package python

import (
	"fmt"
	//"github.com/DataDog/go-python3"
	"github.com/sbinet/go-python"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var pyHelperSrc = `
from __future__ import print_function

import sys

class DummyHelper(object):
    def __init__(self):
#        print('args: %s' % sys.argv)
#        print(dir(sys))
        print('instantiated DummyHelper')
`

func Init() {
	err := python.Initialize()
	if err != nil {
		fmt.Println("failed to initialize python: %s", err)
		os.Exit(1)
	}
	if foo := python.PyRun_SimpleString(pyHelperSrc); foo < 0 {
		fmt.Println("failed to eval python helper code")
		os.Exit(1)
	}
	helperModule := importHelperModule()
	fmt.Printf("helperModule = %#v\n", helperModule)
	testObj := CreateClassInstance(helperModule, `DummyHelper`, nil)
	fmt.Printf("testObj = %#v\n", testObj)
	fmt.Println("initialized python")
}

// TODO: refactor to use PyRun_SimpleString() to eval, PyImport_GetModuleDict(), and
// PyDict_GetItemString(moduleDict, '__main__')
func importHelperModule() *python.PyObject {
	tmpDir, err := ioutil.TempDir("", "tmp.gansible")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tmpDir = %s\n", tmpDir)
	defer os.RemoveAll(tmpDir)
	err = ioutil.WriteFile(path.Join(tmpDir, "helpers.py"), []byte(pyHelperSrc), 0644)
	if err != nil {
		log.Fatal(fmt.Printf("failed to write temp file for python helpers: %s", err))
	}
	python.PySys_SetPath(tmpDir)
	pyModule := python.PyImport_ImportModule("helpers")
	if pyModule == nil {
		log.Fatal(fmt.Printf("failed to import python helpers module"))
	}
	return pyModule
}

func Cleanup() {
	err := python.Finalize()
	if err != nil {
		fmt.Println("failed to finalize python: %s", err)
		os.Exit(1)
	}
	fmt.Println("finalized python")
}

func CreateClassInstance(module *python.PyObject, className string, args []*python.PyObject) *python.PyObject {
	pyClass := module.GetAttrString(className)
	//fmt.Printf("pyClass = %#v, %s, %s\n", pyClass, pyClass.Repr(), pyClass.PyObject_Dir())
	pyObj := pyClass.Call(python.PyTuple_New(0), python.PyDict_New()) //args)
	fmt.Printf("pyObj = %#v\n", pyObj)
	python.PyErr_Print()
	pyClass.DecRef()
	return pyObj
}
