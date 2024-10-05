package contact

import (
	"fmt"
	"bytes"
	"encoding/json"
	"time"
	"github.com/gorilla/websocket"
	"github.com/zhengyichen/spaich/param"
	"github.com/zhengyichen/spaich/authen"
	"github.com/zhengyichen/spaich/player"
	"github.com/zhengyichen/spaich/character"
)

//角色视角下的剧情场景、背景知识、对话人身份信息
func MemoryCreate(pr *player.Player, ch *character.Character) bool {
	memory := param.MemoryCreate{
		AppId: param.AppId,
		AgentId: ch.GetId(),
		InteractionType: "obs", //"cnv"
		Description: " ",
		PlayerInvolved: pr.GetId(),
		AgentInvolved: ch.GetId(),
	}
	jsonData, err := json.Marshal(memory)
	if err != nil {  
		fmt.Printf("error marshaling JSON in interactive memorycreate: %v\n", err)
		return false
	}
	reader := bytes.NewReader(jsonData)
	url := "https://ai-character.xfyun.cn/api/open/interactive/generate"
	body, err := authen.HttpConn(url, "POST", reader)
	if err != nil {  
		fmt.Printf("Error in interactive memorycreate:%v\n", err)  
		return false
	}
	var response param.RespMemory
	err = json.Unmarshal(body, &response)
	if err != nil {  
		fmt.Printf("Error parsing JSON in interactive memorycreate:%v\n", err)  
		return false
	}
	if !response.Success {
		fmt.Printf("Error in interactive memorycreate:%s\n",response.Message)
		return false
	}
	fmt.Println("memory create success!")
	return true
}

func MemoryDelete() bool{
	url := "https://ai-character.xfyun.cn/api/open/interactive/clear-cache"
	url = url + "?appId=" + param.AppId + "&chatId=" + param.ChatId 
	body, err := authen.HttpConn(url, "POST", nil)
	if err != nil {  
		fmt.Printf("Error in memory delete:%v\n", err)  
		return false
	}
	var response param.RespMemory
	err = json.Unmarshal(body, &response)
	if err != nil {  
		fmt.Printf("Error parsing JSON in memory delete:%v\n", err)  
		return false
	}
	if !response.Success {
		fmt.Printf("Error in memory delete:%s\n",response.Message)
		return false
	}
	fmt.Println("memory delete success!")
	return true
}

func Contact(pr *player.Player, ch *character.Character, conn *websocket.Conn, ty int, question string) {
	// 创建一个中断通道，用于优雅地关闭 WebSocket 连接  
	// interrupt := make(chan os.Signal, 1)  
	// signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM) 


	// 创建 done 通道，用于标记读取 goroutine 的完成  
	// 创建 answerchan 通道，用于接收答案
	//done := make(chan struct{})
	answerchan := make(chan string)
	// 启动一个 goroutine 来读取 WebSocket 服务器的消息  
	go func() {  
		//defer close(done)  
		answer := ""
		for {  
			// 设置读取超时，以避免 goroutine 永远阻塞  
			conn.SetReadDeadline(time.Now().Add(time.Second * 120))  
			_, message, err := conn.ReadMessage()  
			if err != nil {  
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {  
					fmt.Printf("Connection closed by server:%v\n", err)  
					return  
				}  
				fmt.Printf("Error reading message in websocket:%v\n", err)  
				return  
			}  
			// 处理接收到的消息
			//fmt.Println("Received message:", string(message))
			var response map[string]interface{}
			err = json.Unmarshal(message, &response)
			if err != nil {
				fmt.Printf("Error parsing JSON in contact:%v\n", err)
				return
			}
			payload := response["payload"].(map[string]interface{})
			choices := payload["choices"].(map[string]interface{})
			header := response["header"].(map[string]interface{})
			code := header["code"].(float64)

			if code != 0 {
				fmt.Println(response["payload"])
				return
			}
			status := choices["status"].(float64)
			//fmt.Println(status)
			text := choices["text"].([]interface{})
			content := text[0].(map[string]interface{})["content"].(string)
			if status != 2 {
				answer += content
			} else {
				answer += content
				answerchan <- answer
				break;
			}
		}  
	}() 

	//创建请求消息
	var texts interface{}
	if ty == 1 {
		texts = []param.Text{
			{Role: pr.GetName(), Content: question},
		}	
	}else{
		texts = []string{}
	}

	req := map[string]interface{}{ 
		"header": map[string]interface{}{ 
			"app_id": param.AppId, 
			"uid": pr.GetId(),
			"agent_id": ch.GetId(),
		},
		"parameter": map[string]interface{}{ 
			"chat": map[string]interface{}{ 
				"chat_id": param.ChatId,
			},
		},
		"payload": map[string]interface{}{ 
			"message": map[string]interface{}{ 
				"text": texts, 
			},
		},
	}  
	reqJSON, err := json.Marshal(req)  
	if err != nil {  
		fmt.Printf("Failed to marshal request to JSON in interactive:%v\n", err)  
	}
	// 发送 JSON 消息到 WebSocket 服务器  
	err = conn.WriteMessage(websocket.TextMessage, reqJSON)  
	if err != nil {  
		fmt.Printf("Failed to send message to WebSocket server:%v\n", err)  
	} 
		
	answer := <-answerchan
	fmt.Println(answer)
	
	// select {  
	// case <-interrupt:  
	// 	fmt.Println("Interrupt signal received, closing connection...") 
	// 	conn.Close()
	// 	os.Exit(0) 
	// // case <-done:  
	// // 	fmt.Println("Reading goroutine finished, exiting...")  
	// }  
}



