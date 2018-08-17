package overrides

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenericOverrides(t *testing.T) {

	Convey("GenericOverrides", t, func() {

		Convey("MergeMaps function", func() {

			Convey("Should merge two maps with non-overlapping keys", func() {
				const m1 = `a:
  b:
    j: "100"
    k: "200"
    l: "300"
`
				const m2 = `p:
  q:
    x1: "1100"
    y1: "2100"
    z1: "3100"
`
				const expected = `a:
  b:
    j: "100"
    k: "200"
    l: "300"
p:
  q:
    x1: "1100"
    y1: "2100"
    z1: "3100"
`
				baseMap, err := ToMap(m1)
				So(err, ShouldBeNil)
				map2, err := ToMap(m2)
				So(err, ShouldBeNil)
				MergeMaps(baseMap, map2)
				res, err := ToYaml(baseMap)
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expected)
			})

			Convey("Should merge two maps with overlapping keys", func() {
				const m1 = `a:
  b:
    j: "100"
    k: "200"
    l: 300
`
				const m2 = `a:
  b:
    i: "1100"
    j: 100
    k:
      x1: foo
      y1:
        z1: bar
    l: "300"

`
				const expected = `a:
  b:
    i: "1100"
    j: 100
    k:
      x1: foo
      y1:
        z1: bar
    l: "300"
`
				baseMap, err := ToMap(m1)
				So(err, ShouldBeNil)
				map2, err := ToMap(m2)
				So(err, ShouldBeNil)
				MergeMaps(baseMap, map2)
				res, err := ToYaml(baseMap)
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expected)
			})

			Convey("Should merge a map with an empty one", func() {
				const m1 = `a:
  b:
    j: "100"
    k: 200
    l: abc
`
				const expected = `a:
  b:
    j: "100"
    k: 200
    l: abc
`
				baseMap, err := ToMap(m1)
				So(err, ShouldBeNil)
				map2, err := ToMap("")
				So(err, ShouldBeNil)
				MergeMaps(baseMap, map2)
				res, err := ToYaml(baseMap)
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expected)
			})
		})

		Convey("ToYaml function", func() {

			Convey("Should not fail for empty map", func() {

				inputMap := map[string]string{}
				res, err := ToYaml(UnflattenMap(inputMap))
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
				res, err := ToYaml(UnflattenMap(inputMap))
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expected)
			})

			Convey("Should handle multi-line string correctly", func() {

				const expected = `a:
  b:
    c: "100"
    d: "200"
    e: |
      300
      400
      500
`
				inputMap := map[string]string{}
				inputMap["a.b.c"] = "100"
				inputMap["a.b.d"] = "200"
				inputMap["a.b.e"] = "300\n400\n500\n"
				res, err := ToYaml(UnflattenMap(inputMap))
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
				res, err := ToYaml(UnflattenMap(inputMap))
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expected)
			})

		})

		Convey("ToMap function", func() {
			Convey("Should unmarshall yaml into a map", func() {
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

		})

		Convey("FlattenMap function", func() {

			Convey("Should flatten the map", func() {
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
	})
}
