#!/bin/bash

set -e

# set working directory to script location
cd "${0%/*}"
echo "using working directory $(pwd)"

# the ip-address:port combination that we will connect to Hadoop
HADOOP_ADDRESS="0.0.0.0:9000"

# the port that will be exposed for accessing the container. if this isn't specified a random unused port is selected
PORT=""

while getopts ":a:p:" opt; do
  case ${opt} in
    a ) 
        HADOOP_ADDRESS=${OPTARG}
      ;;
    p ) 
        PORT=${OPTARG}
    ;;
  esac
done

echo "hadoop address set to: ${HADOOP_ADDRESS}"

docker build -t go-webhdfs . 

function get_unused_port() {
  LPORT=32768;
  UPORT=60999;
  while true; do
      MPORT=$[$LPORT + ($RANDOM % $UPORT)];
      (echo "" >/dev/tcp/127.0.0.1/${MPORT}) >/dev/null 2>&1
      if [ $? -ne 0 ]; then
          echo $MPORT;
          return 0;
      fi
  done
}

if [[ -z "${PORT}" ]]
then
      PORT="$(get_unused_port)"
fi

docker run -e HADOOP_ADDRESS=${HADOOP_ADDRESS} -p ${PORT}:8000 --detach go-webhdfs
echo "starting go-webhdfs"
echo "go-webhdfs is running on port: ${PORT}"