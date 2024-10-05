package param

var(
	AppId = "your appid"
    PlayerId = ""
	ApiSecret = "your apisecret"
	//apiKey = "your apikey"
    ChatId = ""
    PrechatId = ""
)

//Return parameters
type Resp struct {
    Code    int    `json:"code"`
    Data    interface{}  `json:"data"`
    Message string `json:"message"`
    Success bool   `json:"success"`
}

type RespMemory struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Success bool   `json:"success"`
}

// player parameter

type PlayerRegister struct
{
    AppId        	string   `json:"appId"`
    PlayerName   	string   `json:"playerName"`
    Description  	string   `json:"description"`
    PlayerType 		string   `json:"playerType"`
    SenderIdentity  string   `json:"senderIdentity"`
}

type PlayerModify struct
{
    AppId        	string   `json:"appId"`
    PlayerId    	string   `json:"playerId"`
    PlayerName   	string   `json:"playerName"`
    Description  	string   `json:"description"`
    PlayerType 		string   `json:"playerType"`
    SenderIdentity  string   `json:"senderIdentity"`
}

type PlayerDelete struct
{
    AppId        string   `json:"appId"`
    PlayerId    	string   `json:"playerId"`
    PlayerName   	string   `json:"playerName"`
}


//charactor parameter
type CharacterCreate struct
{
    AppId  string `json:"appId"`
    PlayerId string `json:"playerId"`
    AgentName string `json:"agentName"`
    AgentType string `json:"agentType"`
    Description string `json:"description"`
    Identity string `json:"identity"`
    PersonalityDescription string `json:"personalityDescription"`
    Hobby string `json:"hobby"`
   // Speaker string `json:"speaker"`
    KeyPersonality string `json:"keyPersonality"`
    Mission string `json:"mission"`
}

type CharacterEdit struct
{
    AppId  string `json:"appId"`
    PlayerId string `json:"playerId"`
    AgentId string `json:"agentId"`
    AgentName string `json:"agentName"`
    AgentType string `json:"agentType"`
    Description string `json:"description"`
    Identity string `json:"identity"`
    PersonalityDescription string `json:"personalityDescription"`
    Hobby string `json:"hobby"`
   // Speaker string `json:"speaker"`
    KeyPersonality string `json:"keyPersonality"`
    Mission string `json:"mission"`
}

//interactive parameter
type MemoryCreate struct
{
    AppId  string `json:"appId"`
    AgentId string `json:"agentId"`
    InteractionType string `json:"interactionType"`
    Description string `json:"description"`
    PlayerInvolved string `json:"playerInvolved"`
    AgentInvolved string `json:"agentInvolved"`
}

type Text struct {  
	Role string `json:"role"`
    Content string `json:"content"`  
}  
  



