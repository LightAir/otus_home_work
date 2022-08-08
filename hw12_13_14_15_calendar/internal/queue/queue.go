package queue

import "context"

type Queue interface {
	Connect(ctx context.Context) error
	Receive(name string, callback func(body []byte)) error
	Sent(body []byte, name string) error
	Close() error
}
