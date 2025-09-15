package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway-go-server-main/src/utils"
	"net/http"

	"github.com/xpwu/go-httpclient/httpc"
	"github.com/xpwu/go-log/log"
	"github.com/xpwu/go-tinyserver/api"
)

type Suite struct {
	Request *api.Request
}

func (s *Suite) SetUp(ctx context.Context, r *api.Request, apiReq interface{}) bool {
	_, logger := log.WithCtx(ctx)

	if err := json.Unmarshal(r.RawData, apiReq); err != nil {
		logger.Error(err)
		r.Terminate(err)
	}

	s.Request = r

	url := s.Request.Header.Get("x-hook-url")
	method := s.Request.Header.Get("x-hook-method")

	if url != "" && method != "" {
		reqId := s.Request.Header.Get("X-Req-Id")
		api := s.Request.Header.Get("Api")

		go s.performHookRequest(context.Background(), reqId, api, url, method)
	}

	return true
}

// performHookRequest 执行钩子请求
func (s *Suite) performHookRequest(ctx context.Context, reqId, api, url, method string) {
	ctx, logger := log.WithCtx(ctx)

	logger.PushPrefix(fmt.Sprintf("hook:%s, reqid:%s", api, reqId))
	logger.Debug("start")
	defer logger.Debug("end")

	var res []byte
	var resHeader http.Header
	var err error

	if method == "GET" || method == "HEAD" || method == "DELETE" {
		err = httpc.Send(ctx, url, httpc.WithHeader(s.Request.Header),
			httpc.WithMethod(method),
			httpc.WithBytesResponse(&res),
			httpc.WithResponseHeader(&resHeader),
		)
	} else {
		err = httpc.Send(ctx, url, httpc.WithHeader(s.Request.Header),
			httpc.WithMethod(method),
			httpc.WithBytesBody(s.Request.RawData),
			httpc.WithBytesResponse(&res),
			httpc.WithResponseHeader(&resHeader),
		)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("[%s] %s error %+v", method, url, err))
		return
	}

	response, err := utils.DecompressResponse(res, resHeader, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("DecompressResponse error: %+v", err))
		return
	}
	logger.Debug(fmt.Sprintf("response: %s", response))
}

func (s *Suite) TearDown(ctx context.Context, apiRes interface{}, res *api.Response) {
	var err error
	res.RawData, err = json.Marshal(apiRes)
	if err != nil {
		res.Request().Terminate(err)
	}
}

func (s *Suite) MappingPreUri() string {
	return "/"
}

func AddAPI() {
	api.Add(func() api.Suite {
		return &Suite{}
	})
}
