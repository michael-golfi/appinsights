FROM golang:1.9.2 as builder
WORKDIR /go/src/gitlab.com/michael.golfi/appinsights
COPY . /go/src/gitlab.com/michael.golfi/appinsights
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o /usr/bin/appinsights

FROM alpine
RUN mkdir -p /run/docker/plugins
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /usr/bin/appinsights /usr/bin/appinsights
RUN chmod +x /usr/bin/appinsights