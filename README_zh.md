<div align="center">

<pre>
   ___ _             _____ _       _ _
  / __| |_  ___ _  _|_   _| |_  __| | | _
 | (__| ' \/ -_) || | | | | ' \/ _` | || |
  \___|_||_\___|\_, |_| |_|_||_\__,_|_||_|
                |_|
</pre>

**SSH 服务器管理工具 - 为终端用户设计**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

## 📺 效果预览

```bash
$ easyssh server ls

════════════════════════════════════════════════════════════════════════
ID    NAME          GROUP        HOST               USER         STATUS
────────────────────────────────────────────────────────────────────────
[0]   web-master    prod         192.168.1.10      root         ● Connected
[1]   web-slave     prod         192.168.1.11      root         ○ Idle
[2]   db-primary    prod         192.168.1.20     admin         ○ Idle
────────────────────────────────────────────────────────────────────────
[3]   uat-web       uat          10.0.0.10         deploy       ○ Idle
[4]   uat-db        uat          10.0.0.11         deploy       ○ Idle
────────────────────────────────────────────────────────────────────────
Total: 5 servers | prod: 3 | uat: 2

$ easyssh server 0
# 已连接到 root@192.168.1.10
```

---

## ⚡ 快速开始

### 安装

**方式一：下载二进制文件（推荐）**

```bash
# macOS
wget https://github.com/lymboy/easyssh/releases/latest/download/easyssh_darwin_amd64 -O /usr/local/bin/easyssh
chmod +x /usr/local/bin/easyssh

# Linux
wget https://github.com/lymboy/easyssh/releases/latest/download/easyssh_linux_amd64 -O /usr/local/bin/easyssh
chmod +x /usr/local/bin/easyssh
```

**方式二：从源码编译**

```bash
go install github.com/lymboy/easyssh@latest
```

**方式三：交叉编译**

```bash
git clone https://github.com/lymboy/easyssh.git
cd easyssh

# macOS
GOOS=darwin GOARCH=amd64 go build -o easyssh_darwin_amd64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o easyssh_linux_amd64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o easyssh_windows_amd64.exe .

# 或直接编译当前平台
go build -o easyssh .
```

### 配置

```bash
# 1. 创建配置目录
mkdir -p ~/.easyssh

# 2. 复制示例配置
cp easy_config.yaml.example ~/.easyssh/easy_config.yaml

# 3. 编辑配置文件
vim ~/.easyssh/easy_config.yaml
```

### 使用

```bash
# 列出所有服务器
easyssh server ls

# 按索引连接
easyssh server 0

# 按名称连接
easyssh server web-master

# 添加新服务器（交互式）
easyssh add

# 添加新服务器（命令行）
easyssh add -g prod -e web -u root -i "192.168.1.10,192.168.1.11"

# 删除服务器
easyssh remove web-1
```

---

## ✨ 功能特性

| 功能 | 说明 |
|------|------|
| 🔄 **自动保活** | 长时间不操作自动发送心跳，防止断开 |
| 🔐 **双认证支持** | 自动识别 SSH Key 或密码认证 |
| 📋 **灵活连接** | 支持索引号或服务器名称连接 |
| 🎨 **分组管理** | 按环境（prod/uat/dev）分组管理 |
| ➕ **便捷添加** | 交互式添加服务器，支持批量 IP |
| 🔗 **连接复用** | 支持 SSH ControlMaster，多个终端共享连接 |
| ⚡ **轻量快速** | 单个二进制文件，无额外依赖 |
| 🖥️ **纯终端** | 在任何终端中使用，无需 GUI |

---

## 🔗 SSH 连接复用

启用 ControlMaster 后，多个终端可以共享同一个 SSH 连接，快速重连无需重复输入密码。

### 方式一：一键配置（推荐）

```bash
# 自动配置 SSH ControlMaster
$ easyssh setup

  EasySSH Setup - SSH ControlMaster Configuration

  Will add the following to ~/.ssh/config:

    # EasySSH Connection Reuse
    Host *
        ControlMaster auto
        ControlPath ~/.ssh/sockets/%r@%h-%p
        ControlPersist no
        ServerAliveInterval 60

  ✓ SSH ControlMaster configured successfully!

  Created:
    • ~/.ssh/config (updated)
    • ~/.ssh/sockets (created)

  Next steps:
    1. Enable in EasySSH config: use_system_ssh: true
```

### 方式二：手动配置

如果你想手动配置，编辑 `~/.ssh/config`：

```bash
# 添加以下内容到 ~/.ssh/config 文件
Host *
    ControlMaster auto
    ControlPath ~/.ssh/sockets/%r@%h-%p
    ControlPersist no
    ServerAliveInterval 60

# 创建 socket 目录
mkdir -p ~/.ssh/sockets
```

### 启用连接复用

编辑 `~/.easyssh/easy_config.yaml`：

```yaml
ssh:
  use_system_ssh: true
```

### 效果

```bash
# 终端1：连接服务器（需要输入密码）
$ easyssh server 0

# 终端2：再次连接同一服务器（秒连，不需要密码）
$ easyssh server 0

# 终端3：查看状态，显示已连接
$ easyssh server ls
# web-master 显示 "● Connected"
```

---

## 📝 配置文件详解

```yaml
ssh:
  # SSH 私钥文件名（位于 ~/.ssh/ 目录）
  key: "id_rsa"                    # 默认: id_rsa

  # 是否启用保活
  keep_alive: true                 # 默认: true

  # 保活间隔
  keep_alive_interval: "60s"      # 默认: 60s

  # 是否使用系统 SSH（支持 ControlMaster）
  use_system_ssh: true             # 默认: false

server:
  - group: "prod"                  # 分组名（可选）
    name: "web-master"             # 服务器名称（必填）
    host: "192.168.1.10"           # 服务器地址（必填）
    port: 22                       # SSH 端口（默认: 22）
    user: "root"                   # SSH 用户（默认: 当前用户）
    password: ""                   # 密码（可选，建议使用密钥）
    desc: "生产环境Web服务器"       # 描述（可选）
```

---

## 🆚 与其他工具对比

| 特性 | **EasySSH** | **XShell** | **tabby** |
|------|:-----------:|:----------:|:---------:|
| 安装大小 | ~5MB | ~50MB | ~100MB |
| 依赖 | 无 | GUI 库 | Electron |
| 终端原生 | ✅ | ❌ | ❌ |
| Linux/服务器友好 | ✅ | ❌ | ❌ |
| 服务器管理 | ✅ | ✅ | ✅ |
| 免费开源 | ✅ | ❌ | ✅ |

**适合人群：**
- 🖥️ 服务器运维人员
- 🔧 DevOps 工程师
- 👨‍💻 终端重度用户
- 🐧 Linux 用户

---

## 📖 命令参考

### 服务器连接

```bash
easyssh server ls          # 列出所有服务器
easyssh server 0           # 按索引连接
easyssh server web-master  # 按名称连接
```

### 服务器管理

```bash
easyssh add                           # 交互式添加
easyssh add -g prod -e web -u root -i "192.168.1.10,192.168.1.11"
easyssh remove 0                      # 删除第一台服务器
easyssh remove web-1                  # 按名称删除
```

### add 命令参数

| 参数 | 缩写 | 默认值 | 说明 |
|------|------|--------|------|
| `--group` | `-g` | default | 服务器分组 |
| `--env` | `-e` | uat | 环境/服务名 |
| `--user` | `-u` | 当前用户 | SSH 用户 |
| `--ips` | `-i` | - | IP 地址（逗号或空格分隔） |
| `--port` | `-p` | 22 | SSH 端口 |

---

## 🤝 参与贡献

欢迎提交 Issue 和 Pull Request！

---

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

<div align="center">

Made with ❤️ by [lymboy](https://github.com/lymboy)

</div>
