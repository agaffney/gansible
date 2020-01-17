from __future__ import print_function

import concurrent.futures
import grpc
import time

import test_pb2, test_pb2_grpc

class TestServicer(test_pb2_grpc.TestServicer):

    def Ping(self, request, context):
        return test_pb2.PingResponse(pong=True, msg="hi there")


def serve():
    server = grpc.server(concurrent.futures.ThreadPoolExecutor(max_workers=10))
    test_pb2_grpc.add_TestServicer_to_server(TestServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()

serve()

while True:
    time.sleep(60)
