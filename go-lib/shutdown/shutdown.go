package shutdown

//Manage the shutdown of your application.

import (
	"io"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//Shutdown manages the shutdown of your application.
type Shutdown struct {
	wg          *sync.WaitGroup
	sigChannel  chan os.Signal
	callChannel chan interface{}
	timeout     time.Duration
	log         Logger
}

//Logger interface to log.
type Logger interface {
	Error(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
}

//closer is a wrapper to the struct we are going to close with metadata
//to help with debuging close.
type closer struct {
	Index   int
	C       io.Closer
	Name    string
	PKGPath string
}

//NewShutdown Create Shutdown with the signals you want to die from.
func NewShutdown(signals ...os.Signal) (shutdown *Shutdown) {
	shutdown = &Shutdown{timeout: 10 * time.Second,
		sigChannel:  make(chan os.Signal, 1),
		callChannel: make(chan interface{}, 1),
		wg:          &sync.WaitGroup{},
		log:         log.StandardLogger()}
	signal.Notify(shutdown.sigChannel, signals...)
	shutdown.wg.Add(1)
	go shutdown.listenForSignal()
	return shutdown
}

//SetTimeout Overrides the time shutdown is willing to wait for a objects to be closed.
func (d *Shutdown) SetTimeout(t time.Duration) *Shutdown {
	d.timeout = t
	return d
}

//SetLogger Overrides the default logger (seelog)
func (d *Shutdown) SetLogger(l Logger) *Shutdown {
	d.log = l
	return d
}

//WaitForShutdown wait for signal and then kill all items that need to die.
func (d *Shutdown) WaitForShutdown(closable ...io.Closer) {
	d.wg.Wait()
	d.log.Info("Shutdown started...")
	count := len(closable)
	d.log.Debug("Closing ", count, " objects")
	if count > 0 {
		d.closeInMass(closable...)
	}
}

//WaitForShutdownWithFunc allows you to have a single function get called when it's time to
//kill your application.
func (d *Shutdown) WaitForShutdownWithFunc(f func()) {
	d.wg.Wait()
	d.log.Info("Shutdown started...")
	f()
}

//getPkgPath for an io closer.
func getPkgPath(c io.Closer) (name string, pkgPath string) {
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name(), t.PkgPath()
}

//closeInMass Close all the objects at once and wait forr them to finish with a channel.
func (d *Shutdown) closeInMass(closable ...io.Closer) {

	count := len(closable)
	sentToClose := make(map[int]closer)
	//call close async
	doneClosers := make(chan closer, count)
	for i, c := range closable {
		name, pkgPath := getPkgPath(c)
		closer := closer{Index: i, C: c, Name: name, PKGPath: pkgPath}
		go d.closeObjects(closer, doneClosers)
		sentToClose[i] = closer
	}

	//wait on channel for notifications.

	timer := time.NewTimer(d.timeout)
	for {
		select {
		case <-timer.C:
			d.log.Warn(count, " object(s) remaining but timer expired.")
			for _, c := range sentToClose {
				d.log.Error("Failed to close: ", c.PKGPath, "/", c.Name)
			}
			return
		case closer := <-doneClosers:
			delete(sentToClose, closer.Index)
			count--
			d.log.Debug(count, " object(s) left")
			if count == 0 && len(sentToClose) == 0 {
				d.log.Debug("Finished closing objects")
				return
			}
		}
	}
}

//closeObjects and return a bool when finished on a channel.
func (d *Shutdown) closeObjects(closer closer, done chan<- closer) {
	err := closer.C.Close()
	if nil != err {
		d.log.Error(err)
	}
	done <- closer
}

//Start manually initiates the shutdown process.
func (d *Shutdown) Start() {
	select {
	case d.callChannel <- "quit":
	default:
	}
}

//ListenForSignal Manage shutdown of application by signal.
func (d *Shutdown) listenForSignal() {
	defer d.wg.Done()
	for {
		select {
		case <-d.sigChannel:
			return
		case <-d.callChannel:
			return
		}
	}
}

//Expose the manuel trigger for shutting down
func (d *Shutdown) CallChannel() chan interface{} {
	return d.callChannel
}
