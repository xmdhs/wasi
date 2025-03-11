package main

import (
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func BenchmarkGo(b *testing.B) {
	for b.Loop() {
		fibonacci(40)
	}
}

func BenchmarkWasi(b *testing.B) {
	b.StopTimer()
	ctx := b.Context()
	rt := wazero.NewRuntime(ctx)
	wasm, err := os.ReadFile("test.wasm")
	if err != nil {
		panic(err)
	}
	cm, err := rt.CompileModule(ctx, wasm)
	if err != nil {
		panic(err)
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	mod, err := rt.InstantiateModule(ctx, cm, wazero.NewModuleConfig())
	fibonacci := mod.ExportedFunction("fibonacci")
	b.StartTimer()
	for b.Loop() {
		fibonacci.Call(ctx, 40)
	}
}
