from gansible.grpc_gen import inventory_pb2, inventory_pb2_grpc

from ansible.inventory.manager import InventoryManager
from ansible.parsing.dataloader import DataLoader

class InventoryServicer(inventory_pb2_grpc.InventoryServicer):

    _inventory = None

    @classmethod
    def add_to_server(cls, server):
        return inventory_pb2_grpc.add_InventoryServicer_to_server(cls(), server)

    def Load(self, request, context):
        if self._inventory is None:
            loader = DataLoader()
            self._inventory = InventoryManager(loader=loader, sources=request.sources)
        return inventory_pb2.BoolResponse(ret=True)

    def ListHosts(self, request, context):
        hosts = self._inventory.list_hosts(request.pattern)
        hosts = [h.get_name() for h in hosts]
        return inventory_pb2.ListHostsResponse(hosts=hosts)
