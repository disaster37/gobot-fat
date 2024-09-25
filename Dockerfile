FROM golang:1.22-alpine as builder
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
RUN apk add --update git
WORKDIR /go/src/app
COPY . .
RUN \
  CGO_ENABLED=0 go build


FROM alpine:3.15
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
COPY --from=builder /go/src/app/gobot-fat /opt/dfp/bin/dfp
RUN \
  chmod +x /opt/dfp/bin/dfp &&\
  mkdir -p /opt/dfp/config &&\
  mkdir -p /opt/dfp/data &&\
  mkdir -p /opt/dfp/log &&\
  addgroup -g 1000 dfp && \
  adduser -g "DFP user" -D -h /opt/dfp -G dfp -s /bin/sh -u 1000 dfp &&\
  chown -R dfp:dfp /opt/dfp &&\
  apk upgrade &&\
  apk add --update curl bash wget tzdata
RUN   rm -rf /tmp/* /var/cache/apk/*
ENV TZ "Europe/Paris"
WORKDIR "/opt/dfp"
EXPOSE "4040"
VOLUME [ "/opt/dfp/data" ]
CMD [ "/opt/dfp/bin/dfp" ]