package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	pb "github.com/sago35/grpcbuild"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	dir string
}

func (s *server) Init(ctx context.Context, in *pb.InitRequest) (*pb.InitResponse, error) {
	//log.Printf("Init(%#v)\n", in)
	s.dir = in.GetDir()

	fis, err := ioutil.ReadDir(s.dir)
	if os.IsNotExist(err) {
		// skip
	} else if err != nil {
		return nil, err
	} else {
		for _, f := range fis {
			os.RemoveAll(filepath.Join(s.dir, f.Name()))
		}
	}

	os.Mkdir(s.dir, 0644)
	return &pb.InitResponse{}, nil
}

func (s *server) Send(ctx context.Context, in *pb.SendRequest) (*pb.SendResponse, error) {
	//log.Printf("Send(%#v)\n", in)
	err := pb.RetrieveFiles(s.dir, in.Files)
	if err != nil {
		return nil, err
	}
	return &pb.SendResponse{}, nil
}

func (s *server) Exec(ctx context.Context, in *pb.ExecRequest) (*pb.ExecResponse, error) {
	log.Printf("Exec(%#v)\n", in)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	for _, c := range in.GetCmds() {
		so := new(bytes.Buffer)
		se := new(bytes.Buffer)
		mwso := io.MultiWriter(so, os.Stdout)
		mwse := io.MultiWriter(se, os.Stderr)

		cmd := exec.CommandContext(ctx, c.GetPath(), c.GetArgs()[1:]...)
		cmd.Dir = filepath.Join(s.dir, c.GetDir())
		cmd.Stdout = mwso
		cmd.Stderr = mwse
		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("%s : %s", err.Error(), se.String())
		}

		stdout.Write(so.Bytes())
		stderr.Write(se.Bytes())
	}

	f, err := pb.StoreFiles(s.dir, in.GetFiles()...)
	if err != nil {
		return nil, fmt.Errorf(stdout.String() + stderr.String() + err.Error())
	}

	return &pb.ExecResponse{
		Files:    f,
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		ExitCode: 0,
	}, nil
}

var (
	port = flag.Int("port", 12345, "port")
)

func run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.MaxMsgSize(0x7FFFFFFF))
	pb.RegisterGrpcBuildServer(s, &server{})
	reflection.Register(s)

	log.Printf("GrpcBuildServer start")

	return s.Serve(lis)
}

func main() {
	flag.Parse()

	err := run(*port)
	if err != nil {
		log.Fatal(err)
	}
}
