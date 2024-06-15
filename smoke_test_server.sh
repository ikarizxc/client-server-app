#!/bin/bash
sudo dpkg -i server.deb

chmod +x deb_pakage_server/bin/server

deb_pakage_server/bin/server &
server_pid=$!

timeout 10 nc -z 127.0.0.1 8080 > /dev/null 2>&1
if [ $? -eq 0 ]; then
  echo "Server successfully started!"
  kill $server_pid 
  sudo dpkg -r server
  exit 0
else
  echo "Something went wrong while starting server!"
  kill $server_pid
  sudo dpkg -r server
  exit 1
fi

