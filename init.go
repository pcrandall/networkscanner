package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"sync"
	"time"
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

type onlineHostList []onlineHosts

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
