package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"net/http"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

var (
	hkpin string
	hkserial string

	items []DisplayItem
)

type DisplayItem struct {
	Name     string
	HkDevice *accessory.Outlet
}

func main() {
	// Load in flags
	parseFlags()

	// Start Up HomeKit
	fmt.Println("Stating HomeKit Server with Pin: ", hkpin)
	go setupHomekit()

	// Start up web server
	fmt.Println("Setting up pages")
	buildWebHandlers()

	// Run server
	fmt.Println("Stating server on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func parseFlags() {
	// Define Flags
	flag.StringVar(&hkpin, "hkpin", "00102003", "Pin to use for HomeKit Authentication")
	flag.StringVar(&hkserial, "hkserial", "027TC-000001", "Serial number to expose to homekit")
	item_strings := flag.String("items", "Bins,Recycling", "A comma seperated list of items")

	// Parese flags
	flag.Parse()

	// Compute Items
	ia := strings.Split(*item_strings, ",")
	for  _, i := range ia {
		h := accessory.NewOutlet(accessory.Info{Name: i})
		items = append(items, DisplayItem{
			Name:  i,
			HkDevice: h,
		})
		h.Outlet.On.SetValue(false)

	}
}

func setupHomekit() (error) {
	binfo := accessory.Info{
		Name: "Termovision",
		SerialNumber: hkserial,
		Manufacturer: "Tom Cousin",
		Model: "V0.1",
	}

	bacc := accessory.NewBridge(binfo)

	d := []*accessory.Accessory{}
	for _, ad := range items {
		fmt.Println("Adding Accessory: ", ad.Name)
		d = append(d, ad.HkDevice.Accessory)
	}

	config := hc.Config{
		Pin: hkpin,
	}
	t, err := hc.NewIPTransport(config, bacc.Accessory, d...)
	if err != nil {
		fmt.Println("Homekit Error: ", err)
		return err
	}
	hc.OnTermination(func(){
		<-t.Stop()
	})

	t.Start()

	fmt.Println("Homekit started")

	return nil
}

func buildWebHandlers() {
	http.HandleFunc("/", HomeHandler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	out := "["
	for _, d := range items {
		if d.HkDevice.Outlet.On.GetValue() != true {
			continue
		}
		out += "\""
		out += d.Name
		out += "\""
		out += ","
	}
	if out != "[" {
		out = out[:len(out)-1]
	}
	out += "]"
	fmt.Fprintf(w, out)
}
