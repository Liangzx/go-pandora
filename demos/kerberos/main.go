package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/spnego" // 导入SPNEGO包
)

// 示例 1: 客户端基础认证（使用用户名/密码）
func exp1() {
	// 1. 加载 krb5.conf 配置
	cfg, err := config.Load("./krb5.conf")
	if err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 2. 创建客户端实例（使用用户名、领域、密码）
	cl := client.NewWithPassword("liangzx/admin", "LIANGZXGG.COM", "dingjia", cfg, client.DisablePAFXFAST(true)) // 根据KDC配置可能需禁用PA-FX-FAST

	// 3. 执行登录认证
	err = cl.Login()
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	log.Println("Kerberos 登录成功！")
	return

	// 4. 获取特定服务的票据（例如访问HTTP服务）
	spn := "HTTP/host.your.domain.local" // 服务主体名称 (Service Principal Name)
	ticket, key, err := cl.GetServiceTicket(spn)
	if err != nil {
		log.Fatalf("获取服务票据失败: %v", err)
	}
	_ = key
	log.Printf("成功获取服务票据 for %s", spn)
	_ = ticket // 这里ticket可用于后续请求，如设置HTTP头Authorization: Negotiate <ticket>
}

// 示例 2：在 HTTP 客户端中使用 Kerberos 认证
func exp2() {
	cfg, err := config.Load("/path/to/krb5.conf")
	if err != nil {
		log.Fatal(err)
	}

	krb5Client := client.NewWithPassword("username", "YOUR.DOMAIN.LOCAL", "password", cfg, client.DisablePAFXFAST(true))
	err = krb5Client.Login()
	if err != nil {
		log.Fatal(err)
	}

	spn := "HTTP/host.your.domain.local" // 服务主体名称 (Service Principal Name)
	// 目标URL
	url := "http://protected-service.your.domain.local/api/endpoint"

	// 方法一：使用SPNEGO库自动处理协商认证（更推荐）
	// 创建SPNEGO支持的HTTP客户端

	spnClient := spnego.NewClient(krb5Client, &http.Client{}, spn)

	resp, err := spnClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n", body)
}

// 示例 3: 使用 Keytab 文件认证（更安全，适用于服务账户）
func exp3() {
	cfg, err := config.Load("/path/to/your/krb5.conf")
	if err != nil {
		log.Fatal(err)
	}

	// 加载 keytab 文件
	kt, err := keytab.Load("/path/to/your/service.keytab")
	if err != nil {
		log.Fatal(err)
	}

	// 创建使用 keytab 的客户端
	cl := client.NewWithKeytab("service-username", "YOUR.DOMAIN.LOCAL", kt, cfg, client.DisablePAFXFAST(true))

	err = cl.Login()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("使用 Keytab 登录成功！")
	// ... 后续获取服务票据等操作与示例1相同
}

func main() {
	exp1()
	// exp2()
	// exp3()
}
