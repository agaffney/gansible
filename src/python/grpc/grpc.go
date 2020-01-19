package grpc

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
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
	if err != nil {
		fmt.Printf("failed to get stdout pipe: %s\n", err)
		os.Exit(1)
	}
	err = g.pythonCmd.Start()
	if err != nil {
		fmt.Printf("failed to start command: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("started command")
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
			if _, ok := err.(*exec.ExitError); ok {
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
	client := NewTestClient(conn)
	for {
		ret, err := client.Ping(context.Background(), &PingRequest{Ping: true, Msg: "anyone home?"})
		fmt.Printf("ret = %s, err = %#v\n", ret.String(), err)
		if err == nil {
			time.Sleep(2 * time.Second)
			continue
		}
		time.Sleep(100 * time.Millisecond)
	}
	syscall.Kill(-g.pythonCmd.Process.Pid, syscall.SIGKILL)
	os.RemoveAll(dir)
}
