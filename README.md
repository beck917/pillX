# pillX
a simple &amp; powerful network Library written in Go

#目录结构
  gateway 网关服务器
  worker 逻辑服务器
  pillx   pillx网络通信库核心代码
  
#运行原理
  gateway 用于维护客户端连接,并转发客户端的数据给Worker进程处理
  worker 用于处理gateway转发过来的请求,并处理相应的业务逻辑,然后返回结果给gateway
  
#协议
  gateway协议 gateway和worker通信所用的协议
  pill协议 目前pillX网络开发套件,官方推荐使用的二进制协议,特点是间接高效
  text协议 文本协议,主要用于测试
  websocket协议 待开发,用于html5页面通信的协议
  
使用实例
  pillx参考了go官方net库中的http的使用方式
  
	server := &pillx.Server{
		Addr:          ":8080",
		Handler:        nil,
		Protocol:		new(pillx.PillProtocol),
	}
	server.HandleFunc(0x0DDC, helloHandler)
	server.HandleFunc(pillx.SYS_ON_CLOSE, closeHandler)
	fmt.Println("pillX服务端引擎启动")
	server.ListenAndServe()
