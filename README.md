# GOIPAM
This program has been developed to make easy communication with Orion IPAM.

## Build:

```go build .```

It is required to set environment variables ORION_IP, ORION_USER, ORION_PASSWORD. For that please update env_vars.sh file and run below command.

```
source env_vars.sh
```

If Ping function does not work, please run this command in shell.

```
	sudo sysctl -w net.ipv4.ping_group_range="0   2147483647
```

To see more logs, please change logging level by setting following env var.

```
export TF_LOG="DEBUG"
```

## Usage:

### List all ipaddresses in specific vlan
```
    ./goipam -vlan=VLAN100_10.141.16.0m24 -list
    ./goipam -vlan=VLAN_141810.14.18.0m24 -list
```
### Reserve provided ip address
```
    ./goipam -ip=10.141.16.13 -reserve -comment "ala ma kota"
```
### Release provided ip address
```
    ./goipam -ip=10.141.16.13 -release
```
