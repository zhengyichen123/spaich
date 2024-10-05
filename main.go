package main

import (
	"fmt"
	"os"
	"sync"
	"os/signal"
	"syscall"
	"github.com/zhengyichen/spaich/authen"
	"github.com/zhengyichen/spaich/character"
	"github.com/zhengyichen/spaich/interactive"
	"github.com/zhengyichen/spaich/player"
)

/*
该代码封装了星火api的角色模拟功能，详细介绍可以参考
https://www.xfyun.cn/doc/spark/character_simulation/%E8%A7%92%E8%89%B2%E6%A8%A1%E6%8B%9F%E4%BB%8B%E7%BB%8D.html#_0-%E5%BC%95%E8%A8%80
并借助星火大模型实现一个智能情感陪伴机器人的应用
*/


//如何优雅的关闭程序
var (  
	stop     = false       // 停止标志  
	stopMu   sync.Mutex    // 互斥锁，用于保护停止标志  
	interrupt = make(chan os.Signal, 1) // 信号通道  
)  

func init(){
	// 该ai没有记忆，由于上次的会话id不确定
	// 数据初始化...
	// to do
}

func main() {
	// 注册信号通知  
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)  
  
	// 启动一个 goroutine 来处理信号  
	go func() {  
		<-interrupt // 阻塞等待信号  
		stopMu.Lock()  
		stop = true        // 设置停止标志  
		stopMu.Unlock()  
		fmt.Println("Interrupt signal received, preparing to shut down...")  
	}()
	
	var question string
	url := "wss://ai-character.xfyun.cn/api/open/interactivews"
	conn := authen.WebSocketConn(url)
	//fmt.Printf("%s\n",readResp(resp))
	pr := player.NewPlayer("kop", "test13", "test13", "test13")
	first := 0
	for{
		if stop == true{
			break
		}
		if first == 0{
			contact.Contact(pr, character.GetInstance(), conn, 2, "")
			first = 1
		}else{
			fmt.Scanf("%s", &question)
			contact.Contact(pr, character.GetInstance(), conn, 1, question)	
		}
	}
	//pr.IfRegister()
	//demo.Test()
	conn.Close()
	os.Exit(0)
}


	
	//fmt.Println("Hello World!")

