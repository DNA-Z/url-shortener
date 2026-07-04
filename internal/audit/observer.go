package audit

type Observer interface {
	Notify(event AuditEvent) error
}

type IPublisher interface {
	Register(observer Observer)
	Unregister(observer Observer)
	Publish(event AuditEvent)
}
