package proxy

import (
	"context"
	"fmt"
	"gateway-go-server-main/src/utils"
	"net/http"

	"github.com/xpwu/go-httpclient/httpc"
	"github.com/xpwu/go-log/log"
)

func (s *Suite) APIProxy(ctx context.Context, request *string) *string {
	_, logger := log.WithCtx(ctx)
	logger.PushPrefix("Proxy")
	logger.Debug("start")
	logger.Debug(fmt.Sprintf("req: %+v", request))
	defer logger.Debug("end")

	url := s.Request.Header.Get("x-proxy-url")
	if url == "" {
		s.Request.Terminate(fmt.Errorf("x-proxy-url is empty"))
	}

	method := s.Request.Header.Get("x-proxy-method")
	if method == "" {
		s.Request.Terminate(fmt.Errorf("x-proxy-method is empty"))
	}

	pushUrl := s.Request.Header.Get("Pushurl")
	if pushUrl == "" {
		s.Request.Terminate(fmt.Errorf("pushurl is empty"))
	}

	pushToken := s.Request.Header.Get("Pushtoken")
	if pushToken == "" {
		s.Request.Terminate(fmt.Errorf("pushtoken is empty"))
	}

	logger.Debug(fmt.Sprintf("pushUrl: %s", pushUrl))

	var res []byte
	var resHeader http.Header
	var err error

	// 根据 HTTP 方法决定是否发送 body
	if method == "GET" || method == "HEAD" || method == "DELETE" {
		// GET/HEAD/DELETE 请求通常不发送 body
		err = httpc.Send(ctx, url, httpc.WithHeader(s.Request.Header),
			httpc.WithMethod(method),
			httpc.WithBytesResponse(&res),
			httpc.WithResponseHeader(&resHeader),
		)
	} else {
		// POST/PUT/PATCH 等请求发送 body
		err = httpc.Send(ctx, url, httpc.WithHeader(s.Request.Header),
			httpc.WithMethod(method),
			httpc.WithBytesBody([]byte(*request)),
			httpc.WithBytesResponse(&res),
			httpc.WithResponseHeader(&resHeader),
		)
	}
	if err != nil {
		s.Request.Terminate(fmt.Errorf("proxy error: %+v", err))
	}

	response, err := utils.DecompressResponse(res, resHeader, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("DecompressResponse error: %+v", err))
		s.Request.Terminate(err)
		return nil
	}

	logger.Debug(fmt.Sprintf("response: %s", response))
	return &response
}
