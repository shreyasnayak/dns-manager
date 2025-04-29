package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ipBytes), nil
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no non-loopback local IP address found")
}

func getDNSRecordIP(domain string) (string, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found for %s", domain)
	}

	return ips[0], nil
}

func updateGoDaddyDNS(domain, recordType, godaddyAPIKey, godaddyAPISecret string) error {
	currentIP, err := getPublicIP()
	if err != nil {
		return fmt.Errorf("failed to get public IP: %v", err)
	}

	parts := strings.SplitN(domain, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid domain format: %s", domain)
	}
	name := parts[0]
	rootDomain := parts[1]

	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s", rootDomain, recordType, name)

	payload := fmt.Sprintf(`[{"data":"%s","ttl":600}]`, currentIP)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", godaddyAPIKey, godaddyAPISecret))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update GoDaddy DNS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("GoDaddy API error: %s", string(body))
}

func main() {
	getCmd := flag.String("get-ip", "", "Get the IP address, use 'public' or 'local'")
	getDnsIpCmd := flag.Bool("get-dns-ip", false, "Fetch the current DNS IP address for a domain")
	putDnsIpCmd := flag.Bool("put-dns-ip", false, "Update DNS IP address for a domain")
	domain := flag.String("domain", "", "Domain name (example.com)")
	recordType := flag.String("type", "A", "DNS record type (default 'A')")
	domainRegistrar := flag.String("domain-registrar", "", "DNS registrar (GODADDY, CLOUDFLARE, etc.)")
	godaddyAPIKey := flag.String("godaddy-api-key", "", "GoDaddy API Key")
	godaddyAPISecret := flag.String("godaddy-api-secret", "", "GoDaddy API Secret")

	flag.Parse()

	if *getCmd != "" {
		switch strings.ToLower(*getCmd) {
		case "public":
			ip, err := getPublicIP()
			if err != nil {
				fmt.Println("Error fetching public IP:", err)
				os.Exit(1)
			}
			fmt.Println(ip)

		case "local":
			ip, err := getLocalIP()
			if err != nil {
				fmt.Println("Error fetching local IP:", err)
				os.Exit(1)
			}
			fmt.Println(ip)

		default:
			fmt.Println("Invalid option. Use 'public' or 'local'.")
			os.Exit(1)
		}
		return
	} else if *getDnsIpCmd {
		if *domain == "" {
			fmt.Println("Error: --domain is required with --get-dns-ip")
			os.Exit(1)
		}
		ip, err := getDNSRecordIP(*domain)
		if err != nil {
			fmt.Println("Error fetching DNS IP:", err)
			os.Exit(1)
		}
		fmt.Println(ip)
		return
	} else if *putDnsIpCmd {
		if *domain == "" || *domainRegistrar == "" {
			fmt.Println("Error: --domain and --domain-registrar are required with --put-dns-ip")
			os.Exit(1)
		}

		switch strings.ToUpper(*domainRegistrar) {
		case "GODADDY":
			if *godaddyAPIKey == "" || *godaddyAPISecret == "" {
				fmt.Println("Error: --godaddy-api-key and --godaddy-api-secret are required for GoDaddy")
				os.Exit(1)
			}
			err := updateGoDaddyDNS(*domain, *recordType, *godaddyAPIKey, *godaddyAPISecret)
			if err != nil {
				fmt.Println("Failed to update GoDaddy DNS:", err)
				os.Exit(1)
			}
			fmt.Println("Successfully updated GoDaddy DNS record for", *domain)

		default:
			fmt.Println("Error: Unsupported domain registrar:", *domainRegistrar)
			os.Exit(1)
		}
		return

	} else {
		fmt.Println("Usage: dnsmanager [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
