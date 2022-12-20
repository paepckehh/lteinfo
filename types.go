package lteinfo

import (
	"bufio"
	"fmt"
	"os"
)

type display_message struct {
	id  int
	msg string
}

type lte struct {
	feed     *bufio.Scanner
	port     *os.File
	portfile string
}

type lteDEVICE struct {
	Manufacturer     string
	Model            string
	Build            string
	Revision         string
	HWVER            string
	ISO              string
	CFG              string
	IMEI             string
	INI              string
	TEMP             string
	Capabilities     string // network depended capabilities
	Port             string // current command port
	Netlock          string
	Status           string // information status colletion _progress
	Fingerprint      string
	Fingerprint_Done bool
	Done             bool // true if all information successfully parsed
}

type lteSTATUS struct {
	RSSI      string
	HCSQ      string
	DSFLOWRPT string
	HFREQINFO string
	Status    string // information status colletion _progress
}

type lteSESSION struct {
	Uptime  string
	TXSpeed string
	RXSpeed string
	TXTotal string
	RXTotal string
	Status  string // information status colletion _progress
}

type lteNETWORK struct {
	SystemMode         string
	RSSI               string
	lte_BAND           string
	lte_RSSI           string
	lte_RSRP           string
	lte_SINR           string
	lte_RSRQ           string
	UL_FC              string
	UL_BW              string
	UL_CF              string
	DL_FC              string
	DL_BW              string
	DL_CF              string
	lteCAT             string
	DTS_MATCH          string
	DHCP_SPEED         string
	Provider_Name      string
	Provider_Interface string
	Provider_Status    string
	Provider_Roaming   bool
	Provider_GSM       string
	Provider_UMTS      string
	Provider_lte       string
	Auth               string
	SMS                string
	Around             string
	LinkQuality        string // graphic representation of RSSI
	Status             string // information status colletion _progress
}

type lteNIC struct {
	UPLINK     string
	DOWNLINK   string
	STACK      string
	DHCP       string
	IP4        string
	IP4D       string
	IP4NM      string
	IP4GW      string
	IP4GWD     string
	IP4DNS1    string
	IP4DNS2    string
	IP4DNSD    string
	IP4NAT     string
	IP4RFC1918 bool
	Status     string // information status colletion _progress
}

type lteSIM struct {
	spn         string
	ccid        string
	issuer      string
	Card_Status string
	Fingerprint string
	Status      string // information status colletion _progress
}

type cmd struct {
	device_clean [9]string
	device_init  [9]string
	device_setup [9]string
	update       [9]string
	sim          [9]string
	slow         [9]string
	real_slow    [9]string
	smsInit      [9]string
	smsGet       [9]string
}

func getSTATUS() *lteSTATUS {
	status := new(lteSTATUS)
	status.RSSI = _defaults
	status.HCSQ = _defaults
	status.HFREQINFO = _defaults
	status.DSFLOWRPT = _defaults
	status.Status = "active"
	return status
}

func getSIM() *lteSIM {
	sim := new(lteSIM)
	sim.ccid = _defaults
	sim.issuer = _defaults
	sim.Card_Status = _defaults
	sim.Status = "waiting for simcard answer"
	return sim
}

func getNET() *lteNETWORK {
	net := new(lteNETWORK)
	net.SystemMode = _defaults_short
	net.RSSI = " 0dBm "
	net.lte_BAND = _defaults_short
	net.lte_RSSI = _defaults
	net.lte_RSRP = _defaults
	net.lte_SINR = _defaults
	net.lte_RSRQ = _defaults
	net.lteCAT = _defaults_short
	net.DTS_MATCH = _defaults
	net.LinkQuality = _defaults
	net.UL_FC = _defaults_short
	net.UL_BW = _defaults_short
	net.UL_CF = _defaults_short
	net.DL_FC = _defaults_short
	net.DL_BW = _defaults_short
	net.DL_CF = _defaults_short
	net.Provider_Name = _defaults
	net.Provider_Status = _defaults
	net.Provider_Interface = _defaults_short
	net.Provider_Roaming = false
	net.Provider_GSM = _defaults
	net.Provider_lte = _defaults
	net.Provider_UMTS = _defaults
	net.Auth = _defaults
	net.SMS = _defaults
	net.Around = "waiting for cell tower interface answer"
	net.Status = "waiting for cell tower interface answer"
	return net
}

func getSESSION() *lteSESSION {
	session := new(lteSESSION)
	session.Uptime = _defaults
	session.TXSpeed = _defaults
	session.RXSpeed = _defaults
	session.TXTotal = _defaults
	session.RXTotal = _defaults
	session.Status = "waiting for connection request"
	return session
}

func getNIC() *lteNIC {
	nic := new(lteNIC)
	nic.UPLINK = _defaults_short
	nic.DOWNLINK = _defaults_short
	nic.STACK = _defaults
	nic.DHCP = _defaults
	nic.IP4D = _defaults
	nic.IP4NAT = _defaults
	nic.IP4GWD = _defaults
	nic.IP4DNSD = _defaults
	nic.Status = "waiting for device"
	return nic
}

func getDEVICE(port string) *lteDEVICE {
	device := new(lteDEVICE)
	device.Manufacturer = _defaults
	device.Build = _defaults
	device.Model = _defaults
	device.Revision = _defaults
	device.HWVER = _defaults
	device.INI = _defaults
	device.ISO = _defaults
	device.CFG = _defaults
	device.IMEI = _defaults
	device.TEMP = _defaults
	device.Capabilities = _defaults
	device.Port = port
	device.Netlock = _defaults
	device.Fingerprint = _defaults
	device.Fingerprint_Done = false
	device.Done = false
	device.Status = "waiting for device"
	return device
}

func getDISPLAY(lte_device string) *[64]string {
	x := new([64]string)
	nic := getNIC()
	net := getNET()
	sim := getSIM()
	status := getSTATUS()
	session := getSESSION()
	device := getDEVICE(lte_device)
	x[0] = "\n\n\n\n\n\n\n\n"
	x[1] = _section_line_
	x[2] = fmt.Sprintf("%s### DEVICE%s [%s%s%s]", _WHITE, _OFF, _GREY, device.Status, _OFF)
	x[3] = fmt.Sprintf("+ Manufacturer:    %s%s", _BLUE, device.Manufacturer)
	x[4] = fmt.Sprintf("+ Model:           %s%s", _BLUE, device.Model)
	x[5] = fmt.Sprintf("  + Build Date:    %s%s", _BLUE, device.Build)
	x[6] = fmt.Sprintf("  + Software:      %s%s", _BLUE, device.Revision)
	x[7] = fmt.Sprintf("  + Image:         %s%s", _BLUE, device.ISO)
	x[8] = fmt.Sprintf("  + Config:        %s%s", _BLUE, device.CFG)
	x[9] = fmt.Sprintf("  + Custom Init:   %s%s", _BLUE, device.INI)
	x[10] = fmt.Sprintf("  + Hardware:      %s%s", _BLUE, device.HWVER)
	x[11] = fmt.Sprintf("  + Temperature:   %s", device.TEMP)
	x[12] = fmt.Sprintf("  + IMEI:          %s", device.IMEI)
	x[13] = fmt.Sprintf("  + Capabilities:  %s%s", _BLUE, device.Capabilities)
	x[14] = fmt.Sprintf("  + Operator Lock: %s", device.Netlock)
	x[15] = fmt.Sprintf("  + Command Port:  %s%s", _BLUE, device.Port)
	x[16] = _section_line_
	x[17] = fmt.Sprintf("%s### SIM CARD%s [%s%s%s]", _WHITE, _OFF, _GREY, sim.Status, _OFF)
	x[18] = fmt.Sprintf("+ State [PIN] : %s", sim.Card_Status)
	x[19] = fmt.Sprintf("+ Issuer      : %s", sim.issuer)
	x[20] = fmt.Sprintf("  + CCID      : %s", sim.ccid)
	x[21] = _section_line_
	x[22] = fmt.Sprintf("%s### CONNECTED CELLTOWER%s [%s%s%s]", _WHITE, _OFF, _GREY, net.Status, _OFF)
	x[23] = fmt.Sprintf("+ Status: %s%s", _BLUE, net.Provider_Status)
	x[24] = fmt.Sprintf("+ Network Operator: %s", net.Provider_Name)
	x[25] = fmt.Sprintf("+ Authentication: %s%s%s [Roaming %s%v%s]", _BLUE, net.Auth, _OFF, _BLUE, net.Provider_Roaming, _OFF)
	x[26] = fmt.Sprintf("+ System Mode: %s%s%s [RAN AirInterface %s%s%s]", _BLUE, net.SystemMode, _OFF, _BLUE, net.Provider_Interface, _OFF)
	x[27] = fmt.Sprintf("  + LTE Category: %s%s%s Frequency Band %s%s", _BLUE, net.lteCAT, _OFF, _BLUE, net.lte_BAND)
	x[28] = fmt.Sprintf("    + Uplink:   [Carrier %s%v%s] [Bandwidth %s%v%s] [EARFCN %s%v%s]", _BLUE, net.UL_CF, _OFF, _BLUE, net.UL_BW, _OFF, _BLUE, net.UL_FC, _OFF)
	x[29] = fmt.Sprintf("    + Downlink: [Carrier %s%v%s] [Bandwidth %s%v%s] [EARFCN %s%v%s]", _BLUE, net.DL_CF, _OFF, _BLUE, net.DL_BW, _OFF, _BLUE, net.DL_FC, _OFF)
	x[30] = fmt.Sprintf("    + Rx Signal Strength [RSSI]: %s%v", _BLUE, net.lte_RSSI)
	x[31] = fmt.Sprintf("    + Rx Signal Power    [RSRP]: %s%v", _BLUE, net.lte_RSRP)
	x[32] = fmt.Sprintf("    + Rx Signal to Noise [SINR]: %s%v", _BLUE, net.lte_SINR)
	x[33] = fmt.Sprintf("  + Cell SMS Infra: %s", net.SMS)
	x[34] = fmt.Sprintf("  + TimeStamp:      %s", net.DTS_MATCH)
	x[35] = fmt.Sprintf("  + LinkQuality:    %v", net.LinkQuality)
	x[36] = _section_line_
	x[37] = fmt.Sprintf("%s### OTHER VISIBLE CELL INTERFACES%s [%s%s%s]", _WHITE, _OFF, _GREY, net.Around, _OFF)
	x[38] = fmt.Sprintf("+ GSM:         %s%s", _CYAN, net.Provider_GSM)
	x[39] = fmt.Sprintf("+ UTRAN UMTS:  %s%s", _CYAN, net.Provider_UMTS)
	x[40] = fmt.Sprintf("+ E-UTRAN LTE: %s%s", _CYAN, net.Provider_lte)
	x[41] = _section_line_
	x[42] = fmt.Sprintf("%s### INTERFACE NDIS/DHCP%s [%s%s%s]", _WHITE, _OFF, _GREY, nic.Status, _OFF)
	x[43] = fmt.Sprintf("+ Supported Network Stack(s): %s%s", _GREEN, nic.STACK)
	x[44] = fmt.Sprintf("  + IPv4 Mode   : %s", nic.IP4NAT)
	x[45] = fmt.Sprintf("  + IPv4 Address: %s%s", _BLUE, nic.IP4D)
	x[46] = fmt.Sprintf("  + IPv4 Gateway: %s%s", _BLUE, nic.IP4GWD)
	x[47] = fmt.Sprintf("  + IPv4 DNS SRV: %s%s", _BLUE, nic.IP4DNSD)
	x[48] = _section_line_
	x[49] = fmt.Sprintf("%s### CURRENT SESSION%s [%s%s%s]", _WHITE, _OFF, _GREY, session.Status, _OFF)
	x[50] = fmt.Sprintf("+ Uptime: %s%v", _GREEN, session.Uptime)
	x[51] = fmt.Sprintf("  + Tx speed: %s%v", _CYAN, session.TXSpeed)
	x[52] = fmt.Sprintf("  + Rx speed: %s%v", _CYAN, session.RXSpeed)
	x[53] = fmt.Sprintf("  + Tx total: %s%v", _CYAN, session.TXTotal)
	x[54] = fmt.Sprintf("  + Rx total: %s%v", _CYAN, session.RXTotal)
	x[55] = fmt.Sprintf("+ Bandwidth:  [Uplink Channel %s%s%s] [Downlink Channel %s%s%s]", _BLUE, nic.UPLINK, _OFF, _BLUE, nic.DOWNLINK, _OFF)
	x[56] = _section_line_
	x[57] = fmt.Sprintf("%s### RAW SENTENCE%s [%s%s%s]", _WHITE, _OFF, _GREY, status.Status, _OFF)
	x[58] = fmt.Sprintf("+ RSSI:      %s%s", _GREY, status.RSSI)
	x[59] = fmt.Sprintf("+ HCSQ:      %s%s", _GREY, status.HCSQ)
	x[60] = fmt.Sprintf("+ DHCP:      %s%s", _GREY, nic.DHCP)
	x[61] = fmt.Sprintf("+ HFREQINFO: %s%s", _GREY, status.HFREQINFO)
	x[62] = fmt.Sprintf("+ DSFLOWRPT: %s%s", _GREY, status.DSFLOWRPT)
	x[63] = _section_line_
	return x
}
