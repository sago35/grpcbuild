package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/sago35/grpcbuild"
	"github.com/sago35/limichan"

	"google.golang.org/grpc"
)

type worker struct {
	name    string
	conn    *grpc.ClientConn
	client  pb.GrpcBuildClient
	ctx     context.Context
	cancel  context.CancelFunc
	addr    string
	port    int
	threads int
}

func newWorker(addr string, port int, threads int) (*worker, error) {
	w := &worker{
		name:    addr,
		addr:    addr,
		port:    port,
		threads: threads,
		ctx:     context.Background(),
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", w.addr, w.port), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(0x7FFFFFFF)))
	if err != nil {
		return nil, err
	}
	w.conn = conn

	c := pb.NewGrpcBuildClient(conn)
	w.client = c

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	workdir := hostname

	_, err = w.client.Init(context.Background(), &pb.InitRequest{Dir: workdir})
	if err != nil {
		return nil, err
	}

	input, err := makeFileList()
	if err != nil {
		return nil, err
	}

	for _, f := range input {
		sr, err := pb.MkSendRequest(f)
		if err != nil {
			return nil, err
		}

		_, err = w.client.Send(context.Background(), sr)
		if err != nil {
			return nil, err
		}
	}

	return w, nil
}

func (w *worker) Clone(name string) *worker {
	ret := &worker{
		name:    name,
		conn:    w.conn,
		client:  w.client,
		ctx:     w.ctx,
		cancel:  w.cancel,
		addr:    w.addr,
		port:    w.port,
		threads: w.threads,
	}
	return ret
}

func (w *worker) Do(ctx context.Context, lj limichan.Job) error {
	job, ok := lj.(*job)
	if !ok {
		log.Fatal("type error")
	}
	defer close(job.ch)

	start := time.Now()
	job.ch <- fmt.Sprintf("# %-16s %-16s %d", w.name, job.name, time.Now().UnixNano()/1000000)

	for _, f := range job.depFile {
		sr, err := pb.MkSendRequest(f)
		if err != nil {
			return err
		}

		_, err = w.client.Send(ctx, sr)
		if err != nil {
			return err
		}
	}

	er, err := job.GetExecRequest()
	if err != nil {
		return err
	}

	res, err := w.client.Exec(w.ctx, er)
	if err != nil {
		return err
	}

	out := []string{}
	if len(res.GetStdout()) > 0 {
		out = append(out, strings.TrimSpace(string(res.GetStdout())))
	}
	if len(res.GetStderr()) > 0 {
		out = append(out, strings.TrimSpace(string(res.GetStderr())))
	}

	if len(out) > 0 {
		job.ch <- strings.Join(out, "\n")
	}

	err = res.WriteFiles()
	if err != nil {
		return err
	}

	job.ch <- fmt.Sprintf("# %-16s %-16s %d (%dms)", w.name, job.name, time.Now().UnixNano()/1000000, time.Since(start).Nanoseconds()/1000000)

	return nil
}
