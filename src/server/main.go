/*
 The MIT License (MIT)

Copyright © 2023 Kasyanov Nikolay Alexeyevich (Unbewohnte)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Data struct {
	Hostname string            `json:"hostname"`
	Username string            `json:"username"`
	System   map[string]string `json:"system"`
}

const Version string = "server v0.1.0"

var (
	version  *bool   = flag.Bool("version", false, "Print version information and exit")
	port     *uint   = flag.Uint("port", 13370, "Set port to listen on")
	keyFile  *string = flag.String("key", "", "SSL private key file path")
	certFile *string = flag.String("cert", "", "SSL certificate file path")
	verbose  *bool   = flag.Bool("verbose", false, "Print user data-centric messages or not")
)

func main() {
	fmt.Printf(
		`
███████╗██████╗  ██████╗ ██╗     ██╗████████╗███████╗██╗    ██╗ █████╗ ██████╗ ███████╗
██╔════╝██╔══██╗██╔═══██╗██║     ██║╚══██╔══╝██╔════╝██║    ██║██╔══██╗██╔══██╗██╔════╝
███████╗██████╔╝██║   ██║██║     ██║   ██║   █████╗  ██║ █╗ ██║███████║██████╔╝█████╗  
╚════██║██╔═══╝ ██║   ██║██║     ██║   ██║   ██╔══╝  ██║███╗██║██╔══██║██╔══██╗██╔══╝  
███████║██║     ╚██████╔╝███████╗██║   ██║   ███████╗╚███╔███╔╝██║  ██║██║  ██║███████╗
╚══════╝╚═╝      ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝
																																											
`)

	flag.Parse()
	if *version {
		fmt.Printf("by Kasyanov Nikolay Alexeyevich (Unbewohnte) %s\n", Version)
		return
	}

	executablePath, err := os.Executable()
	if err != nil {
		fmt.Printf("[ERROR] Failed to determine executable's path: %s\n", err)
		return
	}
	dataDirPath := filepath.Join(filepath.Dir(executablePath), "collected_data")

	err = os.MkdirAll(dataDirPath, os.ModePerm)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create data directory: %s\n", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if *keyFile != "" && *certFile != "" {
			proto := req.Header.Get("x-forwarded-proto")
			if strings.ToLower(proto) == "http" {
				http.Redirect(w, req, fmt.Sprintf("https://%s%s", req.Host, req.URL), http.StatusMovedPermanently)
				return
			}
		}

		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content is not application/json", http.StatusBadRequest)
			return
		}

		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			// Internal error, - no need to tell about it to the other side
			if *verbose {
				fmt.Printf("[ERROR] Failed to read request's body from %s: %s\n", req.RemoteAddr, err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			if *verbose {
				fmt.Printf("[ERROR] Failed to split host and port from %s: %s\n", req.RemoteAddr, err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		var clientData Data
		err = json.Unmarshal(body, &clientData)
		if err != nil {
			if *verbose {
				fmt.Printf("[ERROR] Failed to unmarshal data from %s: %s\n", req.RemoteAddr, err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		clientDir := filepath.Join(dataDirPath, fmt.Sprintf("%s_%s_%s", clientData.Hostname, clientData.Username, ip))
		err = os.MkdirAll(clientDir, os.ModePerm)
		if err != nil {
			if *verbose {
				fmt.Printf("[ERROR] Failed to create directory for %s: %s\n", ip, err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		clientFilePath := filepath.Join(clientDir, fmt.Sprintf("data_%d.json", time.Now().Unix()))
		clientFile, err := os.OpenFile(clientFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			if *verbose {
				fmt.Printf("[ERROR] Failed to create file for %s: %s\n", ip, err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		defer clientFile.Close()

		_, err = clientFile.Write(body)
		if err != nil {
			if *verbose {
				fmt.Printf("[ERROR] Failed to append data to %s: %s\n", clientFilePath, err)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusOK)
		if *verbose {
			fmt.Printf("[INFO] Received data from %s_%s_%s\n", clientData.Hostname, clientData.Username, ip)
		}
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	if *keyFile != "" && *certFile != "" {
		fmt.Printf("[INFO] Using TLS\n")
		fmt.Printf("[INFO] Server listening on port %d\n", *port)
		server.ListenAndServeTLS(*certFile, *keyFile)
	} else {
		fmt.Printf("[INFO] Not using TLS\n")
		fmt.Printf("[INFO] Server listening on port %d\n", *port)
		server.ListenAndServe()
	}
}
