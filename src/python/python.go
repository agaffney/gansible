package python

import (
	"fmt"
	"github.com/agaffney/gansible/python/grpc"
)

func Init() {
	g := grpc.New()
	g.Start()
	//helperModule := loadHelperModule()
	//testObj := CreateClassInstance(helperModule, `DummyHelper`, nil)
	//fmt.Printf("testObj = %#v, %s\n", testObj, testObj.Repr())
	fmt.Println("initialized python")
}
