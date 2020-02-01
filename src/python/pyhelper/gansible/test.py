from gansible.grpc_gen import test_pb2, test_pb2_grpc

class TestServicer(test_pb2_grpc.TestServicer):

    @classmethod
    def add_to_server(cls, server):
        return test_pb2_grpc.add_TestServicer_to_server(cls(), server)

    def Ping(self, request, context):
        return test_pb2.PingResponse(pong=True, msg="hi there")
