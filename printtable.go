package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

func (p onlineHostList) Len() int      { return len(p) }
func (p onlineHostList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p onlineHostList) Less(i, j int) bool {
	// 192.168.1.1  192.168.1.8
	// 19216811  19216818
	// return 19216811 < 19216818
	if len(p[i].addr) == len(p[j].addr) {
		ii, _ := strconv.Atoi(strings.ReplaceAll(p[i].addr, ".", ""))
		jj, _ := strconv.Atoi(strings.ReplaceAll(p[j].addr, ".", ""))
		return ii < jj
	}
	return len(p[i].addr) < len(p[j].addr)
}

func printOnlineHosts(pingedAddr map[string]*onlineHosts) {
	// make sortable list
	pingedAddrList := make(onlineHostList, len(pingedAddr))
	i := 0
	for _, v := range pingedAddr {
		pingedAddrList[i] = *v
		i++
	}
	sort.Sort(pingedAddrList)

	// initialize tabwriter
	w := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 4, '\t', 0)

	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "IP address", "Response time", "MAC address", "Bytes")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "--------------", "-------------", "-----------------", "-----")
	for _, val := range pingedAddrList {
		if val.bytes > 0 { // console output and write to file
			fmt.Fprintf(w, "\n %s\t%s\t%s\t%d\t ", val.addr, val.time, val.mac, val.bytes)
		}
	}
	w.Flush()
}
