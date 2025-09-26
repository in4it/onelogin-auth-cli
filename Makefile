BINARY = onelogin-auth

all: build tests

tests:
	go test -v ./...

build:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY}-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY}-darwin-arm64 main.go

clean:
	rm -f ${BINARY}-linux-${GOARCH}
