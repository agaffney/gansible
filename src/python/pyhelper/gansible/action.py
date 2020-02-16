from gansible.grpc_gen import action_pb2, action_pb2_grpc
from gansible.util import exception_wrapper
from gansible.variable import Value

from ansible.plugins.loader import action_loader

class ActionServicer(action_pb2_grpc.ActionServicer):

    @classmethod
    def add_to_server(cls, server):
        return action_pb2_grpc.add_ActionServicer_to_server(cls(), server)

    @exception_wrapper
    def Init(self, request, context):
        # TODO: instantiate all the action plugins
        return action_pb2.InitEmpty()

    @exception_wrapper
    def Run(self, request, context):
        # TODO: run the named action plugin
        return action_pb2.RunResponse(result=Value({'foo': 'bar', 'baz': 123}))
