docker run -d --name "example-logger" --log-driver michael-golfi/appinsights --log-opt insights-url=https://dc.services.visualstudio.com --log-opt insights-key=b11d730f-995c-4eda-ac8a-79093fcace6d ubuntu bash -c 'while true; do echo "{\"msg\": \"something\", \"time\": \"`date +%s`\"}"; sleep 2; done;'

