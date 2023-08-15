package interfaces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/internal/biz"
	"ws/internal/service"
)

type ChatUseCase struct {
	chatService *service.ChatUseCase
	log         *log.Helper
}

func NewChatUseCase(ch *service.ChatUseCase, logger log.Logger) *ChatUseCase {
	return &ChatUseCase{
		chatService: ch,
		log:         log.NewHelper(logger),
	}
}

type Req struct {
	ActionType string `json:"action_type"` // 1-发送消息
	Content    string `json:"content"`     // 消息内容
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

// 处理WebSocket跨域
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wc = make(map[string]*websocket.Conn)

func (cc *ChatUseCase) Send(c *gin.Context) {
	action_type := c.Query("action_type")
	content := c.Query("content")
	to_user_id := c.Query("to_user_id")
	token := c.Query("token")
	Message := &Req{
		ActionType: action_type,
		Content:    content,
		ToUserID:   to_user_id,
		Token:      token,
	}

	fmt.Printf("Message: %v\n", Message)
	fmt.Printf("wc: %v\n", wc)
	fmt.Printf("upgrade: %v\n", upgrade)
	// 升级为WebSocket协议 s
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 400,
			"status_msg":  fmt.Sprintf("升级为WebSocket协议失败: %v", err),
		})
		return
	}
	Chat(c, conn, wc, Message)
}

func Chat(c *gin.Context, conn *websocket.Conn, ws map[string]*websocket.Conn, msg *Req) {
	for {
		// 创建消息JSON结构体, 保存消息与额外的信息
		message := new(biz.ChatReq)
		// ReadJSON
		// message: 需要读取的消息对象(一条消息一般包含多个属性用于其他用途)
		err := conn.ReadJSON(message)
		if err != nil {
			panic(err)
		}

		// 消息结构体
		msg := &biz.ChatReq{
			ActionType: msg.ActionType,
			Content:    msg.Content,
			ToUserID:   msg.ToUserID,
			Token:      msg.Token,
		}

		fmt.Printf("消息内容: %v", msg)

		// 根据连接绑定用户id
		ws[msg.ToUserID] = conn

		for _, cc := range ws {
			// WriteMessage
			// 1 消息类型: websocket.TextMessage文本
			// 2 传输类型: []byte二进制
			if cc.WriteJSON(msg) != nil {
				break
			}
		}
	}
}
