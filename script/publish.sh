#!/bin/bash

#DSTHOST=pi@hummerd.dyndns.org
#DSTHOST=dima@192.168.100.204
DSTHOST=pi@192.168.1.51

#DSTPATH=/home/dima/fihttp
DSTPATH=/home/pi/fihttp


ssh $DSTHOST "sudo /etc/init.d/fihttpd.sh stop; sudo rm -rf $DSTPATH; mkdir $DSTPATH"

rsync -rav -e ssh \
--exclude='*.pyc' \
--exclude='*.go' \
--exclude='*.log' \
--exclude='*.md' \
--exclude='*.txt' \
--exclude='/test' \
--exclude='/.git' \
--exclude='**/publish.sh' \
--exclude='**/goprepare.sh' \
--exclude='**/.*' \
--exclude='**/__pycache__' \
../ $DSTHOST:$DSTPATH/

ssh $DSTHOST "sudo /etc/init.d/fihttpd.sh start"
