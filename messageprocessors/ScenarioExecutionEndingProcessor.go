package messageprocessors

import (
	m "github.com/ezbuy/gauge-go/gauge_messages"
	t "github.com/ezbuy/gauge-go/testsuit"
)

type ScenarioExecutionEndingProcessor struct{}

func (r *ScenarioExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetScenarioExecutionEndingRequest().GetCurrentExecutionInfo().GetCurrentScenario().GetTags()
	hooks := context.GetHooks(t.AFTERSCENARIO, tags)
	return executeHooks(hooks, msg)
}
