FROM golang:latest
WORKDIR /go/src/app
COPY . .
RUN \
  go build CGO_ENABLED=1 -o dfp
