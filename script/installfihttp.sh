#!/bin/bash

sudo cp ../fihttp.conf /etc
sudo chmod +x ../com.dima.fihttp

if [ "$1" == "first" ]; then
  sudo cp ./fihttpd.sh /etc/init.d
  sudo chmod +x /etc/init.d/fihttp.sh

  sudo update-rc.d fihttpd.sh defaults
  sudo update-rc.d fihttpd.sh enable
fi
