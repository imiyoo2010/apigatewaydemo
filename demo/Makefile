#程序生成的执行文件
BINARYC=api-gateway

BINARYS=api-server

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BINARYC} main.go

	GOOS=darwin GOARCH=amd64 go build -o ${BINARYS} server/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARYC} main.go

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARYS} server/main.go



clean:
	rm -rf ${BINARYC}
	rm -rf ${BINARYS}