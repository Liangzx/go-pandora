[进程间如何传递fd](https://mp.weixin.qq.com/s/ts_qsUee6mUFv0FpykaOXQ)
[进程间传递文件描述符fd](https://www.cnblogs.com/aquester/p/9891633.html)

```text
 进程A监听9087端口对外提供服务
 进程B监听/tmp/sock.sock文件，接受其他进程（这里只有进程A）传递过来已经listen文件描述符，继续进行accept处理
 进程A监听Interrupt信号，收到信号将监听9087端口的文件描述符传递给B进程，然后退出
 最后B进程监听9087端口对外提供服务
```