// +build linux bsd darwin

package shutdown

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/cihub/seelog"
)

type Unhashable map[string]interface{}

func (u Unhashable) Close() error {
	return nil
}

func TestShutdown(t *testing.T) {
	defer seelog.Flush()

	Convey("Validate shutdown handles unhashable types", t, func() {
		u := make(Unhashable)
		shutdown := NewShutdown(syscall.SIGTERM)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		shutdown.WaitForShutdown(u)
	})

	Convey("Validate shutdown happens cleanly", t, func() {
		shutdown := NewShutdown(syscall.SIGTERM)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
		shutdown.WaitForShutdown()

	})

	Convey("Validate shutdown happens with other signals", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		shutdown.WaitForShutdown(closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate shutdown happens with a manual call", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		closeMe := &CloseMe{}
		shutdown.Start()
		shutdown.WaitForShutdown(closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate multiple sword falls do not block", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		closeMe := &CloseMe{}
		shutdown.Start()
		shutdown.Start()
		shutdown.WaitForShutdown(closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate multiple sword falls do not block even after we have exited waitForDeath", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		closeMe := &CloseMe{}
		shutdown.Start()
		shutdown.Start()
		shutdown.WaitForShutdown(closeMe)
		shutdown.Start()
		shutdown.Start()
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Validate shutdown gives up after timeout", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		shutdown.SetTimeout(10 * time.Millisecond)
		neverClose := &neverClose{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		shutdown.WaitForShutdown(neverClose)

	})

	Convey("Validate shutdown uses new logger", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		closeMe := &CloseMe{}
		logger := &MockLogger{}
		shutdown.SetLogger(logger)

		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		shutdown.WaitForShutdown(closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
		So(logger.Logs, ShouldNotBeEmpty)
	})

	Convey("Close multiple things with one that fails the timer", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		shutdown.SetTimeout(10 * time.Millisecond)
		neverClose := &neverClose{}
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		shutdown.WaitForShutdown(neverClose, closeMe)
		So(closeMe.Closed, ShouldEqual, 1)
	})

	Convey("Close with anonymous function", t, func() {
		shutdown := NewShutdown(syscall.SIGHUP)
		shutdown.SetTimeout(5 * time.Millisecond)
		closeMe := &CloseMe{}
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
		shutdown.WaitForShutdownWithFunc(func() {
			closeMe.Close()
			So(true, ShouldBeTrue)
		})
		So(closeMe.Closed, ShouldEqual, 1)
	})

}

type MockLogger struct {
	Logs []interface{}
}

func (l *MockLogger) Info(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

func (l *MockLogger) Debug(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

func (l *MockLogger) Error(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

func (l *MockLogger) Warn(v ...interface{}) {
	for _, log := range v {
		l.Logs = append(l.Logs, log)
	}
}

type neverClose struct {
}

func (n *neverClose) Close() error {
	time.Sleep(2 * time.Minute)
	return nil
}

type CloseMe struct {
	Closed int
}

func (c *CloseMe) Close() error {
	c.Closed++
	return errors.New("I have been closed")
}
