/*
 The MIT License (MIT)

Copyright © 2023 Kasyanov Nikolay Alexeyevich (Unbewohnte)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package osutil

import (
	"os/exec"
	"runtime"
	"strings"
)

func GetSystemInfo() map[string]string {
	var sysInfo map[string]string = make(map[string]string)

	switch runtime.GOOS {
	case "windows":
		{
			output, err := exec.Command("systeminfo").Output()
			if err == nil {
				sysInfo["systeminfo"] = string(output)
			}

			output, err = exec.Command("wmic", "partition", "get", "name,size,type").Output()
			if err == nil {
				sysInfo["wmic"] = string(output)
			}
		}
	default:
		{
			output, err := exec.Command("lscpu").Output()
			if err == nil {
				sysInfo["lscpu"] = string(output)
			}

			output, err = exec.Command("lsblk").Output()
			if err == nil {
				sysInfo["lsblk"] = string(output)
			}

			output, err = exec.Command("lspci").Output()
			if err == nil {
				sysInfo["lspci"] = string(output)
			}

			output, err = exec.Command("uname", "-a").Output()
			if err == nil {
				sysInfo["uname"] = string(output)
			}

			output, err = exec.Command("hostname").Output()
			if err == nil {
				sysInfo["hostname"] = string(output)
			}
		}
	}

	return sysInfo
}

func GetHostname() string {
	output, err := exec.Command("hostname").Output()
	if err != nil {
		return "unknown_hostname"
	}

	return strings.TrimSpace(string(output))
}

func GetUsername() string {
	output, err := exec.Command("whoami").Output()
	if err != nil {
		return "unknown_username"
	}

	return strings.TrimSpace(string(output))
}
