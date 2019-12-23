## Usage:
### List all ipaddresses in specified vlan
```
    ./goipam -vlan=VLAN100_10.141.16.0m24
```
### Reserve provided ip address
```
    ./goipam -ip=10.141.16.13 -reserve -comment "ala ma kota"
```
### Release provided ip address
```
    ./goipam -ip=10.141.16.13 -release
```
