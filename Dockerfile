FROM golang:1.17 AS builder

COPY . /go/src/app
WORKDIR /go/src/app

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux make build

RUN git log -1 --oneline > version.txt



VOLUME ["/go/src/analytics-file-export-service/files"]

FROM alpine:latest

RUN mkdir -p /go/src/analytics-file-export-service/files

WORKDIR /go/src/analytics-file-export-service

COPY --from=builder /go/src/app/analytics-file-export-service .
COPY --from=builder /go/src/app/version.txt .

ENTRYPOINT ["./analytics-file-export-service"]


