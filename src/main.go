package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

var (
	hkpin string
	hkserial string
)

func main() {
	// Load in flags
	parseFlags()

	// Start Up HomeKit
	err := setupHomekit()
	if err != nil {
		return
	}
}

func parseFlags() {
	// Define Flags
	flag.StringVar(&hkpin, "hkpin", "00102003", "Pin to use for HomeKit Authentication")
	flag.StringVar(&hkserial, "hkserial", "027TC-000001", "Serial number to expose to homekit")

	// Parese flags
	flag.Parse()
}

func setupHomekit() (error) {
	binfo := accessory.Info{
		Name: "Termovision",
		SerialNumber: hkserial,
		Manufacturer: "Tom Cousin",
		Model: "V0.1",
	}

	bacc := accessory.NewBridge(binfo)

	config := hc.Config{
		Pin: hkpin,
	}
	t, err := hc.NewIPTransport(config, bacc.Accessory)
	if err != nil {
		fmt.Println("Homekit Error: ", err)
		return err
	}
	hc.OnTermination(func(){
		<-t.Stop()
	})

	t.Start()

	return nil
}
