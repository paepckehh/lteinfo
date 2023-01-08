package lteinfo

func getCommand(in string) *cmd {
	switch in {
	case "HUAWEI_E3372":
		return &cmd{
			update:       [9]string{"w", "x", "x", "x", "x", "x", "x", "x", "x"},
			device_clean: [9]string{"^CURC=0", "x", "x", "x", "x", "x", "x", "x", "x"},
			device_init:  [9]string{"I", "W", "W", "^VERSION?", "W", "W", "x", "x", "x"},
			device_setup: [9]string{"^CURC=1", "^CREG=1", "x", "x", "x", "x", "x", "x", "x"},
			sim:          [9]string{"^SPN=1", "w", "^ICCID?", "w", "+CPIN?", "W", "^CARDLOCK?", "w", "+CSMS?"},
			slow:         [9]string{"^HFREQINFO?", "+CREG?", "^DHCP?", "w", "^NWTIME?", "w", "^CHIPTEMP?", "w", "+COPS?"},
			real_slow:    [9]string{"w", "^lteCAT?", "W", "^NDISSTATQRY?", "W", "+COPS=?", "W", "W"},
			smsInit:      [9]string{"+CMGF=1", "x", "x", "x", "x", "x", "x", "x", "x"},
			smsGet:       [9]string{"+CMGR=", "x", "x", "x", "x", "x", "x", "x", "x"}, // add slot number (0-30)
		}
	default:
		return &cmd{
			update:       [9]string{"x", "x", "x", "x", "x", "x", "x", "x", "x"},
			device_clean: [9]string{"x", "x", "x", "x", "x", "x", "x", "x", "x"},
			device_init:  [9]string{"I", "W", "W", "I", "W", "W", "^VERSION?", "W", "^VERSION?"},
			device_setup: [9]string{"x", "x", "x", "x", "x", "x", "x", "x", "x"},
			sim:          [9]string{"x", "x", "x", "w", "x", "x", "x", "x", "x"},
			slow:         [9]string{"x", "x", "x", "x", "x", "x", "x", "x", "x"},
			real_slow:    [9]string{"x", "x", "x", "x", "x", "x", "x", "x", "x"},
			smsInit:      [9]string{"+CMGF=1", "x", "x", "x", "x", "x", "x", "x", "x"},
			smsGet:       [9]string{"+CMGR=", "x", "x", "x", "x", "x", "x", "x", "x"}, // add slot number (0-30)
		}
	}
}
