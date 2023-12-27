package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
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

	// Scan file by new line
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

		// Check if host exists and update

		if containsHost(hosts, Host{
			IP: fields[0],
		}) {
			// Host exists, update it

			// Update hostname
			for _, host := range hosts {
				if host.IP == fields[0] {
					fields[0] = host.IP
					fields[1] = host.Hostname
					fields[2] = host.Comment
					break
				}
			}

			// Write updated line to file
			_, err := file.WriteString(strings.Join(fields, " ") + "\n")
			if err != nil {
				slog.Error("[App]: [writeHostFile] Error writing hosts to file", "error", err.Error())
				return err
			}
		} else {
			// Host does not exist, Append it

			// Write line to file
			_, err := file.WriteString(line + "\n")
			if err != nil {
				slog.Error("[App]: [writeHostFile] Error writing hosts to file", "error", err.Error())
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("[App]: [writeHostFile] Error scanning hosts file", "error", err.Error())
		return err
	}

	// Write hosts to file

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

func updateHost(hosts []Host, host Host) []Host {
	for i, h := range hosts {
		if h.IP == host.IP {
			hosts[i] = host
			break
		}
	}

	return hosts
}
