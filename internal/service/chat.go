package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/internal/biz"
)

type ChatUseCase struct {
	sh  *biz.ChatUseCase
	log *log.Helper
}

func NewChatUseCase(sh *biz.ChatUseCase) *ChatUseCase {
	return &ChatUseCase{sh: sh}
}

// 处理WebSocket跨域
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var wc = make(map[string]*websocket.Conn)

func (cc *ChatUseCase) Send(c *gin.Context, req *biz.ChatReq) (*biz.ChatRes, error) {
	// 升级为WebSocket协议 s
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	cc.log.WithContext(c).Infof("Send: %v", req)
	// return &biz.ChatRes(ctx, req), nil
	result, err := cc.sh.Send(c, conn, wc, req)
	return &biz.ChatRes{
		StatusCode: 0,
		StatusMsg:  result.StatusMsg,
	}, nil
}
