package overrides

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenericOverrides(t *testing.T) {
	Convey("GenericOverrides", t, func() {

		Convey("Should not fail for empty map", func() {

			overridesMap := map[string]string{}
			res, err := mapToYaml(overridesMap)
			So(err, ShouldBeNil)
			So(res, ShouldBeBlank)
		})

		Convey("Should merge several entries into one yaml", func() {

			const expected = `a:
  b:
    c: "100"
    d: "200"
    e: "300"
`
			overridesMap := map[string]string{}
			overridesMap["a.b.c"] = "100"
			overridesMap["a.b.d"] = "200"
			overridesMap["a.b.e"] = "300"
			res, err := mapToYaml(overridesMap)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, expected)
		})

		Convey("Should handle global values", func() {

			const expected = `a:
  b:
    c: "100"
    d: "200"
    e: "300"
global:
  foo: bar
h:
  o:
    o: xyz
`
			overridesMap := map[string]string{}
			overridesMap["a.b.c"] = "100"
			overridesMap["a.b.d"] = "200"
			overridesMap["a.b.e"] = "300"
			overridesMap["global.foo"] = "bar"
			overridesMap["h.o.o"] = "xyz"
			res, err := mapToYaml(overridesMap)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, expected)
		})
	})
}
