package rk_query


import "go.uber.org/zap"

// A helper function for easy use of EventData
type EventHelper struct {
	factory    *EventFactory
	TimeSource TimeSource
}

func NewEventHelperWithLogger(appName string, logger *zap.Logger) *EventHelper {
	factory, _ := NewEventFactory(logger)
	factory.AppName = appName

	return &EventHelper{factory, &RealTimeSource{}}
}

func (helper *EventHelper) Start(operationName string) Event {
	event := helper.factory.CreateEvent()

	event.SetOperation(operationName)
	event.SetStartTimeMS(helper.TimeSource.CurrentTimeMS())
	return event
}

func (helper *EventHelper) Finish(event Event) {
	event.SetEndTimeMS(helper.TimeSource.CurrentTimeMS())
	event.WriteLog()
}

func (helper *EventHelper) FinishWithCond(event Event, success bool) {
	if success {
		event.SetCounter("success", 1)
	} else {
		event.SetCounter("failure", 1)
	}

	helper.Finish(event)
}

func (helper *EventHelper) FinishWithError(event Event, err error) {
	if err == nil {
		helper.FinishWithCond(event, true)
	}

	helper.FinishWithCond(event, false)
}
