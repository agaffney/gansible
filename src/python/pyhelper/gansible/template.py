import json

from gansible.grpc_gen import template_pb2, template_pb2_grpc
from gansible.util import exception_wrapper

from ansible.errors import AnsibleUndefinedVariable
from ansible.parsing.dataloader import DataLoader
from ansible.template import Templar

class TemplateServicer(template_pb2_grpc.TemplateServicer):

    @classmethod
    def add_to_server(cls, server):
        return template_pb2_grpc.add_TemplateServicer_to_server(cls(), server)

    @exception_wrapper
    def Render(self, request, context):
        loader = DataLoader()
        templar = Templar(loader)
        try:
            result = templar.template(request.template)
        except Exception as e:
            error_type = template_pb2.ErrorType.OTHER
            if isinstance(e, AnsibleUndefinedVariable):
                error_type = template_pb2.ErrorType.UNDEFINED
            return template_pb2.TemplateResponse(errorType=error_type, error=str(e))
        return template_pb2.TemplateResponse(result=result)
