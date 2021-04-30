FROM golang:1.16.3-alpine3.13

RUN \
    apk update && \
    apk add --no-cache make && \
    rm -rf '/var/cache/apk/*'

WORKDIR /app
