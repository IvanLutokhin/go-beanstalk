package beanstalk

import (
	"errors"
	"sync/atomic"
)

type Pool struct {
	newClient func() (*Client, error)
	clients   chan *Client
	closed    int32
}

func NewPool(newClient func() (*Client, error), capacity int, open bool) (*Pool, error) {
	if newClient == nil {
		return nil, errors.New("beanstalk: pool: factory function is nil")
	}

	if capacity < 0 {
		return nil, errors.New("beanstalk: pool: capacity should be greater than or equal to 1")
	}

	p := &Pool{
		newClient: newClient,
		clients:   make(chan *Client, capacity),
		closed:    1,
	}

	if open {
		if err := p.Open(); err != nil {
			return p, err
		}
	}

	return p, nil
}

func NewDefaultPool(address string, capacity int, open bool) (*Pool, error) {
	return NewPool(func() (*Client, error) { return Dial(address) }, capacity, open)
}

func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

func (p *Pool) Open() error {
	if !p.IsClosed() {
		return errors.New("beanstalk: pool: already opened")
	}

	atomic.StoreInt32(&p.closed, 0)

	for i := 0; i < cap(p.clients); i++ {
		client, err := p.newClient()
		if err != nil {
			return err
		}

		p.clients <- client
	}

	return nil
}

func (p *Pool) Close() error {
	if p.IsClosed() {
		return errors.New("beanstalk: pool: already closed")
	}

	atomic.StoreInt32(&p.closed, 1)

	close(p.clients)

	for client := range p.clients {
		client.Close()
	}

	return nil
}

func (p *Pool) Get() (*Client, error) {
	if p.IsClosed() {
		return nil, errors.New("beanstalk: pool: closed")
	}

	select {
	case client := <-p.clients:
		return client, nil

	default:
		return p.newClient()
	}
}

func (p *Pool) Put(client *Client) error {
	if p.IsClosed() {
		return errors.New("beanstalk: pool: closed")
	}

	select {
	case p.clients <- client:
		return nil

	default:
		return client.Close()
	}
}

func (p *Pool) Len() int {
	return len(p.clients)
}
