package shutdown

import (
	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetPkgPath(t *testing.T) {
	defer log.Flush()

	Convey("Give pkgPath a ptr", t, func() {
		c := &Closer{}
		name, pkgPath := getPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "bitbucket.org/be-mobile/go-lib/shutdown")

	})

	Convey("Give pkgPath a interface", t, func() {
		var closable Closable
		closable = Closer{}
		name, pkgPath := getPkgPath(closable)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "bitbucket.org/be-mobile/go-lib/shutdown")
	})

	Convey("Give pkgPath a copy", t, func() {
		c := Closer{}
		name, pkgPath := getPkgPath(c)
		So(name, ShouldEqual, "Closer")
		So(pkgPath, ShouldEqual, "bitbucket.org/be-mobile/go-lib/shutdown")
	})
}

type Closable interface {
	Close() error
}
type Closer struct {
}

func (c Closer) Close() error {
	return nil
}
