package lteinfo

import (
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"strconv"
	"strings"
)

// hruIEC converts value to hru IEC 60027 units
func hruIEC(i uint64, u string) string {
	return hru(i, 1024, u)
}

// hruSI converts value to hru SI units
func hruSI(i uint64, u string) string {
	return hru(i, 1000, u)
}

// hru [human readable units] backend
func hru(i, unit uint64, u string) string {
	if i < unit {
		return fmt.Sprintf("%d %s", i, u)
	}
	div, exp := unit, 0
	for n := i / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	switch u {
	case "":
		return fmt.Sprintf("%.3f %c", float64(i)/float64(div), "kMGTPE"[exp])
	case "bit":
		return fmt.Sprintf("%.0f %c%s", float64(i)/float64(div), "kMGTPE"[exp], u)
	case "bytes", "bytes/sec":
		return fmt.Sprintf("%.1f %c%s", float64(i)/float64(div), "kMGTPE"[exp], u)
	}
	return fmt.Sprintf("%.3f %c%s", float64(i)/float64(div), "kMGTPE"[exp], u)
}

// hex2IP4 converts hex representation to default ip4 notation
func hex2IP4(hex string) string {
	var s [ipv4_OCTETS]uint64
	for i := 0; i < ipv4_OCTETS; i++ {
		s[i], _ = strconv.ParseUint(hex[i*2:i*2+2], 16, 0)
	}
	return fmt.Sprintf("%d.%d.%d.%d", s[3], s[2], s[1], s[0])
}

// s2b converts a generic string to string-bool
func s2b(in string) string {
	switch in {
	case "0":
		return _ALERT + "false" + _OFF
	case "1":
		return _GREEN + "true" + _OFF
	}
	panic("bool can only 0 = false  / 1 = true")
}

// s2s converts a generic string to status
func s2s(in string) string {
	if in == "0" {
		return "n/a"
	}
	return in
}

// id2netlock converts network enum to msg
func id2netlock(in string) string {
	switch in {
	case "1":
		return _alert + _ALERT + "NETLOCK ACTIVE]" + _OFF
	case "2":
		return _GREEN + "[NETLOCK DEACTIVATED]" + _OFF
	case "3":
		return _alert + _ALERT + "[PERSISTENT NETLOCK ACTIVE]" + _OFF
	}
	return _ALERT + "unknow / unclear status code" + _OFF
}

// id2airInterface converts id -> network type
func id2airInterface(in string) string {
	switch in {
	case "0":
		return "GSM"
	case "1":
		return "GSM Compact"
	case "2":
		return "UTRAN UMTS"
	case "3":
		return "GSM GPRS"
	case "4":
		return "UTRAN HSDPA"
	case "5":
		return "UTRAN HSUPA"
	case "6":
		return "UTRAN HSDAP & HSUPA"
	case "7":
		return "E-UTRAN LTE"
	}
	return "ERROR - unknow air interface"
}

// id2provider converts a provider status resolver
func id2provider(in string) (string, bool) {
	switch in {
	case "0":
		return _alert + _ALERT + "[MT not registered] [Not Searching!]" + _OFF, false
	case "1":
		return _GREEN + "[MT registered] [complete] [home network]" + _OFF, false
	case "2":
		return _ALERT + "[MT not registered] [Searching ...]" + _OFF, false
	case "3":
		return _ALERT + "[MT access denied]" + _OFF, false
	case "5":
		return _GREEN + "[MT registered] [complete] [roaming]" + _OFF, true
	}
	return _ALERT + "[MT status unknow]" + _OFF, false
}

// temp2color converts a temperatur to color / _alert indicator
func temp2color(in, alert uint64) string {
	switch {
	case in == alert:
		return _CYAN
	case in > alert:
		return _ALERT
	}
	return _GREEN
}

// int2indicator converts to a bar graph ui
func int2indicator(in, alert uint64, indicator string) string {
	if in < alert {
		return _ALERT + strings.Repeat(indicator, int(in/2)) + _OFF
	}
	return _ALERT_G + strings.Repeat(indicator, int(in/2)) + _OFF
}

// s2b32 converts a string to b32 fingerprint display rep
func s2b32f(in string) string {
	h := s2b32(sha512.Sum512_224([]byte(in)))
	return fmt.Sprintf("%s-%s-%s-%s", h[:4], h[4:6], h[6:8], h[8:12])
}

// s2b32 converts a slice[28] to base32
func s2b32(in [28]byte) string { return base32.StdEncoding.EncodeToString(in[:]) }
