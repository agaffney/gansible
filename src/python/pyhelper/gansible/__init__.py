# Suppress warnings about setuptools on py2 and old versions of cryptography
import warnings
warnings.filterwarnings("ignore", module='pkg_resources.py2_warn')
warnings.filterwarnings("ignore", module='requests')

import grpc
import logging
import os
import sys

from gansible.callback import CallbackServicer
from gansible.inventory import InventoryServicer
from gansible.template import TemplateServicer
from concurrent import futures

# The Go parent process opens up a pipe on FD 3 for this process to communicate
# its listening port for gRPC
PORT_PIPE_FD = 3

class Gansible(object):

    def __init__(self):
        logging.basicConfig()

    def start(self):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        InventoryServicer.add_to_server(server)
        CallbackServicer.add_to_server(server)
        TemplateServicer.add_to_server(server)
        port = server.add_insecure_port('[::]:0')
        if port == 0:
            logging.error('failed to bind to port')
            sys.exit(1)
        os.write(3, "PORT=%d\n" % port)
        server.start()
        server.wait_for_termination()
