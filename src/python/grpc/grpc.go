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
	// Create pipe for communicating listening port
	portReader, portWriter, _ := os.Pipe()
	g.pythonCmd.ExtraFiles = []*os.File{portWriter}
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
	g.pythonCmd.Stdout = os.Stdout
	g.pythonCmd.Stderr = os.Stderr
	err := g.pythonCmd.Start()
	if err != nil {
		fmt.Printf("failed to start command: %s\n", err)
		os.Exit(1)
	}
	var pythonPort string
	portReaderBuf := bufio.NewReader(portReader)
	portLine, _ := portReaderBuf.ReadString('\n')
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
				fmt.Printf("Unexpected exit of python process\n")
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

	templateClient := grpc_gen.NewTemplateClient(conn)
	ret5, err := templateClient.Render(context.Background(), &grpc_gen.TemplateRequest{Template: `foo {{ 'bar' }} baz, {{ {'foo':'bar'} | dict2items }}`})
	fmt.Printf("ret5 = %s, err = %#v\n", ret5.String(), err)

	actionClient := grpc_gen.NewActionClient(conn)
	ret6, err := actionClient.Run(context.Background(), &grpc_gen.RunRequest{Action: `template`})
	fmt.Printf("ret6 = %s, err = %#v\n", ret6.String(), err)

	syscall.Kill(-g.pythonCmd.Process.Pid, syscall.SIGKILL)
	os.RemoveAll(dir)
}
