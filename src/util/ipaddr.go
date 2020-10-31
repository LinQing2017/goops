package util

import (
	"net"
	"strconv"
	"strings"
)

// 判断IP 是否属于cidr网段
func IpaddrBelong(ip, cidr string) bool {

	ipAddr := strings.Split(ip, `.`)
	if len(ipAddr) < 4 {
		return false
	}
	cidrArr := strings.Split(cidr, `/`)
	if len(cidrArr) < 2 {
		return false
	}

	var tmp = make([]string, 0)

	maskBits, _ := strconv.Atoi(cidrArr[1])
	cidrMask := net.CIDRMask(maskBits, 32) // 获取掩码

	for key, value := range cidrMask {
		iint := int(value)

		iint2, _ := strconv.Atoi(ipAddr[key])

		tmp = append(tmp, strconv.Itoa(iint&iint2))
	}
	return strings.Join(tmp, `.`) == cidrArr[0]
}
