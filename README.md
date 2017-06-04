## About DataBag

Package Data Bag offers a simple implementation similar to data bag attribute from Symfony2/3 framework.

## Installation

Install in the usual way:

    go get -u github.com/vcraescu/databag

## Usage

Accessors example:
```go
bag := NewDataBag()
bag.Set("a.b.c.d", "this is some value")

fmt.Println(bag.Get("a.b.c.d"))
// Output: this is some value
```

Merge 2 data bags:
```go
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
```	


Please refer to the [GoDoc API](https://godoc.org/github.com/vcraescu/databag) listing for a summary of the API. 
