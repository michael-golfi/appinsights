#!/bin/bash
docker build -t rootfsimage .
id=$(docker create rootfsimage true)
mkdir -p plugin/rootfs
docker export "$id" | tar -x -C plugin/rootfs

docker rm -vf "$id"
docker rmi rootfsimage

cp config.json ./plugin/