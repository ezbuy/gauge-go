package messageprocessors

import (
	m "github.com/ezbuy/gauge-go/gauge_messages"
	t "github.com/ezbuy/gauge-go/testsuit"
)

type SpecExecutionStartingRequestProcessor struct{}

func (r *SpecExecutionStartingRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetSpecExecutionStartingRequest().GetCurrentExecutionInfo().GetCurrentSpec().GetTags()
	hooks := context.GetHooks(t.BEFORESPEC, tags)

	return executeHooks(hooks, msg)
}
