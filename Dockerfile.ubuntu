FROM golang:1.9.2 as builder
WORKDIR /go/src/gitlab.com/michael.golfi/appinsights
COPY . /go/src/gitlab.com/michael.golfi/appinsights
RUN go build --ldflags '-extldflags "-static"' -o /usr/bin/appinsights

FROM ubuntu
RUN mkdir -p /run/docker/plugins
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/bin/appinsights /usr/bin/appinsights
RUN chmod +x /usr/bin/appinsights