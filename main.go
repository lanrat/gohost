package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/digineo/ripego"
)

var (
	addr_port = "0.0.0.0:8181"
)

func initSettings() {
	env_port, set := os.LookupEnv("LISTEN_ADDR")
	if set {
		addr_port = env_port
	}
}

func AllHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIpAddress(r)
	name := getHostname(ip)
	log.Printf("IP: %s\n", ip)
	fmt.Fprintf(w, "IP: %s\n", ip)
	fmt.Fprintf(w, "DNS: %s\n", name)

	wr, err := ripego.IPLookup(ip)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "WHO: %s\n", wr.Organization)

}

func IpHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIpAddress(r)
	log.Printf("IP: %s\n", ip)
	fmt.Fprintf(w, "%s\n", ip)
}

func HostHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIpAddress(r)
	name := getHostname(ip)
	log.Printf("IP: %s\n", ip)
	fmt.Fprintf(w, "%s\n", name)
}

func getHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	return strings.Join(names, " ")
}

// https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
func getIpAddress_merge(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() { //|| isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

func getIpAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIp := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIp == "" && hdrForwardedFor == "" {
		hdrRealIp, _, _ := net.SplitHostPort(r.RemoteAddr)
		return hdrRealIp
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIp
}

func init() {
	initSettings()
	fmt.Printf("Started server at %v.\n", addr_port)
	http.HandleFunc("/", AllHandler)
	http.HandleFunc("/ip", IpHandler)
	http.HandleFunc("/host", HostHandler)
	http.ListenAndServe(addr_port, nil)
}

func main() {}
