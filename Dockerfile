FROM golang:latest
WORKDIR /go/src/app
COPY . .
RUN \
  CGO_ENABLED=0 go build -o dfp
