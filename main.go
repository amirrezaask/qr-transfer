package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mdp/qrterminal"
)

func replaceChar(fileName string) string {
	s := strings.Replace(fileName, " ", "%20", -1)
	//s = strings.Replace(fileName, ".", "%2E", -1)
	return s
}

func getIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("could not get interface addrs :%v", err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("could not find ip")
}

func main() {
	//get file name
	if len(os.Args) < 2 {
		log.Fatalln("need a file name ")
	}
	fileName := os.Args[1]
	fileName = replaceChar(fileName)
	//get ip
	ip, err := getIP()
	if err != nil {
		log.Fatalf("could not get ip:%v", err)
	}
	//start a file server
	fs := http.FileServer(http.Dir("."))
	url := fmt.Sprintf("http://%s:7777/%s", ip, fileName)
	color.Cyan("🎊 %s 🎊 ", url)
	qrterminal.Generate(url, qrterminal.L, os.Stdout)
	//start webserver
	http.Handle("/", http.StripPrefix("/", fs))
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatalf("could not create web server :%v", err)
	}
}
