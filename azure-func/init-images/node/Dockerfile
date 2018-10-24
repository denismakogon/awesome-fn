FROM golang:1.11-stretch as build-stage

RUN go get github.com/fnproject/fdk-go
ADD func.go /go/src/func/func.go
WORKDIR /go/src/func
RUN go build -o func-binary

################################################
# ordinary dockerfile for azure serverless app #
################################################
FROM mcr.microsoft.com/azure-functions/node:2.0

ENV AzureWebJobsScriptRoot=/home/site/wwwroot
COPY . /home/site/wwwroot
################################################
##################### end ######################
################################################

COPY --from=build-stage /go/src/func/func-binary /func
ADD run_host.sh /run_host.sh
ENTRYPOINT ["/func"]
