package led

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/water78813/iot/manager"

	"gobot.io/x/gobot/platforms/firmata"
)

type ledModule struct {
	ledState  int
	pin       string
	ip        string
	funcState string
	stopCh    chan struct{}
}

//
func LedHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		s, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		mng := manager.GetMng()
		led, err := (*mng).GetMod("led")
		if err != nil {
			lm := &ledModule{
				funcState: "init",
				ledState:  0,
				pin:       "12",
				ip:        "192.168.2.113:3030",
				stopCh:    make(chan struct{}, 1),
			}
			mng.AddMod("led", lm)
			if led, err = (*mng).GetMod("led"); err != nil {
				w.WriteHeader(200)
				w.Write([]byte("app fail"))
			}
		}

		if string(s) == "on" {
			(*led).Start()
			mng.ModReload()
		} else if string(s) == "off" {
			(*led).Stop()
		} else if string(s) == "remove" {
			mng.RemoveMod("led")
		}
		w.WriteHeader(200)
	} else {
		respContext := []byte("post is the only valid method")
		w.Write(respContext)
	}
}

func (lm *ledModule) Start() {
	lm.funcState = "start"
}

func (lm *ledModule) Run() {
	interval := time.Duration(time.Second)
	next := time.Now()
	adaptor := firmata.NewTCPAdaptor(lm.ip)
	err := adaptor.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		//before stop the program set the pin state be origin
		err := adaptor.DigitalWrite(lm.pin, 0)
		if err != nil {
			log.Fatal(err)
		}
		//wait for esp8266 getting the signal
		time.Sleep(time.Second)
		fmt.Println("End")
		adaptor.Disconnect()
	}()
	for {
		if lm.ledState == 0 {
			err := adaptor.DigitalWrite(lm.pin, 1)
			if err != nil {
				log.Fatal(err)
			}
			lm.ledState = 1
		} else {
			err := adaptor.DigitalWrite(lm.pin, 0)
			if err != nil {
				log.Fatal(err)
			}
			lm.ledState = 0
		}
		if interval > 0 {
			now := time.Now()
			next = next.Add(interval)
			if next.Before(now) {
				next = now.Add(interval)
			}
			select {
			case <-lm.stopCh:
				return
			case <-time.After(next.Sub(now)):
			}
		}
	}
}

func (lm *ledModule) Stop() {
	lm.funcState = "stop"
	lm.stopCh <- struct{}{}
}

func (lm *ledModule) GetFuncState() string {
	return lm.funcState
}

func (lm *ledModule) SetFuncState(s string) {
	lm.funcState = s
}
