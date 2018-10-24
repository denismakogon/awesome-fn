FROM alpine:latest

RUN apk add tar
COPY node/ /function/
WORKDIR /function
RUN tar -cvf init.tar *
RUN apk del tar
CMD ["cat", "/function/init.tar"]
