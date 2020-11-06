BINARY = app
MAIN_PATH = .

linux-build:
	GOOS=linux go build -o ${BINARY} ${MAIN_PATH}/main.go

windows-build:
	GOOS=windows go build -o ${BINARY}.exe ${MAIN_PATH}/main.go

mac-build:
	GOOS=darwin go build -o ${BINARY} ${MAIN_PATH}/main.go

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@printf "%-30s %-100s\n" "make" "默认编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make linux-build" "编译 linux 平台的二进制文件"
	@printf "%-30s %-100s\n" "make windows-build" "编译 windows 平台的二进制文件"
	@printf "%-30s %-100s\n" "make mac-build" "编译 mac 平台的二进制文件"
	@printf "%-30s %-100s\n" "make clean" "清理编译生成的二进制文件"
