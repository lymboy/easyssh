package ssh

import (
	"easyssh/config"
	"easyssh/util"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Cli struct {
	IP         string
	Username   string
	Password   string
	Port       int
	client     *ssh.Client
	LastResult string
}

func New(server *config.Server) *Cli {
	cli := new(Cli)
	cli.IP = server.GetHost()
	cli.Username = server.GetUser()
	cli.Password = server.GetPassword()
	cli.Port = server.GetPort()
	return cli
}

func getPublicKeyPath() string {
	pubKey := config.GetConf().GetSSHConfig().GetKey()
	pubKeyPath := filepath.Join(util.GetHomeDir(), ".ssh", pubKey)
	if !util.IsFile(pubKeyPath) {
		_, _ = os.Create(pubKeyPath)
	}
	return pubKeyPath
}

func getPrivateKeySign() ssh.Signer {
	// 读取私钥文件
	privateBytes, err := os.ReadFile(getPublicKeyPath())
	if err != nil {
		log.Fatalf("Failed to load private key (%s)", err)
	}

	// 解析私钥
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key (%s)", err)
	}
	return private
}

func getPassword() string {
	//var password string
	//fmt.Print("Enter password: ")
	//_, _ = fmt.Scanln(&password)
	//return password

	fmt.Print("Enter password: ")
	// 获取当前的终端
	oldTerm, _ := term.GetState(syscall.Stdin)

	// 设置为原始模式，以便我们可以读取每个按键
	_, _ = term.MakeRaw(syscall.Stdin)

	// 读取用户输入
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	// 将读取的字节转换为字符串并输出
	password := string(bytePassword)
	// 恢复原来的终端状态
	_ = term.Restore(syscall.Stdin, oldTerm)
	fmt.Println() // 换行打印一个换行符
	return password
}

func (c *Cli) connect() error {
	authList := make([]ssh.AuthMethod, 0)
	sign := getPrivateKeySign()
	if sign != nil {
		authList = append(authList, ssh.PublicKeys(sign))
	}
	authList = append(authList, ssh.Password(c.Password))

	clientConfig := ssh.ClientConfig{
		User:            c.Username,
		Auth:            authList,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略服务器的 HostKey 验证
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	var sshClient *ssh.Client
	var err error
	const retry = 3
	for i := 0; i < retry; i++ {
		sshClient, err = ssh.Dial("tcp", addr, &clientConfig)
		if err != nil {
			authList = authList[:len(authList)-1]
			authList = append(authList, ssh.Password(getPassword()))
			clientConfig.Auth = authList
		}
	}
	if err != nil {
		log.Printf("Failed to dial (attempt %d times): %s", retry, err)
		return err
	}
	c.client = sshClient
	return nil
}

func (c *Cli) RunTerminal(stdout, stderr io.Writer) error {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}
	ticker := time.NewTicker(3 * time.Second) // 每3分钟发送一次探活消息
	defer ticker.Stop()
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer func(fd int, oldState *terminal.State) {
		err := terminal.Restore(fd, oldState)
		if err != nil {

		}
	}(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	// 处理终端大小变化
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGWINCH)
	go func() {
		for range signals {
			if termWidth, termHeight, err = terminal.GetSize(fd); err == nil {
				_ = session.WindowChange(termHeight, termWidth)
			}
		}
	}()

	// 开启探活
	if config.GetConf().GetSSHConfig().KeepAlive {
		// 设置 keepalive 参数
		ticker := time.NewTicker(5 * time.Second) // 每分钟发送一次探活消息
		defer ticker.Stop()
		// 发送探活消息
		go func() {
			for range ticker.C {
				// 向服务器发送一个全局请求，比如 keepalive@openssh.com
				_, err := session.SendRequest("keepalive@openssh.com", true, nil)
				if err != nil {
					log.Printf("error sending keep-alive message: %s", err)
					// 发生错误时，可以选择重新建立连接或者执行其他处理逻辑
				}
			}
		}()
	}

	_ = session.Shell()
	_ = session.Wait()
	return nil
}
