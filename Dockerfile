FROM golang:latest AS swizzle

COPY . /go/src/github.com/opsline/swizzle
WORKDIR /go/src/github.com/opsline/swizzle
RUN go get -v github.com/opsline/swizzle/...

FROM 253379484728.dkr.ecr.us-east-1.amazonaws.com/opsline/chalk:latest as chalk

FROM debian:stretch
COPY --from=chalk /usr/local/bin/chalk /usr/local/bin/chalk
COPY --from=swizzle /go/bin/swizzle /usr/local/bin/swizzle
ENTRYPOINT /usr/local/bin/swizzle
