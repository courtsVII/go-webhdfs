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

MIN_PORT=30000;
MAX_PORT=60000;
RANGE=$(($MAX_PORT-$MIN_PORT+1));
RANDOM_PORT=$RANDOM;
let "RANDOM_PORT %= $RANGE";
PORT=$(($RANDOM_PORT+$MIN_PORT));

docker run -p ${PORT}:8000 go-webhdfs  &
echo "go-webhdfs is running on port: ${PORT}"