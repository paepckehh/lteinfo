// package main
package main

// import
import (
	"syscall"

	"paepcke.de/lteinfo"
)

// main ...
func main() {
	c := lteinfo.Config{
		DevicePort:       "/dev/lte0",
		DeviceModel:      "HUAWEI_E3372",
		DeviceKnownList:  "4BD2-XO-TK-WFQR,RQJD-MX-NW-UFO7,5TU4-C5-HA-7J4J",
		SimCardKnownList: "YMI5-7T-J7-RW3N,TYYD-6S-OC-VPV5,LLNX-AL-AQ-6AE2",
		EcoMode:          "normal",
	}
	switch {
	case isEnv("SMS"):
		c.Sms()
	default:
		c.Stats()
	}
}

// isEnv
func isEnv(in string) bool {
	if _, ok := syscall.Getenv(in); ok {
		return true
	}
	return false
}
