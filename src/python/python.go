package python

import (
	"fmt"
	//"github.com/DataDog/go-python3"
	"github.com/sbinet/go-python"
	"os"
)

func Init() {
	err := python.Initialize()
	if err != nil {
		fmt.Println("failed to initialize python: %s", err)
		os.Exit(1)
	}
	fmt.Println("initialized python")
}

func Cleanup() {
	err := python.Finalize()
	if err != nil {
		fmt.Println("failed to finalize python: %s", err)
		os.Exit(1)
	}
	fmt.Println("finalized python")
}
