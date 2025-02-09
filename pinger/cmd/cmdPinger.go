package cmd

import (
	"os"
	"os/exec"
	messaging "pinger/amqp"
	"strconv"
	"sync"
	"time"
)

type Pinger struct {
	addressesToPing map[string]bool
	mux             sync.Mutex
}

func (_ *Pinger) pingHost(address string, packageCnt int) bool {
	err := exec.Command("ping", address, "-c", strconv.Itoa(packageCnt)).Run()
	if err != nil {
		return false
	}
	return true
}
func (pinger *Pinger) doPing(address string, packageCnt int, wg *sync.WaitGroup) {
	res := pinger.pingHost(address, packageCnt)
	messaging.SendToAddPing(address, time.Now(), res)
	wg.Done()
}

func (pinger *Pinger) AddService(address string) {
	pinger.mux.Lock()
	if pinger.addressesToPing == nil {
		pinger.addressesToPing = make(map[string]bool)
	}
	pinger.addressesToPing[address] = true
	pinger.mux.Unlock()
}
func (pinger *Pinger) DeleteService(address string) {
	pinger.mux.Lock()
	if pinger.addressesToPing == nil {
		pinger.addressesToPing = make(map[string]bool)
	}
	pinger.addressesToPing[address] = false
	pinger.mux.Unlock()
}
func (pinger *Pinger) DoPinging() {
	cnt, err := strconv.Atoi(os.Getenv("PING_PACKAGES_CNT"))
	if err != nil {
		cnt = 3
	}
	pingingTime, err := strconv.Atoi(os.Getenv("PING_TIME"))
	if err != nil {
		pingingTime = 10
	}
	var wg sync.WaitGroup
	pinger.mux.Lock()
	if pinger.addressesToPing == nil {
		pinger.addressesToPing = make(map[string]bool)
	}
	pinger.mux.Unlock()
	for {
		pinger.mux.Lock()
		for key, val := range pinger.addressesToPing {
			if val {
				wg.Add(1)
				go pinger.doPing(key, cnt, &wg)
			}
		}
		wg.Wait()
		pinger.mux.Unlock()
		time.Sleep(time.Duration(pingingTime) * time.Second)
	}
}
