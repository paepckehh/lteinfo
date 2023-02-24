# OVERVIEW
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/lteinfo.svg)](https://pkg.go.dev/paepcke.de/lteinfo) [![Go Report Card](https://goreportcard.com/badge/paepcke.de/lteinfo)](https://goreportcard.com/report/paepcke.de/lteinfo) [![Go Build](https://github.com/paepckehh/lteinfo/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/lteinfo/actions/workflows/golang.yml)

[paepcke.de/lteinfo](https://paepcke.de/lteinfo/)

- Show and decode information from your LTE Modem.
- Hardware Revision, Firmware Version, Operational Status (core,env temp)
- All Provider and Visible Network Infrastructure (capabilities eg, cell-broadcast,...)
- SIM, IMEI, MVNO, authentication mode, used Frequencies/Transponder ....
- Focus on small embedded systems (debugging) on restricted resources.
- Focus onpower saving parser (NO CLEAN idomatic go code for hot loop, no clean full state maschine, quick hack)
- 100 % pure go, stdlib only, no external dependencies, use as app or api (see api.go)

# Supported Devices 

- Huawei E3372h / E3372v153 / E5572C
- PRs welcome 

# INSTALL
```
go install paepcke.de/lteinfo/cmd/lteinfo@latest
```

### DOWNLOAD (prebuild)

[github.com/paepckehh/lteinfo/releases](https://github.com/paepckehh/lteinfo/releases)

# SHOWTIME 

```Shell
lteinfo /dev/lte0
```
# DOCS

[pkg.go.dev/paepcke.de/lteinfo](https://pkg.go.dev/paepcke.de/lteinfo)

# CONTRIBUTION

Yes, Please! PRs Welcome! 
