package player

import (
	"fmt"
	"encoding/json"
	"bytes"
	"github.com/zhengyichen/spaich/authen"
	"github.com/zhengyichen/spaich/param"
)

type Player struct {
	appId          string
	playerId       string
	playerName     string
	playerType     string
	description    string
	senderIdentity string
}

func NewPlayer( playerName string, playerType string, description string, senderIdentity string) *Player {
	pr := &Player{
		appId:          param.AppId,
		playerId:       "",
		playerName:     playerName,
		playerType:     playerType,
		description:    description,
		senderIdentity: senderIdentity,
	}
	
	res := pr.IfRegister()
	if res {
		pr.register()
		return pr
	}
	return nil
}

//判断玩家是否注册, true为未注册, false为已注册或发生错误
func (pr *Player) IfRegister() bool {
	
	url := "https://ai-character.xfyun.cn/api/open/player/if-register"
	url = url + "?appId=" + pr.appId + "&playerName=" + pr.playerName
	body, err := authen.HttpConn(url, "GET", nil)
	if err != nil {
		fmt.Printf("Error in player IfRegister:%v\n" , err)
		return false
	}

	//解析数据
	var response param.Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON in player IfRegister:%v\n", err)
		return false
	}
	if !response.Data.(bool) {
		fmt.Println("The player is not registered")
		return true
	}else{
		fmt.Println("The player is registered")
		return false
	}
}

func (pr *Player) register() bool {
	player := param.PlayerRegister{  
		AppId:        param.AppId,
		PlayerName:   pr.playerName,
		PlayerType:   pr.playerType,
		Description:  pr.description,
		SenderIdentity: pr.senderIdentity,
	}
	
	jsonData, err := json.Marshal(player)  
	if err != nil {  
		fmt.Printf("error marshaling JSON in player Register: %v\n", err)
		return false
	} 
	reader := bytes.NewReader(jsonData)
	url := "https://ai-character.xfyun.cn/api/open/player/register"
	
	body, err := authen.HttpConn(url, "POST", reader)
	if err != nil {  
		fmt.Printf("Error in player Register:%v\n", err)  
		return false
	}
	
	//解析数据
	var response param.Resp  
	err = json.Unmarshal(body, &response)  
	if err != nil {  
		fmt.Printf("Error parsing JSON in player Register:%v\n", err)  
		return false
	} 
	
	if !response.Success {
		fmt.Printf("Error in player register:%v\n", response.Message)
		return false
	}else{
		fmt.Println("player register sucess")
		pr.playerId = response.Data.(string)
		param.PlayerId = pr.playerId //只设置一个玩家
		return true	
	}
}

func (pr *Player) modify() bool {
	player := param.PlayerModify{  
		AppId:        param.AppId,
		PlayerId:     pr.playerId,
		PlayerName:   pr.playerName,
		PlayerType:   pr.playerType,
		Description:  pr.description,
		SenderIdentity: pr.senderIdentity,
	}
	jsonData, err := json.Marshal(player)
	if err != nil {  
		fmt.Printf("error marshaling JSON in player Modify: %v\n", err)
		return false
	}
	reader := bytes.NewReader(jsonData)
	url := "https://ai-character.xfyun.cn/api/open/player/modify"
	body, err := authen.HttpConn(url, "POST", reader)
	if err != nil {  
		fmt.Printf("Error in modify:%v\n", err)  
		return false
	}
	var response param.Resp
	err = json.Unmarshal(body, &response)
	if err != nil {  
		fmt.Printf("Error parsing JSON in player Modify:%v\n", err)  
		return false
	}
	if response.Data.(bool) {
		fmt.Println("player modify sucess")
		return true
	}else{
		fmt.Printf("Error in player modify:%v\n", response.Message)
		return false
	}
}

func (pr *Player) delete() bool {
	player := param.PlayerDelete{  
		AppId:        param.AppId,
		PlayerId:     pr.playerId,
		PlayerName:   pr.playerName,
	}
	
	jsonData, err := json.Marshal(player)  
	if err != nil {  
		fmt.Printf("error marshaling JSON in player Delete: %v\n", err)
	}
	reader := bytes.NewReader(jsonData)
	url := "https://ai-character.xfyun.cn/api/open/player/delete"
	body, err := authen.HttpConn(url, "POST", reader)
	if err != nil {
		fmt.Printf("Error in player Delete:%v\n", err)
	}
	var response param.Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON in player Delete:%v\n", err)
	}
	if response.Data.(bool) {
		fmt.Println("player delete sucess")
		return true
	}else{
		fmt.Printf("Error in player delete:%v\n", response.Message)
		return false
	}
}

func (pr *Player)GetId() string {
	return pr.playerId
}

func (pr *Player)GetName() string {
	return pr.playerName
}
