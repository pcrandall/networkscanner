package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-ping/ping"
)

func Ping(addr string, ctx *context.Context, host *host, pingedAddr *onlineHosts) {

	defer wg.Done()

	pinger, err := ping.NewPinger(addr)

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	if err != nil {
		panic(err)
	}

	pinger.Count = host.Count
	pinger.Timeout = host.PingTime
	pinger.Interval = host.Interval
	pinger.OnRecv = func(pkt *ping.Packet) {
		// got reply remove from available addresses
		if pingedAddr.bytes == 0 {
			openAddr.delAddr(addr)
		}

		pingedAddr.addr = addr
		pingedAddr.bytes += pkt.Nbytes
		pingedAddr.time = pkt.Rtt
		if debug {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
		}
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		if debug {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}
	}

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		fmt.Printf("Failed to ping target host: %s", err)
	}
}
