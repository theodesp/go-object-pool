package go_object_pool

import (
	"testing"

	"bytes"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FixedPoolSuite struct{}

var _ = Suite(&FixedPoolSuite{})

type ByteBufferObject struct {
	buffer *bytes.Buffer
}

func (b *ByteBufferObject) Reset() {
	b.buffer.Reset()
}

type ByteBufferFactory struct{}

func (f ByteBufferFactory) Create() (PooledObject, error) {
	return &ByteBufferObject{
		buffer: bytes.NewBuffer(make([]byte, 1024)),
	}, nil
}

func (s *FixedPoolSuite) TestConstructor(c *C) {
	factory := &ByteBufferFactory{}
	pool := NewFixedPool(16, factory)

	c.Assert(len(pool.inUse), Equals, 0)
	c.Assert(len(pool.available), Equals, 0)
}

func (s *FixedPoolSuite) TestGet(c *C) {
	factory := &ByteBufferFactory{}
	pool := NewFixedPool(16, factory)

	obj1, _ := pool.Get()
	obj2, _ := pool.Get()

	c.Assert(len(pool.available), Equals, 0)
	c.Assert(len(pool.inUse), Equals, 2)
	c.Assert(obj1, Not(Equals), obj2)
}

func (s *FixedPoolSuite) TestGetFull(c *C) {
	factory := &ByteBufferFactory{}
	pool := NewFixedPool(1, factory)

	_, err1 := pool.Get()
	_, err2 := pool.Get()

	c.Assert(err1, IsNil)
	c.Assert(err2, ErrorMatches, "fixed Pool reached maximum capacity")
}

func (s *FixedPoolSuite) TestReturn(c *C) {
	factory := &ByteBufferFactory{}
	pool := NewFixedPool(2, factory)

	obj1, _ := pool.Get()
	obj2, _ := pool.Get()

	pool.Return(obj1)
	pool.Return(obj2)

	c.Assert(len(pool.available), Equals, 2)
	c.Assert(len(pool.inUse), Equals, 0)
}

func (s *FixedPoolSuite) TestReturnInvalidObj(c *C) {
	factory := &ByteBufferFactory{}
	pool := NewFixedPool(2, factory)

	obj1, _ := pool.Get()
	obj2, _ := factory.Create()

	err1 := pool.Return(obj1)
	err2 := pool.Return(obj2)

	c.Assert(err1, IsNil)
	c.Assert(err2, ErrorMatches, "unrecognized pooled object returned")
}
