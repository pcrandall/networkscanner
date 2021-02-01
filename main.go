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
	"text/tabwriter"
	"time"

	"github.com/go-ping/ping"
	"github.com/mostlygeek/arp"

	// "github.com/pcrandall/networkscanner/arp"
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
	mac   string
	bytes int
	time  time.Duration
}

type netTable map[string]string

type availAddr struct {
	table netTable
	mutex sync.Mutex
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

	flag.IntVar(&_timeout, "t", 1000, "Timeout in milliseconds -t 500")
	flag.IntVar(&_interval, "i", 200, "Ping interval -i 200")
	flag.IntVar(&count, "c", 2, "Ping count -c 2")

	flag.BoolVar(&debug, "d", false, "Debug -d=true")
	flag.BoolVar(&verbose, "v", true, "No output to console -v=false")
	flag.BoolVar(&write, "w", false, "Write output to availableIPS.txt in current directory -w=true")
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
			if mac != "" && mac != "00:00:00:00:00:00" { // remove from available only if valid mac address was returned
				if debug {
					fmt.Println(ip, "from arp.search() mac: ", mac)
				}
				pingedAddr[ip].mac = mac
				openAddr.delAddr(ip)
			}
		}(ip)
	}
	// wait for all go routines to finish
	wg.Wait()

	if verbose {
		// initialize tabwriter
		w := new(tabwriter.Writer)
		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 4, '\t', 0)
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "IP address", "Response time", "MAC address", "Bytes")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "--------------", "-------------", "-----------------", "-----")
		for _, val := range pingedAddr {
			if val.bytes > 0 { // console output and write to file
				fmt.Fprintf(w, "\n %s\t%s\t%s\t%d\t ", val.addr, val.time, val.mac, val.bytes)
			}
		}
		w.Flush()
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

func (a *availAddr) addAddr(ip string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.table[ip] = ip
}

func (a *availAddr) delAddr(ip string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, ok := a.table[ip]
	if ok {
		delete(a.table, ip)
	}
}
