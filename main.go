package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/pcrandall/networkscanner/arp"
	"github.com/pcrandall/networkscanner/network"
)

type host struct {
	Address  string
	Count    int
	PingTime time.Duration
	Interval time.Duration
	Ctx      *context.Context
}

type onlineHosts struct {
	addr  string
	bytes int
	time  time.Duration
}

type netTable map[string]string

type availAddr struct {
	// sync.RWMutex
	table netTable
}

var (
	wg        sync.WaitGroup
	count     int
	_interval int
	_timeout  int

	debug   bool
	write   bool
	verbose bool

	err error
	ip  string
	e   string

	exclude []string
	output  *os.File

	openAddr = &availAddr{
		table: make(netTable),
	}
)

func init() {
	// Initalize flags
	flag.StringVar(&ip, "ip", "192.168.1.1/24", "Addresses to scan. Only CIDR format supported. -ip 192.168.1.1/24")
	flag.StringVar(&e, "e", "", "Addresses to exclude from available list; seperated by comma.  -e 192.168.1.0,192,168.1.255")

	flag.IntVar(&_timeout, "t", 500, "Timeout in milliseconds -t 500 ")
	flag.IntVar(&_interval, "i", 200, "Ping interval -i 200")
	flag.IntVar(&count, "c", 2, "Ping count -c 2")

	flag.BoolVar(&debug, "d", false, "Debug -d=true ")
	flag.BoolVar(&verbose, "v", false, "Output to console -v=true ")
	flag.BoolVar(&write, "w", false, "Write output to availableIPS.txt in current directory -w=false")

	flag.Parse()

	if e != "" {
		exclude = strings.Split(e, ",")
	}

	if write {
		output, err = os.Create("availableIPS.txt")
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	timeout := time.Duration(_timeout) * time.Millisecond
	interval := time.Duration(_interval) * time.Millisecond
	hostConn := make(map[string]*host)
	pingedAddr := make(map[string]*onlineHosts)

	_, IPrange := network.CalculateCIDR(ip)

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
			if mac != "" { // remove from available only if valid mac address was returned
				if verbose {
					fmt.Println(ip, " mac: ", mac)
				}
				openAddr.delAddr(ip)
			}
		}(ip)
	}
	wg.Wait() // wait for all go routines to finish

	for _, val := range pingedAddr {
		if val.bytes > 0 && verbose { // console output and write to file
			fmt.Println(val.addr, "is taken\tbytes rec:", val.bytes, "\ttime:", val.time)
		}
	}

	// delete excluded addresses
	for _, ip := range exclude {
		openAddr.delAddr(ip)
	}

	//write remaining available addresses to file.
	if write {
		for _, val := range openAddr.table {
			io.WriteString(output, val+"\n")
		}
	}
}

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

	pinger.OnRecv = func(pkt *ping.Packet) {
		// got reply remove from available addresses
		if pingedAddr.bytes == 0 {
			openAddr.delAddr(addr)
		}

		pingedAddr.addr = addr
		pingedAddr.bytes = pkt.Nbytes
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
		return
	}

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		fmt.Printf("Failed to ping target host: %s", err)
	}
}

func (a *availAddr) addAddr(ip string) {
	a.table[ip] = ip
}

func (a *availAddr) delAddr(ip string) {
	_, ok := a.table[ip]
	if ok {
		delete(a.table, ip)
	}
}
