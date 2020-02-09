import json

from gansible.grpc_gen import callback_pb2, callback_pb2_grpc
from gansible.util import exception_wrapper

from ansible.plugins.loader import callback_loader
from ansible.executor.task_queue_manager import TaskQueueManager
from ansible.playbook.task import Task
from ansible.executor.task_result import TaskResult
from ansible.inventory.host import Host

class CallbackServicer(callback_pb2_grpc.CallbackServicer):

    # Vars used by the load_callbacks() method from TQM
    _callbacks_loaded = False
    _callback_plugins = None
    _run_additional_callbacks = False
    _run_tree = False
    _stdout_callback = None

    @classmethod
    def add_to_server(cls, server):
        return callback_pb2_grpc.add_CallbackServicer_to_server(cls(), server)

    @exception_wrapper
    def Init(self, request, context):
        # Call function to load callbacks from TQM using ourself as the class instance
        _load_callbacks = TaskQueueManager.__dict__['load_callbacks']
        _load_callbacks(self)
        return callback_pb2.Empty()

    def _call_plugins(self, fn_name, *args, **kwargs):
        for callback_plugin in self._callback_plugins:
            fn = getattr(callback_plugin, fn_name)
            fn(*args, **kwargs)

    def _create_task_obj(self, request):
        pass

    def _create_result_obj(self, request):
        return TaskResult(Host("testhost"), Task(), dict(foo="bar", baz=123))

    def _create_play_obj(self, request):
        pass

    def _create_playbook_obj(self, request):
        pass

    @exception_wrapper
    def RunnerOnOk(self, request, context):
        self._call_plugins('v2_runner_on_ok', self._create_result_obj(request))
        return callback_pb2.Empty()
