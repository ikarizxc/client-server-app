#!/bin/bash
sudo dpkg -i client.deb

chmod +x deb_pakage_client/bin/client

deb_pakage_client/bin/client &
client_pid=$!

sleep 1

if ps -p $client_pid > /dev/null; then
  echo "Client successfully started!"
  kill $client_pid 
  exit 0
else
  echo "Something went wrong while starting client!"
  exit 1
fi