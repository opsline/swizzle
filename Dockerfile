FROM golang:latest AS swizzle

COPY . /go/src/github.com/opsline/swizzle/
WORKDIR /go/src/github.com/opsline/swizzle
RUN go get -v github.com/opsline/swizzle/...

FROM opsline/echo-debian:latest

COPY --from=swizzle /go/bin/swizzle /usr/local/bin/swizzle

COPY docker/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

