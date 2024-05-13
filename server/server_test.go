package main_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServerIntegration(t *testing.T) {
	Convey("Testing The Test", t, func() {
		So(1, ShouldEqual, 1)
	})
}
