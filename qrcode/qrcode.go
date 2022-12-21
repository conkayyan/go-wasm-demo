package main

import (
	"encoding/base64"
	"fmt"
	"github.com/skip2/go-qrcode"
	"syscall/js"
)

func getDocumentDom() js.Value {
	return js.Global().Get("document")
}

func getElementByID(id string) js.Value {
	return getDocumentDom().Call("getElementById", id)
}

func registerCallbackFunc() {
	generate := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		str := getElementByID("str").Get("value").String()
		png, err := qrcode.Encode(str, qrcode.Medium, 256)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		qr := getElementByID("qrcode")
		qr.Set("src", fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(png)))
		return nil
	})

	getElementByID("btn").Call("addEventListener", "click", generate)

}

func main() {
	registerCallbackFunc()
	select {}
}

// $ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static
// GOARCH=wasm GOOS=js go build -o static/qrcode.wasm qrcode.go
