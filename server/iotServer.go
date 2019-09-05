package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/water78813/iot/manager"
	"github.com/water78813/iot/module/led"
	"github.com/water78813/iot/module/light"
)

func main() {
	go manager.IotFuncMng()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/led", ledHandler)
	http.HandleFunc("/light", light.LightHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("home.html"))
	if err := t.ExecuteTemplate(w, "home.html", nil); err != nil {
		log.Fatal(err)
	}
}

func ledHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("in ledhander")
	if req.Method == "POST" {
		fmt.Println("in ledhander post")
		err := req.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		host := req.Form.Get("host")
		pin := req.Form.Get("pin")
		status := req.Form.Get("status")
		interval := req.Form.Get("interval")
		m := map[string]string{
			"host":     host,
			"pin":      pin,
			"status":   status,
			"interval": interval,
		}
		fmt.Printf("map is %v¥r¥n", m)
		if err := led.LedAccessor(m); err != nil {
			log.Fatal(err)
		}
		t := template.Must(template.ParseFiles("led.html"))
		if err := t.ExecuteTemplate(w, "led.html", m); err != nil {
			log.Fatal(err)
		}
	}
}
