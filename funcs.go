package main

import (
	"bufio"
	"net"
	"os"
	"runtime"
	"strings"
)

var filePath string

func isIP(ip string) bool {
	// Check if IP is IPv4 or IPv6
	if strings.Contains(ip, ":") {
		return false
	}

	// Check if IP is valid
	if net.ParseIP(ip) == nil {
		return false
	}

	return true
}

func parseHostFile(file *os.File) (hosts []Host, err error) {
	// Read file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Split line by whitespace
		fields := strings.Fields(line)

		// Skip if line is empty
		if len(fields) == 0 {
			continue
		}

		// Skip if line doesn't have an IP address
		if !isIP(fields[0]) {
			continue
		}

		// Check if line has comment and include it. Remove # from comment
		// if len(fields) > 2 {
		// 	hosts = append(hosts, Host{
		// 		Hostname: fields[1],
		// 		IP:       fields[0],
		// 		Comment:  strings.TrimPrefix(fields[2], "#"),
		// 	})
		// 	continue
		// }

		// // Add host to hosts slice
		// hosts = append(hosts, Host{
		// 	Hostname: fields[1],
		// 	IP:       fields[0],
		// 	Comment:  fields[2],
		// })

		// Check if line has comment and include it. Remove # from comment
		if len(fields) > 2 {
			hosts = append(hosts, Host{
				Hostname: fields[1],
				IP:       fields[0],
				Comment:  strings.TrimPrefix(strings.TrimSpace(strings.Join(fields[2:], " ")), "#"),
			})
			continue
		}

		// Add host to hosts slice
		hosts = append(hosts, Host{
			Hostname: fields[1],
			IP:       fields[0],
			Comment:  fields[2],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func readHostFile() (hosts []Host, err error) {
	// check if windows or linux

	switch runtime.GOOS {
	case "windows":
		filePath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		filePath = "/etc/hosts"
	default:
		filePath = "/etc/hosts"
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Read file
	hosts, err = parseHostFile(file)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}
