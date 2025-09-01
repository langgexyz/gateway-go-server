package api

import (
	"context"

	"github.com/xpwu/go-log/log"
)

type PingRequest struct {
}

type PingResponse struct {
}

func (s *Suite) APIPing(ctx context.Context, request *PingRequest) *PingResponse {
	_, logger := log.WithCtx(ctx)
	logger.PushPrefix("Ping")
	logger.Debug("start")
	defer logger.Debug("end")

	return &PingResponse{}
}
