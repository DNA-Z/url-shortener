package audit

import "sync"

type Publisher struct {
	observers map[Observer]struct{}
	mu        sync.RWMutex
}

func NewPublisher() IPublisher {
	return &Publisher{
		observers: make(map[Observer]struct{}),
	}
}

func (p *Publisher) Register(observer Observer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observers[observer] = struct{}{}
}

func (p *Publisher) Unregister(observer Observer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.observers, observer)

	if fo, ok := observer.(*FileObserver); ok {
		fo.Close()
	}
}

func (p *Publisher) Publish(event AuditEvent) {
	p.mu.RLock()
	observers := make([]Observer, 0, len(p.observers))
	for observer := range p.observers {
		observers = append(observers, observer)
	}
	p.mu.RUnlock()

	for _, observer := range observers {
		go func(o Observer, e AuditEvent) {
			_ = o.Notify(e)
		}(observer, event)
	}
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for observer := range p.observers {
		if fo, ok := observer.(*FileObserver); ok {
			fo.Close()
		}
	}
	p.observers = make(map[Observer]struct{})
}
