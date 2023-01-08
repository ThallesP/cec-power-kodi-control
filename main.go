package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
)

type Rpc struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
}

func main() {
	target, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)

			if err != nil {
				panic(err)
			}

			var rpc Rpc

			err = json.Unmarshal(body, &rpc)

			if err == nil {
				if rpc.Method == "System.Shutdown" {
					cmd := exec.Command("bash", "-c", "echo 'standby 0' | cec-client -s")

					err := cmd.Run()

					if err != nil {
						fmt.Println(err)
					}

					fmt.Println("tv turned off")
					return
				}
			}

			req.Body.Close()
			req.Body = ioutil.NopCloser(bytes.NewReader(body))

			req.URL.Host = target.Host
			req.URL.Scheme = target.Scheme
		},
	}

	server := &http.Server{
		Addr:    ":8081",
		Handler: &proxy,
	}

	log.Println("Starting proxy server on port 8080")

	go listenForWOL()

	log.Fatal(server.ListenAndServe())
}

func listenForWOL() {
	conn, err := net.ListenPacket("udp4", ":9")

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	for {
		conn.ReadFrom(buf)

		cmd := exec.Command("bash", "-c", "echo 'on 0' | cec-client -s")

		err := cmd.Run()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("tv turned on")
	}
}
