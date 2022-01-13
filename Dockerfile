FROM golang:1.17-alpine3.13 as builder

RUN apk add make

WORKDIR /app/

COPY . .

RUN make download && make build

FROM alpine:3.15

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

WORKDIR /app/

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/config/app.yaml.example /app/config/app.yaml
COPY --from=builder /app/bin/app /app/bin/app

CMD ["./bin/app"]
