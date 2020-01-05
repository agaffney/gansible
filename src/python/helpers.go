package python

import (
	"github.com/sbinet/go-python"
)

var pyHelperSrc = `
from __future__ import print_function

class DummyHelper(object):
    def __init__(self, *args, **kwargs):
        print('self = %s, args = %s, kwargs = %s' % (self, args, kwargs))
        print('instantiated DummyHelper')
`

func loadHelperModule() *python.PyObject {
	helperModule := importModuleFromString(buildModuleName(`helpers`), pyHelperSrc)
	return helperModule
}
