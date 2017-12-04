# multilogger - a logger docker plugin

**Consider this currently as a work in progress, alpha/beta plugin. While it current works, I haven't riggerously tested it, so there could be performance or other edge cases that have not been addressed yet.**

This implements a logger as a docker plugin that allows for logging to splunk and a json file that is internal to the plugin. The benefit to this is all the logs end up in splunk, but the command `docker logs` still works.

## Credits

Thanks to [cpuguy83](http://github.com/cpuguy83/) for the basis of this plugin, without him this would not exist. Thanks to [docker](https://github.com/docker/) for creating docker plugins!

## Installation

```bash
docker plugin install ekristen/multilogger
```

## Usage

After installation you can reference `ekristen/multilogger` as your logging driver, from there, all the jsonfile logger options and the splunk logger options are all valid and accepted. The only required options are `splunk-url` and `splunk-token` everything else is optional. 

### Example Usage

```bash
docker run -d --name "example-logger" \
  --log-driver michael-golfi/appinsights
  --log-opt splunk-url=https://dc.visualstudio.com \
  --log-opt splunk-token=b11d730f-995c-4eda-ac8a-79093fcace6d \
  ubuntu bash -c 'while true; do echo "{\"msg\": \"something\", \"time\": \"`date +%s`\"}"; sleep 2; done;'


```

## Supported Drivers

* json-file
* splunk

### JSON File Logging Driver

Documentation: https://docs.docker.com/engine/admin/logging/json-file/#options

### Splunk Logging Driver

Documentation: https://docs.docker.com/engine/admin/logging/splunk/#splunk-options

## Building

This plugin is using `govendor` to manage all it's dependencies along with multi-stage docker builds. You'll need docker 17.05-ce or later to support the build process.

```bash
bash build.sh
```

## License

[Apache License](LICENSE.md)

