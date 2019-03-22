#!/bin/bash

set -e

HADOOP_ADDRESS="0.0.0.0:9000"

while getopts "a:" opt; do
  case ${opt} in
    a ) 
        HADOOP_ADDRESS=${OPTARG}
      ;;
  esac
done

set -x

echo "hadoop address set to: ${HADOOP_ADDRESS}"

docker build -t go-webhdfs --build-arg DOCKER_HADOOP_ADDRESS=${HADOOP_ADDRESS} . 

function get_unused_port() {
  for port in $(seq 4444 65000);
  do
    [ $? -eq 1 ] && echo "$port" && break;
  done
}
PORT="$(get_unused_port)"

docker run -p ${PORT}:8000 --detach go-webhdfs  
echo "go-webhdfs is running on port: ${PORT}"