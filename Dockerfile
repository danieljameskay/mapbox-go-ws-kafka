FROM golang:1.10.3-alpine3.7

RUN apk update && apk upgrade && apk add --no-cache bash \
      libressl \
      tar \
      git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils

WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN /root/librdkafka/configure
RUN make
RUN make install
WORKDIR /go/src/app
RUN git clone https://github.com/danieljameskay/mapbox-go-ws-kafka.git
WORKDIR /go/src/app/mapbox-go-ws-kafka
RUN go get -d -v ./...
#RUN go install -v ./...
RUN go build
CMD ["./mapbox-go-ws-kafka"]
