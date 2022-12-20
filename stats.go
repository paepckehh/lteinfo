// package lteinfo paepcke.de/lteinfo (2022)
// Currently supported devices:
//
//	HUAWEI (full)
//	+ E3372*      [ USB alternate profile vendor: 0x12d1 / product: 0x1505 / NCM Mode / FBSD: MSC_EJECT_HUAWEI2 ]
//	HUAWEI (minimal support via E3372 mode)
//	+ E5573CS-322 [ USB alternate profile vendor: 0x12d1 / product: 0x155e / NCM Mode / FBSD: MSC_EJECT_HUAWEI4 ]
//	+ E5573CS-322 [ USB alternate profile vendor: 0x12d1 / product: 0x1442 / TTY Mode / FBSD: MSC_EJECT_HUAWEI  ]
package lteinfo // import paepcke.de/lteinfo
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (

	// some internal _defaults
	ipv4_OCTETS  int           = 4
	one          time.Duration = time.Duration(1 * time.Second)
	nwTimeDiffOK time.Duration = time.Duration(2650 * time.Millisecond)
	nwTimeLayout string        = "06/01/02,15:04:05"
)

func (c *Config) displayStats() {
	var (
		line, indicator  string
		u, uu, uuu       uint64
		txspeed, rxspeed uint64
		txtotal, rxtotal uint64
		m                display_message
		ts               time.Time
		diff, uptime     time.Duration
	)

	// eco mode config
	eco_sleep, display_refresh_rate := ecomode(c.EcoMode)

	// command and inbound sentence channel
	channelDisplay := make(chan display_message, 200)
	channelParser := make(chan string, 100)
	channelCmd := make(chan [9]string, 1)
	channelOut := make(chan string, 10)
	channelFingerprint := make(chan bool, 1)

	// spin up background bufio fetcher process
	go func() {
		go func() {
			var (
				l             int
				devRead       lte
				sentence_line string
			)
			// open / buffer [bufio] the device stream
			devRead = c.getLTE()
			for devRead.feed.Scan() {
				sentence_line = devRead.feed.Text()
				l = len(sentence_line)
				if l < 12 || l > 128 {
					continue
				}
				if !strings.Contains(sentence_line, ":") {
					continue
				}
				channelParser <- sentence_line
				fmt.Printf("%s#%s", _GREY, _OFF)
				time.Sleep(eco_sleep)
			}
		}()
	}()

	// spin up the channelCmd pacemaker process
	go func() {
		var eco_cmd_tasker time.Duration
		devWrite := c.getLTE()
		for cmd := range channelCmd {
			for i := 0; i < 9; i++ {
				switch cmd[i] {
				case "", "x":
					continue
				case "w":
					time.Sleep(one)
					continue
				case "W":
					time.Sleep(one * 2)
					continue
				}
				fmt.Printf("%s# AT%s #%s", _GREY, cmd[i], _OFF)
				if _, err := devWrite.port.WriteString("AT" + cmd[i] + "\r\n"); err != nil {
					fmt.Printf("\n****** ERROR: command [%v] -> [%v] ********\n", i, err)
					devWrite.port.Close()
					fmt.Print("############ LTE DEVICE GONE!  ################\n")
					fmt.Print("############     [EXIT]        ################\n\n")
					panic("")
				}
				time.Sleep(one)
			}
			time.Sleep(one*1 + eco_cmd_tasker)
			eco_cmd_tasker = eco_cmd_tasker + (eco_sleep * 3)
		}
	}()

	// spin up background command scheduler process
	go func() {
		cmd := getCommand(c.DeviceModel)
		var (
			fingerprint_done  bool
			catch_all         uint
			slow_updates      uint = 10
			real_slow_updates uint = 10
		)
		channelCmd <- cmd.device_clean
		channelCmd <- cmd.device_init
		channelCmd <- cmd.sim
		channelCmd <- cmd.device_setup
		for {
			channelCmd <- cmd.update
			slow_updates++
			if slow_updates > 3 {
				channelCmd <- cmd.slow
				slow_updates = 0
				real_slow_updates++
				if real_slow_updates > 6 {
					channelCmd <- cmd.real_slow
					real_slow_updates = 0
					catch_all++
					if catch_all > 8 {
						if !fingerprint_done {
							channelCmd <- cmd.device_clean
							time.Sleep(one * 2)
							channelCmd <- cmd.device_init
							fingerprint_done = <-channelFingerprint
						}
						channelCmd <- cmd.sim
						time.Sleep(one * 2)
						channelCmd <- cmd.device_setup
						time.Sleep(one * 2)
						catch_all = 0
					}
				}
			}
			time.Sleep(one * 2)
		}
	}()

	// setup container _defaults
	nic := getNIC()
	net := getNET()
	sim := getSIM()
	status := getSTATUS()
	session := getSESSION()
	device := getDEVICE(c.DevicePort)

	// spin up display frame builder engine
	go func() {
		var (
			frame   uint64
			inbound display_message
			b       strings.Builder
		)
		d := getDISPLAY(c.DevicePort)
		r := time.Now()
		for inbound = range channelDisplay {
			d[inbound.id] = inbound.msg
			if time.Now().Unix() > r.Add(display_refresh_rate).Unix() {
				frame++
				for i := 0; i < len(d); i++ {
					b.WriteString(d[i])
					b.WriteString(_clean_newline)
				}
				fmt.Fprintf(&b, "%s[Frame %v] ", _GREY, frame)
				channelOut <- b.String()
				b.Reset()
				r = time.Now()
			}
			time.Sleep(eco_sleep)
		}
	}()

	// spin up display out writer
	go func() {
		for s := range channelOut {
			os.Stdout.Write([]byte(s))
		}
	}()

	// main parser loop
	var l int
	for line = range channelParser {
		l = len(line)
		switch {
		case line[:2] == "^D":
			switch {
			case line[:11] == "^DSFLOWRPT:":
				status.DSFLOWRPT = line[11:]
				m.id = 62
				m.msg = fmt.Sprintf("+ DSFLOWRPT: %s%s", _GREY, status.DSFLOWRPT)
				channelDisplay <- m
				s := strings.Split(status.DSFLOWRPT, ",")
				if len(s) == 7 {
					// decode DSFLOWRPT for current session uptime
					u, _ = strconv.ParseUint(s[0], 16, 64)
					uptime = time.Duration(u) * time.Second
					session.Uptime = uptime.String()
					m.id = 50
					m.msg = fmt.Sprintf("+ Uptime: %s%v", _GREEN, session.Uptime)
					channelDisplay <- m
					// decode DSFLOWRPT components for LinkSpeed
					txspeed, _ = strconv.ParseUint(s[1], 16, 64)
					session.TXSpeed = hruIEC(txspeed, "bytes/sec")
					m.id = 51
					m.msg = fmt.Sprintf("  + Tx speed: %s%v", _CYAN, session.TXSpeed)
					channelDisplay <- m
					rxspeed, _ = strconv.ParseUint(s[2], 16, 64)
					session.RXSpeed = hruIEC(rxspeed, "bytes/sec")
					m.id = 52
					m.msg = fmt.Sprintf("  + Rx speed: %s%v", _CYAN, session.RXSpeed)
					channelDisplay <- m
					// decode DSFLOWRPT components for Transfer Volume
					txtotal, _ = strconv.ParseUint(s[3], 16, 64)
					session.TXTotal = hruIEC(txtotal, "bytes")
					m.id = 53
					m.msg = fmt.Sprintf("  + Tx total: %s%v", _CYAN, session.TXTotal)
					channelDisplay <- m
					rxtotal, _ = strconv.ParseUint(s[4], 16, 64)
					session.RXTotal = hruIEC(rxtotal, "bytes")
					m.id = 54
					m.msg = fmt.Sprintf("  + Rx total: %s%v", _CYAN, session.RXTotal)
					channelDisplay <- m
					session.Status = "active"
					m.id = 49
					m.msg = fmt.Sprintf("%s### CURRENT SESSION%s [%s%s%s]", _WHITE, _OFF, _GREY, session.Status, _OFF)
					channelDisplay <- m
				}
			case line[:7] == "^DHCP: ":
				nic.DHCP = line[7:]
				m.id = 60
				m.msg = fmt.Sprintf("+ DHCP:      %s%s", _GREY, nic.DHCP)
				channelDisplay <- m
				s := strings.Split(nic.DHCP, ",")
				if len(s) == 8 {
					u, _ = strconv.ParseUint(s[6], 10, 64)
					nic.UPLINK = hruSI(u, "bit")
					m.id = 55
					m.msg = fmt.Sprintf("+ Bandwidth:  [Uplink Channel %s%s%s] [Downlink Channel %s%s%s]", _BLUE, nic.UPLINK, _OFF, _BLUE, nic.DOWNLINK, _OFF)
					channelDisplay <- m
					u, _ = strconv.ParseUint(s[7], 10, 64)
					nic.DOWNLINK = hruSI(u, "bit")
					nic.IP4 = hex2IP4(s[0])
					if strings.HasPrefix(nic.IP4, "192") || strings.HasPrefix(nic.IP4, "10") {
						nic.IP4RFC1918 = true
						nic.IP4NAT = fmt.Sprintf("%s[Provider NAT] [RFC1918 %v]%s", _BLUE, nic.IP4RFC1918, _OFF)
					} else {
						nic.IP4RFC1918 = false
						nic.IP4NAT = fmt.Sprintf("%s[direct]%s", _BLUE, _OFF)
					}
					m.id = 44
					m.msg = fmt.Sprintf("  + IPv4 Mode   : %s", nic.IP4NAT)
					channelDisplay <- m
					nic.IP4NM = hex2IP4(s[1])
					nic.IP4GW = hex2IP4(s[2])
					nic.IP4DNS1 = hex2IP4(s[4])
					nic.IP4DNS2 = hex2IP4(s[5])
					nic.IP4GWD = fmt.Sprintf("[%s]", nic.IP4GW)
					m.id = 46
					m.msg = fmt.Sprintf("  + IPv4 Gateway: %s%s", _BLUE, nic.IP4GWD)
					channelDisplay <- m
					nic.IP4D = fmt.Sprintf("[%s] [MASK %s]", nic.IP4, nic.IP4NM)
					m.id = 45
					m.msg = fmt.Sprintf("  + IPv4 Address: %s%s", _BLUE, nic.IP4D)
					channelDisplay <- m
					nic.IP4DNSD = fmt.Sprintf("[%s] [%s]", nic.IP4DNS1, nic.IP4DNS2)
					m.id = 47
					m.msg = fmt.Sprintf("  + IPv4 DNS SRV: %s%s", _BLUE, nic.IP4DNSD)
					channelDisplay <- m
					nic.Status = "active"
					m.id = 42
					m.msg = fmt.Sprintf("%s### INTERFACE NDIS/DHCP%s [%s%s%s]", _WHITE, _OFF, _GREY, nic.Status, _OFF)
					channelDisplay <- m
				}
			}
		case line[:2] == "^H":
			switch {
			case line[:6] == "^HCSQ:":
				status.HCSQ = strings.ReplaceAll(line[6:], "\"", "")
				m.id = 59
				m.msg = fmt.Sprintf("+ HCSQ:      %s%s", _GREY, status.HCSQ)
				channelDisplay <- m
				s := strings.Split(status.HCSQ, ",")
				if len(s) > 0 {
					net.SystemMode = "[" + s[0] + "]"
					m.id = 26
					m.msg = fmt.Sprintf("+ System Mode: %s%s%s [RAN AirInterface %s%s%s]", _BLUE, net.SystemMode, _OFF, _BLUE, net.Provider_Interface, _OFF)
					channelDisplay <- m
					// we got an lte mode sentence
					if len(s) == 5 {
						// decode HCSQ to RSSI
						u, _ = strconv.ParseUint(s[1], 10, 64)
						indicator = int2indicator(120-u, 60, "|")
						net.lte_RSSI = fmt.Sprintf("%s %s[-%v dBm]%s", indicator, _CYAN, 120-u, _OFF)
						m.id = 30
						m.msg = fmt.Sprintf("    + Rx Signal Strength [RSSI]: %s%v%s", _BLUE, net.lte_RSSI, _OFF)
						channelDisplay <- m
						// decode HCSQ to RSRP
						u, _ = strconv.ParseUint(s[2], 10, 64)
						indicator = int2indicator(140-u, 60, "|")
						net.lte_RSRP = fmt.Sprintf("%s %s[-%v dBm]%s", indicator, _CYAN, 140-u, _OFF)
						m.id = 31
						m.msg = fmt.Sprintf("    + Rx Signal Power    [RSRP]: %s%v%s", _BLUE, net.lte_RSRP, _OFF)
						channelDisplay <- m
						// decode HCSQ to SINR
						u, _ = strconv.ParseUint(s[4], 10, 64)
						indicator = int2indicator(20-(u/2), 10, "|||")
						net.lte_SINR = fmt.Sprintf("%s %s[-%v dB]%s", indicator, _CYAN, 20-(u/2), _OFF)
						m.id = 32
						m.msg = fmt.Sprintf("    + Rx Signal to Noise [SINR]: %s%v%s", _BLUE, net.lte_SINR, _OFF)
						channelDisplay <- m
					}
					net.Status = "active"
					m.id = 22
					m.msg = fmt.Sprintf("%s### CONNECTED CELLTOWER%s [%s%s%s]", _WHITE, _OFF, _GREY, net.Status, _OFF)
					channelDisplay <- m
				}
			case line[:11] == "^HFREQINFO:":
				status.HFREQINFO = line[11:]
				m.id = 61
				m.msg = fmt.Sprintf("+ HFREQINFO: %s%s%s", _GREY, status.HFREQINFO, _OFF)
				channelDisplay <- m
				s := strings.Split(status.HFREQINFO, ",")
				if len(s) == 9 {
					if s[1] == "6" {
						net.SystemMode = "[LTE]"
						m.id = 26
						m.msg = fmt.Sprintf("+ System Mode: %s%s%s [RAN AirInterface %s%s%s]", _BLUE, net.SystemMode, _OFF, _BLUE, net.Provider_Interface, _OFF)
						channelDisplay <- m
						net.lte_BAND = "[" + s[2] + "]"
						m.id = 27
						m.msg = fmt.Sprintf("  + LTE Category: %s%s%s Frequency Band %s%s%s", _BLUE, net.lteCAT, _OFF, _BLUE, net.lte_BAND, _OFF)
						channelDisplay <- m
						// uplink Frequencies
						net.DL_FC = s[3]
						u, _ = strconv.ParseUint(s[4], 10, 64)
						net.DL_CF = hruSI(u*100000, "Hz")
						u, _ = strconv.ParseUint(s[5], 10, 64)
						net.DL_BW = hruSI(u*1000, "Hz")
						m.id = 29
						m.msg = fmt.Sprintf("    + Downlink: [Carrier %s%v%s] [Bandwidth %s%v%s] [EARFCN %s%v%s]", _BLUE, net.DL_CF, _OFF, _BLUE, net.DL_BW, _OFF, _BLUE, net.DL_FC, _OFF)
						channelDisplay <- m
						// downlink Frequencies
						net.UL_FC = s[6]
						u, _ = strconv.ParseUint(s[7], 10, 64)
						net.UL_CF = hruSI(u*100000, "Hz")
						u, _ = strconv.ParseUint(s[8], 10, 64)
						net.UL_BW = hruSI(u*1000, "Hz")
						m.id = 28
						m.msg = fmt.Sprintf("    + Uplink:   [Carrier %s%v%s] [Bandwidth %s%v%s] [EARFCN %s%v%s]", _BLUE, net.UL_CF, _OFF, _BLUE, net.UL_BW, _OFF, _BLUE, net.UL_FC, _OFF)
						channelDisplay <- m
					}
				}
			}
		case line[:6] == "^RSSI:":
			status.RSSI = line[6:]
			m.id = 58
			m.msg = fmt.Sprintf("+ RSSI:      %s%s", _GREY, status.RSSI)
			channelDisplay <- m
			// decode RSSI as graphic LinkQuality status indicator
			u, _ = strconv.ParseUint(status.RSSI, 10, 8)
			// net.RSSI = fmt.Sprintf(" -%v dBm", 140-u)
			net.LinkQuality = fmt.Sprintf("%s %s[-%v dB]%s", int2indicator(140-u/2, 60, "|"), _CYAN, 140-u, _OFF)
			m.id = 35
			m.msg = fmt.Sprintf("  + LinkQuality:    %v", net.LinkQuality)
			channelDisplay <- m
		case line[:2] == "^N":
			switch {
			case line[:8] == "^NWTIME:":
				dts := time.Now()
				s := strings.Split(line[8:], "+")
				ts, _ = time.Parse(nwTimeLayout, s[0])
				diff = dts.Sub(ts)
				if diff < nwTimeDiffOK {
					net.DTS_MATCH = fmt.Sprintf("%s [Celltower vs. Reference: %s%s%s]", _ok, _BLUE, dts.Format(time.RFC3339), _OFF)
				} else {
					net.DTS_MATCH = fmt.Sprintf("%s %sDIFF %v", _alert, _ALERT, diff)
				}
				m.id = 34
				m.msg = fmt.Sprintf("  + TimeStamp:      %s", net.DTS_MATCH)
				channelDisplay <- m
			case line[:13] == "^NDISSTATQRY:":
				s := strings.Split(strings.ReplaceAll(line[13:], "\"", ""), ",")
				if len(s) == 4 {
					nic.STACK = "[" + s[3] + "]"
					m.id = 43
					m.msg = fmt.Sprintf("+ Supported Network Stack(s): %s%s", _GREEN, nic.STACK)
					channelDisplay <- m
				}
			}
		case line[:8] == "^LTECAT:":
			net.lteCAT = "[" + line[8:] + "]"
			m.id = 27
			m.msg = fmt.Sprintf("  + LTE Category: %s%s%s Frequency Band %s%s", _BLUE, net.lteCAT, _OFF, _BLUE, net.lte_BAND)
			channelDisplay <- m
		case line[:11] == "^CARDLOCK: ":
			s := strings.Split(line[11:], ",")
			if len(s) > 2 {
				device.Netlock = fmt.Sprintf("%s %s[ATTEMPTS OPEN %v] [OPERATORCODE %v]%s", id2netlock(s[0]), _GREY, s[1], s2s(s[2]), _OFF)
				m.id = 14
				m.msg = fmt.Sprintf("  + Operator Lock: %s", device.Netlock)
				channelDisplay <- m
			}
		case line[:7] == "^ICCID:":
			if len(line) == 28 {
				s := line[7:]
				sim.ccid = fmt.Sprintf("MajorID %s%s%s CountryCode %s%s%s IssuerID %s%s%s AccountID %s%s%s CD %s%s%s", _BLUE, s[0:3], _OFF, _BLUE, s[3:5], _OFF, _BLUE, s[5:7], _OFF, _BLUE, s[7:19], _OFF, _BLUE, s[19:], _OFF)
				m.id = 20
				m.msg = fmt.Sprintf("  + CCID      : %s", sim.ccid)
				channelDisplay <- m
				sim.Fingerprint = s2b32f(s)
				if strings.Contains(c.SimCardKnownList, sim.Fingerprint) {
					m.msg = fmt.Sprintf("%s### SIM CARD %s[OK]%s [%s]", _WHITE, _ALERT_G, _GREY, sim.Fingerprint)
				} else {
					m.msg = fmt.Sprintf("%s### SIM CARD %sUNLISTED [%s]", _WHITE, _RED, sim.Fingerprint)
				}
				m.id = 17
				channelDisplay <- m
			}
		case line[:5] == "^SPN:":
			s := strings.Split(strings.ReplaceAll(line[5:], "\"", ""), ",")
			if strings.HasPrefix(s[2], "FF") {
				sim.spn = "n/a"
				sim.issuer = _ALERT + "[" + sim.spn + "]" + _OFF
			} else {
				sim.spn = s[2]
				sim.issuer = _BLUE + "[" + sim.spn + "]" + _OFF
			}
			m.id = 19
			m.msg = fmt.Sprintf("+ Issuer      : %s", sim.issuer)
			channelDisplay <- m
			switch sim.spn {
			case net.Provider_Name:
				net.Auth = fmt.Sprintf("[ESP] [Network Operator SPN %s]", sim.spn)
			default:
				if net.Provider_Roaming {
					net.Auth = fmt.Sprintf("[ESP/MVNO] [Roaming Access via SPN %s]", sim.spn)
				} else {
					net.Auth = fmt.Sprintf("[MVNO] [Virtual Operator via SPN %s]", sim.spn)
				}
			}
			m.id = 25
			m.msg = fmt.Sprintf("+ Authentication: %s%s%s [Roaming %s%v%s]", _BLUE, net.Auth, _OFF, _BLUE, net.Provider_Roaming, _OFF)
			channelDisplay <- m
		case line[0] == '+':
			switch {
			case line[:6] == "+COPS:":
				s := strings.Split(strings.ReplaceAll(line[6:], "\"", ""), ",")
				if len(s) == 4 {
					net.Provider_Name = _CYAN + "[" + s[2] + "]" + _OFF
					m.id = 24
					m.msg = fmt.Sprintf("+ Network Operator: %s", net.Provider_Name)
					channelDisplay <- m
					net.Provider_Interface = id2airInterface(s[3])
				} else {
					net.Provider_lte = ""
					net.Provider_GSM = ""
					net.Provider_UMTS = ""
					a := strings.TrimPrefix(line, "+COPS: ")
					b := strings.ReplaceAll(a, "\"", "")
					c := strings.ReplaceAll(b, "(", "")
					s = strings.Split(c, "),")
					for ii := 0; ii < len(s); ii++ {
						ss := strings.Split(s[ii], ",")
						if len(ss) == 5 {
							switch ss[4] {
							case "0":
								net.Provider_GSM = net.Provider_GSM + "[" + ss[1] + "] "
							case "2":
								net.Provider_UMTS = net.Provider_UMTS + "[" + ss[1] + "] "
							case "7":
								net.Provider_lte = net.Provider_lte + "[" + ss[1] + "] "
							}
						}
						net.Around = "active"
						m.id = 37
						m.msg = fmt.Sprintf("%s### OTHER VISIBLE CELL INTERFACES%s [%s%s%s]", _WHITE, _OFF, _GREY, net.Around, _OFF)
						channelDisplay <- m
					}
					if net.Provider_lte == "" {
						net.Provider_lte = _defaults_short
					}
					m.id = 40
					m.msg = fmt.Sprintf("+ E-UTRAN LTE: %s%s", _CYAN, net.Provider_lte)
					channelDisplay <- m
					if net.Provider_GSM == "" {
						net.Provider_GSM = _defaults_short
					}
					m.id = 38
					m.msg = fmt.Sprintf("+ GSM:         %s%s", _CYAN, net.Provider_GSM)
					channelDisplay <- m
					if net.Provider_UMTS == "" {
						net.Provider_UMTS = _defaults_short
					}
					m.id = 39
					m.msg = fmt.Sprintf("+ UTRAN UMTS:  %s%s", _CYAN, net.Provider_UMTS)
					channelDisplay <- m
				}
			case line[:6] == "+CREG:":
				s := strings.Split(line[6:], ",")
				if len(s) == 2 {
					net.Provider_Status, net.Provider_Roaming = id2provider(s[1])
					m.id = 23
					m.msg = fmt.Sprintf("+ Status: %s%s", _BLUE, net.Provider_Status)
					channelDisplay <- m
				}
			case line[:7] == "+CPIN: ":
				switch line[7:] {
				case "READY":
					sim.Card_Status = _GREEN + "[UNLOCKED]" + _OFF
				default:
					sim.Card_Status = _ALERT + "[LOCKED]" + _OFF
				}
				m.id = 18
				m.msg = fmt.Sprintf("+ State [PIN] : %s", sim.Card_Status)
				channelDisplay <- m
			case line[:7] == "+CSMS: ":
				s := strings.Split(line[7:], ",")
				if len(s) == 4 {
					net.SMS = fmt.Sprintf("[GSM %s] [Send %s] [Receive %s] [Cell Broadcast %s]", s2b(s[0]), s2b(s[1]), s2b(s[2]), s2b(s[3]))
					m.id = 33
					m.msg = fmt.Sprintf("  + Cell SMS Infra: %s", net.SMS)
					channelDisplay <- m
				}
			}
		case line[:11] == "^CHIPTEMP: ":
			s := strings.Split(line[11:], ",")
			if len(s) == 5 {
				u, _ = strconv.ParseUint(s[0], 10, 64)
				uu, _ = strconv.ParseUint(s[1], 10, 64)
				uuu, _ = strconv.ParseUint(s[3], 10, 64)
				device.TEMP = fmt.Sprintf("CPU %s%v%sC PSU %s%v%sC CASE %s%v%sC", temp2color(u/10, 42), u/10, _OFF, temp2color(uu/10, 42), uu/10, _OFF, temp2color(uuu, 39), s[3], _OFF)
				m.id = 11
				m.msg = fmt.Sprintf("  + Temperature:   %s", device.TEMP)
				channelDisplay <- m
			}
		} // end main parser loop switch
		// hw/fw/sim related response (should) never change during a session, so stop parsing for it
		if !device.Fingerprint_Done {
			switch {
			case line[:2] == "^V":
				if l > 12 {
					switch {
					case line[:13] == "^VERSION:INI:":
						device.INI = line[13:]
						if device.INI == "" {
							device.INI = fmt.Sprintf("%s[OK]%s [NONE]", _GREEN, _OFF)
						} else {
							device.INI = fmt.Sprintf("%s %s%s%s", _alert, _ALERT, device.INI, _OFF)
						}
						m.id = 9
						m.msg = fmt.Sprintf("  + Custom Init:   %s%s", _BLUE, device.INI)
						channelDisplay <- m
					default:
					}
				}
				if l > 13 {
					switch {
					case line[:13] == "^VERSION:CFG:":
						device.CFG = line[13:]
						m.id = 8
						m.msg = fmt.Sprintf("  + Config:        %s%s", _BLUE, device.CFG)
						channelDisplay <- m
					case line[:13] == "^VERSION:BDT:":
						device.Build = line[13:]
						m.id = 5
						m.msg = fmt.Sprintf("  + Build Date:    %s%s", _BLUE, device.Build)
						channelDisplay <- m
					}
				}
				if l > 14 {
					switch {
					case line[:14] == "^VERSION:EXTU:":
						device.Model = line[14:]
						m.id = 4
						m.msg = fmt.Sprintf("+ Model:           %s%s", _BLUE, device.Model)
						channelDisplay <- m
					case line[:14] == "^VERSION:EXTH:":
						device.HWVER = line[14:]
						m.id = 10
						m.msg = fmt.Sprintf("  + Hardware:      %s%s", _BLUE, device.HWVER)
						channelDisplay <- m
					case line[:14] == "^VERSION:EXTD:":
						device.ISO = line[14:]
						m.id = 7
						m.msg = fmt.Sprintf("  + Image:         %s%s", _BLUE, device.ISO)
						channelDisplay <- m
					case line[:14] == "^VERSION:EXTS:":
						device.Revision = line[14:]
						m.id = 6
						m.msg = fmt.Sprintf("  + Software:      %s%s", _BLUE, device.Revision)
						channelDisplay <- m
					}
				}
			case line[:5] == "IMEI:":
				s := line[5:]
				if len(s) > 14 {
					device.IMEI = fmt.Sprintf("TAC %s%s-%s%s SNR %s%s%s CD %s%s%s", _BLUE, s[0:2], s[2:8], _OFF, _BLUE, s[8:14], _OFF, _BLUE, s[14:], _OFF)
					m.id = 12
					m.msg = fmt.Sprintf("  + IMEI:          %s", device.IMEI)
					channelDisplay <- m
				}
			case line[:7] == "+GCAP: ":
				device.Capabilities = strings.ReplaceAll(line[7:], "+", "")
				m.id = 13
				m.msg = fmt.Sprintf("  + Capabilities:  %s%s", _BLUE, device.Capabilities)
				channelDisplay <- m
			}
			if l > 14 {
				switch {
				case line[:14] == "Manufacturer: ":
					device.Manufacturer = strings.ToUpper(line[14:])
					device.Status = "active"
					m.id = 2
					m.msg = fmt.Sprintf("%s### DEVICE%s [%s%s%s]", _WHITE, _OFF, _GREY, device.Status, _OFF)
					channelDisplay <- m
					m.id = 3
					m.msg = fmt.Sprintf("+ Manufacturer:    %s%s", _BLUE, device.Manufacturer)
					channelDisplay <- m
				default:
				}
			}
			if device.Manufacturer != _defaults && device.Model != _defaults && device.Build != _defaults && device.Revision != _defaults && device.ISO != _defaults && device.CFG != _defaults && device.HWVER != _defaults && device.IMEI != _defaults && device.Capabilities != _defaults && device.INI != _defaults {
				s := device.Manufacturer + device.Model + device.Build + device.Revision + device.ISO + device.CFG + device.HWVER + device.IMEI + device.Capabilities + device.INI
				device.Fingerprint = s2b32f(s)
				if strings.Contains(c.DeviceKnownList, device.Fingerprint) {
					m.msg = fmt.Sprintf("%s### DEVICE %s[OK]%s [%s]", _WHITE, _ALERT_G, _GREY, device.Fingerprint)
				} else {
					m.msg = fmt.Sprintf("%s### DEVICE %sUNLISTED [%s]", _WHITE, _RED, device.Fingerprint)
				}
				m.id = 2
				channelDisplay <- m
				channelFingerprint <- true
				device.Fingerprint_Done = true
			}
			continue
		}
		time.Sleep(eco_sleep)
	} // end parser loop
}

// open / probe / get device
func (c *Config) getLTE() lte {
	for {
		port, err := os.OpenFile(c.DevicePort, os.O_RDWR, 0o660)
		if err != nil {
			fmt.Printf("### DEVICE ERROR: lte device -> %s not ready! Waiting for device ... \n", c.DevicePort)
			time.Sleep(one * 5)
			continue
		}
		return lte{
			feed:     bufio.NewScanner(port),
			port:     port,
			portfile: c.DevicePort,
		}
	}
}
