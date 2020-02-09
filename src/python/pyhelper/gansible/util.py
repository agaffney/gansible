from __future__ import print_function

import functools
import traceback

def exception_wrapper(func):
    '''
    Wrapper function that catches exceptions and re-throws them with the entire stack
    trace in the exception message. This is necessary because gRPC only returns the
    exception message, which makes debugging harder.
    '''
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except Exception as e:
            msg = traceback.format_exc()
            raise Exception(msg)
    return wrapper
