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

func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIpAddress(r)
	log.Printf("header request from : %s\n", ip)
	for k, v:= range r.Header {
		fmt.Fprintf(w, "%s:%v\n", k, v)
	}
}

func getHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	return strings.Join(names, " ")
}

func getIpAddress(r *http.Request) string {
	for _, h := range []string{"Cf-Connecting-Ip", "X-Forwarded-For", "X-Real-Ip"} {
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
	realIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	return realIP
}

func init() {
	initSettings()
	fmt.Printf("Started server at %v.\n", addr_port)
	http.HandleFunc("/", AllHandler)
	http.HandleFunc("/ip", IpHandler)
	http.HandleFunc("/host", HostHandler)
	http.HandleFunc("/headers", HeaderHandler)
	http.ListenAndServe(addr_port, nil)
}

func main() {}
