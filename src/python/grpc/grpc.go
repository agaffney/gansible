package grpc

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	grpc_gen "github.com/agaffney/gansible/python/grpc/generated"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

type Grpc struct {
	pythonCmd *exec.Cmd
	port      int
}

func New() *Grpc {
	g := &Grpc{}
	return g
}

func (g *Grpc) Start() {
	dir, _ := ioutil.TempDir("", "tmp.gansible")
	zipContent, _ := base64.StdEncoding.DecodeString(pyGrpcZipContent)
	zipFile := path.Join(dir, "python_grpc.zip")
	_ = ioutil.WriteFile(zipFile, zipContent, 0400)
	g.pythonCmd = exec.Command("python", zipFile)
	// Put child process in its own process group to make it easier to kill with
	/// all its children
	g.pythonCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// Setup signal handler
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	signalReceived := false
	go func() {
		<-s
		signalReceived = true
		fmt.Println("Killing python process...")
		syscall.Kill(-g.pythonCmd.Process.Pid, syscall.SIGKILL)
		os.Exit(0)
	}()
	stdout, err := g.pythonCmd.StdoutPipe()
	//g.pythonCmd.Stderr = os.Stderr
	if err != nil {
		fmt.Printf("failed to get stdout pipe: %s\n", err)
		os.Exit(1)
	}
	err = g.pythonCmd.Start()
	if err != nil {
		fmt.Printf("failed to start command: %s\n", err)
		os.Exit(1)
	}
	var pythonPort string
	stdoutBuf := bufio.NewReader(stdout)
	portLine, _ := stdoutBuf.ReadString('\n')
	if strings.HasPrefix(portLine, "PORT=") {
		pythonPort = portLine[5 : len(portLine)-1]
	} else {
		fmt.Printf("Malformed output: %s", portLine)
		os.Exit(1)
	}
	go func() {
		err := g.pythonCmd.Wait()
		if err != nil {
			if signalReceived {
				return
			}
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("exitErr = %s\n", exitErr.Error())
				fmt.Printf("Unexpected exit of python process with the following output:\n\n%%s\n")
			} else {
				fmt.Printf("Unexpected error running python process: %s\n", err)
			}
			os.Exit(1)
		}
	}()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", pythonPort), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to dial: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := grpc_gen.NewTestClient(conn)
	ret, err := client.Ping(context.Background(), &grpc_gen.PingRequest{Ping: true, Msg: "anyone home?"})
	fmt.Printf("ret = %s, err = %#v\n", ret.String(), err)
	if err != nil {
		fmt.Printf("failed to Ping: %s\n", err)
		os.Exit(1)
	}
	inventoryClient := grpc_gen.NewInventoryClient(conn)
	ret1, err := inventoryClient.Load(context.Background(), &grpc_gen.LoadRequest{Sources: []string{"/etc/ansible/hosts"}})
	fmt.Printf("ret1 = %s, err = %#v\n", ret1.String(), err)
	ret2, err := inventoryClient.ListHosts(context.Background(), &grpc_gen.ListHostsRequest{Pattern: "all"})
	fmt.Printf("ret2 = %s, err = %#v\n", ret2.String(), err)
	callbackClient := grpc_gen.NewCallbackClient(conn)
	ret3, err := callbackClient.Init(context.Background(), &grpc_gen.Empty{})
	fmt.Printf("ret3 = %s, err = %#v\n", ret3.String(), err)
	ret4, err := callbackClient.RunnerOnOk(context.Background(), &grpc_gen.TaskResult{})
	fmt.Printf("ret4 = %s, err = %#v\n", ret4.String(), err)
	syscall.Kill(-g.pythonCmd.Process.Pid, syscall.SIGKILL)
	os.RemoveAll(dir)
}
