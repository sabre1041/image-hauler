image-hauler
===============

Tool to Transfer Docker Images Between Registries Without A Docker Client

## Overview

Docker as a format defines a mechanism to build an image once and store it in any number of registries. The transfer of images typically requires the use of the Docker client to facilitate the transfer. On some systems, installing docker can be cumbersome or impossible due to system limitations. 

*image-hauler* is a tool to move docker images between registries without needing to install or configure Docker on the local machine

## Usage

```
$ ./image-hauler --help
image-hauler is a Tool to Transfer Docker Images Between Registries

Usage:
  image-hauler [flags]

Flags:
      --destination-image string      Destination Docker Image
      --destination-password string   Password for the Destination Registry
      --destination-registry string   Address of the Destination Docker Registry
      --destination-tag string        Destination Docker Tag (default "latest")
      --destination-username string   Username for the Destination Registry
      --source-image string           Source Docker Image
      --source-password string        Password for the Source Registry
      --source-registry string        Address of the Source Docker Registry
      --source-tag string             Source Docker Tag (default "latest")
      --source-username string        Username for the Source Registry
      --storage-dir string            Location to Store Docker Images During Execution (default "/tmp/image-hauler")
```

### Example Executions

```
./image-hauler --source-registry=http://<docker_ip>:5000 --source-image=docker.io/busybox --destination-registry=http://<docker_ip>:5000 --destination-image=sabre1041/busybox --destination-tag=1.0
```