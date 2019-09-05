package led

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/water78813/iot/manager"

	"gobot.io/x/gobot/platforms/firmata"
)

type ledModule struct {
	ledState  int
	pin       string
	ip        string
	interval  string
	funcState string
	stopCh    chan struct{}
}

//
func LedAccessor(m map[string]string) error {
	pin := m["pin"]
	host := m["host"]
	interval := m["interval"]
	status := m["status"]
	mng := manager.GetMng()
	led, err := (*mng).GetMod("led")
	if err != nil {
		lm := &ledModule{
			funcState: "init",
			ledState:  0,
			pin:       pin,
			ip:        host,
			interval:  interval,
			stopCh:    make(chan struct{}, 1),
		}
		mng.AddMod("led", lm)
		if led, err = (*mng).GetMod("led"); err != nil {
			return err
		}
	}
	fmt.Println("in accessor")
	if status == "on" {
		(*led).Start()
		mng.ModReload()
	} else if status == "off" {
		(*led).Stop()
	} else if status == "remove" {
		if (*led).GetFuncState() != "on" {
			(*led).Stop()
		}
		mng.RemoveMod("led")
	}
	return nil
}

func (lm *ledModule) Start() {
	fmt.Println("start")
	lm.funcState = "start"
}

func (lm *ledModule) Run() {
	interg, _ := strconv.Atoi(lm.interval)
	interval := time.Duration(interg) * time.Second
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
