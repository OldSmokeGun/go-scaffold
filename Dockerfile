FROM golang:1.17 as builder

ARG PROTOC_VERSION=3.19.1
ARG PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip

RUN apt-get update && apt-get install -y unzip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP} \
    && unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc \
    && rm -f ${PROTOC_ZIP} \
    && apt-get autoclean && apt-get clean

WORKDIR /app/

COPY . .

RUN make download && make proto && make build

FROM alpine:3.14

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

WORKDIR /app/

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/.air.conf.example /app/.air.conf
COPY --from=builder /app/etc/config.yaml.example /app/etc/config.yaml
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/bin/app /app/bin/app

CMD ["./bin/app"]
