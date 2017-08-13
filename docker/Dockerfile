FROM alpine:3.6

EXPOSE 53 53/udp

ENV HOST_FILE /etc/hosts.private

ENV RESOLV_FILE /etc/dnsmasq.resolv.conf

ENV GRPC_PORT 9999

RUN apk --no-cache add dnsmasq

RUN touch $HOST_FILE

RUN echo nameserver 8.8.8.8 > $RESOLV_FILE

RUN echo addn-hosts=$HOST_FILE > /etc/dnsmasq.d/private.hosts.conf

RUN echo resolv-file=$RESOLV_FILE > /etc/dnsmasq.d/private.resolv.conf

RUN apk add --no-cache --virtual .build-deps go git gcc openssl musl-dev

RUN go get github.com/nokamoto/private-dns/dnssupervisor

RUN go get github.com/nokamoto/private-dns/dnscli

RUN apk del .build-deps

COPY run.sh run.sh

RUN chmod u+x run.sh

CMD ./run.sh
