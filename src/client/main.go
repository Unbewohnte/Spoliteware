/*
 The MIT License (MIT)

Copyright © 2023 Kasyanov Nikolay Alexeyevich (Unbewohnte)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"Unbewohnte/spolitewareClient/osutil"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	Version string = "client v0.1.1"
)

var (
	version    *bool   = flag.Bool("version", false, "Print version information and exit")
	serverAddr *string = flag.String("server", "http://localhost:13370/", "Set scheme://addr:port for receiving server")
)

// The structure of the generalised per computer data that is identical to the server's
// representation
type Data struct {
	Hostname string            `json:"hostname"`
	Username string            `json:"username"`
	System   map[string]string `json:"system"`
}

// Prints notification on the state of things
func greeting() {
	fmt.Printf(
		`I'm sorry to inform you, but you've (intentionally or not) launched a spyware on your machine!
But rest assured, I will do no harm to your computer and am to withdraw if you would wish so.
No data will be collected and sent without your permission.  

`)
}

// Asks for permission to collect all kinds of data on the machine
func askForAllPerms() bool {
	fmt.Printf(`Would you grant me a permission to [system information; files lookup; ]
(If no -> (optional) specify separate permissions afterwards) y/N: `)
	var input string = "n"
	fmt.Scanf("%s", &input)
	input = strings.ToLower(input)

	if strings.HasPrefix(input, "y") {
		fmt.Printf("Your approval is appreciated; I am to treat your system with care\n")
		return true
	} else {
		return false
	}
}

// Asks for permission to collect system information
func askForSystemInfo() bool {
	fmt.Printf("\nWould you allow me to look around and collect some information about your computer ? [y/N]: ")
	var input string = "n"
	fmt.Scanf("%s", &input)
	input = strings.ToLower(input)

	if strings.HasPrefix(input, "y") {
		fmt.Printf("System information scan allowed\n")
		return true
	} else {
		fmt.Printf("As you wish\n")
		return false
	}
}

// Asks whether a local copy of the whole collected data is needed
func localCopyNeeded() bool {
	fmt.Printf("\nDo you want to save a local copy as well ? [y/N]: ")
	var input string = "n"
	fmt.Scanf("%s", &input)
	input = strings.ToLower(input)

	if strings.HasPrefix(input, "y") {
		return true
	} else {
		return false
	}
}

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

	var data Data
	data.System = nil

	// Greet and ask for permissions
	greeting()
	if askForAllPerms() {
		data.System = osutil.GetSystemInfo()
	} else {
		if askForSystemInfo() {
			data.System = osutil.GetSystemInfo()
		}
	}

	if data.System == nil {
		// NOTICE! add new checks for new fields when they're added
		fmt.Printf("\nNothing to send. Bailing out\n")
		return
	}

	data.Hostname = osutil.GetHostname()
	data.Username = osutil.GetUsername()

	dataJson, err := json.MarshalIndent(&data, "", " ")
	if err != nil {
		return
	}

	// Try to send collected information to the server
	postBody := bytes.NewBuffer(dataJson)
	var retries uint8 = 0
	var response *http.Response
	for {
		if retries == 5 {
			fmt.Printf("\nFailed to send data\n")
			return
		}

		response, err = http.Post(*serverAddr, "application/json", postBody)
		if err != nil || response == nil {
			// Try to resend
			time.Sleep(time.Second * 5)
			retries++
			continue
		}

		if response.StatusCode != http.StatusOK {
			// Try to resend
			time.Sleep(time.Second * 5)
			retries++
			continue
		}
		break
	}

	// Successfully sent
	fmt.Printf("\nSuccesfully sent data\n")

	// Print thanks from the server
	defer response.Body.Close()
	thanks, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("\nFailed to read thanks message from the server. Better luck next time. Thank you !\n")
	}
	fmt.Printf("\n%s\n", thanks)

	// Create local copy if needed
	if localCopyNeeded() {
		wdir, err := os.Getwd()
		if err != nil {
			wdir = "./"
		}

		copyPath := filepath.Join(wdir, "spoliteware_scan_copy.txt")
		f, err := os.Create(copyPath)
		if err != nil {
			fmt.Printf("Failed to create a local copy: %s\n", err)
		} else {
			_, err = f.Write(dataJson)
			if err != nil {
				fmt.Printf("Failed to write to a file to save a local copy: %s\n", err)
			}
		}
		f.Close()
		fmt.Printf("Saved to %s\n", copyPath)
	}

	// Don't close the console window for windows users
	if runtime.GOOS == "windows" {
		fmt.Scanln()
	}
}
