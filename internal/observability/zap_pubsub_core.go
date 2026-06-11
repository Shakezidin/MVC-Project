package observability

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type PubSubCore struct {
	zapcore.Core
	publisher *PubSubPublisher
	service   string
}

func NewPubSubCore(
	core zapcore.Core,
	publisher *PubSubPublisher,
	service string,
) zapcore.Core {

	return &PubSubCore{
		Core:      core,
		publisher: publisher,
		service:   service,
	}
}

func (c *PubSubCore) Enabled(
	level zapcore.Level,
) bool {
	return true
}

func (c *PubSubCore) With(
	fields []zapcore.Field,
) zapcore.Core {
	return c
}

func (c *PubSubCore) Check(
	entry zapcore.Entry,
	checkedEntry *zapcore.CheckedEntry,
) *zapcore.CheckedEntry {

	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}

	return checkedEntry
}

func (c *PubSubCore) Write(
	entry zapcore.Entry,
	fields []zapcore.Field,
) error {

	logEntry := LogEntry{
		Service:   c.service,
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err := c.publisher.Publish(logEntry)
	if err != nil {
		return err
	}

	return nil
}

func (c *PubSubCore) Sync() error {
	return nil
}
