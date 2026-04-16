package ssh

import (
	"easyssh/config"
	"easyssh/util"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
	"io"
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

func getPrivateKeySign() (ssh.Signer, error) {
	// 读取私钥文件
	privateBytes, err := os.ReadFile(getPublicKeyPath())
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// 解析私钥
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return private, nil
}

func getPassword() string {
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

	sign, err := getPrivateKeySign()
	if err == nil && sign != nil {
		authList = append(authList, ssh.PublicKeys(sign))
	}

	if c.Password != "" {
		authList = append(authList, ssh.Password(c.Password))
	}

	clientConfig := ssh.ClientConfig{
		User:            c.Username,
		Auth:            authList,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	var sshClient *ssh.Client
	const retry = 3

	for i := 0; i < retry; i++ {
		sshClient, err = ssh.Dial("tcp", addr, &clientConfig)
		if err == nil {
			break
		}
		// Try password prompt on failure
		if i < retry-1 {
			fmt.Printf("\n\033[33mAuthentication failed for %s@%s\033[0m\n", c.Username, c.IP)
			fmt.Printf("Enter password: ")
			password := getPassword()
			authList = []ssh.AuthMethod{ssh.Password(password)}
			clientConfig.Auth = authList
		}
	}

	if err != nil {
		return fmt.Errorf("failed to connect to %s@%s:%d after %d attempts: %w", c.Username, c.IP, c.Port, retry, err)
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
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}
	defer func(fd int, oldState *terminal.State) {
		_ = terminal.Restore(fd, oldState)
	}(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
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

	// Keep-alive logic
	if config.GetConf().GetSSHConfig().KeepAlive {
		keepAliveInterval := config.GetConf().GetSSHConfig().GetKeepAliveInterval()
		keepAliveTicker := time.NewTicker(keepAliveInterval)
		defer keepAliveTicker.Stop()

		go func() {
			for range keepAliveTicker.C {
				_, err := session.SendRequest("keepalive@openssh.com", true, nil)
				if err != nil {
					// Connection likely dead, will be handled by session.Wait()
				}
			}
		}()
	}

	_ = session.Shell()
	_ = session.Wait()
	return nil
}
