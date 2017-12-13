package messageprocessors

import (
	m "github.com/ezbuy/gauge-go/gauge_messages"
	t "github.com/ezbuy/gauge-go/testsuit"
)

type MessageProcessor interface {
	Process(*m.Message, *t.GaugeContext) *m.Message
}

type ProcessorDictionary map[m.Message_MessageType]MessageProcessor
