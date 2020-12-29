FROM golang:1.15.3 as build

WORKDIR /go/app/

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download \
    && make linux-build

FROM scratch

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

WORKDIR /go/app/

COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/app/server /go/app/server

EXPOSE 9527

CMD ["./server"]
