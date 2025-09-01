package proxy

import (
	"context"
	"fmt"

	"github.com/xpwu/go-tinyserver/api"
)

type Suite struct {
	Request *api.Request
}

func (s *Suite) SetUp(ctx context.Context, r *api.Request, apiReq interface{}) bool {
	if str, ok := apiReq.(*string); ok {
		*str = string(r.RawData)
	} else {
		s.Request.Terminate(fmt.Errorf("request is not *string"))
	}

	s.Request = r
	return true
}

func (s *Suite) TearDown(ctx context.Context, apiRes interface{}, res *api.Response) {
	if str, ok := apiRes.(*string); ok {
		res.RawData = []byte(*str)
	} else {
		s.Request.Terminate(fmt.Errorf("response is not string"))
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
