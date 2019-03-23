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
    print ("workdir has been set to {}".format(dname))
    os.chdir(dname)


parser=argparse.ArgumentParser()
parser.add_argument('--hadoop_address', '-a', help='the address:port combination used to access your hadoop cluster', type=str, default="0.0.0.0:9000")
parser.add_argument('--port', '-p', help='port to expose your container on', type=int, default=get_free_port())
parser.add_argument('--detach', '-d', help='run container in detached mode or not', action='store_true')
parser.add_argument('--tag', '-t', help='tag for your container', type=str, default="go-webhdfs")
args=parser.parse_args()

print ("tag has been set to {}".format(args.tag)) 
print ("port has been set to {}".format(args.port)) 
print ("detach has been set to {}".format(args.detach)) 
print ("hadoop_address has been set to {}".format(args.hadoop_address)) 

chdir_to_scripts_folder()

docker_client = docker.from_env()
print ("building docker image {}".format(args.tag))
docker_client.images.build(path=".", tag=args.tag, buildargs={'DOCKER_HADOOP_ADDRESS' : args.hadoop_address})
print("running go-webhdfs on port: {}".format(args.port))
docker_client.containers.run("go-webhdfs", detach=args.detach, ports =  {'8000': args.port})

