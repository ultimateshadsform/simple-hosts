package main

import (
	"bufio"
	"fmt"
	"log/slog"
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

		if len(fields) > 2 {
			hosts = append(hosts, Host{
				Hostname: fields[1],
				IP:       fields[0],
				Comment:  strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.Join(fields[2:], " ")), "#")),
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

func writeHostFile(hosts []Host) error {
	// Only replace the line and not the whole file
	// Leave everything untouched and only replace the line needed
	// This is to avoid losing comments and other lines that are not managed by this app

	slog.Info("[App]: [writeHostFile] Writing hosts to file", "hosts", fmt.Sprintf("%+v", hosts))

	// Open file for writing

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	// Check if host exists and update else append it

	// Loop through hosts in the file and check if they exist in the hosts slice
	// If they exist, update them

	for _, host := range hosts {
		// Check if host exists
		if containsHost(hosts, host) {
			// Host exists, update it
			slog.Info("[App]: [writeHostFile] Host exists, updating it", "host", fmt.Sprintf("%+v", host))

			// Write host to file
			if _, err := file.WriteString(fmt.Sprintf("%s\t%s\t# %s\n", host.IP, host.Hostname, host.Comment)); err != nil {
				return err
			}
		} else {
			// Host does not exist, append it
			slog.Info("[App]: [writeHostFile] Host does not exist, appending it", "host", fmt.Sprintf("%+v", host))

			// Write host to file
			if _, err := file.WriteString(fmt.Sprintf("%s\t%s\t# %s\n", host.IP, host.Hostname, host.Comment)); err != nil {
				return err
			}
		}
	}

	return nil
}

func containsHost(hosts []Host, host Host) bool {
	for _, h := range hosts {
		if h.IP == host.IP {
			return true
		}
	}

	return false
}
