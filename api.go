// package lteinfo ...
package lteinfo

// Config ...
type Config struct {
	DevicePort       string // device port, eg /dev/lte0
	DeviceModel      string // device mode, eg HUAWEI_E3372
	DeviceKnownList  string // known devices, hashlist
	SimCardKnownList string // known simcards, hashlist
	EcoMode          string // eco mode
}

// Stats ...
func (c *Config) Stats() {
	c.displayStats()
}

// Sms ...
func (c *Config) Sms() {
	c.sms()
}
