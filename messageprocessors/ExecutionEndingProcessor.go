package messageprocessors

import (
	m "github.com/ezbuy/gauge-go/gauge_messages"
	t "github.com/ezbuy/gauge-go/testsuit"
)

type ExecutionEndingProcessor struct{}

func (r *ExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetExecutionEndingRequest().GetCurrentExecutionInfo().GetCurrentScenario().GetTags()
	hooks := context.GetHooks(t.AFTERSUITE, tags)

	return executeHooks(hooks, msg)
}
