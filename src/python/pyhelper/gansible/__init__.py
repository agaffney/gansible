# Suppress warnings about setuptools on py2
import warnings
warnings.filterwarnings("ignore", module='pkg_resources.py2_warn')

import grpc
import logging
import sys

from gansible.callback import CallbackServicer
from gansible.inventory import InventoryServicer
from gansible.template import TemplateServicer
from gansible.test import TestServicer
from concurrent import futures

class Gansible(object):

    def __init__(self):
        logging.basicConfig()

    def start(self):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        TestServicer.add_to_server(server)
        InventoryServicer.add_to_server(server)
        CallbackServicer.add_to_server(server)
        TemplateServicer.add_to_server(server)
        port = server.add_insecure_port('[::]:0')
        if port == 0:
            logging.error('failed to bind to port')
            sys.exit(1)
        sys.stderr.write("PORT=%d\n" % port)
        sys.stderr.flush()
        server.start()
        server.wait_for_termination()
