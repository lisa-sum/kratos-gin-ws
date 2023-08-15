package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

type ChatReq struct {
	ActionType string `json:"action_type"` // 1-发送消息
	Content    string `json:"content"`     // 消息内容
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

type ChatRes struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type Chat struct{}

type ChatRepo interface {
	Send(c *gin.Context, conn *websocket.Conn, ws map[string]*websocket.Conn, req *ChatReq) (*ChatRes, error)
}

type ChatUseCase struct {
	repo ChatRepo
	log  *log.Helper
}

func NewChatUseCase(repo ChatRepo, logger log.Logger) *ChatUseCase {
	return &ChatUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *ChatUseCase) Send(c *gin.Context, conn *websocket.Conn, ws map[string]*websocket.Conn, req *ChatReq) (*ChatRes, error) {
	result, err := cc.repo.Send(c, conn, ws, req)
	if err != nil {
		return nil, err
	}

	return &ChatRes{
		StatusCode: 0,
		StatusMsg:  result.StatusMsg,
	}, nil
}
