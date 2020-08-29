package utils

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"math"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const SERVICE_HOST = "SERVICE_HOST"

func GetExecFilename() string {
	_, filename, _, _ := runtime.Caller(1)
	cwdPath, _ := os.Getwd()
	return cwdPath + filename
}

func GetHost() string {
	ret := getIpFromEnv()
	if len(ret) > 0 {
		return ret
	}

	ret = getIpFromSpecialName("eth0", "em1")
	if len(ret) > 0 {
		return ret
	}

	ret = getFirstIp()
	if len(ret) > 0 {
		return ret
	}

	panic(fmt.Errorf("ip not found"))
	return ""
}

func getIpFromEnv() string {
	return os.Getenv(SERVICE_HOST)
}

func getIpFromSpecialName(name ...string) string {
	ifaces, e := net.Interfaces()
	if e != nil {
		panic(e)
	}

	match := func(n string) bool {
		for _, v := range name {
			if v == n {
				return true
			}
		}
		return false
	}

	for _, v := range ifaces {
		if match(v.Name) {
			return _getIpByFace(v)
		}
	}
	return ""
}

func _getIpByFace(iface net.Interface) string {
	if iface.Flags&net.FlagUp == 0 {
		return ""
	}

	if iface.Flags&net.FlagLoopback != 0 {
		return ""
	}

	// ignore docker and warden bridge
	if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
		return ""
	}

	addrs, e := iface.Addrs()
	if e != nil {
		return ""
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}

		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}
		return ip.String()
	}
	return ""
}

func getFirstIp() string {
	ips := make([]string, 0)
	ifaces, e := net.Interfaces()
	if e != nil {
		panic(e)
	}
	for _, iface := range ifaces {
		ipStr := _getIpByFace(iface)
		if len(ipStr) > 0 {
			ips = append(ips, ipStr)
		}
	}
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

func GetFreePort() (int, error) {
	ports, err := GetFreePorts(1)
	if err != nil {
		return 0, err
	}
	if len(ports) == 0 {
		return 0, fmt.Errorf("no useful port")
	}
	return ports[0], nil
}

func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

func ValueInject(start int, count int) string {
	var ret bytes.Buffer
	ret.WriteByte('(')
	for i := 1; i <= count; i++ {
		ret.WriteByte('$')
		ret.WriteString(Int64ToString(int64(start + i)))
		if i != count {
			ret.WriteByte(',')
		}
	}
	ret.WriteByte(')')
	return ret.String()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type Page struct {
	Start uint64
	Size  uint64
}

func SplitPage(pageSize, count uint64) []*Page {
	allPage := uint64(math.Ceil(float64(count) / float64(pageSize)))
	ret := make([]*Page, 0)
	if pageSize >= count {
		ret = append(ret, &Page{
			Start: 0,
			Size:  count,
		})
	} else {
		for i := uint64(0); i < allPage-1; i++ {
			ret = append(ret, &Page{
				Start: i * pageSize,
				Size:  pageSize,
			})
		}
		ret = append(ret, &Page{
			Start: (allPage - 1) * pageSize,
			Size:  count - (allPage-1)*pageSize,
		})
	}
	return ret
}

func UInt64ToString(n uint64) string {
	return strconv.FormatUint(uint64(n), 10)
}

func Decimal(value float64) float64 {
	return math.Round(value*1000000) / 1000000
}
func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func StringToUint64(s string) uint64 {
	ret, _ := strconv.ParseUint(s, 10, 64)
	return ret
}

func StringToInt64(s string) int64 {
	ret, _ := strconv.ParseInt(s, 10, 64)
	return ret
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func StringToFloat(s string) float64 {
	ret, _ := strconv.ParseFloat(s, 64)
	return ret
}

func LnglatValid(lng, lat float64) bool {
	return ValidLng(lng) && ValidLat(lat)
}

func ValidLng(lng float64) bool {
	return lng >= -180.0 && lng <= 180.0
}

func ValidLat(lat float64) bool {
	return lat <= 90.0 && lat >= -90.0
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func CondSql(first bool) string {
	if first {
		return " where"
	}
	return " and"
}

func ArrayToSqlIn(ids ...string) string {
	var buffer bytes.Buffer
	for _, v := range ids {
		buffer.WriteString("'")
		buffer.WriteString(v)
		buffer.WriteString("'")
		buffer.WriteString(",")
	}
	temp := buffer.String()
	return temp[:len(temp)-1]
}

func ReadFileContent(filename string) ([]byte, error) {
	obj, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(obj)
	_ = obj.Close()
	return buf, err
}

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	errc <- terminateError
}

func ExactStringArrayRequests(req string, seq string) ([]string, error) {
	args := make([]string, 0)
	_ = jsoniter.ConfigFastest.Unmarshal([]byte(req), &args)
	if len(args) > 0 {
		return args, nil
	}
	arr := strings.Split(req, seq)
	for i := range arr {
		args = append(args, arr[i])
	}
	if len(args) == 0 {
		args = append(args, req)
	}
	return args, nil
}
