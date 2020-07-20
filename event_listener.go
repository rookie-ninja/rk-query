package rk_query

type eventEntryListener interface {
	notify(eventEntry)
}