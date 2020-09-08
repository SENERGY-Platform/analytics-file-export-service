FROM golang:1.14

COPY . /go/src/analytics-file-export-service
WORKDIR /go/src/analytics-file-export-service

RUN mkdir /data

VOLUME ["/go/src/analytics-file-export-service/files"]

RUN make build

CMD ./analytics-file-export-service