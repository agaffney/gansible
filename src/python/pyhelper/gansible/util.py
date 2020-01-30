from __future__ import print_function

import sys

# Python buffers output to a pipe by default, so we use a helper function
# to flush stdout after printing
def print_flush(*args, **kwargs):
    print(*args, **kwargs)
    sys.stdout.flush()
