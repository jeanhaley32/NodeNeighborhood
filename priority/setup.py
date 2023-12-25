import os

curr_dir=os.getcwd()

lib_dir='/'.join(curr_dir.split('/'))+'/src'
lib_export='export PYTHONPATH=$PYTHONPATH:"{0}"'.format(lib_dir)
print('\n'.join([lib_export]))