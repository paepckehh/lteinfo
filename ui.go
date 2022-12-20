package lteinfo

const (
	// ansi terminal color definitions
	_OFF     = "\033[0m"
	_RED     = "\033[2;31m"
	_GREEN   = "\033[2;32m"
	_YELLOW  = "\033[2;33m"
	_BLUE    = "\033[2;34m"
	_MAGENTA = "\033[2;35m"
	_CYAN    = "\033[2;36m"
	_WHITE   = "\033[2;37m"
	_GREY    = "\033[2;90m"
	_ALERT   = "\033[1;31m" // red alert
	_ALERT_G = "\033[1;32m" // grean alert

	// some generic terminal ui _defaults
	_clean_newline  = "\n" + _OFF
	_alert          = _ALERT + "[-= ***! ALERT !*** =-]" + _OFF
	_ok             = _ALERT_G + "[OK]" + _OFF
	_progress       = _OFF + "[" + _GREY + "-= information colletion in _progress =-" + _OFF + "]"
	_defaults       = _OFF + "[" + _GREY + "-= information not (yet) emitted from device =-" + _OFF + "]"
	_defaults_short = _GREY + "n/a" + _OFF
	_section_line_  = _WHITE + "###########################################################################################################"
)
