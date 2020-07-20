package rk_query


import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

type eventEntryImpl struct{}

func (entry *eventEntryImpl) FormatAsRk(event Event) string {
	eventImpl := event.(*EventImpl)

	eventImpl.addDefaultKvs()

	builder := bytes.Buffer{}

	builder.WriteString(EventScopeDelimiter + "\n")

	entry.addPredifinedValues(&builder, eventImpl)

	builder.WriteString(EventEOE)
	return builder.String()
}

func (entry *eventEntryImpl) FormatAsRkMin(event Event) string {
	eventImpl := event.(*EventImpl)

	eventImpl.addDefaultKvs()

	builder := bytes.Buffer{}

	builder.WriteString(EventScopeDelimiter + "\n")

	builder.WriteString("timing=")
	eventImpl.appendTimers(&builder)
	builder.WriteString("\n")

	builder.WriteString("counters=")
	eventImpl.appendCounters(&builder)
	builder.WriteString("\n")

	builder.WriteString("kvs=")
	eventImpl.appendKvs(&builder)
	builder.WriteString("\n")

	builder.WriteString("errors=")
	eventImpl.appendErrs(&builder)
	builder.WriteString("\n")

	builder.WriteString(EventEOE)
	return builder.String()
}

func (entry *eventEntryImpl) FormatAsJson(event Event) *zap.Logger {
	logger := event.GetZapLogger()
	impl := event.(*EventImpl)
	fields := make([]zapcore.Field, 0)

	// EndTime
	fields = append(fields, zap.Time("end_time", time.Unix(0, event.GetEndTimeMS()*1000000)))
	// StartTime
	fields = append(fields, zap.Time("start_time", time.Unix(0, event.GetStartTimeMS()*1000000)))
	// Timers
	fields = append(fields, entry.getTimerAsZapFields(impl)...)
	// Counters
	fields = append(fields, entry.getCounterAsZapFields(impl)...)
	// KVs
	fields = append(fields, entry.getKvAsZapFields(impl)...)
	// Remote Address
	if len(event.GetRemoteAddr()) > 0 {
		fields = append(fields, zap.String("remote_addr", event.GetRemoteAddr()))
	}
	// App
	if len(event.GetAppName()) > 0 {
		fields = append(fields, zap.String("app", event.GetAppName()))
	}
	// Hostname
	if len(event.GetHostName()) > 0 {
		fields = append(fields, zap.String("host_name", event.GetHostName()))
	}
	// Status
	if len(event.GetEventStatus().String()) > 0 {
		fields = append(fields, zap.String("event_status", event.GetEventStatus().String()))
	}
	// History
	if impl.producesHistory() {
		builder := &bytes.Buffer{}
		event.GetEventHistory().appendTo(builder)

		fields = append(fields, zap.String("history", builder.String()))
	}

	return logger.With(fields...)
}

func (entry *eventEntryImpl) FormatAsJsonMin(event Event) *zap.Logger {
	logger := event.GetZapLogger()
	impl := event.(*EventImpl)

	fields := make([]zapcore.Field, 0)

	// Timers
	fields = append(fields, entry.getTimerAsZapFields(impl)...)
	// Counters
	fields = append(fields, entry.getCounterAsZapFields(impl)...)
	// Name values
	fields = append(fields, entry.getKvAsZapFields(impl)...)

	return logger.With(fields...)
}

func (entry *eventEntryImpl) getKvAsZapFields(event *EventImpl) []zapcore.Field {
	fields := make([]zapcore.Field, 0)

	for k, v := range event.kvs {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}

func (entry *eventEntryImpl) getErrAsZapFields(event *EventImpl) []zapcore.Field {
	fields := make([]zapcore.Field, 0)

	for k, v := range event.errors {
		fields = append(fields, zap.Int64(k, v))
	}

	return fields
}

func (entry *eventEntryImpl) getCounterAsZapFields(event *EventImpl) []zapcore.Field {
	fields := make([]zapcore.Field, 0)

	for k, v := range event.counters {
		fields = append(fields, zap.Int64(k, v))
	}

	return fields
}

func (entry *eventEntryImpl) getTimerAsZapFields(event *EventImpl) []zapcore.Field {
	fields := make([]zapcore.Field, 0)

	for _, v := range event.tracker {
		f, _ := v.ToZapFieldsWithTimeSource(event.timeSource)
		fields = append(fields, f...)
	}

	return fields
}

func (entry *eventEntryImpl) addPredifinedValues(builder *bytes.Buffer, event *EventImpl) {
	// EndTime
	builder.WriteString("end_time=" + time.Unix(0, event.GetEndTimeMS()*1000000).Format(time.RFC3339) + "\n")

	// StartTime
	builder.WriteString("start_time=" + time.Unix(0, event.GetStartTimeMS()*1000000).Format(time.RFC3339) + "\n")

	// Time
	builder.WriteString("time=" + strconv.FormatInt(event.GetEndTimeMS()-event.GetStartTimeMS(), 10) + "\n")

	// Hostname
	builder.WriteString("hostname=" + event.GetHostName() + "\n")

	// Timing
	if len(event.tracker) > 0 {
		builder.WriteString("timing=")
		event.appendTimers(builder)
		builder.WriteString("\n")
	}

	// Counters
	if event.hasCounters() {
		builder.WriteString("counters=")
		event.appendCounters(builder)
		builder.WriteString("\n")
	}

	// kvs
	if event.hasKvs() {
		builder.WriteString("kvs=")
		event.appendKvs(builder)
		builder.WriteString("\n")
	}

	// err
	if event.hasCounters() {
		builder.WriteString("errors=")
		event.appendErrs(builder)
		builder.WriteString("\n")
	}

	// Remote address
	if len(event.GetRemoteAddr()) > 0 {
		builder.WriteString("remote_addr=" + event.GetRemoteAddr() + "\n")
	}

	// Program
	builder.WriteString("app=" + event.GetAppName() + "\n")

	// Operation
	builder.WriteString("operation=" + event.GetOperation() + "\n")

	// History
	if event.producesHistory() && event.GetEventHistory().builder.Len() > 0 {
		builder.WriteString("history=")
		event.GetEventHistory().appendTo(builder)
		builder.WriteString("\n")
	}
}