#!/usr/bin/env python

import argparse, sys, docker, os

def get_free_port():
    import socket
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("",0))
    free_port = s.getsockname()[1]
    s.close
    return free_port

def chdir_to_scripts_folder():
    abspath = os.path.abspath(__file__)
    dname = os.path.dirname(abspath)
    os.chdir(dname)


parser=argparse.ArgumentParser()

parser.add_argument('--hadoop_address', '-a', help='the address:port combination used to access your hadoop cluster', type=str, default="0.0.0.0:9000")
parser.add_argument('--port', '-p', help='port to expose your container on', type=int, default=get_free_port())
parser.add_argument('--detach', '-d', help='run container in detached mode or not', type=bool, default=True)
parser.add_argument('--tag', '-t', help='tag for your container', type=str, default="go-webhdfs")

args=parser.parse_args()

chdir_to_scripts_folder()


print ("hadoop_address has been set to {}".format(args.hadoop_address)) 
print ("port has been set to {}".format(args.port)) 

docker_client = docker.from_env()
tag = "go-webhdfs"
buildargs = {'DOCKER_HADOOP_ADDRESS' : args.hadoop_address}
print("buiding go-webhdfs image. tag = {}, buildargs = {}".format(tag, buildargs))
docker_client.images.build(path=".", tag=tag, buildargs=buildargs)

print("running go-webhdfs container. detach = {} ports = {}".format(args.detach, args.port))
docker_client.containers.run("go-webhdfs", detach=args.detach, ports =  {'8000': args.port})

print("go-webhdfs is running on port: {}".format(args.port))
