#!/bin/bash

### BEGIN INIT INFO
# Provides:	   start
# Required-Start: $remote_fs $syslog
# Required-Stop:  $remote_fs $syslog
# Should-Start:	  $local_fs
# Should-Stop:	  $local_fs
# Default-Start:  2 3 4 5
# Default-Stop:	  0 1 6
# Short-Description: fihttp
# Description:	     fihttp
### END INIT INFO

HOME_DIR=/home/pi
DAEMON_DIR=$HOME_DIR/fihttp
NAME=com.dima.fihttp
DAEMON=$DAEMON_DIR/$NAME
DAEMON_OPT="/etc/fihttp.conf"
DESC="fihttp daemon"

test -f $DAEMON || exit 0
set -e


case "$1" in
  start)
        echo -n "Starting $DESC: "
        start-stop-daemon --start -m -b --pidfile /var/run/$NAME.pid --startas $DAEMON -- $DAEMON_OPT 
        echo "$NAME."
        ;;
  stop)
        echo -n "Stopping $DESC: "
        start-stop-daemon --stop --quiet --oknodo \
            --pidfile /var/run/$NAME.pid
        rm -f /var/run/$NAME.pid
        echo "$NAME."
        ;;
  restart)
        echo -n "Restarting $DESC: "
        start-stop-daemon --stop --quiet --oknodo \
            --pidfile /var/run/$NAME.pid
        rm -f /var/run/$NAME.pid
        start-stop-daemon --start -m -b --pidfile /var/run/$NAME.pid --startas $DAEMON -- $DAEMON_OPT
        echo "$NAME."
esac

exit 0
