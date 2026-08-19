//go:build !js || !wasm

package main

func browserMode() string {
	return "electric-circuit"
}
