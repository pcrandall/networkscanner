package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/mostlygeek/arp"
	"github.com/pcrandall/networkscanner/network"
)

func main() {
	// initialize vars
	timeout := time.Duration(_timeout) * time.Millisecond
	interval := time.Duration(_interval) * time.Millisecond
	hostConn := make(map[string]*host)
	pingedAddr := make(map[string]*onlineHosts)
	_, IPrange := network.CalculateCIDR(ip)

	// ping addresses
	for _, ip := range IPrange {
		openAddr.addAddr(ip) // build list of possible addresses
		wg.Add(1)
		ctx := context.Background()
		hostConn[ip] = &host{
			Address:  ip,
			Ctx:      &ctx,
			Count:    count,
			PingTime: timeout,
			Interval: interval,
		}
		pingedAddr[ip] = &onlineHosts{}
		go Ping(ip, hostConn[ip].Ctx, hostConn[ip], pingedAddr[ip])
	}

	// arp lookup
	for ip, _ := range arp.Table() {
		go func(ip string) {
			mac := arp.Search(ip)
			if mac != "" && mac != "00:00:00:00:00:00" { // remove from available only if valid mac address was returned
				if debug {
					fmt.Println(ip, "from arp.search() mac: ", mac)
				}
				pingedAddr[ip].mac = mac
				openAddr.delAddr(ip)
			}
		}(ip)
	}

	// wait for all gofuncs to finish
	wg.Wait()

	// print & write online host table
	if verbose {
		printOnlineHosts(pingedAddr)
	}
	// delete excluded addresses
	for _, ip := range exclude {
		openAddr.delAddr(ip)
	}
	//write remaining available addresses to file.
	if write {
		for _, val := range openAddr.table {
			io.WriteString(availableHostsFile, val+"\n")
		}
	}
}
