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
)

func main() {

	ip := flag.String("ip", "192.168.1.1/24", "-ip=192.168.1.1/24 CIDRblock to scan.\nOnly CIDR format supported at this time.")
	t := flag.Int("t", 500, "-t=500 timeout in milliseconds")
	c := flag.Int("c", 1, "-c=1 number of times to ping")
	silent := flag.Bool("s", false, "-s=true no output to console")
	write := flag.Bool("w", false, "-w=true write output to ips.txt")

	flag.Parse()

	count := c
	timeout := time.Duration(*t) * time.Millisecond
	fmt.Println(count)

	peerConn := make(map[string]*Peer)
	unavailable := make(map[string]*online)

	_, IPrange := network.CalculateCIDR(*ip)

	for _, ip := range IPrange {

		wg.Add(1)

		ctx := context.Background()

		peerConn[ip] = &Peer{
			Address:  ip,
			Ctx:      &ctx,
			Count:    4,
			PingTime: timeout,
		}

		unavailable[ip] = &online{}

		go Ping(ip, peerConn[ip].Ctx, peerConn[ip], unavailable[ip])
	}

	wg.Wait()

	var output *os.File
	var err error

	if *write {
		// Check if file exists
		if _, err = os.Stat("ips.txt"); os.IsNotExist(err) {
			output, err = os.Create("ips.txt")
		} else { // If file exists, open it
			os.Remove("ips.txt") // remove file and create new one
			output, err = os.Create("ips.txt")
			output, err = os.Open("ips.txt")
		}
	}

	for _, val := range unavailable {
		if val.bytes > 0 && !*silent && *write { // console output and write to file
			fmt.Println(val.addr, "\tbytes rec:", val.bytes, "\ttime:", val.time)
			// io.WriteString(output, val.addr+"\n")
		} else if val.bytes > 0 && *silent && *write { // write to file
			// io.WriteString(output, val.addr+"\n")
		} else if val.bytes > 0 && !*silent && !*write { // conole output
			fmt.Println(val.addr, "\tbytes rec:", val.bytes, "\ttime:", val.time)
		}
	}

	if *write { // console output and write to file
		for _, val := range available {
			io.WriteString(output, val)
		}
	}

	fmt.Println(available)

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
		unavailable.addr = _ip
		unavailable.bytes = pkt.Nbytes
		unavailable.time = pkt.Rtt
		// fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
		// 	pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		// if stats.PacketLoss < 100 {
		// 	fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		// 	fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
		// 		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		// }

		if stats.PacketLoss > 99.9 {
			// fmt.Printf("\n%s Available ---%d packets transmitted, %d packets received, %v%% packet loss\n", stats.Addr, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			s := fmt.Sprintf("%0f", stats.PacketLoss)
			r := fmt.Sprintf("%x", stats.PacketsSent)
			t := fmt.Sprintf("%x", stats.PacketsRecv)
			available = append(available, string(stats.Addr)+" Available --- packets transmitted: "+r+" packets received: "+t+" packet loss "+s+"\n")
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		}
		// fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		// 	stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		return
	}

	err = pinger.Run() // Blocks until finished.

	if err != nil {
		fmt.Printf("Failed to ping target host: %s", err)
	}
}
