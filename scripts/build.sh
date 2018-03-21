#!/bin/bash
rm -rf plugin

docker build -t rootfsimage .
id=$(docker create rootfsimage true)
mkdir -p plugin/rootfs
docker export "$id" | tar -x -C plugin/rootfs

docker rm -vf "$id"
docker rmi rootfsimage
cp config.json ./plugin/

docker plugin disable michaelgolfi/appinsights
docker plugin rm michaelgolfi/appinsights
docker plugin create michaelgolfi/appinsights ./plugin

rm -rf plugin