# Microsoft Application Insights Docker Logging Plugin

<p align="center">
  <a href="https://gitlab.com/michael.golfi/appinsights-backup/commits/master"><img alt="pipeline status" src="https://gitlab.com/michael.golfi/appinsights-backup/badges/master/pipeline.svg" /></a>
  <a href="https://gitlab.com/michael.golfi/appinsights-backup/commits/master"><img alt="coverage report" src="https://gitlab.com/michael.golfi/appinsights-backup/badges/master/coverage.svg" /></a>
  <a href="https://goreportcard.com/report/gitlab.com/michael.golfi/appinsights"><img alt="coverage report" src="https://goreportcard.com/badge/gitlab.com/michael.golfi/appinsights" /></a>
</p>

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
