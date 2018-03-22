FROM alpine:3.7
#RUN apk --no-cache add ca-certificates
COPY appinsights /bin/
WORKDIR /bin
ENTRYPOINT [ "/bin/appinsights" ]