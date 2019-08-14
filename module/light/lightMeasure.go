package light

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/water78813/iot/manager"

	"gobot.io/x/gobot/platforms/firmata"
)

type lightModule struct {
	pin       string
	ip        string
	funcState string
	stopCh    chan struct{}
}

//
func LightHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		s, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		mng := manager.GetMng()
		light, err := (*mng).GetMod("light")
		if err != nil {
			lm := &lightModule{
				funcState: "init",
				pin:       "0",
				ip:        "192.168.2.113:3030",
				stopCh:    make(chan struct{}, 1),
			}
			mng.AddMod("light", lm)
			if light, err = (*mng).GetMod("light"); err != nil {
				w.WriteHeader(200)
				w.Write([]byte("app fail"))
			}
		}

		if string(s) == "on" {
			(*light).Start()
			mng.ModReload()
		} else if string(s) == "off" {
			(*light).Stop()
		} else if string(s) == "remove" {
			mng.RemoveMod("light")
		}
		w.WriteHeader(200)
	} else {
		respContext := []byte("post is the only valid method")
		w.Write(respContext)
	}
}

func (lm *lightModule) Start() {
	lm.funcState = "start"
}

func (lm *lightModule) Run() {
	interval := time.Duration(time.Second)
	next := time.Now()
	adaptor := firmata.NewTCPAdaptor(lm.ip)
	err := adaptor.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("End")
		adaptor.Disconnect()
	}()
	for {
		light, err := adaptor.AnalogRead(lm.pin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(light)
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

func (lm *lightModule) Stop() {
	lm.funcState = "stop"
	lm.stopCh <- struct{}{}
}

func (lm *lightModule) GetFuncState() string {
	return lm.funcState
}

func (lm *lightModule) SetFuncState(s string) {
	lm.funcState = s
}
