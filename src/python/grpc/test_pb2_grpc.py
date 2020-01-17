# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import test_pb2 as test__pb2


class TestStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Ping = channel.unary_unary(
        '/test.Test/Ping',
        request_serializer=test__pb2.PingRequest.SerializeToString,
        response_deserializer=test__pb2.PingResponse.FromString,
        )


class TestServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Ping(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_TestServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Ping': grpc.unary_unary_rpc_method_handler(
          servicer.Ping,
          request_deserializer=test__pb2.PingRequest.FromString,
          response_serializer=test__pb2.PingResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'test.Test', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
