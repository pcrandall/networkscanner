package network

import (
	"log"
	"net"
)

func CalculateCIDR(cidr string) (int64, []string) {
	IP, IPnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Println("Error in parsing CIDR")
	}

	newIP := IP.Mask(IPnet.Mask)

	begin, _ := IPnet.Mask.Size()
	count := int64(0)
	// IPRange := make([]net.IP, 1<<(32-uint(begin)))
	var IPRange []string

	for i := 0; i < 1<<(32-uint(begin)); i++ {
		count++
		IPRange = append(IPRange, newIP.String())
		if newIP[3] < 0xff {
			newIP[3] = newIP[3] + 1
			continue
		}
		if newIP[2] < 0xff {
			newIP[2] = newIP[2] + 1
			continue
		}
		if newIP[1] < 0xff {
			newIP[1] = newIP[1] + 1
			continue
		}
		if newIP[0] < 0xff {
			newIP[0] = newIP[0] + 1
			continue
		}

	}
	return count, IPRange
}
