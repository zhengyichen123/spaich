package character

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zhengyichen/spaich/authen"
	"github.com/zhengyichen/spaich/param"
)

type Character struct {
	appId                  string
	playerId               string
	agentId                string
	agentName              string
	agentType              string
	description            string
	identity               string
	personalityDescription string
	hobby                  string
	//speaker                string
	keyPersonality string
	mission        string
}

var (
	instance *Character
	once     sync.Once
)

func GetInstance() *Character {
	once.Do(func() {
		instance = &Character{
			appId:                  param.AppId,
			playerId:               param.PlayerId,
			agentId:                "",
			agentName:              "xx",
			agentType:              "xx",
			description:            "xx",
			identity:               "xx是女大学生",
			personalityDescription: "xx心思细腻，待人温柔体贴",
			hobby:                  "xx喜欢粘着朋友和朋友分享快乐，善于排解对方的负面情绪，喜欢音乐、电影、阅读、运动、舞蹈",
			//speaker: "",
			keyPersonality: "善良，热情，天真可爱，说话有点傻乎乎",
			mission:        "xx在微信上认识了男孩,发现他内心有点忧愁，希望给他一丝慰藉",
		}
		instance.create()
	})
	return instance
}

func NewCharacter(appId string, playerId string, agentName string, agentType string, description string, identity string, personalityDescription string, hobby string, speaker string, keyPersonality string, mission string) *Character {
	character := &Character{
		appId:                  appId,
		playerId:               playerId,
		agentId:                "",
		agentName:              agentName,
		agentType:              agentType,
		description:            description,
		identity:               identity,
		personalityDescription: personalityDescription,
		hobby:                  hobby,
		//speaker:                speaker,
		keyPersonality: keyPersonality,
		mission:        mission,
	}

	return character
}

func (ch *Character) create() bool {
	character := param.CharacterCreate{
		AppId:                  ch.appId,
		PlayerId:               ch.playerId,
		AgentName:              ch.agentName,
		AgentType:              ch.agentType,
		Description:            ch.description,
		Identity:               ch.identity,
		PersonalityDescription: ch.personalityDescription,
		Hobby:                  ch.hobby,
		//Speaker: ch.speaker,
		KeyPersonality: ch.keyPersonality,
		Mission:        ch.mission,
	}
	jsonData, err := json.Marshal(character)
	if err != nil {
		fmt.Printf("error marshaling JSON: %v\n", err)
		return false
	}
	reader := bytes.NewReader(jsonData)
	url := "https://ai-character.xfyun.cn/api/open/agent/edit-character"
	body, err := authen.HttpConn(url, "POST", reader)
	if err != nil {
		fmt.Printf("Error in character create:%v\n", err)
		return false
	}
	var response param.Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON in character create:%v\n", err)
		return false
	}
	if !response.Success {
		fmt.Printf("Error in character create:%v\n", response.Message)
		return false
	} else {
		fmt.Println("character create success")
		ch.agentId = response.Data.(string)
		return true
	}
}

func (ch *Character) edit() bool {
	character := param.CharacterEdit{
		AppId:                  ch.appId,
		PlayerId:               ch.playerId,
		AgentId:                ch.agentId,
		AgentName:              ch.agentName,
		AgentType:              ch.agentType,
		Description:            ch.description,
		Identity:               ch.identity,
		PersonalityDescription: ch.personalityDescription,
		Hobby:                  ch.hobby,
		//Speaker: ch.speaker,
		KeyPersonality: ch.keyPersonality,
		Mission:        ch.mission,
	}
	jsonData, err := json.Marshal(character)
	if err != nil {
		fmt.Printf("error marshaling JSON in character edit: %v\n", err)
		return false
	}
	reader := bytes.NewReader(jsonData)
	url := "https://ai-character.xfyun.cn/api/open/agent/edit-character"
	body, err := authen.HttpConn(url, "POST", reader)
	if err != nil {
		fmt.Printf("Error in character edit:%v\n", err)
		return false
	}
	var response param.Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON in character edit:%v\n", err)
		return false
	}
	if !response.Success {
		fmt.Printf("Error in character edit:%v\n", response.Message)
		return false
	} else {
		fmt.Println("character edit success")
		return true
	}
}

// other
func (ch *Character) GetId() string {
	return ch.agentId
}
