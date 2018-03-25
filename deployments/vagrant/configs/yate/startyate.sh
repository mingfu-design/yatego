#! /bin/bash

export LD_LIBRARY_PATH=/opt/yate/lib
#exec /opt/yate/bin/yate $@
exec /opt/yate/bin/yate -ds -r -p /var/run/yate.pid -l /var/log/yate/messages -vvv
