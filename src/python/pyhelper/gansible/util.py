from __future__ import print_function

import functools
import sys
import traceback

# Python buffers output to a pipe by default, so we use a helper function
# to flush stdout after printing
def print_flush(*args, **kwargs):
    print(*args, **kwargs)
    sys.stdout.flush()

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
