FROM golang:1.14

COPY . /go/src/analytics-file-export-service
WORKDIR /go/src/analytics-file-export-service

VOLUME ["/data"]

RUN make build

CMD ./analytics-file-export-service