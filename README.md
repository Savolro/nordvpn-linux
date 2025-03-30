# NordVPN-Linux for the best OS
![icon](./assets/icon.svg)
---
This branch is created to provide possibility to use the best NordVPN client on
the greatest OS that ever existed.

## THIS SECTION IS CAPS-LOCKED
PLEASE DO NOT USE THIS IN PRODUCTION. THIS IS NOT TESTED AND SOME SECURITY
FEATURES WERE PURPOSEFULLY REMOVED IN ORDER TO ACHIEVE POC STATE. SUCH AS:
- SERVICE LISTENS ON `tcp:localhost:6969` AND HAS NO AUTHENTICATION WHATSOEVER.
  THIS MEANS THAT ANY PROCESS IN YOUR DEVICE COULD MANAGE YOUR VPN SETUP OR
  MODIFY YOUR MESHNET.
- THIS WAS DEVELOPED WITHIN A FEW HOURS ON MY MACHINE TO RECEIVE A SUCCESSFUL
  VPN CONNECTION AND HAD NO FURTHER TESTING. IT COULD CRASH AT ANY MINUTE
  WHILIST LEAKING YOUR PRECIOUS DATA.
- ANY EXTRA SECURITY FEATURES (SUCH AS NETWORK MONITORING, FIREWALL
  CONFIGURATION, KILLSWITCH, ETC.) WERE CUT OUT. IF TUNNEL IS DEAD - YOU ARE
  LEAKING.

# Current state of this
## What works on my machine
- VPN connections (quic connect, connect to country, group, etc.)
- Basic meshnet connection. Incoming connections are enabled by default. Traffic
  routing to this device will not work due to missing firewall implementations.
  Traffic to another device may work but it was not tested.
- `nrodvpn-linux-gui` worked with windows build type added and gRPC address
  changed to `tcp:localhost:6969`.

## What will definitely not work as of now
- Any firewall configuration (killswitch, firewall settings)
- DNS nameservers configuration (even for regular connections, TP Lite, custom
  DNS)
- NordVPN user service, nordtray and push notifications
- Meshnet File Sharing
- OpenVPN
- NordWhisper
- Full OAuth flow due to not handling mime types (use `nrodvpn login --callback '<link_from_browser>'`
  to work this around)
- Probably something else :) Expect some crashes as well.

# Prerequisites
- Go <= 1.23. Note: It has to be smaller as uniffi-bindgen-go does not work with
  latest Go as of now.
- Everything mentioned in: https://github.com/NordSecurity/libtelio;

# Building
Some extra steps may be requried but in general it should be something like this
```shell
# Build libtelio. See https://github.com/NordSecurity/libtelio regarding build dependencies
mkdir -p build
cd build
git clone --branch v5.1.9 https://github.com/NordSecurity/libtelio 
$env:BYPASS_LLT_SECRETS=1
cd libtelio
cargo build
cd ..\..

# Copy telio.dll to a correct directory
mkdir -p .\bin\deps\lib\libtelio\current\amd64
cp .\build\libtelio\target\debug\telio.dll .\bin\deps\lib\libtelio\current\amd64\telio.dll

# Build go binaries
go build -o nordvpnd.exe -tags telio .\cmd\daemon
go build -o nordvpn.exe .\cmd\cli
```

# Running
Starting the nordvpnd:
1. Run `Windows PowerShell` as an Administrator in this directory
1. Execute the following script in this terminal:
```script
$env:Path += ';.\bin\deps\lib\libtelio\current\amd64'
.\nordvpnd.exe
```
Note: Sometimes the process freezes until focused again or receiving keyboard
input. Probably some tweaks could be done for this as well.

Running the cli:
1. Simply run `.\nordvpn.exe --help` or any other commands.

Logging in:
1. Execute `.\nordvpn.exe login`
1. Open given URL in your browser and complete the login flow.
1. Copy the link under `Continue` button in browser.
1. Execute `.\nordvpn.exe login --callback '<link_from_browser>'