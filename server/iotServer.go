package main

import (
	"log"
	"net/http"

	"github.com/water78813/iot/manager"
	"github.com/water78813/iot/module/led"
	"github.com/water78813/iot/module/light"
)

func main() {
	go manager.IotFuncMng()
	http.HandleFunc("/led", led.LedHandler)
	http.HandleFunc("/light", light.LightHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
