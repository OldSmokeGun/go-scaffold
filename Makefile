.PHONY: build linux-build windows-build mac-build download clean test help

BINARY_PATH = bin/httpserver
HTTPSERVER_HTTPSERVER_MAIN_DIR = cmd/httpserver

build:
	go generate -x ./...
ifeq (${OS}, Windows_NT)
	set CGO_ENABLED=0
	set GOOS=windows
	go build -o ${BINARY_PATH}.exe ${HTTPSERVER_HTTPSERVER_MAIN_DIR}/main.go
else
	CGO_ENABLED=0 go build -o ${BINARY_PATH} ${HTTPSERVER_HTTPSERVER_MAIN_DIR}/main.go
endif

linux-build:
	go generate -x ./...
	CGO_ENABLED=0 GOOS=linux go build -o ${BINARY_PATH} ${HTTPSERVER_HTTPSERVER_MAIN_DIR}/main.go

windows-build:
	go generate -x ./...
	set CGO_ENABLED=0
	set GOOS=windows
	go build -o ${BINARY_PATH}.exe ${HTTPSERVER_HTTPSERVER_MAIN_DIR}/main.go

mac-build:
	go generate -x ./...
	CGO_ENABLED=0 GOOS=darwin go build -o ${BINARY_PATH} ${HTTPSERVER_HTTPSERVER_MAIN_DIR}/main.go

download:
	go env -w GOPROXY=https://goproxy.cn,direct; go mod download

clean:
	@if [ -f ${BINARY_PATH} ] ; then rm ${BINARY_PATH} ; fi

test:
	go test -v ./...

help:
	@printf "%-30s %-100s\n" "make" "默认自动根据平台编译二进制文件"
	@printf "%-30s %-100s\n" "make build" "自动根据平台编译二进制文件"
	@printf "%-30s %-100s\n" "make linux-build" "编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make windows-build" "编译 windows 平台的二进制文件"
	@printf "%-30s %-100s\n" "make mac-build" "编译 mac 平台的二进制文件"
	@printf "%-30s %-100s\n" "make download" "下载编译所需的依赖包"
	@printf "%-30s %-100s\n" "make clean" "清理编译生成的二进制文件"
	@printf "%-30s %-100s\n" "test" "单元测试"
