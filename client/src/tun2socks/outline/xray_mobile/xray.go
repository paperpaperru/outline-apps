package xrayMobile

import (
	"fmt"
	"github.com/xtls/libxray"
	"os"
)

func writeConfig() {
	config := []byte("") // config here
	os.WriteFile("config.json", config, 0644)
}

func StartXrayServer() {
	libXray.LoadGeoData(".")
	writeConfig()
	fmt.Println("writeConfig done")
	libXray.RunXray(".", "config.json", 13*1000*1000)
}

func StopXrayServer() {
	libXray.StopXray()
}
