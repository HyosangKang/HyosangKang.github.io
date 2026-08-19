//go:build js && wasm

package main

import "syscall/js"

func browserMode() string {
	body := js.Global().Get("document").Get("body")
	mode := body.Get("dataset").Get("gameMode").String()
	if mode == "" {
		return "electric-circuit"
	}
	return mode
}
