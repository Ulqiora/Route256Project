package printer

import (
	"context"
	"fmt"
	"io"

	"homework/pkg/engine"
	"homework/pkg/engine/engine_response"
)

type SyncPrinter struct {
	output io.Writer
}

func NewSyncPrinter(file io.Writer) SyncPrinter {
	return SyncPrinter{
		output: file,
	}
}

func (p *SyncPrinter) Start(_ context.Context, _ ...context.CancelFunc) {
	return
}

func (p *SyncPrinter) Print(response engine.Response) {
	if status, err := response.StatusCode(); status == engine_response.StatusFail && err != nil {
		_, _ = p.output.Write([]byte(fmt.Sprintf("Fail: %s", err.Error())))
		return
	}
	bytesAll, _ := io.ReadAll(response.GetBody())
	_, _ = p.output.Write(bytesAll)
}
