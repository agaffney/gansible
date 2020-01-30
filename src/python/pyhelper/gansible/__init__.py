from __future__ import print_function

import grpc
import logging
import sys
import time

from gansible.grpc_gen import test_pb2, test_pb2_grpc
from gansible.util import print_flush
from concurrent import futures

class TestServicer(test_pb2_grpc.TestServicer):

    def Ping(self, request, context):
        return test_pb2.PingResponse(pong=True, msg="hi there")


class Gansible(object):

    def __init__(self):
        logging.basicConfig()

    def start(self):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        test_pb2_grpc.add_TestServicer_to_server(TestServicer(), server)
        port = server.add_insecure_port('[::]:0')
        if port == 0:
            logging.error('failed to bind to port')
            sys.exit(1)
        print_flush('PORT=%d' % port)
        server.start()
        server.wait_for_termination()
