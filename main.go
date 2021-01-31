package main

import (
	"context"
	"flag"
	"sync"

	"fmt"
	"runtime"
	"time"

	"github.com/go-ping/ping"

	"github.com/workit/network-scanner/network"
)

type Peer struct {
	Address  string
	Status   bool
	Count    int
	PingTime time.Duration
	Interval time.Duration
	Ctx      *context.Context
}

var wg sync.WaitGroup

func main() {

	cidr := flag.String("cidr", "10.136.18.0/27", "CIDR block to scan")

	flag.Parse()

	peerConn := make(map[string]*Peer)

	_, IPrange := network.CalculateCIDR(*cidr)

	for _, ip := range IPrange {

		wg.Add(1)

		ctx := context.Background()

		peerConn[ip] = &Peer{
			Address:  ip,
			Ctx:      &ctx,
			Count:    1,
			PingTime: 1000 * time.Millisecond,
			Interval: 1000 * time.Millisecond,
		}

		go Ping(ip, peerConn[ip].Ctx, peerConn[ip])
	}

	wg.Wait()
}

func Ping(address string, ctx *context.Context, peer *Peer) {

	defer wg.Done()

	pinger, err := ping.NewPinger(address)
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}
	if err != nil {
		panic(err)
	}

	pinger.Count = peer.Count

	pinger.Interval = peer.Interval

	pinger.Timeout = peer.PingTime

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%s Online, bytes received: %d: \n", pkt.IPAddr, pkt.Nbytes)
		// fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
		// 	pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
		// if pkt.Nbytes > 0 {
		// 	peer.Status = true
		// 	peer.PingTime = pkt.Rtt
		// }
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		// fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		// fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
		// 	stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		// fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		// 	stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		return
	}

	err = pinger.Run() // Blocks until finished.

	if err != nil {
		fmt.Printf("Failed to ping target host: %s", err)
	}

}
