TINYDFS_BINARY_NAME=tinydfs

build-tinydfs:
	GOARCH=amd64 GOOS=windows go build -o build/${TINYDFS_BINARY_NAME}-amd64-windows ./tinydfs
	GOARCH=amd64 GOOS=darwin go build -o build/${TINYDFS_BINARY_NAME}-amd64-darwin ./tinydfs
	GOARCH=amd64 GOOS=linux go build -o build/${TINYDFS_BINARY_NAME}-amd64-linux ./tinydfs
	GOARCH=arm64 GOOS=darwin go build -o build/${TINYDFS_BINARY_NAME}-arm64-darwin ./tinydfs

clean:
	go clean
	rm ./build/${TINYDFS_BINARY_NAME}-amd64-windows
	rm ./build/${TINYDFS_BINARY_NAME}-amd64-darwin
	rm ./build/${TINYDFS_BINARY_NAME}-amd64-linux
	rm ./build/${TINYDFS_BINARY_NAME}-arm64-darwin
	
deps:
	go mod download

test:
	go test ./...