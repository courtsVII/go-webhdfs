FROM golang:1.11


ARG DOCKER_HADOOP_ADDRESS=0.0.0.0:9000
ENV HADOOP_ADDRESS=${DOCKER_HADOOP_ADDRESS}
WORKDIR /go/src/webhdfs
COPY . .

RUN go get github.com/colinmarc/hdfs && \
    go get -d -v ./... && \
    go install -v ./...

CMD controller
