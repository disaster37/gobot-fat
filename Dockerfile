FROM arm32v6/alpine:latest as builder

WORKDIR /root

RUN \
  apk add --update git curl go &&\
  git clone https://github.com/disaster37/gobot-fat.git &&\
  cd gobot-fat &&\
  go build

FROM arm32v6/alpine:latest


COPY --from=builder /root/gobot-fat/gobot-fat /opt/dfp/bin/dfp

RUN \
  chmod +x /opt/dfp/bin/dfp &&\
  mkdir -p /opt/dfp/config &&\
  mkdir -p /opt/dfp/data &&\
  mkdir -p /opt/dfp/log &&\
  addgroup -g 1000 dfp && \
  adduser -g "DFP user" -D -h /opt/dfp -G dfp -s /bin/sh -u 1000 dfp &&\
  chown -R dfp:dfp /opt/dfp &&\
  apk upgrade &&\
  apk add --update curl bash &&\
  rm -rf /tmp/* /var/cache/apk/*

WORKDIR "/opt/dfp"

EXPOSE "4040"

VOLUME [ "/opt/dfp/data" ]

CMD [ "bin/dfp" ]