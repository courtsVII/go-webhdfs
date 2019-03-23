FROM golang:1.11

WORKDIR /go/src/webhdfs
COPY . .

RUN go get github.com/colinmarc/hdfs && \
    go get -d -v ./... && \
    go install -v ./...

CMD controller
