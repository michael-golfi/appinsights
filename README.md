# Microsoft Application Insights Docker Logging Plugin

This project implements a Docker logging driver that will allow Docker to stream logs to local JSON files and to also redirect log streams to Microsoft Application Insights. This project supports the `docker logs` command.

## Installation

```bash
docker plugin install michael-golfi/appinsights
```

## Building

This plugin uses godep for vendoring. 
- Run `make install` to install all dependencies. 
- Run `./build.sh` to build the plugin.

## Usage

```bash
docker run -d --name "example-logger" \
  --log-driver michael-golfi/appinsights
  --log-opt insights-token=b11d730f-995c-4eda-ac8a-79093fcace6d \
  ubuntu bash -c 'while true; do echo "{\"msg\": \"something\", \"time\": \"`date +%s`\"}"; sleep 2; done;'
```

## Supported Log Endpoints

* json-file
* appinsights

### JSON File Logging Driver

Documentation: https://docs.docker.com/engine/admin/logging/json-file/#options

### Microsoft Application Insights Logging Driver

The documentation for this plugin are a work in progress currently.