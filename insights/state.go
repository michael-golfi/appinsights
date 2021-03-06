package insights

import (
	"sync"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
)

func (l *insightsLogger) worker() {
	timer := time.NewTicker(l.postMessagesFrequency)
	var messages []*contracts.Envelope
	for {
		select {
		case message, open := <-l.stream:
			if !open {
				l.postMessages(messages, true)
				l.lock.Lock()

				l.transport.CloseIdleConnections()
				l.closed = true
				l.closedCond.Signal()

				l.lock.Unlock()
				return
			}
			messages = append(messages, message)
			// Only sending when we get exactly to the batch size,
			// This also helps not to fire postMessages on every new message,
			// when previous try failed.
			if len(messages)%l.postMessagesBatchSize == 0 {
				messages = l.postMessages(messages, false)
			}
		case <-timer.C:
			messages = l.postMessages(messages, false)
		}
	}
}

func (l *insightsLogger) Close() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.closedCond == nil {
		l.closedCond = sync.NewCond(&l.lock)
		close(l.stream)
		for !l.closed {
			l.closedCond.Wait()
		}
	}
	return nil
}
