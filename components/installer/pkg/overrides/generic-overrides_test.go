package overrides

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenericOverrides(t *testing.T) {
	Convey("GenericOverrides", t, func() {

		Convey("Should not fail for empty map", func() {

			inputMap := map[string]string{}
			res, err := ToYaml(flatMapToOverridesMap(inputMap))
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
			inputMap := map[string]string{}
			inputMap["a.b.c"] = "100"
			inputMap["a.b.d"] = "200"
			inputMap["a.b.e"] = "300"
			res, err := ToYaml(flatMapToOverridesMap(inputMap))
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
			inputMap := map[string]string{}
			inputMap["a.b.c"] = "100"
			inputMap["a.b.d"] = "200"
			inputMap["a.b.e"] = "300"
			inputMap["global.foo"] = "bar"
			inputMap["h.o.o"] = "xyz"
			res, err := ToYaml(flatMapToOverridesMap(inputMap))
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
			res, err := ToMap(value)
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

		Convey("flatten the map", func() {
			const value = `
a:
  b:
    c: "100"
    d: "200"
    e: "300"
`
			oMap, err := ToMap(value)
			So(err, ShouldBeNil)
			res := FlattenMap(oMap)
			So(len(res), ShouldEqual, 3)
			So(res["a.b.c"], ShouldEqual, "100")
			So(res["a.b.d"], ShouldEqual, "200")
			So(res["a.b.e"], ShouldEqual, "300")
		})
	})
}
