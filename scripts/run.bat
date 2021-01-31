@ECHO off

networkscanner.exe -ip=10.136.18.0/27 -c=2 -i=250 -w=true -t=1000
:: networkscanner.exe -ip=10.136.18.0/27 -c=2 -i=250 -w=true -t=2000
:: networkscanner.exe -ip=10.136.18.0/27 -c=4 -i=250 -w=false -t=2000
:: networkscanner.exe -ip=10.136.18.0/27 -w=true -s=true -t=3000
:: networkscanner.exe -ip=10.136.18.0/27
:: networkscanner.exe -ip=10.136.18.0/27 -w=true
:: networkscanner.exe -ip=10.136.18.0/27 -t=1000

pause
