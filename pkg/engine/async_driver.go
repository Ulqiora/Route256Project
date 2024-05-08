package engine

import (
	"context"
	"errors"

	"homework/pkg/engine/engine_request"
	"homework/pkg/engine/engine_response"
	"homework/pkg/logger"
)

type FunctionHandler func(w engine_response.ResponseEngine, r engine_request.RequestEngine)

type AsyncDriver struct {
	readerRequest  engine_request.Reader
	writerResponse engine_response.Writer
	logObj         logger.LogEngine
	headers        map[engine_request.Path]FunctionHandler
}

func NewAsyncDriver(reader engine_request.Reader,
	writer engine_response.Writer,
	obj logger.LogEngine) AsyncDriver {
	return AsyncDriver{
		readerRequest:  reader,
		writerResponse: writer,
		logObj:         obj,
		headers:        make(map[engine_request.Path]FunctionHandler),
	}
}

func (d *AsyncDriver) SetFunctionHandler(name engine_request.Path, handler FunctionHandler) {
	d.headers[name] = handler
}

func (d *AsyncDriver) StartDriver(ctx context.Context, cancels ...context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			for _, cancel := range cancels {
				cancel()
			}
			d.logObj.Info("Engine", "the work of the parser is completed")
			return
		case request := <-d.readerRequest.Read():
			var response engine_response.ResponseRWType
			handler, ok := d.headers[request.GetPath()]
			if !ok {
				response.SetStatusCode(engine_response.StatusFail, errors.New("incorrect path"))
				d.writerResponse.Write(&response)
			}
			handler(&response, request)
			d.writerResponse.Write(&response)
		}
	}
}
