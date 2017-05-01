package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func main() {
	const (
		ledpath         = "/sys/class/leds"
		triggerfilename = "/sys/class/leds/led0/trigger"
		ledfilename     = "/sys/class/leds/led0/brightness"
		triggerOff      = "none"
		triggerOn       = "mmc0"
	)
	log.Println("starting...")
	log.Printf("LED list: %v\n", listLed(ledpath))

	// disable trigger in OS (standard use of LED)
	trigger(triggerfilename, []byte(triggerOff))

	// Pi Zero has inversed brightness logic, and also only
	// on and off
	log.Println("turning LED off")
	brightness(ledfilename, 1)

	log.Println("toggling...")
	toggle := 0
	for i := 0; i < 25; i++ {
		brightness(ledfilename, toggle)
		if toggle == 0 {
			toggle = 1
		} else if toggle == 1 {
			toggle = 0
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func brightness(f string, mode int) {
	buf := []byte(strconv.Itoa(mode))
	if err := ioutil.WriteFile(f, buf, 0644); err != nil {
		log.Fatalf("cannot write to %v: %v\n", f, err)
	}
}

// Pi zero has one LED but reports led0 and led1
func listLed(ledpath string) (leds []string) {
	fis, err := ioutil.ReadDir(ledpath)
	if err != nil {
		log.Fatalf("error listing %v: %v\n", ledpath, err)
	}
	for _, fi := range fis {
		leds = append(leds, fi.Name())
	}
	return
}

// trigger sets the new trigger and returns previous setting
func trigger(f string, msg []byte) []byte {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("cannot read %v: %v\n", f, err)
	}
	if err := ioutil.WriteFile(f, msg, 0644); err != nil {
		log.Fatalf("cannot write to %v: %v\n", f, err)
	}
	return buf
}
