package overrides

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenericOverrides(t *testing.T) {
	Convey("GenericOverrides", t, func() {

		Convey("Should not fail for empty map", func() {

			overridesMap := map[string]string{}
			res, err := flatMapToYaml(overridesMap)
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
			res, err := flatMapToYaml(overridesMap)
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
			res, err := flatMapToYaml(overridesMap)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, expected)
		})

		Convey("Unmarshall yaml into a map", func() {
			const value = `
a:
  b:
    c: "100"
    d: "200"
    e: "300"
`
			res, err := unmarshallToNestedMap(value)
			So(err, ShouldBeNil)

			a, ok := res["a"].(map[string]interface{})
			So(ok, ShouldBeTrue)

			b, ok := a["b"].(map[string]interface{})
			So(ok, ShouldBeTrue)

			c, ok := b["c"].(string)
			So(ok, ShouldBeTrue)
			So(c, ShouldEqual, "100")

			d, ok := b["d"].(string)
			So(ok, ShouldBeTrue)
			So(d, ShouldEqual, "200")

			e, ok := b["e"].(string)
			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, "300")
		})
	})
}
