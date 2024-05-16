package engine

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/pkg/engine/engine_request"
	"github.com/Ulqiora/Route256Project/pkg/engine/engine_response"
	"github.com/Ulqiora/Route256Project/pkg/logger"
)

type SyncDriver struct {
	logObj  logger.LogEngine
	headers map[engine_request.Path]FunctionHandler
}

func NewSyncDriver(obj logger.LogEngine) SyncDriver {
	return SyncDriver{
		logObj:  obj,
		headers: make(map[engine_request.Path]FunctionHandler),
	}
}

func (d *SyncDriver) SetFunctionHandler(name engine_request.Path, handler FunctionHandler) {
	d.headers[name] = handler
}

func (d *SyncDriver) StartDriver(_ context.Context, _ ...context.CancelFunc) {
	return
}

func (d *SyncDriver) CallFunc(name engine_request.Path, r engine_request.RequestEngine) (Response, error) {
	if _, ok := d.headers[name]; ok {
		return nil, fmt.Errorf("duplicate function handler for %s", name)
	}
	var w engine_response.ResponseRWType
	d.headers[name](&w, r)
	return &w, nil
}
