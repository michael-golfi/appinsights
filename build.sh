#!/bin/bash

rm -rf ./plugin

docker build -t rootfsimage .
id=$(docker create rootfsimage true) # id was cd851ce43a403 when the image was created
mkdir -p plugin/rootfs
docker export "$id" | tar -x -C plugin/rootfs
docker rm -vf "$id"
docker rmi rootfsimage
cp config.json ./plugin/

docker plugin disable michael-golfi/appinsights
docker plugin rm michael-golfi/appinsights
docker plugin create michael-golfi/appinsights ./plugin

rm -rf ./plugin
