.PHONY: build linux-build windows-build mac-build download clean test help

APP_BIN_PATH = bin/app
APP_MAIN_DIR = cmd/app
API_DOC_SCAN_DIR = internal/app
API_DOC_SCAN_ENTRY = app.go
API_DOC_OUT_DIR = internal/app/http/api/docs

build:
	go generate -x ./...
ifeq (${OS}, Windows_NT)
	set CGO_ENABLED=0
	set GOOS=windows
	go build -tags=jsoniter -o ${APP_BIN_PATH}.exe ${APP_MAIN_DIR}/main.go
else
	CGO_ENABLED=0 go build -tags=jsoniter -o ${APP_BIN_PATH} ${APP_MAIN_DIR}/main.go
endif

linux-build:
	go generate -x ./...
	CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -o ${APP_BIN_PATH}_linux ${APP_MAIN_DIR}/main.go

windows-build:
	go generate -x ./...
	set CGO_ENABLED=0
	set GOOS=windows
	go build -tags=jsoniter -o ${APP_BIN_PATH}_windows.exe ${APP_MAIN_DIR}/main.go

mac-build:
	go generate -x ./...
	CGO_ENABLED=0 GOOS=darwin go build -tags=jsoniter -o ${APP_BIN_PATH}_mac ${APP_MAIN_DIR}/main.go

download:
	go env -w GOPROXY=https://goproxy.cn,direct; go mod download; \
	go install github.com/cosmtrek/air@latest; \
	go install github.com/swaggo/swag/cmd/swag@v1.7.8; \
	go install github.com/golang/mock/mockgen@latest; \
	go install entgo.io/ent/cmd/ent@latest;

clean:
	@if [ -f ${APP_BIN_PATH} ] ; then rm ${APP_BIN_PATH} ; fi

test:
	go test -gcflags=-l -v ./...

doc:
	swag fmt -d ${API_DOC_SCAN_DIR} -g ${API_DOC_SCAN_ENTRY}
	swag init -d ${API_DOC_SCAN_DIR} -g ${API_DOC_SCAN_ENTRY} -o ${API_DOC_OUT_DIR} --parseInternal

help:
	@printf "%-30s %-100s\n" "make" "默认自动根据平台编译二进制文件"
	@printf "%-30s %-100s\n" "make build" "自动根据平台编译二进制文件"
	@printf "%-30s %-100s\n" "make linux-build" "编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make windows-build" "编译 windows 平台的二进制文件"
	@printf "%-30s %-100s\n" "make mac-build" "编译 mac 平台的二进制文件"
	@printf "%-30s %-100s\n" "make download" "下载编译所需的依赖包"
	@printf "%-30s %-100s\n" "make clean" "清理编译生成的二进制文件"
	@printf "%-30s %-100s\n" "make test" "单元测试"
	@printf "%-30s %-100s\n" "make doc" "生成文档"
