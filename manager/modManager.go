package manager

import "fmt"

type IotFunc interface {
	Start()
	Run()
	Stop()
	GetFuncState() string
	SetFuncState(string)
}

type modMng struct {
	iotModules map[string]IotFunc
	reloadCh   chan struct{}
}

var manager = modMng{
	iotModules: make(map[string]IotFunc),
	reloadCh:   make(chan struct{}, 1),
}

// Return the Iot Module Manager
func GetMng() *modMng {
	return &manager
}

func (mm *modMng) GetMod(name string) (*IotFunc, error) {
	iot, ok := mm.iotModules[name]
	if !ok {
		return nil, fmt.Errorf("%s is not exisit", name)
	}
	return &iot, nil
}

func (mm *modMng) ModReload() {
	mm.reloadCh <- struct{}{}
}

func (mm *modMng) AddMod(name string, mod IotFunc) {
	mm.iotModules[name] = mod
}
func (mm *modMng) RemoveMod(name string) {
	mm.iotModules[name] = nil
}

//
func IotFuncMng() {
	for {
		for _, iotmod := range manager.iotModules {
			if iotmod.GetFuncState() == "start" {
				go iotmod.Run()
				iotmod.SetFuncState("running")
			}
		}
		<-manager.reloadCh
		close(manager.reloadCh)
		manager.reloadCh = make(chan struct{}, 1)
	}
}
