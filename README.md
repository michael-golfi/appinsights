# Microsoft Application Insights Docker Logging Plugin

[![pipeline status](https://gitlab.com/michael.golfi/appinsights/badges/master/pipeline.svg)](https://gitlab.com/michael.golfi/appinsights/commits/master)
[![coverage report](https://gitlab.com/michael.golfi/appinsights/badges/master/coverage.svg)](https://gitlab.com/michael.golfi/appinsights/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/michael.golfi/appinsights)](https://goreportcard.com/report/gitlab.com/michael.golfi/appinsights)

This project implements a Docker logging driver that will allow Docker to stream logs to local JSON files and Microsoft Application Insights. 
This log plugin supports the `docker logs` command.

## Installation

```bash
docker plugin install --alias appinsights michaelgolfi/appinsights
```

## Usage

```bash
docker run -d --name "example-logger" \
  --log-driver appinsights \
  --log-opt token=$AppInsightsToken \
  ubuntu bash -c 'while true; do echo "{\"msg\": \"something\", \"time\": \"`date +%s`\"}"; sleep 2; done;'
```

### Log Options

| Option               | Default                                         |
|----------------------|-------------------------------------------------|
| endpoint             | "https://dc.services.visualstudio.com/v2/track" |
| token                |                                                 |
| verify-connection    | "true"                                          |
| insecure-skip-verify | "false"                                         |
| gzip                 | "false"                                         |
| gzip-level           | "0"                                             |
| batch-size           | "1024"                                          |
| batch-interval       | "5s"                                            |

## Building

This plugin uses godep for vendoring. 
- Run `make install` to install all dependencies. 
- Run `./scripts/build.sh` to build the plugin.

## References

### JSON File Logging Driver

Documentation: https://docs.docker.com/engine/admin/logging/json-file/#options
