FROM golang:1.16.0 as build

WORKDIR /app/

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download \
    && make linux-build

FROM scratch

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

WORKDIR /app/

COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/config/config.yaml.example /app/config/config.yaml
COPY --from=build /app/bin/server /app/bin/server

EXPOSE 9527

CMD ["./bin/server"]
