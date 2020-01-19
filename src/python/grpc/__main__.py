from __future__ import print_function

import grpc
import logging
import sys
import time

import test_pb2, test_pb2_grpc

from concurrent import futures

class TestServicer(test_pb2_grpc.TestServicer):

    def Ping(self, request, context):
        return test_pb2.PingResponse(pong=True, msg="hi there")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    test_pb2_grpc.add_TestServicer_to_server(TestServicer(), server)
    port = server.add_insecure_port('[::]:0') #50051')
    if port == 0:
        logging.error('failed to bind to port')
        sys.exit(1)
    print_flush('PORT=%d' % port)
    server.start()
    server.wait_for_termination()


# Python buffers output to a pipe by default, so we use a helper function
# to flush stdout after printing
def print_flush(*args, **kwargs):
    print(*args, **kwargs)
    sys.stdout.flush()


if __name__ == '__main__':
    logging.basicConfig()
    serve()
