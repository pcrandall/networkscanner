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
	"github.com/pcrandall/networkscanner/arp"
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
	bytes int
	time  time.Duration
}

var (
	wg        sync.WaitGroup
	available []string
	debug     *bool
	write     *bool
	silent    *bool
	err       error
)

func main() {
	// flags
	ip := flag.String("ip", "192.168.1.1/24", "default:-ip=192.168.1.1/24 CIDRblock to scan.\nOnly CIDR format supported at this time.")
	t := flag.Int("t", 500, "default:-t=500 Timeout in milliseconds)")
	c := flag.Int("c", 2, "default:-c=2 Number of times to ping")
	i := flag.Int("i", 200, "default:-i=200 Ping interval")
	// interval := time.Duration(*flag.Duration("i", 250, "default:-i=250 Ping interval")) * time.Millisecond
	debug = flag.Bool("d", false, "default:false Debug")
	silent = flag.Bool("s", false, "default:-s=false No output to console")
	write = flag.Bool("w", false, "default:-w=false Write output to ips.txt in current directory")
	flag.Parse()

	output, err := os.Create("ips.txt")
	if err != nil {
		panic(err)
	}

	count := *c
	timeout := time.Duration(*t) * time.Millisecond
	interval := time.Duration(*i) * time.Millisecond
	peerConn := make(map[string]*Peer)
	unavailable := make(map[string]*online)

	_, IPrange := network.CalculateCIDR(*ip)

	// arp lookup
	for ip, _ := range arp.Table() {
		go func(ip string) {
			mac := arp.Search(ip)
			if mac != "" { // remove from available only if valid mac address was returned
				if !*silent {
					fmt.Println(ip, " mac: ", mac)
				}
				removeAvailable(ip)
			}
		}(ip)
	}

	for _, ip := range IPrange {
		available = append(available, ip) // create slice to remove all unavailable addresses later
		wg.Add(1)
		ctx := context.Background()
		peerConn[ip] = &Peer{
			Address:  ip,
			Ctx:      &ctx,
			Count:    count,
			PingTime: timeout,
			Interval: interval,
		}
		unavailable[ip] = &online{}
		go Ping(ip, peerConn[ip].Ctx, peerConn[ip], unavailable[ip])
	}

	wg.Wait()

	for _, val := range unavailable {
		if val.bytes > 0 && !*silent { // console output and write to file
			fmt.Println(val.addr, "is taken\tbytes rec:", val.bytes, "\ttime:", val.time)
		}
	}

	//write remaining available addresses to file.
	for _, val := range available {
		io.WriteString(output, val+"\n")
	}

	if !*write {
		err = os.Remove("ips.txt")
		if err != nil {
			panic(err)
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
	pinger.Timeout = peer.PingTime

	pinger.OnRecv = func(pkt *ping.Packet) {
		// got reply remove from available addresses
		if unavailable.bytes == 0 {
			removeAvailable(_ip)
		}

		unavailable.addr = _ip
		unavailable.bytes = pkt.Nbytes
		unavailable.time = pkt.Rtt
		if *debug {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
		}
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		if *debug {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}
		return
	}

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		fmt.Printf("Failed to ping target host: %s", err)
	}
}

func removeAvailable(ip string) {
	if *debug {
		fmt.Println("in remove available", available)
	}
	for idx, val := range available {
		if val == ip {
			tmp := available[:idx]
			available = available[idx+1 : len(available)-1]
			tmp = append(tmp, available...)
			available = tmp
		}
	}
}
