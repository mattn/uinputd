package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	"github.com/gvalkov/golang-evdev"
)

type Config struct {
	Device    string `yaml:"device"`
	Press     string `yaml:"press"`
	Release   string `yaml:"release"`
	LongPress string `yaml:"longpress"`
	verbose   bool
}

func (c *Config) exec(n, s string) {
	if c.verbose {
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

func main() {
	var cfg Config

	var device, config string
	var verbose bool
	flag.StringVar(&config, "c", "", "config file")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	if config != "" {
		f, err := os.Open(config)
		if err != nil {
			log.Fatal(err)
		}
		yaml.NewDecoder(f).Decode(&cfg)
		f.Close()
	}
	cfg.verbose = verbose

	devs, err := evdev.ListInputDevices()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		for _, dev := range devs {
			dev.File.Close()
		}
	}()

	var input *evdev.InputDevice
	if device != "" {
		for _, dev := range devs {
			if verbose {
				log.Println(dev.Name)
			}
			if input == nil && dev.Name == device {
				input = dev
			}
		}
	} else if len(devs) > 0 {
		input = devs[0]
	}

	if input == nil {
		log.Fatal("cannot open uinput device")
	}

	last := int64(0)
	for {
		ev, err := input.ReadOne()
		if err != nil {
			log.Fatal(err)
		}
		if ev.Type != evdev.EV_KEY {
			continue
		}
		switch ev.Value {
		case 1:
			cfg.press()
			last = ev.Time.Nano()
		case 0:
			if last != 0 && ev.Time.Nano()-last > 10000000 {
				cfg.release()
			} else {
				cfg.longpress()
			}
			last = 0
		}
	}
}
