package api

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/xpwu/go-log/log"
)

type UnsubscribeRequest struct {
	Cmd []string `json:"cmd"`
}

type UnsubscribeResponse struct {
	ErrMsg string `json:"errMsg,omitempty"`
}

func (s *Suite) APIUnsubscribe(ctx context.Context, request *UnsubscribeRequest) *UnsubscribeResponse {
	_, logger := log.WithCtx(ctx)
	logger.PushPrefix("Unsubscribe")
	logger.Debug("start")
	logger.Debug(fmt.Sprintf("req: %+v", request))
	defer logger.Debug("end")

	if len(request.Cmd) == 0 {
		return &UnsubscribeResponse{ErrMsg: "cmd cannot be empty"}
	}

	pushToken := s.Request.Header.Get("Pushtoken")
	if pushToken == "" {
		s.Request.Terminate(fmt.Errorf("pushtoken is empty"))
	}

	// 检查是否有空的cmd
	if slices.Contains(request.Cmd, "") {
		return &UnsubscribeResponse{ErrMsg: "cmd cannot be empty"}
	}

	// 从指定的cmd列表中移除客户端
	for _, cmd := range request.Cmd {
		if clientsInterface, ok := channelClientsMap.Load(cmd); ok {
			clients := clientsInterface.(*sync.Map) // 该命令的客户端映射
			clients.Delete(pushToken)               // 从该命令移除指定客户端
			logger.Debug(fmt.Sprintf("Unsubscribe pushToken=%s from cmd=%s", pushToken, cmd))
		}
	}

	return &UnsubscribeResponse{}
}
