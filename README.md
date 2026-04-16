<div align="center">

<pre>
   ___ _             _____ _       _ _
  / __| |_  ___ _  _|_   _| |_  __| | | _
 | (__| ' \/ -_) || | | | | ' \/ _` | || |
  \___|_||_\___|\_, |_| |_|_||_\__,_|_||_|
                |_|
</pre>

**SSH server management for terminal users**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/lymboy/easyssh.svg)](https://github.com/lymboy/easyssh/releases)

</div>

---

## 📺 Preview

> **A picture is worth a thousand words**

```bash
$ easyssh server ls

════════════════════════════════════════════════════════════════════════════════

ID    NAME               GROUP        HOST               USER         STATUS
──────────────────────────────────────────────────────────────────────────────
[0]   web-master         prod         192.168.1.10      root         ○ Idle
[1]   web-slave          prod         192.168.1.11      root         ○ Idle
[2]   db-primary         prod         192.168.1.20     admin         ○ Idle
──────────────────────────────────────────────────────────────────────────────
[3]   dev-server         dev          10.0.0.5          dev          ○ Idle
[4]   test-server        dev          10.0.0.6          dev          ○ Idle
──────────────────────────────────────────────────────────────────────────────
Total: 5 servers | prod: 3 | dev: 2

$ easyssh server 0
# Connected to root@192.168.1.10
```

---

## ⚡ Quick Start

### Install

**From binary (recommended):**
```bash
# macOS
wget https://github.com/lymboy/easyssh/releases/latest/download/easyssh_darwin_amd64 -O /usr/local/bin/easyssh
chmod +x /usr/local/bin/easyssh

# Linux
wget https://github.com/lymboy/easyssh/releases/latest/download/easyssh_linux_amd64 -O /usr/local/bin/easyssh
chmod +x /usr/local/bin/easyssh
```

**From source:**
```bash
go install github.com/lymboy/easyssh@latest
```

**Build from source (cross-platform):**
```bash
git clone https://github.com/lymboy/easyssh.git
cd easyssh

# macOS
GOOS=darwin GOARCH=amd64 go build -o easyssh_darwin_amd64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o easyssh_linux_amd64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o easyssh_windows_amd64.exe .

# Or build for your current platform
go build -o easyssh .
```

### Configure

```bash
mkdir -p ~/.easyssh
cat > ~/.easyssh/easy_config.yaml << 'EOF'
ssh:
  key: "id_rsa"
  keep_alive: true
  keep_alive_interval: "60s"

server:
  - group: "prod"
    name: "web-master"
    host: "192.168.1.10"
    user: "root"
  - group: "prod"
    name: "web-slave"
    host: "192.168.1.11"
    user: "root"
  - group: "dev"
    name: "dev-server"
    host: "10.0.0.5"
    user: "dev"
EOF
```

### Use

```bash
easyssh server ls     # List all servers
easyssh server 0      # Connect by index
easyssh server web-master  # Connect by name
```

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🔄 **Auto Keep-Alive** | Maintains connection during long idle periods |
| 🔐 **Dual Auth** | Automatic detection of SSH Key or password |
| 📋 **Flexible Connection** | Connect by index number or server name |
| 🎨 **Group Management** | Organize servers by environment |
| ⚡ **Lightweight** | Single binary, no dependencies |
| 🖥️ **Terminal Native** | Works in any terminal, no GUI needed |

---

## 🆚 Why EasySSH?

| | **EasySSH** | **XShell** | **tabby** |
|-|:-----------:|:----------:|:---------:|
| **Size** | ~5MB | ~50MB | ~100MB |
| **Dependencies** | None | GUI libs | Electron |
| **Terminal native** | ✅ | ❌ | ❌ |
| **Linux/Server friendly** | ✅ | ❌ | ❌ |
| **Server management** | ✅ | ✅ | ✅ |

**Perfect for:**
- 🖥️ Server administrators
- 🔧 DevOps engineers
- 👨‍💻 Developers who live in the terminal
- 🐧 Linux users

---

## 📄 Configuration

<details>
<summary>Click to expand full configuration reference</summary>

```yaml
ssh:
  key: "id_rsa"              # SSH private key filename (default: id_rsa)
  keep_alive: true           # Enable keep-alive (default: true)
  keep_alive_interval: "60s" # Keep-alive interval (default: 60s)

server:
  - group: "prod"           # Group name (optional)
    name: "web-master"       # Server name (required)
    host: "192.168.1.10"     # Server host (required)
    port: 22                 # SSH port (default: 22)
    user: "root"             # SSH user (default: current user)
    password: ""             # Password (optional, prefer SSH key)
    desc: "Production web"   # Description (optional)
```

</details>

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

Made with ❤️ by [lymboy](https://github.com/lymboy)

</div>
