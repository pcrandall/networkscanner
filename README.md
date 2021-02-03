# Network Scanner
Check for online hosts/ available addresses.

Use build tags to build for proper os

With the -w=true flag set scannedIPs directory will be created in current directory. Available addresses will be written to "scannedIPs/availableIPS.txt", and Reserved addresses will be written to 'scannedIPs/reservedIPS.txt'

Flag syntax
```
-flag
-flag=x
-flag x  // non-boolean flags only
```
Flag descriptions:
```
./networkscanner -h

```
Examples:

Windows: from cmd.exe
```
networkscanner.exe -ip 192.168.1.0/27 -c 4 -w=true -t 1500 -e 192.168.1.0,192.168.1.1,192.168.1.2,192.168.1.30,192.168.1.31
networkscanner.exe -ip 192.168.1.0/27 -c 2 -i 250 -w=true -t 1000 -e 192.168.1.0,192.168.1.1,192.168.1.2,192.168.1.30,192.168.1.31
networkscanner.exe -ip=192.168.1.0/27 -c 2 -i 250 -w=true -t 1000
networkscanner.exe -ip=192.168.1.0/27 -c 2 -i 250 -w=true -t 2000
networkscanner.exe -ip=192.168.1.0/27 -c 4 -i 250 -w=false -t 2000
networkscanner.exe -ip=192.168.1.0/27 -w=true -v=false -t 3000
networkscanner.exe -ip=192.168.1.0/27
networkscanner.exe -ip=192.168.1.0/27 -w=true
networkscanner.exe -ip=192.168.1.0/27 -t 1000
```

Linux/Osx:
```
./networkscanner -ip=192.168.1.0/27
./networkscanner -ip 192.168.1.0/27 -c 4 -w=true -t 1500 -e 192.168.1.0,192.168.1.1,192.168.1.2,192.168.1.30,192.168.1.31
./networkscanner -ip 192.168.1.0/27 -c 2 -i 250 -w=true -t 1000 -e 192.168.1.0,192.168.1.1,192.168.1.2,192.168.1.30,192.168.1.31
./networkscanner -ip=192.168.1.0/27 -c 2 -i 250 -w=true -t 1000
./networkscanner -ip=192.168.1.0/27 -c 4 -i 250 -w=false -t 2000
./networkscanner -ip=192.168.1.0/27 -w=true -v=false -t 3000
./networkscanner -ip=192.168.1.0/27 -w=true
./networkscanner -ip=192.168.1.0/27 -t 1000
```
