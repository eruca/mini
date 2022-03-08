package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

const token = "nbk2022"

func main() {
	router := gin.Default()

	router.GET("/", wxCheckSignature)
	router.POST("/", WXMsgReceive)

	err := CreateMenu()
	if err != nil {
		panic(err)
	}

	log.Fatalln(router.Run(":8080"))
}

func wxCheckSignature(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	if checkSignature(signature, timestamp, nonce, token) {
		log.Println("check signature successed")
		c.String(http.StatusOK, "%s", echostr)
		return
	}
	log.Println("un successed")
}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
}

// CheckSignature 微信公众号签名检查
func checkSignature(signature, timestamp, nonce, token string) bool {
	h := sha1.New()

	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)

	for i := 0; i < len(sl); i++ {
		h.Write([]byte(sl[i]))
	}

	return hex.EncodeToString(h.Sum(nil)) == signature
}

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}

	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)

	// 对接收的消息进行被动回复
	WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName)
}

// WXRepTextMsg 微信回复文本消息结构体
type WXRepTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "谢谢你的支持",
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	c.Writer.Write(msg)
}
