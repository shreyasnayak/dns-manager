

![Main-CI](https://github.com/shreyasnayak/dns-manager/actions/workflows/ci.yml/badge.svg) ![GitHub release](https://img.shields.io/github/v/release/shreyasnayak/dns-manager) ![License](https://img.shields.io/github/license/shreyasnayak/dns-manager)


**DNS Manager** is a simple command-line tool that helps users manage DNS records for their domain names. It streamlines the process of adding, updating, or deleting DNS records without needing to navigate complex registrar interfaces.

---

## 🚀 Why DNS Manager?

This tool was inspired by a common pain point when using services like AWS EC2. When an EC2 instance is stopped and restarted, its public IP changes—requiring manual DNS updates. DNS Manager automates this task, saving time and effort.

---

## 🤔 Why not just use cURL, Python, or another language?

While it’s possible to script these operations using existing tools, the ultimate goal of DNS Manager is to provide a secure and user-friendly experience. This includes:

- Safely storing API keys and secrets using encryption.
- Automatically decrypting and using them when needed.
- Building a dedicated CLI for clean, repeatable operations.

This project is just the beginning—expect more security and automation features in future releases.

---

## 📦 Usage

Here are some common commands:

```bash
# Get the local IP address
dnsmanager --get-ip local

# Get the public IP address
dnsmanager --get-ip public

# Get the current IP address associated with a DNS record
dnsmanager --get-dns-ip --domain www.example.com

# Update a DNS record using GoDaddy API
dnsmanager --put-dns-ip --domain www.example.com \
           --domain-registrar GODADDY \
           --godaddy-api-key <KEY> \
           --godaddy-api-secret <TOKEN>
```