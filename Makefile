.PHONY: build linux-build windows-build mac-build clean test help

BINARY_PATH = bin/httpserver
HTTPSERVER_HTTPSERVER_MAIN_DIR = cmd/httpserver


build:
	go generate -x ./...
	@binaryPath=${BINARY_PATH}; \
	os=`go env GOOS`; \
	echo "os: $${os}"; \
	if [ $${os} == "windows" ]; then binaryPath=$${binaryPath}.exe; fi; \
	CGO_ENABLED=0 go build -o $${binaryPath} ${HTTPSERVER_HTTPSERVER_MAIN_DIR}/main.go

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
	@printf "%-30s %-100s\n" "make" "默认编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make linux-build" "编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make windows-build" "编译 windows 平台的二进制文件"
	@printf "%-30s %-100s\n" "make mac-build" "编译 mac 平台的二进制文件"
	@printf "%-30s %-100s\n" "make clean" "清理编译生成的二进制文件"
