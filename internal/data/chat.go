package data

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	"ws/internal/biz"
)

type chatRepo struct {
	data *Data
	log  *log.Helper
}

func (cc chatRepo) Send(c *gin.Context, conn *websocket.Conn, ws map[string]*websocket.Conn, req *biz.ChatReq) (*biz.ChatRes, error) {
	for {
		// 创建消息JSON结构体, 保存消息与额外的信息
		message := new(biz.ChatReq)
		// ReadJSON
		// message: 需要读取的消息对象(一条消息一般包含多个属性用于其他用途)
		err := conn.ReadJSON(message)
		if err != nil {
			return nil, err
		}

		// 消息结构体
		msg := &biz.ChatReq{
			ActionType: req.ActionType,
			Content:    req.Content,
			ToUserID:   req.ToUserID,
			Token:      req.Token,
		}

		// 根据连接绑定用户id
		ws[msg.Token] = conn

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

func NewUserRepo(data *Data, logger log.Logger) biz.ChatRepo {
	return &chatRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
