package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/xpwu/go-log/log"
	"github.com/xpwu/go-stream/push/core"
)

type PublishRequest struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type PublishResponse struct {
	ErrMsg string `json:"errMsg,omitempty"`
}

func (s *Suite) APIPublish(ctx context.Context, request *PublishRequest) *PublishResponse {
	_, logger := log.WithCtx(ctx)
	logger.PushPrefix("Publish")
	logger.Debug("start")
	logger.Debug(fmt.Sprintf("req: %+v", request))
	defer logger.Debug("end")

	if request.Cmd == "" {
		return &PublishResponse{ErrMsg: "cmd cannot be empty"}
	}

	// 获取指定cmd的已订阅客户端
	// 1. 从全局映射中获取该命令的客户端集合
	var pushTokens []string
	if clientsInterface, ok := channelClientsMap.Load(request.Cmd); ok {
		clients := clientsInterface.(*sync.Map) // 该命令的客户端映射
		// 2. 遍历该命令的所有客户端，收集pushToken
		clients.Range(func(key, value interface{}) bool {
			if token, ok := key.(string); ok {
				pushTokens = append(pushTokens, token)
			}
			return true
		})
	}

	if len(pushTokens) == 0 {
		// 没有订阅者是正常现象，不返回错误
		return &PublishResponse{}
	}

	logger.Debug(fmt.Sprintf("Found %d clients to publish to cmd '%s'", len(pushTokens), request.Cmd))

	// 获取请求header用于追踪
	reqId := s.Request.Header.Get("X-Req-Id")

	// 准备发送的数据，格式：{Cmd: string, Data: string, Header: {X-Req-Id: string}}
	dataBytes, err := json.Marshal(map[string]interface{}{
		"cmd":  request.Cmd,
		"data": request.Data,
		"header": map[string]string{
			"X-Req-Id": reqId,
		},
	})
	if err != nil {
		return &PublishResponse{ErrMsg: fmt.Sprintf("failed to marshal data: %v", err)}
	}
	dataToSend := string(dataBytes)

	// 异步推送数据到所有订阅的客户端
	for _, pushToken := range pushTokens {
		go func(token string) {
			// 直接使用核心推送，避免HTTP开销
			clientConn, st := core.GetClientConn(ctx, token)
			if st != core.Success {
				// 连接获取失败，从该命令移除失效的客户端
				if clientsInterface, ok := channelClientsMap.Load(request.Cmd); ok {
					clients := clientsInterface.(*sync.Map)
					clients.Delete(token)
				}
				return
			}

			// 推送数据到客户端
			st = core.PushDataToClient(ctx, clientConn, []byte(dataToSend))
			if st != core.Success {
				// 推送失败，从该命令移除失效的客户端
				if clientsInterface, ok := channelClientsMap.Load(request.Cmd); ok {
					clients := clientsInterface.(*sync.Map)
					clients.Delete(token)
				}
			}
		}(pushToken)
	}

	return &PublishResponse{}
}
