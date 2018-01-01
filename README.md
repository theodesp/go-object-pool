go-object-pool
---
<a href="https://godoc.org/github.com/theodesp/go-object-pool">
<img src="https://godoc.org/github.com/theodesp/go-object-pool?status.svg" alt="GoDoc">
</a>

<a href="https://opensource.org/licenses/MIT" rel="nofollow">
<img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="License"/>
</a>

<a href="https://travis-ci.org/theodesp/go-object-pool" rel="nofollow">
<img src="https://travis-ci.org/theodesp/go-object-pool.svg?branch=master" />
</a>

<a href="https://codecov.io/gh/theodesp/go-object-pool">
  <img src="https://codecov.io/gh/theodesp/go-object-pool/branch/master/graph/badge.svg" />
</a>

<a href="https://ci.appveyor.com/project/theodesp/go-object-pool" rel="nofollow">
  <img src="https://ci.appveyor.com/api/projects/status/7yiwtn68qmcj71xy?svg=true" />
</a>

The object pool pattern is a software **creational 
design pattern** that uses a set of initialized objects 
kept ready to use – a "pool" – rather than allocating and 
destroying them on demand. 
A client of the pool will request an object from the pool 
and perform operations on the returned object. 
When the client has finished, it returns the object 
to the pool rather than destroying it; 
this can be done manually or automatically.

## Installation
```bash
$ go get -u github.com/theodesp/go-object-pool
```

## Usage

Provide the necessary interface implementations first and then
create the Pool

```go
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

factory := &ByteBufferFactory{}
pool := NewFixedPool(16, factory)

obj, _ := pool.Get()
pool.Return(obj)
```


## LICENCE
Copyright © 2017 Theo Despoudis MIT license