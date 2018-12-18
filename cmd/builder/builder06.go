package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	pb "github.com/sago35/grpcbuild"
	"github.com/sago35/limichan"
	"github.com/sago35/ochan"

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
	job, _ := lj.(*job)
	defer close(job.ch)

	cmd := job.cmd

	// もし依存ファイルがあるなら送る
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

	// 実際に処理を行う
	er := &pb.ExecRequest{
		Files: job.outFile,
	}
	er.Cmds = append(er.Cmds, &pb.Cmd{
		Path: cmd.Path,
		Args: cmd.Args[1:],
		Dir:  cmd.Dir,
	})

	res, err := w.client.Exec(w.ctx, er)
	if err != nil {
		return err
	}

	if len(res.GetStdout()) > 0 {
		job.ch <- strings.TrimSpace(string(res.GetStdout()))
	}

	// リモートで処理したファイルをローカルに保存する
	err = res.WriteFiles()
	if err != nil {
		return err
	}

	return nil
}

func makeFileList() ([]string, error) {
	ret := []string{}

	b, err := exec.Command(`git`, `ls-files`, `testdata`).Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, nil
}

func build06(cmds []*exec.Cmd) error {
	outCh := make(chan string, 10000)
	defer close(outCh)

	oc := ochan.NewOchan(outCh, 100)
	go func() {
		for ch := range outCh {
			fmt.Println(ch)
		}
	}()

	l, _ := limichan.New(context.Background())
	w, err := newWorker(`127.0.0.1`, 12345, *threads)
	if err != nil {
		return err
	}
	for i := 0; i < *threads; i++ {
		l.AddWorker(w.Clone(fmt.Sprintf("grpcWorker(%d)", i)))
	}

	for _, cmd := range cmds {
		cmd := cmd

		if cmd.Path != dummyCc {
			// コンパイラではない時は、直前までのコンパイルが終わるのを待つ
			fmt.Println(1, cmd.Path)
			oc.Wait()
		}

		ch := oc.GetCh()
		j := &job{
			cmd: cmd,
			ch:  ch,
		}

		l.Do(j)
	}

	return nil
}
