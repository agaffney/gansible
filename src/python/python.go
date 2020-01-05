package python

import (
	"fmt"
	// I'm using the py2.7 library because it will build easily against the
	// python libs provided by my Ubuntu 16.04 system. The v3 version of the
	// python library requires py3.7, and only py3.5 is provided via apt.
	//"github.com/DataDog/go-python3"
	"github.com/sbinet/go-python"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const PY_MODULE_PREFIX = `gansible_`

var origSysPath string

func Init() {
	err := python.Initialize()
	if err != nil {
		fmt.Println("failed to initialize python: %s", err)
		os.Exit(1)
	}
	origSysPath = getSysPath()
	helperModule := loadHelperModule()
	testObj := CreateClassInstance(helperModule, `DummyHelper`, nil)
	fmt.Printf("testObj = %#v, %s\n", testObj, testObj.Repr())
	fmt.Println("initialized python")
}

func getSysPath() string {
	sysModule := getModule("sys")
	sysPath := getAttr(sysModule, "path")
	tmp := make([]string, 0)
	for _, val := range ToGoType(sysPath).([]interface{}) {
		tmp = append(tmp, val.(string))
	}
	return strings.Join(tmp, ":")
}

func setSysPath(path string) {
	python.PySys_SetPath(path)
}

func sysPathPrepend(path string) {
	sysModule := getModule("sys")
	sysPath := getAttr(sysModule, "path")
	python.PyList_Insert(sysPath, 0, String(path))
}

func sysPathReset() {
	setSysPath(origSysPath)
}

func buildModuleName(name string) string {
	return fmt.Sprintf("%s%s", PY_MODULE_PREFIX, name)
}

func importModule(moduleName string) *python.PyObject {
	pyModule := python.PyImport_ImportModule(moduleName)
	if pyModule == nil {
		printException()
		log.Fatal(fmt.Sprintf("failed to import python module '%s'", moduleName))
	}
	return pyModule
}

func importModuleFromString(moduleName string, moduleSrc string) *python.PyObject {
	tmpDir, err := ioutil.TempDir("", "tmp.gansible")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create temp dir: %s", err))
	}
	defer os.RemoveAll(tmpDir)
	err = ioutil.WriteFile(path.Join(tmpDir, fmt.Sprintf("%s.py", moduleName)), []byte(moduleSrc), 0644)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to write temp file: %s", err))
	}
	sysPathPrepend(tmpDir)
	module := importModule(moduleName)
	sysPathReset()
	return module
}

func getModule(moduleName string) *python.PyObject {
	d := python.PyImport_GetModuleDict()
	module := python.PyDict_GetItemString(d, moduleName)
	return module
}

func getAttr(obj *python.PyObject, key string) *python.PyObject {
	return obj.GetAttrString(key)
}

func printException() {
	python.PyErr_Print()
}

func Cleanup() {
	err := python.Finalize()
	if err != nil {
		fmt.Println("failed to finalize python: %s", err)
		os.Exit(1)
	}
	fmt.Println("finalized python")
}

func CreateClassInstance(module *python.PyObject, className string, args []interface{}) *python.PyObject {
	pyClass := getAttr(module, className)
	pyObj := pyClass.Call(Tuple(args), Dict(nil))
	pyClass.Clear()
	return pyObj
}
