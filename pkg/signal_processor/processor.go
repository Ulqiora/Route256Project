package signal_processor

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type Processor struct {
	ctxWithCancel context.Context
	syncObject    atomic.Bool
}

func New(signals ...os.Signal) (*Processor, context.Context) {
	parentContext := context.Background()
	ctx, _ := signal.NotifyContext(parentContext, signals...)
	processor := &Processor{
		ctxWithCancel: ctx,
	}
	return processor, ctx
}

func (p *Processor) IsCanceled() bool {
	return p.syncObject.Load()
}

func (p *Processor) Start(canFunc ...context.CancelFunc) {
	for {
		select {
		case <-p.ctxWithCancel.Done():
			fmt.Println("Программа завершается, я завершаю работу...")
			p.syncObject.Store(true)
			time.Sleep(500 * time.Millisecond)
			for i := range canFunc {
				canFunc[i]()
			}
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}

	}
}
