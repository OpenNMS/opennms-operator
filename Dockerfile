#Builder
FROM golang:1.17 as builder

WORKDIR app/

COPY Makefile Makefile

RUN make all

#Runner
FROM alpine:3.14

WORKDIR app/

CMD echo "Hello world!"