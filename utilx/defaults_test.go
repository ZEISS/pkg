package utilx

import (
	"time"

	v1 "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.

type DefaultsSuite struct{}

var _ = v1.Suite(&DefaultsSuite{})

type Parent struct {
	Children []Child
}

type Child struct {
	Name string
	Age  int `default:"10"`
}

type ExampleBasic struct {
	Bool       bool    `default:"true"`
	Integer    int     `default:"33"`
	Integer8   int8    `default:"8"`
	Integer16  int16   `default:"16"`
	Integer32  int32   `default:"32"`
	Integer64  int64   `default:"64"`
	UInteger   uint    `default:"11"`
	UInteger8  uint8   `default:"18"`
	UInteger16 uint16  `default:"116"`
	UInteger32 uint32  `default:"132"`
	UInteger64 uint64  `default:"164"`
	String     string  `default:"foo"`
	Bytes      []byte  `default:"bar"`
	Float32    float32 `default:"3.2"`
	Float64    float64 `default:"6.4"`
	Struct     struct {
		Bool    bool `default:"true"`
		Integer int  `default:"33"`
	}
	Duration         time.Duration `default:"1s"`
	Children         []Child
	Second           time.Duration `default:"1s"`
	StringSlice      []string      `default:"[1,2,3,4]"`
	IntSlice         []int         `default:"[1,2,3,4]"`
	IntSliceSlice    [][]int       `default:"[[1],[2],[3],[4]]"`
	StringSliceSlice [][]string    `default:"[[1],[]]"`

	DateTime  string `default:"{{date:1,-10,0}} {{time:1,-5,10}}"`
	HumanSize string `default:"1KiB"`
}

func (s *DefaultsSuite) TestSetDefaultsBasic(c *v1.C) {
	foo := &ExampleBasic{}
	SetDefaults(foo)

	s.assertTypes(c, foo)
}

type ExampleNested struct {
	Struct ExampleBasic
}

func (s *DefaultsSuite) TestSetDefaultsNested(c *v1.C) {
	foo := &ExampleNested{}
	SetDefaults(foo)

	s.assertTypes(c, &foo.Struct)
}

func (s *DefaultsSuite) assertTypes(c *v1.C, foo *ExampleBasic) {
	c.Assert(foo.Bool, v1.Equals, true)
	c.Assert(foo.Integer, v1.Equals, 33)
	c.Assert(foo.Integer8, v1.Equals, int8(8))
	c.Assert(foo.Integer16, v1.Equals, int16(16))
	c.Assert(foo.Integer32, v1.Equals, int32(32))
	c.Assert(foo.Integer64, v1.Equals, int64(64))
	c.Assert(foo.UInteger, v1.Equals, uint(11))
	c.Assert(foo.UInteger8, v1.Equals, uint8(18))
	c.Assert(foo.UInteger16, v1.Equals, uint16(116))
	c.Assert(foo.UInteger32, v1.Equals, uint32(132))
	c.Assert(foo.UInteger64, v1.Equals, uint64(164))
	c.Assert(foo.String, v1.Equals, "foo")
	c.Assert(string(foo.Bytes), v1.Equals, "bar")
	c.Assert(foo.Float32, v1.Equals, float32(3.2))
	c.Assert(foo.Float64, v1.Equals, 6.4)
	c.Assert(foo.Struct.Bool, v1.Equals, true)
	c.Assert(foo.Duration, v1.Equals, time.Second)
	c.Assert(foo.Children, v1.IsNil)
	c.Assert(foo.Second, v1.Equals, time.Second)
	c.Assert(foo.StringSlice, v1.DeepEquals, []string{"1", "2", "3", "4"})
	c.Assert(foo.IntSlice, v1.DeepEquals, []int{1, 2, 3, 4})
	c.Assert(foo.IntSliceSlice, v1.DeepEquals, [][]int{{1}, {2}, {3}, {4}})
	c.Assert(foo.StringSliceSlice, v1.DeepEquals, [][]string{{"1"}, {}})
	c.Assert(foo.DateTime, v1.Equals, "2020-08-10 12:55:10")
	c.Assert(foo.HumanSize, v1.Equals, "1KiB")
}

func (s *DefaultsSuite) TestSetDefaultsWithValues(c *v1.C) {
	foo := &ExampleBasic{
		Integer:  55,
		UInteger: 22,
		Float32:  9.9,
		String:   "bar",
		Bytes:    []byte("foo"),
		Children: []Child{{Name: "alice"}, {Name: "bob", Age: 2}},
	}

	SetDefaults(foo)

	c.Assert(foo.Integer, v1.Equals, 55)
	c.Assert(foo.UInteger, v1.Equals, uint(22))
	c.Assert(foo.Float32, v1.Equals, float32(9.9))
	c.Assert(foo.String, v1.Equals, "bar")
	c.Assert(string(foo.Bytes), v1.Equals, "foo")
	c.Assert(foo.Children[0].Age, v1.Equals, 10)
	c.Assert(foo.Children[1].Age, v1.Equals, 2)
}

func (s *DefaultsSuite) BenchmarkLogic(c *v1.C) {
	for i := 0; i < c.N; i++ {
		foo := &ExampleBasic{}
		SetDefaults(foo)
	}
}
