FROM golang:1.15 as builder
WORKDIR /go/src/app
COPY . .
RUN \
  CGO_ENABLED=0 go build

FROM alpine:3.12
COPY --from=builder /go/src/app/gobot-fat /opt/dfp/bin/dfp
RUN \
  chmod +x /opt/dfp/bin/dfp &&\
  mkdir -p /opt/dfp/config &&\
  mkdir -p /opt/dfp/data &&\
  mkdir -p /opt/dfp/log &&\
  addgroup -g 1000 dfp && \
  adduser -g "DFP user" -D -h /opt/dfp -G dfp -s /bin/sh -u 1000 dfp &&\
  chown -R dfp:dfp /opt/dfp &&\
  apk upgrade
RUN apk add curl

# Break buildx on github
#RUN apk add --update bash
#RUN apk add --update wget
#RUN apk add tzdata
RUN   rm -rf /tmp/* /var/cache/apk/*
ENV TZ "Europe/Paris"
WORKDIR "/opt/dfp"
EXPOSE "4040"
VOLUME [ "/opt/dfp/data" ]
CMD [ "/opt/dfp/bin/dfp" ]