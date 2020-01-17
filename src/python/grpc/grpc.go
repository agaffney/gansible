package grpc

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Grpc struct {
}

func New() *Grpc {
	g := &Grpc{}
	return g
}

func (g *Grpc) Start() {
	dir, _ := ioutil.TempDir("", "tmp.gansible")
	zipContent, _ := base64.StdEncoding.DecodeString(pyGrpcZipContent)
	_ = ioutil.WriteFile(path.Join(dir, "python_grpc.zip"), zipContent, 0400)
	pythonCmd := exec.Command("python", path.Join(dir, "python_grpc.zip"))
	output, err := pythonCmd.CombinedOutput()
	os.RemoveAll(dir)
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Unexpected exit of python process with the following output:\n\n%s\n", output)
		} else {
			fmt.Printf("Unexpected error running python process: %s\n", err)
		}
		os.Exit(1)
	}
}
