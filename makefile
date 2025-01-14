KIWI_VERSION := "v0.20.3"

build:
	@C_INCLUDE_PATH=$(pwd)/kiwi/include \
	LIBRARY_PATH=$(pwd)/kiwi/lib \
	LD_LIBRARY_PATH=$(pwd)/kiwi/lib \
	go build -o main ./cmd/main.go
