.PHONY: linux-build windows-build mac-build clean test help

BINARY = bin/server
MAIN_PATH = core/cmd/server

linux-build:
	go generate -x ./...
	CGO_ENABLED=0 GOOS=linux go build -o ${BINARY} ${MAIN_PATH}/main.go

windows-build:
	go generate -x ./...
	set CGO_ENABLED=0
	set GOOS=windows
	go build -o ${BINARY}.exe ${MAIN_PATH}/main.go

mac-build:
	go generate -x ./...
	CGO_ENABLED=0 GOOS=darwin go build -o ${BINARY} ${MAIN_PATH}/main.go

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

test:
	go test -v ./...

help:
	@printf "%-30s %-100s\n" "make" "默认编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make linux-build" "编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make windows-build" "编译 windows 平台的二进制文件"
	@printf "%-30s %-100s\n" "make mac-build" "编译 mac 平台的二进制文件"
	@printf "%-30s %-100s\n" "make clean" "清理编译生成的二进制文件"
