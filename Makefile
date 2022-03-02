.PHONY: build linux-build windows-build mac-build download clean test generate doc proto help

APP_BIN_PATH = bin/app
APP_MAIN_DIR = cmd/app
API_SWAGGER_SCAN_DIR = internal/app/transport/http/handler
API_SWAGGER_SCAN_ENTRY = handler.go
API_SWAGGER_OUT_DIR = internal/app/transport/http/handler/docs
API_PROTO_FILES=$(shell find internal/app/api -name *.proto)

build:
	make generate
ifeq (${OS}, Windows_NT)
	set CGO_ENABLED=0
	set GOOS=windows
	go build -tags=jsoniter -o ${APP_BIN_PATH}.exe ${APP_MAIN_DIR}/main.go ${APP_MAIN_DIR}/wire_gen.go
else
	CGO_ENABLED=0 go build -tags=jsoniter -o ${APP_BIN_PATH} ${APP_MAIN_DIR}/main.go ${APP_MAIN_DIR}/wire_gen.go
endif

linux-build:
	make generate
	CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -o ${APP_BIN_PATH}_linux ${APP_MAIN_DIR}/main.go ${APP_MAIN_DIR}/wire_gen.go

windows-build:
	make generate
	set CGO_ENABLED=0
	set GOOS=windows
	go build -tags=jsoniter -o ${APP_BIN_PATH}_windows.exe ${APP_MAIN_DIR}/main.go ${APP_MAIN_DIR}/wire_gen.go

mac-build:
	make generate
	CGO_ENABLED=0 GOOS=darwin go build -tags=jsoniter -o ${APP_BIN_PATH}_mac ${APP_MAIN_DIR}/main.go ${APP_MAIN_DIR}/wire_gen.go

download:
	go env -w GOPROXY=https://goproxy.cn,direct; go mod download; \
	go get -u github.com/davecgh/go-spew/spew; \
	go get github.com/google/wire/cmd/wire@v0.5.0; \
	go install github.com/google/wire/cmd/wire@latest; \
	go install github.com/cosmtrek/air@latest; \
	go install github.com/swaggo/swag/cmd/swag@v1.7.8; \
	go install github.com/golang/mock/mockgen@latest; \
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest; \
	go install github.com/envoyproxy/protoc-gen-validate@latest;

clean:
	@if [ -f ${APP_BIN_PATH} ] ; then rm ${APP_BIN_PATH} ; fi

test:
	go test -gcflags=-l -v ./...

generate:
	go generate ./...

doc:
	swag fmt -d ${API_SWAGGER_SCAN_DIR} -g ${API_SWAGGER_SCAN_ENTRY}
	swag init -d ${API_SWAGGER_SCAN_DIR} -g ${API_SWAGGER_SCAN_ENTRY} -o ${API_SWAGGER_OUT_DIR} --parseInternal

proto:
	kratos proto client --proto_path=./proto internal/app/api/v1/greet/greet.proto
	kratos proto client --proto_path=./proto internal/app/api/v1/user/user.proto

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
