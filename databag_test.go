package databag

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestNamespaceSpliter(t *testing.T) {
	parts := namespaceSplitter("this.is.a.test", ".")

	assert.Len(t, parts, 4)
}

func TestDataBagGet(t *testing.T) {
	data := map[interface{}]interface{}{
		"a": map[interface{}]interface{}{
			"b": map[interface{}]interface{} {
				"c": 10,
			},
		},

		"b": 20,
		"c": map[interface{}]interface{}{
			"d": 30,
		},
	}

	bag := NewDataBagFrom(data)

	{
		r, ok := bag.Get("a.b.c")
		assert.True(t, ok)
		assert.Equal(t, 10, r)
	}

	{
		r, ok := bag.Get("b")
		assert.True(t, ok)
		assert.Equal(t, 20, r)
	}

	{
		r, ok := bag.Get("c.d")
		assert.True(t, ok)
		assert.Equal(t, 30, r)
	}


	{
		r, ok := bag.Get("c")
		assert.True(t, ok)
		assert.Equal(t, map[interface{}]interface{}{"d": 30}, r)
	}
}

func TestDataBagSet(t *testing.T) {
	bag := NewDataBag()

	{
		bag.Set("a.b.c", 10)
		r, ok := bag.Get("a.b.c")
		assert.True(t, ok)
		assert.Equal(t, 10, r)
	}

	{
		bag.Set("c", 10)
		r, ok := bag.Get("c")
		assert.True(t, ok)
		assert.Equal(t, 10, r)
	}

	{
		bag.Set("c.d", 10)
		r, ok := bag.Get("c.d")
		assert.True(t, ok)
		assert.Equal(t, 10, r)
	}
}

func TestDataBagAll(t *testing.T) {
	bag := NewDataBag()

	{
		bag.Set("a.b.c", 10)
		r := bag.All()
		assert.Equal(t, r, map[interface{}]interface{}{
			"a": map[interface{}]interface{}{
				"b": map[interface{}]interface{}{
					"c": 10,
				},
			},
		})
	}
}

func TestDataBagMerge(t *testing.T) {
	a := NewDataBagFrom(map[interface{}]interface{}{
		"messages": map[interface{}]interface{}{
			"foo": map[interface{}]interface{} {
				"bar": "This is test",
			},
			"boo": "This is boo value",
		},
	})

	b := NewDataBagFrom(map[interface{}]interface{}{
		"messages": map[interface{}]interface{}{
			"foo": map[interface{}]interface{} {
				"bar": "This is another test",
			},
			"bar": map[interface{}]interface{} {
				"foo": "This is a test",
			},
		},
		"validations": map[interface{}]interface{}{
			"invalid": map[interface{}]interface{}{
				"foo": "Invalid value",
			},
		},
	})

	a.Merge(b)

	{
		r, ok := a.Get("messages.foo.bar")
		if assert.True(t, ok) {
			assert.Exactly(t, r, "This is another test")
		}
	}
	{
		r, ok := a.Get("messages.boo")
		if assert.True(t, ok) {
			assert.Exactly(t, r, "This is boo value")
		}
	}
	{
		r, ok := a.Get("messages.bar.foo")
		if assert.True(t, ok) {
			assert.Exactly(t, r, "This is a test")
		}
	}
	{
		r, ok := a.Get("validations.invalid.foo")
		if assert.True(t, ok) {
			assert.Exactly(t, r, "Invalid value")
		}
	}
}

func TestDeepCopyMap(t *testing.T) {
	src := map[interface{}]interface{}{
		"foo": map[interface{}]interface{} {
			"bar": map[interface{}]interface{} {
				"a": "This is A value",
				"b": "This is B value",
			},
		},
		"boo": map[interface{}]interface{} {
			"c": "This is C value",
			"d": map[interface{}]interface{} {
				"e": "This is E value",
				"f": "This is F value",
			},
		},
		"doo": "This is doo value",
		"moo": "This is moo value",
	}

	dst := make(map[interface{}]interface{})
	deepCopyMap(dst, src)

	if assert.Exactly(t, src, dst) {
		dst["foo"] = "changed"
		assert.NotEqual(t, src["foo"], dst["foo"])

		dst["foo"] = "changed"
		assert.NotEqual(t, src["foo"], dst["foo"])

		(dst["boo"].(map[interface{}]interface{}))["d"].(map[interface{}]interface{})["e"] = "changed"
		assert.NotEqual(t, src, dst)
	}
}

func TestDeepMergeMap(t *testing.T) {
	a := map[interface{}]interface{}{
		"foo": map[interface{}]interface{} {
			"bar": map[interface{}]interface{} {
				"a": "This is A value",
				"b": "This is B value",
			},
		},
		"boo": map[interface{}]interface{} {
			"c": "This is C value",
			"d": map[interface{}]interface{} {
				"e": "This is E value",
				"f": "This is F value",
			},
		},
		"doo": "This is doo value",
		"moo": "This is moo value",
	}

	b := map[interface{}]interface{}{
		"foo": map[interface{}]interface{} {
			"bar": map[interface{}]interface{} {
				"a": "This is another A value",
				"b": "This is another B value",
			},
		},
		"boo": map[interface{}]interface{} {
			"d": map[interface{}]interface{} {
				"f": "This is another F value",
			},
		},
		"coo": map[interface{}]interface{} {
			"c": "This is C value",
			"d": map[interface{}]interface{} {
				"e": "This is E value",
				"f": "This is F value",
			},
		},
		"moo": "This is another moo value",
	}

	c := map[interface{}]interface{} {
		"boo": map[interface{}]interface{} {
			"d": map[interface{}]interface{} {
				"f": "This is TOTALLY another F value",
			},
		},
	}

	r := deepMergeMap(a, b, c)

	assert.Exactly(
		t,
		(r["boo"].(map[interface{}]interface{}))["d"].(map[interface{}]interface{})["f"],
		"This is TOTALLY another F value",
	)
	assert.Exactly(
		t,
		r["moo"],
		"This is another moo value",
	)
}

func ExampleDataBag_Get() {
	bag := NewDataBag()
	bag.Set("a.b.c.d", "this is some value")

	fmt.Println(bag.Get("a.b.c.d"))
	// Output: this is some value
}

func ExampleDataBag_Set() {
	bag := NewDataBag()
	bag.Set("a.b.c.d", "this is some value")

	fmt.Println(bag.Get("a.b.c.d"))
	// Output: this is some value
}

func ExampleDataBag_All() {
	bag := NewDataBag()
	bag.Set("a.b.c.d", "this is some value")

	fmt.Println(bag.All())
	// Output: map[a:map[b:map[c:map[d:this is some value]]]]
}

func ExampleDataBag_Merge() {
	a := NewDataBag()
	a.Set("a.b.c.d", "this is d value")
	a.Set("a.b.c.f", "this is f value")
	a.Set("foo.bar", "this is bar value")

	b := NewDataBag()
	b.Set("a.b.c.f", "this is the other f value")
	b.Set("foo.bar", "this is the other bar value")

	a.Merge(b)
	fmt.Println(a.All())
	// Output: map[a:map[b:map[c:map[f:this is the other f value d:this is d value]]] foo:map[bar:this is the other bar value]]
}
