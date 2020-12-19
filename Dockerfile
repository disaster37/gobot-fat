FROM golang:latest
WORKDIR /go/src/app
COPY . .
RUN \
  CGO_ENABLED=1 go build -o dfp
