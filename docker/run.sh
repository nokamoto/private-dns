#!/bin/sh

# Start the first process
dnsmasq --log-facility=- &
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start dnsmasq: $status"
  exit $status
fi

# Start the second process
/root/go/bin/dnssupervisor -hostsfile $HOST_FILE -port $GRPC_PORT -apikey $API_KEY &
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start dnssupervisor: $status"
  exit $status
fi

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container will exit with an error
# if it detects that either of the processes has exited.
# Otherwise it will loop forever, waking up every 60 seconds

while /bin/true; do
  ps aux | grep dnsmasq | grep -q -v grep
  PROCESS_1_STATUS=$?
  ps aux |grep dnssupervisor |grep -q -v grep
  PROCESS_2_STATUS=$?
  # If the greps above find anything, they will exit with 0 status
  # If they are not both 0, then something is wrong
  if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
    echo "One of the processes has already exited."
    exit -1
  fi
  kill -1 `pgrep dnsmasq`
  sleep 60
done
