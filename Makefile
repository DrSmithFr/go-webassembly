install:
	./aliases.sh copy_wasm_script

build: install
	GOOS=js GOARCH=wasm go build -o public/assets/main.wasm src/main.go

serve: install build
	go run server.go

