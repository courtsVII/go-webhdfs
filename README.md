# go-webhdfs

web service for interacting with the Hadoop FileSystem

## setup

the service can be run as is or run in a container. deploy the service in a node that can communicate with your Hadoop namenode then run commands against your Hadoop infrastructure using the exposed HTTP endpoints

### run as is

optionally set the GO_WEBHDFS_ADDRESS (ip:port combination to expose your service at, default is 0.0.0.0:8000) and HADOOP_ADDRESS (ip:port combination to connect to your Hadoop namenode, default is 0.0.0.0:9000) environment variables in your shell environment then execute `go run controller/controller.go` to expose the service and connect to your Hadoop cluster

### run as docker container

this can be done manually or by executing either the run.sh or run.py script

#### run.sh

##### run.sh can take the following parameters

`-a` | the address of your Hadoop namenode. defaults to 0.0.0.0:9000
`-p` | the port to expose the service on in your machine. If not specified random unassigned port is selected and printed to stdout

#### run.py

##### run.py can take the following parameters

`--env`  | the environment variables passed to the Docker file. (HADOOP_ADDRESS environment variable which specifies the address of your Hadoop namenode should be provided here)
`--port` | the port to expose the service on in your machine. If not specified random unassigned port is selected and printed to stdout
`--detach` | whether or not to run the container in detached mode
`--tag` | tag for your container

## endpoints

### `mv`

move file from source to destination

`src` | HDFS source path
`dst` | HDFS destination path

### `cp`

copy file from source to destination

`src` | HDFS source path
`dst` | HDFS destination path

### `mkdir`

create a directory at specified location

`path` | HDFS path to create folder at

### `rm`

remove file at specific location

#### `params`

`path` | HDFS path to file to remove

### `rmall`

remove all files at specific location recursively

`path` | HDFS path to file to remove

### `ls`

ls files at specific location

`path`| HDFS path to location to ls
`recursive` | boolean value which specifies whether or not to do a recursive ls

### `createfile`

create empty file at a specified destination

`path` | HDFS path to create the empty file at

### `writefile`

copy a file from the client into a file in HDFS

`path` | HDFS path to create the replica file at
`file` | multipart/form-data which represents the source file you wish to upload (e.g. curl -F 'file=@example.txt' -L  http://go-webhdfs-service:8000/v1/hdfs/writefile\?path\=/my/hdfs/location/example.txt)

### `write`

copy arbritary data from the HTTP request body into a file in HDFS

`path` | HDFS path to create the file at

### `readfile`

read or copy data from a file  in HDFS

`path` | HDFS path to the file to be read
(e.g. curl -L  http://go-webhdfs-service:8000/v1/hdfs/readfile\?path\=/my/hdfs/location/example.txt > copyofexample.txt)


### `getcontentsummary`

get contents summary of a file  in HDFS

`path` | HDFS path to the file from which to getcontentsummary
  
### `chmod`

chmod a file  in HDFS

`path` | HDFS path to the file to chmod
`mask` | mask to apply to file

### `chown`

chown a file  in HDFS

`path` | HDFS path to the file to chown
`user` | user to own file
`group` | group to own file
