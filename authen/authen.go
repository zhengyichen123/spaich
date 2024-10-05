//鉴权连接
package authen

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"math/rand"
	"time"
	"strconv"
	"github.com/gorilla/websocket"
	"github.com/zhengyichen/spaich/param"
)

//生成会话id
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
   b := make([]byte, length)
   for i := range b {
      b[i] = charset[seededRand.Intn(len(charset))]
   }
   return string(b)
}

func generateMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func generateHmacSHA1(text, secret string) string {  
	h := hmac.New(sha1.New, []byte(secret))  
	h.Write([]byte(text))  
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func GenerateSig() (int64,string) {
	timestamp := time.Now().UnixMilli() //时间戳

	//用appid和timestamp通过md5生成随机授权参数auth
	auth := generateMD5(param.AppId + fmt.Sprintf("%d", timestamp))
	//用auth和apisecret通过hmacSHA1生成签名signature
	signature := generateHmacSHA1(auth, param.ApiSecret)
	//fmt.Printf("signature: %s\n", signature)
	
	return timestamp, signature
}

func GenerateLoginUrl(hostUrl string) string {
	
	timestamp,signature := GenerateSig()
	param.ChatId = generateRandomString(16)
	param.PrechatId = param.ChatId

	/* 打印结果  
	fmt.Printf("Auth: %s\n", auth)  
	fmt.Printf("Signature: %s\n", signature) */
	return fmt.Sprintf(hostUrl + "/{%s}?appId=%s&timestamp=%d&signature=%s", param.ChatId, param.AppId, timestamp, signature)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("read body error: %s", err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
func WebSocketConn(hostUrl string) *websocket.Conn {
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立连接
	conn, resp, err := d.Dial(GenerateLoginUrl(hostUrl), nil)

	if err != nil {
		fmt.Println(readResp(resp) + err.Error())
		return nil
	}else if(resp.StatusCode != 101){
		fmt.Println(readResp(resp))
		return nil
	}
	
	fmt.Println("success connected !!!")

	return conn
}

func HttpConn(url string, rtype string, data interface{}) ([]byte, error){
	
	var req *http.Request
	var err error

	if reader, ok := data.(io.Reader); ok {  
		req, err = http.NewRequest(rtype, url, reader)   
	} else {  
		req, err = http.NewRequest(rtype, url, nil)  
	}
	
	if err != nil {  
		fmt.Printf("error creating request: %v\n", err)
		return nil, err
	} 
	
	//添加请求头
	timestamp, signature := GenerateSig()
	req.Header.Set("appId", param.AppId)
	req.Header.Set("timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("signature", signature)
	// 创建HTTP客户端并发送请求  
	client := &http.Client{}  
	resp, err := client.Do(req)  
	if err != nil { 
		fmt.Printf("error sending request: %v\n", err) 
		return nil, err  
	}  
	defer resp.Body.Close()
		
	// 读取响应体  
	body, err := io.ReadAll(resp.Body)  
	if err != nil { 
		fmt.Printf("error reading response body: %v\n", err) 
		return nil, err 
	} 
	//fmt.Printf("%s\n", string(body))

	return body, nil
}
