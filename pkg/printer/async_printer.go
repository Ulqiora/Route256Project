package printer

import (
	"context"
	"fmt"
	"io"

	"homework/pkg/engine"
	"homework/pkg/engine/engine_response"
)

type AsyncPrinter struct {
	responseReader *engine.ResponseReader
	reader         engine_response.Reader
	output         io.Writer
}

func New(reader engine_response.Reader,
	file io.Writer) AsyncPrinter {
	return AsyncPrinter{
		responseReader: engine.NewResponseReader(reader),
		reader:         reader,
		output:         file,
	}
}

func (p *AsyncPrinter) Start(ctx context.Context, cancelFuncs ...context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Printer is ended")
			for _, cancel := range cancelFuncs {
				cancel()
			}
			return
		default:
			response := p.responseReader.Read()
			if status, err := response.StatusCode(); status == engine_response.StatusFail && err != nil {
				_, _ = p.output.Write([]byte(fmt.Sprintf("Fail: %s", err.Error())))
				continue
			}
			bytesAll, _ := io.ReadAll(response.GetBody())
			_, _ = p.output.Write(bytesAll)
		}
	}
}
