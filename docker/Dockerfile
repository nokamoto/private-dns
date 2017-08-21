FROM alpine:3.6

EXPOSE 53 53/udp

ENV HOST_FILE /etc/hosts.private

ENV GRPC_PORT 9999

ENV API_KEY default

RUN apk --no-cache add dnsmasq

RUN touch $HOST_FILE

RUN echo addn-hosts=$HOST_FILE > /etc/dnsmasq.d/private.hosts.conf

RUN apk add --no-cache --virtual .build-deps go git gcc openssl musl-dev

RUN go get github.com/nokamoto/private-dns/dnssupervisor

RUN go get github.com/nokamoto/private-dns/dnscli

RUN apk del .build-deps

COPY run.sh run.sh

RUN chmod u+x run.sh

ENV PATH $PATH:/root/go/bin

CMD ./run.sh
