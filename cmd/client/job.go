package main

import (
	"os/exec"

	pb "github.com/sago35/grpcbuild"
)

type job struct {
	Cmds    []exec.Cmd
	name    string
	ch      chan string
	outFile []string
	depFile []string
}

func newJob() *job {
	return &job{}
}

func (j *job) GetExecRequest() (*pb.ExecRequest, error) {
	ret := &pb.ExecRequest{}
	ret.Files = j.outFile

	for _, c := range j.Cmds {
		ret.Cmds = append(ret.Cmds, &pb.Cmd{
			Path: c.Path,
			Args: c.Args[1:],
			Env:  []string{},
			Dir:  c.Dir,
		})
	}

	return ret, nil
}
