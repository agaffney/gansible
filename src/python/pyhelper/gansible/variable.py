import ansible.module_utils.six as six
from ansible.module_utils.common._collections_compat import Mapping

from gansible.grpc_gen import variable_pb2

def Value(value):
    kwargs = None
    if isinstance(value, Mapping):
        dictItems = []
        for key in value.keys():
            dictItems.append(variable_pb2.DictEntry(key=Value(key), value=Value(value[key])))
        kwargs = dict(type=variable_pb2.ValueType.DICT, dictValue=variable_pb2.Dict(items=dictItems))
    elif isinstance(value, list):
        listItems = []
        for item in value:
            listItems.append(Value(item))
        kwargs = dict(type=variable_pb2.ValueType.LIST, listValue=variable_pb2.List(items=listItems))
    elif isinstance(value, six.string_types):
        kwargs = dict(type=variable_pb2.ValueType.STRING, stringValue=value)
    elif isinstance(value, bool):
        kwargs = dict(type=variable_pb2.ValueType.BOOL, boolValue=value)
    elif isinstance(value, float):
        kwargs = dict(type=variable_pb2.ValueType.FLOAT, floatValue=value)
    elif isinstance(value, int):
        kwargs = dict(type=variable_pb2.ValueType.INT, intValue=value)
    else:
        raise Exception("unknown type for: %s" % repr(value))
    return variable_pb2.Value(**kwargs)
