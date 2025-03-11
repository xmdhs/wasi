SET CGO_ENABLED=0
SET GOOS=wasip1
SET GOARCH=wasm
go build -trimpath -ldflags "-w -s" -buildmode=c-shared -o test.wasm