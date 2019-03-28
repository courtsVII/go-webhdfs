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

def parse_environment_variables():
    env_list = args.env
    env_dict = dict()
    for env_variable in env_list:
        entry = env_variable.split("=")
        if len(entry) != 2:
            print("env variable '{}' is invalid".format(env_variable))
            print("env variables must take the form key=value")
            quit()
        key = entry[0]
        val = entry[1]
        env_dict[key] = val
    return env_dict

parser=argparse.ArgumentParser()
parser.add_argument('--env', '-e', help='docker environment variables passed as key=value (multiple can be provided)', type=str, action='append')
parser.add_argument('--port', '-p', help='port to expose your container on', type=int, default=get_free_port())
parser.add_argument('--detach', '-d', help='run container in detached mode or not', action='store_true')
parser.add_argument('--tag', '-t', help='tag for your container', type=str, default="go-webhdfs")
args=parser.parse_args()

chdir_to_scripts_folder()
env_vars = parse_environment_variables()

print ("tag has been set to {}".format(args.tag)) 
print ("port has been set to {}".format(args.port)) 
print ("detach has been set to {}".format(args.detach)) 
print ("environment has been set to {}".format(env_vars))

docker_client = docker.from_env()
print ("building docker image {}".format(args.tag))
docker_client.images.build(path=".", tag=args.tag)
print("running go-webhdfs on port: {}".format(args.port))
docker_client.containers.run("go-webhdfs", detach=args.detach, ports =  {'8000': args.port}, environment=env_vars)

