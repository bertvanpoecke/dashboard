# Shutdown [![Build Status](http://jenkins.be-mobile-ops.net/job/tolling/job/go-lib/job/master/badge/icon)](http://jenkins.be-mobile-ops.net/job/tolling/job/go-lib/job/master/)

<p>Simple library to make it easier to manage the 'graceful' shutdown of your application.</p>

## Based on

Version | Go Get URL | source | doc | Notes |
--------|------------|--------|-----|-------|
2.x     | [gopkg.in/vrecan/death.v2](https://gopkg.in/vrecan/death.v2)| [source]() | [doc]() | This supports loggers who _do not_ return an error from their `Error` and `Warn` functions like [logrus](https://github.com/sirupsen/logrus)



Example
```bash
go get bitbucket.org/be-mobile/go-lib/shutdown
```
## Use The Library

```go
package main

import (
	"bitbucket.org/be-mobile/go-lib/shutdown"
	SYS "syscall"
)

func main() {
	sample_shutdown := shutdown.NewShutdown(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	//when you want to block for shutdown signals
	sample_shutdown.WaitForShutdown() // this will finish when a signal of your type is sent to your application
}
```

### Close Other Objects On Shutdown
<p>One simple feature of shutdown is that it can also close other objects when shutdown starts</p>

```go
package main

import (
	"bitbucket.org/be-mobile/go-lib/shutdown"
	SYS "syscall"
	"io"
)

func main() {
	sample_shutdown := shutdown.NewShutdown(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	objects := make([]io.Closer, 0)

	objects = append(objects, &NewType{}) // this will work as long as the type implements a Close method

	//when you want to block for shutdown signals
	sample_shutdown.WaitForShutdown(objects...) // this will finish when a signal of your type is sent to your application
}

type NewType struct {
}

func (c *NewType) Close() error {
	return nil
}
```

### Or close using an anonymous function

```go
package main

import (
	"bitbucket.org/be-mobile/go-lib/shutdown"
	SYS "syscall"
)

func main() {
	sample_shutdown := shutdown.NewShutdown(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application
	//when you want to block for shutdown signals
	sample_shutdown.WaitForShutdownWithFunc(func(){ 
		//do whatever you want on shutdown
	}) 
}
```

## go-samples

Some samples are available at the [go-samples repository](https://bitbucket.org/be-mobile/go-samples/src)

* graceful_shutdown_lib