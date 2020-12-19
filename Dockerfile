FROM golang:latest as builder
WORKDIR /go/src/app
COPY . .
RUN \
  go build -o dfp


FROM alpine:latest
COPY --from=builder /go/src/app/dfp /opt/dfp/bin/dfp
RUN \
  chmod +x /opt/dfp/bin/dfp &&\
  mkdir -p /opt/dfp/config &&\
  mkdir -p /opt/dfp/data &&\
  mkdir -p /opt/dfp/log &&\
  addgroup -g 1000 dfp && \
  adduser -g "DFP user" -D -h /opt/dfp -G dfp -s /bin/sh -u 1000 dfp &&\
  chown -R dfp:dfp /opt/dfp &&\
  apk upgrade &&\
  apk add --update curl bash wget tzdata &&\
  rm -rf /tmp/* /var/cache/apk/*
ENV TZ "Europe/Paris"
WORKDIR "/opt/dfp"
EXPOSE "4040"
VOLUME [ "/opt/dfp/data" ]
CMD [ "/opt/dfp/bin/dfp" ]