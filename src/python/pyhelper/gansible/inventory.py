import json

from gansible.grpc_gen import inventory_pb2, inventory_pb2_grpc
from gansible.util import exception_wrapper

from ansible.inventory.manager import InventoryManager
from ansible.parsing.dataloader import DataLoader

class InventoryServicer(inventory_pb2_grpc.InventoryServicer):

    _inventory = None

    @classmethod
    def add_to_server(cls, server):
        return inventory_pb2_grpc.add_InventoryServicer_to_server(cls(), server)

    @exception_wrapper
    def Load(self, request, context):
        if self._inventory is None:
            loader = DataLoader()
            self._inventory = InventoryManager(loader=loader, sources=request.sources)
        return inventory_pb2.BoolResponse(ret=True)

    @exception_wrapper
    def ListHosts(self, request, context):
        hosts = self._inventory.list_hosts(request.pattern)
        ret = []
        for host in hosts:
            tmp_groups = []
            for group in host.get_groups():
                tmp_groups.append(inventory_pb2.Group(name=group.get_name(), vars=json.dumps(group.get_vars()), hosts=[h.get_name() for h in group.get_hosts()]))
            ret.append(inventory_pb2.Host(name=host.get_name(), groups=tmp_groups, vars=json.dumps(host.get_vars())))
        return inventory_pb2.ListHostsResponse(hosts=ret)
