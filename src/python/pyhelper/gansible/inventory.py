from gansible.grpc_gen import inventory_pb2, inventory_pb2_grpc

class InventoryServicer(inventory_pb2_grpc.InventoryServicer):

    @classmethod
    def add_to_server(cls, server):
        return inventory_pb2_grpc.add_InventoryServicer_to_server(cls(), server)

    def Parse(self, request, context):
        return inventory_pb2.BoolResponse(ret=True)
