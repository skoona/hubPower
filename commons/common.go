/**
 * commons
 * is the collector of common utilities used
 */

package commons

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	AppIcon         = "apcupsd"
	PreferencesIcon = "preferences"
	TrustedIpKey    = "TRUSTED_IP"
	HubIpAddressKey = "HUBITAT_IP"
	DebugKey        = "SKN_DEBUG"
)

// ShutdownSignals alternate panic() implementation, causes an orderly shutdown
var ShutdownSignals chan os.Signal
var DebugLoggingEnabled = ("true" == os.Getenv(DebugKey)) // "true" / "false"
var logs = log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)

func DebugLog(args ...any) {
	if DebugLoggingEnabled {
		_ = logs.Output(2, fmt.Sprint(args...))
	}
}

// Keys returns the keys of the map m.
// The keys will be an indeterminate order.
// alternate reflect based: reflect.ValueOf(m).MapKeys()
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// ChangeTimeFormat converts APC timestamp to something more human readable
// time.RFC1123, time.RFC3339 are good choices
// returns local time version of value
func ChangeTimeFormat(timeString string, format string) string {
	if format == "" {
		format = time.RFC1123
	}
	if timeString == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(timeString))
	if err != nil {
		DebugLog("ApcService::ChangeTimeFormat() Time Parse Error, src: ", timeString, ", err: ", err.Error())
	}
	return t.Format(format)
}

// RemoveIndexFromSlice remove the given index from any type of slice
func RemoveIndexFromSlice[K comparable](index int, slice []K) []K {
	var idx int

	if len(slice) == 0 {
		return slice
	}

	if index > len(slice) {
		idx = len(slice) - 1
	} else if index < 0 {
		idx = 0
	} else {
		idx = index
	}
	return append(slice[:idx], slice[idx+1:]...)
}

// ShiftSlice drops index 0 and append newData to any type of slice
func ShiftSlice[K comparable](newData K, slice []K) []K {
	idx := 0
	if len(slice) == 0 {
		return append(slice, newData)
	}
	shorter := append(slice[:idx], slice[idx+1:]...)
	shorter = append(shorter, newData)
	return shorter
}

func DefaultIp() string {
	var currentIP string

	// use env if found
	currentIP = os.Getenv(TrustedIpKey)
	if currentIP != "" {
		return currentIP
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		DebugLog(err.Error())
	}

	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				DebugLog("Current IP address : ", ipnet.IP.String())
				currentIP = ipnet.IP.String()
				break // take the first one
			}
		}
	}
	return currentIP
}

func HubitatIP() string {
	ip := os.Getenv(HubIpAddressKey)
	if ip == "" {
		ip = "10.100.1.41"
	}
	return ip
}
