package pool

import (
	"sync/atomic"

	"google.golang.org/grpc"
)

type Conn interface {
	Value() *grpc.ClientConn
	Close() error
}

type conn struct {
	cc    *grpc.ClientConn
	pool  *pool
	once  bool
	count atomic.Int32
}

var _ Conn = (*conn)(nil)

func (c *conn) Value() *grpc.ClientConn {
	return c.cc
}

func (c *conn) Close() error {
	c.pool.decRef()
	if c.once {
		return c.reset()
	}
	return nil
}

func (c *conn) reset() error {
	cc := c.cc
	c.cc = nil
	c.once = false
	if cc != nil {
		return cc.Close()
	}
	return nil
}
