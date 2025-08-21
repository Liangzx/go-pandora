package main

import (
	"fmt"
	"net"
	"os"

	"github.com/golang-auth/go-gssapi/v2"
	_ "github.com/golang-auth/go-gssapi/v2/krb5" // 注册Kerberos机制
)

func exp1() {
	// 1. 创建GSS上下文
	ctx := gssapi.NewMech("kerberos_v5")
	service := "ldap/ldap.example.com" // 服务Principal格式：primary/instance@REALM[6](@ref)

	// 2. 初始化上下文（客户端）
	flags := gssapi.ContextFlagInteg | gssapi.ContextFlagConf | gssapi.ContextFlagMutual
	err := ctx.Initiate(service, flags, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "初始化上下文失败:", err)
		os.Exit(1)
	}

	// 3. 连接服务端
	conn, err := net.Dial("tcp", "ldap.example.com:389")
	if err != nil {
		fmt.Fprintln(os.Stderr, "连接服务端失败:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 4. 安全上下文协商循环
	var inToken, outToken []byte
	for !ctx.IsEstablished() {
		outToken, err = ctx.Continue(inToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, "协商失败:", err)
			break
		}
		if len(outToken) > 0 {
			if _, err = conn.Write(outToken); err != nil {
				fmt.Fprintln(os.Stderr, "发送令牌失败:", err)
				break
			}
		}
		if !ctx.IsEstablished() {
			inToken = make([]byte, 4096)
			n, err := conn.Read(inToken)
			if err != nil {
				fmt.Fprintln(os.Stderr, "接收令牌失败:", err)
				break
			}
			inToken = inToken[:n]
		}
	}

	// 5. 发送加密数据
	msg := "Hello, secure world!"
	sealed := true // 启用加密
	wrapped, err := ctx.Wrap([]byte(msg), sealed)
	if err != nil {
		fmt.Fprintln(os.Stderr, "加密数据失败:", err)
		return
	}
	if _, err = conn.Write(wrapped); err != nil {
		fmt.Fprintln(os.Stderr, "发送加密数据失败:", err)
	}
}

func main() {
	exp1()
}
