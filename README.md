# Overview

[paepche.de/lteinfo](https://paepcke.de/lteinfo)

- Show and decode information from your LTE Modem.
- Hardware Revision, Firmware Version, Operational Status (core,env temp)
- Signal, Provider, Visible Network Infra capabilities around, ... 
- Sim, IMEI, MVNO, authentication mode, Frequencies, ....
- Focus on small embedded systems (debugging) on restricted resources.
- Focus onpower saving parser (NO CLEAN idomatic go code for hot loop, no clean full state maschine, quick hack)
- 100 % pure go, stdlib only, no external dependencies, use as app or api (see api.go)

# Supported Devices 

- Huawei E3372h / E3372v153 / E5572C
- PRs welcome 

## Install 
```
go install paepcke.de/lteinfo/cmd/lteinfo@latest
```

# Showtime 

```Shell
lteinfo /dev/lte0
```
