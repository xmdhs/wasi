package main

import (
	"context"
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

	// 注册宿主函数，供 WASM 中的 Go 代码调用
	_, err = rt.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, a, b int32) int32 {
			return a + b
		}).
		Export("hostAdd").
		Instantiate(ctx)
	if err != nil {
		panic(err)
	}

	mod, err := rt.InstantiateModule(ctx, cm, wazero.NewModuleConfig())
	if err != nil {
		panic(err)
	}
	fibonacci := mod.ExportedFunction("fibonacci")
	for i := 0; i < 1000; i++ {
		fibonacci.Call(ctx, 40)
	}
	b.StartTimer()
	b.ResetTimer()
	for b.Loop() {
		fibonacci.Call(ctx, 40)
	}
}

func TestCallHostAdd(t *testing.T) {
	ctx := t.Context()
	rt := wazero.NewRuntime(ctx)
	defer rt.Close(ctx)

	wasmBytes, err := os.ReadFile("test.wasm")
	if err != nil {
		t.Fatal(err)
	}
	cm, err := rt.CompileModule(ctx, wasmBytes)
	if err != nil {
		t.Fatal(err)
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	// 注册宿主函数
	_, err = rt.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, a, b int32) int32 {
			return a + b
		}).
		Export("hostAdd").
		Instantiate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mod, err := rt.InstantiateModule(ctx, cm, wazero.NewModuleConfig())
	if err != nil {
		t.Fatal(err)
	}

	callHostAdd := mod.ExportedFunction("callHostAdd")
	results, err := callHostAdd.Call(ctx, 3, 7)
	if err != nil {
		t.Fatal(err)
	}
	if int32(results[0]) != 10 {
		t.Fatalf("expected 10, got %d", results[0])
	}
}
