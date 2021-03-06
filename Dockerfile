FROM golang:1.16-alpine3.13

RUN \
    apk update && \
    apk add --no-cache bash build-base && \
    rm -rf '/var/cache/apk/*'

WORKDIR /app
