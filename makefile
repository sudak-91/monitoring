wasm:
	@GOOS="js" GOARCH="wasm" go build -o ./template/main.wasm ./template/main.go 