---
title: Subscription Format
description: Subscription format options and examples
---

Xray Checker supports four different formats for proxy configuration. Use the [environment variable](/configuration/envs#subscription_url) `SUBSCRIPTION_URL` for setup.

For information about how proxies are verified, see [check methods](/configuration/check-methods).

### 1. Subscription URL (Default)

Standard subscription URL returning Base64 encoded list of proxy links.

Example:

```bash
SUBSCRIPTION_URL="https://example.com/subscription"
```

Requirements:

- HTTPS URL
- Returns Base64 encoded content
- Content is newline-separated proxy URLs
- Supports standard User-Agent headers

Headers sent:

```
Accept: */*
User-Agent: Xray-Checker
```

### 2. Base64 String

Direct Base64 encoded string containing proxy configuration links.

Example:

```bash
SUBSCRIPTION_URL="dmxlc3M6Ly91dWlkQGV4YW1wbGUuY29tOjQ0MyVlbmNyeXB0aW9uPW5vbmUmc2VjdXJpdHk9dGxzI3Byb3h5MQ=="
```

Content format (before encoding):

```
vless://uuid@example.com:443?encryption=none&security=tls#proxy1
trojan://password@example.com:443?security=tls#proxy2
vmess://base64encodedconfig
ss://base64encodedconfig
```

### 3. V2Ray JSON File

Single JSON configuration file in V2Ray/Xray format.

Example:

```bash
SUBSCRIPTION_URL="file:///path/to/config.json"
```

File format:

```json
{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "example.com",
            "port": 443,
            "users": [
              {
                "id": "uuid",
                "encryption": "none"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "tls"
      }
    }
  ]
}
```

### 4. Configuration Folder

Directory containing multiple V2Ray/Xray JSON configuration files.

Example:

```bash
SUBSCRIPTION_URL="folder:///path/to/configs"
```

Requirements:

- Directory must contain .json files
- Each file follows V2Ray JSON format
- Files are processed in alphabetical order
- Invalid files are skipped with warning
