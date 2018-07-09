package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/mdp/qrterminal"
)

func main() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	myIP := ""
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				myIP = ipnet.IP.String()
			}
		}
	}
	if len(os.Args) < 2 {
		log.Fatalln("need a file name ")
	}
	fileName := os.Args[1]
	fs := http.FileServer(http.Dir("."))
	url := fmt.Sprintf("http://%s:7777/%s", myIP, fileName)
	qrterminal.Generate(url, qrterminal.L, os.Stdout)
	http.Handle("/", http.StripPrefix("/", fs))
	log.Println("File URL: " + url)
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatalf("could not create web server :%v", err)
	}
}
