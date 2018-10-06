FROM golang:1.11.1-stretch as build-stage

ADD . /go/src/func
WORKDIR /go/src/func
RUN go build -o func

FROM debian:stretch

RUN mkdir /shared_libs
RUN apt-get update && apt-get install file -qy

COPY --from=build-stage /go/src/func/sample/linux_invoker.so /shared_libs/linux_invoker.so
COPY --from=build-stage /go/src/func/func /func

ENTRYPOINT ["./func"]
