package api

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/xpwu/go-log/log"
)

type SubscribeRequest struct {
	Cmd []string `json:"cmd"`
}

type SubscribeResponse struct {
	ErrMsg string `json:"errMsg,omitempty"`
}

// 全局变量：按命令管理客户端订阅关系
// 数据结构：sync.Map[string, *sync.Map[string, bool]]
// - 外层key: cmd名称 (string)
// - 外层value: 该命令的客户端映射 (*sync.Map)
//   - 内层key: pushToken (string)
//   - 内层value: 固定为true (bool)，表示已订阅
//
// 示例：
//
//	channelClientsMap["news"] -> sync.Map{"token1": true, "token2": true}
//	channelClientsMap["chat"] -> sync.Map{"token3": true, "token4": true}
var channelClientsMap sync.Map

func (s *Suite) APISubscribe(ctx context.Context, request *SubscribeRequest) *SubscribeResponse {
	_, logger := log.WithCtx(ctx)
	logger.PushPrefix("Subscribe")
	logger.Debug("start")
	logger.Debug(fmt.Sprintf("req: %+v", request))
	defer logger.Debug("end")

	if len(request.Cmd) == 0 {
		return &SubscribeResponse{ErrMsg: "cmd cannot be empty"}
	}

	pushUrl := s.Request.Header.Get("Pushurl")
	if pushUrl == "" {
		s.Request.Terminate(fmt.Errorf("pushurl is empty"))
	}

	pushToken := s.Request.Header.Get("Pushtoken")
	if pushToken == "" {
		s.Request.Terminate(fmt.Errorf("pushtoken is empty"))
	}

	// 检查是否有空的cmd
	if slices.Contains(request.Cmd, "") {
		return &SubscribeResponse{ErrMsg: "cmd cannot be empty"}
	}

	// 批量订阅多个cmd
	for _, cmd := range request.Cmd {
		// 获取或创建该cmd的客户端映射
		// 1. 从全局映射中获取或创建该命令的客户端集合
		clientsInterface, _ := channelClientsMap.LoadOrStore(cmd, &sync.Map{})
		clients := clientsInterface.(*sync.Map) // 该命令的客户端映射

		// 2. 将客户端添加到该命令（value固定为true表示已订阅）
		clients.Store(pushToken, true)
		logger.Debug(fmt.Sprintf("Subscribe pushToken=%s to cmd=%s", pushToken, cmd))
	}

	return &SubscribeResponse{}
}
