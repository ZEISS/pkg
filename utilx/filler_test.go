package utilx

import (
	"fmt"
	"reflect"

	v1 "gopkg.in/check.v1"
)

type FillerSuite struct{}

var _ = v1.Suite(&FillerSuite{})

type FixtureTypeInt int

func (s *FillerSuite) TestFuncByNameIsEmpty(c *v1.C) {
	calledA := false
	calledB := false

	f := &Filler{
		FuncByName: map[string]FillerFunc{
			"Foo": func(_ *FieldData) {
				calledA = true
			},
		},
		FuncByKind: map[reflect.Kind]FillerFunc{
			reflect.Int: func(_ *FieldData) {
				calledB = true
			},
		},
	}

	f.Fill(&struct{ Foo int }{})
	c.Assert(calledA, v1.Equals, true)
	c.Assert(calledB, v1.Equals, false)
}

func (s *FillerSuite) TestFuncByTypeIsEmpty(c *v1.C) {
	calledA := false
	calledB := false

	t := GetTypeHash(reflect.TypeOf(new(FixtureTypeInt)))
	f := &Filler{
		FuncByType: map[TypeHash]FillerFunc{
			t: func(_ *FieldData) {
				calledA = true
			},
		},
		FuncByKind: map[reflect.Kind]FillerFunc{
			reflect.Int: func(_ *FieldData) {
				calledB = true
			},
		},
	}

	f.Fill(&struct{ Foo FixtureTypeInt }{})
	c.Assert(calledA, v1.Equals, true)
	c.Assert(calledB, v1.Equals, false)
}

func (s *FillerSuite) TestFuncByKindIsNotEmpty(c *v1.C) {
	called := false
	f := &Filler{FuncByKind: map[reflect.Kind]FillerFunc{
		reflect.Int: func(_ *FieldData) {
			called = true
		},
	}}

	f.Fill(&struct{ Foo int }{Foo: 42})
	c.Assert(called, v1.Equals, false)
}

func (s *FillerSuite) TestFuncByKindSlice(_ *v1.C) {
	fmt.Println(GetTypeHash(reflect.TypeOf(new([]string)))) // nolint:forbidigo
}

func (s *FillerSuite) TestFuncByKindTag(c *v1.C) {
	var called string
	f := &Filler{Tag: "foo", FuncByKind: map[reflect.Kind]FillerFunc{
		reflect.Int: func(field *FieldData) {
			called = field.TagValue
		},
	}}

	f.Fill(&struct {
		Foo int `foo:"qux"`
	}{})
	c.Assert(called, v1.Equals, "qux")
}

func (s *FillerSuite) TestFuncByKindIsEmpty(c *v1.C) {
	called := false
	f := &Filler{FuncByKind: map[reflect.Kind]FillerFunc{
		reflect.Int: func(_ *FieldData) {
			called = true
		},
	}}

	f.Fill(&struct{ Foo int }{})
	c.Assert(called, v1.Equals, true)
}
