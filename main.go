package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/gvalkov/golang-evdev"
)

var (
	verbose bool
)

type Config struct {
	Device    string `yaml:"device"`
	Press     string `yaml:"press"`
	Release   string `yaml:"release"`
	LongPress string `yaml:"longpress"`
	uinput    *evdev.InputDevice
}

func (c *Config) exec(n, s string) {
	if verbose {
		log.Print(n + " triggered")
	}
	if s == "" {
		return
	}
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Print(err)
	}
}

func (c *Config) press() {
	c.exec("press", c.Press)
}

func (c *Config) release() {
	c.exec("release", c.Release)
}

func (c *Config) longpress() {
	c.exec("long press", c.LongPress)
}

func (c *Config) watch(wg *sync.WaitGroup) {
	defer wg.Done()

	last := int64(0)
	for {
		ev, err := c.uinput.ReadOne()
		if err != nil {
			log.Print(err)
			continue
		}
		if ev.Type != evdev.EV_KEY {
			continue
		}
		switch ev.Value {
		case 1:
			c.press()
			last = ev.Time.Nano()
		case 0:
			if last != 0 && ev.Time.Nano()-last > 10000000 {
				c.release()
			} else {
				c.longpress()
			}
			last = 0
		}
	}
}

func main() {
	var cfg []*Config

	var config string
	flag.StringVar(&config, "c", filepath.Join(os.Getenv("HOME"), ".config", "uinputd"), "config file")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	f, err := os.Open(config)
	if err != nil {
		log.Print(err)
	} else {
		yaml.NewDecoder(f).Decode(&cfg)
		f.Close()
	}

	devs, err := evdev.ListInputDevices()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		for _, dev := range devs {
			dev.File.Close()
		}
	}()

	var wg sync.WaitGroup
	n := 0

	for _, dev := range devs {
		if verbose {
			log.Println("found uinput device", dev.Name)
		}
		for _, c := range cfg {
			if dev.Phys == c.Device {
				c.uinput = dev
				wg.Add(1)
				n++
				if verbose {
					log.Println("connected", dev.Name)
				}
				go c.watch(&wg)
			}
		}

	}

	if n == 0 {
		log.Fatal("cannot open uinput device")
	}

	wg.Wait()
}
