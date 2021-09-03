FROM golang:1.16.0 as build

WORKDIR /app/

COPY . .

RUN make download && make build

FROM scratch

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

WORKDIR /app/

COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/config/web.yaml.example /app/config/web.yaml
COPY --from=build /app/bin/web /app/bin/web

EXPOSE 9527

CMD ["./bin/web"]
