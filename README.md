# CEC Power Kodi Control
This basically sits on front of the kodi rpc and forwards the traffic, but if it detects the method being "shutdown" it actually shut down my TV.  
And to turn it on, it uses WOL (Wake on Lan).

It is used with Kore (Kodi remote control app) which have WOL and the shutdown option.

This is for my internal use, some bugs may occur to you.

# Build
```bash
git clone https://github.com/thallesp/cec-power-kodi-control
cd cec-power-kodi-control
go build -o cec

# Cross compile for arm 32 bit?
# env GOOS=linux GOARCH=arm go build -o cec
```

# Systemd
```
[Install]
WantedBy=multi-user.target

[Unit]
Description=CEC Power Kodi Control

[Service]
RemainAfterExit=yes
ExecStart=/storage/cec
StartLimitInterval=0
```