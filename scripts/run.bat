@ECHO off

networkscanner.exe -ip 10.136.18.0/27 -c 4 -w=true -t 1500 -e 10.136.18.0,10.136.18.1,10.136.18.2,10.136.18.30,10.136.18.31
:: networkscanner.exe -ip 10.136.18.0/27 -c 2 -i 250 -w=true -t 1000 -e 10.136.18.0,10.136.18.1,10.136.18.2,10.136.18.30,10.136.18.31
:: networkscanner.exe -ip=10.136.18.0/27 -c=2 -i=250 -w=true -t=1000
:: networkscanner.exe -ip=10.136.18.0/27 -c=2 -i=250 -w=true -t=2000
:: networkscanner.exe -ip=10.136.18.0/27 -c=4 -i=250 -w=false -t=2000
:: networkscanner.exe -ip=10.136.18.0/27 -w=true -v=false -t=3000
:: networkscanner.exe -ip=10.136.18.0/27
:: networkscanner.exe -ip=10.136.18.0/27 -w=true
:: networkscanner.exe -ip=10.136.18.0/27 -t=1000

pause
