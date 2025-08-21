package http_client

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

const (
	spn = "HTTP/liangzxgg.com"
	url = "http://liangzxgg.com:9080/world"
	// negotiate = "Negotiate YIICbAYGKwYBBQUCoIICYDCCAlygDTALBgkqhkiG9xIBAgKiggJJBIICRWCCAkEGCSqGSIb3EgECAgEAboICMDCCAiygAwIBBaEDAgEOogcDBQAAAAAAo4IBY2GCAV8wggFboAMCAQWhDxsNTElBTkdaWEdHLkNPTaIgMB6gAwIBAaEXMBUbBEhUVFAbDWxpYW5nenhnZy5jb22jggEfMIIBG6ADAgESoQMCAQKiggENBIIBCYe9kY0z5/Re8VPFfcM5aL3dC3eqxxqPX6UVSmB9A2bZi5fYKY/ufhU28XShdFLW3PQX5pMRXQL+KxYhInbJl940BK5r+qyLtSxwnbdl7UqTZzwoer1EpRKiP/j1mPD7f1disIGHSA4AnvBXUrzTRGeWz2Thsf1qja1YUAa0VgbhLd0WxTXIl1o6slicCOFd2S2+ka1PAnspysqEUjCaY3y3QT7SHETitDdVt2n07YqandtSZxl+2pHn5IYz7ko/InPis4zSAIAm3D7VWRiD2SNRUaZ4c1q43+Q7bKMh5i6R+Gj1EetkgpqKXfwhh3IMuANPt2riH10W/ckJH0fmpnPiy/usZnws1pSkga8wgaygAwIBEqEDAgECooGfBIGciaV29ZEWV+z1NJKHS6XdBZw6GCk0ECdrppy44H6NDcPpxuqHIZdjFe2olm5fvZmqpn5Y4UE67fVgbb7oDFIyleiD0hTrO1Vrh5CmWnUOouEkgH9aCcnk+J/GB00/wOOOLuJQkzCjk1N5aAg6V/th41eNqoVyJ6UzCAJzttKUUNmcildXi0KI0Yr0A8ZE1Y+cgZ32e3/XrfrUbnb7"
	//    YYIBSjCCAUagAwIBBaEPGw1MSUFOR1pYR0cuQ09NoiAwHqADAgEBoRcwFRsESFRUUBsNbGlhbmd6eGdnLmNvbaOCAQowggEGoAMCARKhAwIBAqKB+QSB9hkzd2SjBfKkv1RA0HTC3JY9dOXD6w9nTR5nl10wUeR4AkrqYWNGZ+RQyLLE7W7NWMMKyUrq3qnFqeSRlt3GG/fKzp486asHV/8ayfnc2WM5esCBdGIGnhDoFPPZxqnRPgogRKLa6yuoqq+pi8TIXy7zl+Uv+b+gAoNt1pF8iO1ZPaa2x2C8dshghj12HFY/kkGw9t42YdYyy3Tqw6hN7FflPPmzY0MFSQvFD2CWK8IKiOe4qOQ7WHmdXXmkLeKTgday7oVV+jtata8/4UgVujRimoIkB+nRjwLeVexaFCr8ePUfxQGzpeDip0FtuJGz8zeu0fqaSg==
)

func SpnegoRequest() {
	// Load the client krb5 config
	cfg, err := config.Load("./krb5.conf")
	// Create the client with the pwd
	cl := client.NewWithPassword("liangzx/admin", "LIANGZXGG.COM", "dingjia", cfg, client.DisablePAFXFAST(true))
	// Create the client with keytab
	/*
		kt, err := keytab.Load("./liangzx_admin.keytab")
		if err != nil {
			log.Fatal(err)
		}
		cl := client.NewWithKeytab("liangzx/admin", "LIANGZXGG.COM", kt, cfg, client.DisablePAFXFAST(true))
	*/
	// Log in the client
	err = cl.Login()
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	log.Println("Kerberos 登录成功！")

	spnegoCl := spnego.NewClient(cl, nil, spn)

	// Make the request
	resp, err := spnegoCl.Get(url)
	if err != nil {
		log.Fatalf("error making request: %v", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v", err)
	}

	fmt.Println(string(b))
}

func HttpRequest() {
	cfg, err := config.Load("./krb5.conf")
	cl := client.NewWithPassword("liangzx/admin", "LIANGZXGG.COM", "dingjia", cfg, client.DisablePAFXFAST(true))
	err = cl.Login()
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	log.Println("Kerberos 登录成功！")
	req, err := http.NewRequest("GET", url, nil)
	ticket, key, err := cl.GetServiceTicket(spn)
	_ = key
	ticketByte, _ := ticket.Marshal()
	negotiate := base64.StdEncoding.EncodeToString(ticketByte)
	log.Println("ticketByte", negotiate)
	// "Negotiate YIICbAYGKwYBBQUCoIICYDCCAlygDTALBgkqhkiG9xIBAgKiggJJBIICRWCCAkEGCSqGSIb3EgECAgEAboICMDCCAiygAwIBBaEDAgEOogcDBQAAAAAAo4IBY2GCAV8wggFboAMCAQWhDxsNTElBTkdaWEdHLkNPTaIgMB6gAwIBAaEXMBUbBEhUVFAbDWxpYW5nenhnZy5jb22jggEfMIIBG6ADAgESoQMCAQKiggENBIIBCYe9kY0z5/Re8VPFfcM5aL3dC3eqxxqPX6UVSmB9A2bZi5fYKY/ufhU28XShdFLW3PQX5pMRXQL+KxYhInbJl940BK5r+qyLtSxwnbdl7UqTZzwoer1EpRKiP/j1mPD7f1disIGHSA4AnvBXUrzTRGeWz2Thsf1qja1YUAa0VgbhLd0WxTXIl1o6slicCOFd2S2+ka1PAnspysqEUjCaY3y3QT7SHETitDdVt2n07YqandtSZxl+2pHn5IYz7ko/InPis4zSAIAm3D7VWRiD2SNRUaZ4c1q43+Q7bKMh5i6R+Gj1EetkgpqKXfwhh3IMuANPt2riH10W/ckJH0fmpnPiy/usZnws1pSkga8wgaygAwIBEqEDAgECooGfBIGciaV29ZEWV+z1NJKHS6XdBZw6GCk0ECdrppy44H6NDcPpxuqHIZdjFe2olm5fvZmqpn5Y4UE67fVgbb7oDFIyleiD0hTrO1Vrh5CmWnUOouEkgH9aCcnk+J/GB00/wOOOLuJQkzCjk1N5aAg6V/th41eNqoVyJ6UzCAJzttKUUNmcildXi0KI0Yr0A8ZE1Y+cgZ32e3/XrfrUbnb7"
	req.Header.Add("Authorization", "Negotiate "+negotiate)
	httpCli := &http.Client{}
	// 发送请求
	resp, err := httpCli.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	fmt.Println(string(body))
}
