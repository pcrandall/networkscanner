package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/mostlygeek/arp"
	"github.com/pcrandall/networkscanner/network"
)

type Peer struct {
	Address  string
	Status   bool
	Count    int
	PingTime time.Duration
	Interval time.Duration
	Ctx      *context.Context
}

type online struct {
	addr  string
	mac   string
	bytes int
	time  time.Duration
	// pinger ping.Pinger
}

var (
	wg sync.WaitGroup
)

func main() {

	ip := flag.String("ip", "192.168.1.1/24", "-ip=<192.168.1.1/24> ip block to scan")

	flag.Parse()

	output, err := os.Create("ips.txt")
	if err != nil {
		panic(err)
	}

	peerConn := make(map[string]*Peer)
	unavailable := make(map[string]*online)

	_, IPrange := network.CalculateCIDR(*ip)

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

		unavailable[ip] = &online{}

		go Ping(ip, peerConn[ip].Ctx, peerConn[ip], unavailable[ip])
	}

	wg.Wait()

	for ip, _ := range arp.Table() {
		unavailable[ip].mac = arp.Search(ip)
	}

	for _, val := range unavailable {
		if val.bytes > 0 {
			if val.mac == "" {
				fmt.Println(val.addr, "\tmac: 00:00:00:00:00:00", "\tbytes rec:", val.bytes, "\ttime:", val.time)
				io.WriteString(output, val.addr+"\n")
			} else {
				fmt.Println(val.addr, "\tmac:", val.mac, "\tbytes rec:", val.bytes, "\ttime:", val.time)
				io.WriteString(output, val.addr+"\n")
			}
		}
	}
}

func Ping(_ip string, ctx *context.Context, peer *Peer, unavailable *online) {

	defer wg.Done()

	pinger, err := ping.NewPinger(_ip)

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
		unavailable.addr = _ip
		unavailable.bytes = pkt.Nbytes
		unavailable.time = pkt.Rtt
		// fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
		// 	pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
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
