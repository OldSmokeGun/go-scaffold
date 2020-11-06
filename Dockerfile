FROM golang:1.15.3 as build

WORKDIR /go/app/

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download \
    && make build

FROM scratch

ENV TZ=Asia/Shanghai
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

WORKDIR /go/

COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build /go/app/bin/app /go/

EXPOSE 9527

CMD ["./app"]
