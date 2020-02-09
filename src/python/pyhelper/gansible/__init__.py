import grpc
import logging
import sys

from gansible.callback import CallbackServicer
from gansible.inventory import InventoryServicer
from gansible.test import TestServicer
from gansible.util import print_flush
from concurrent import futures

class Gansible(object):

    def __init__(self):
        logging.basicConfig()

    def start(self):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        TestServicer.add_to_server(server)
        InventoryServicer.add_to_server(server)
        CallbackServicer.add_to_server(server)
        port = server.add_insecure_port('[::]:0')
        if port == 0:
            logging.error('failed to bind to port')
            sys.exit(1)
        print_flush('PORT=%d' % port)
        server.start()
        server.wait_for_termination()
