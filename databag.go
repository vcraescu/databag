package databag

import (
	"strings"
	"reflect"
)

const defaultNamespaceSep = "."

type Bag interface {
	Get(name string) (interface{}, bool)
	Set(name string, value interface{})
	All() map[interface{}]interface{}
	Merge(b Bag)
}

type DataBag struct {
	data map[interface{}]interface{}
	namespaceSep string
}

func NewDataBag() *DataBag {
	return NewDataBagFrom(make(map[interface{}]interface{}))
}

func NewDataBagFrom(data map[interface{}]interface{}) *DataBag {
	return &DataBag{
		data: data,
		namespaceSep: defaultNamespaceSep,
	}
}

func namespaceSplitter(name string, sep string) []string {
	return strings.Split(name, sep)
}

func (d *DataBag) Get(name string) (interface{}, bool) {
	keys := namespaceSplitter(name, d.namespaceSep)

	count := len(keys)
	if count == 1 {
		v, ok := d.data[name]
		return v, ok
	}

	cp := d.data
	for i, key := range keys {
		v, ok := cp[key]
		if !ok || i == count-1 {
			return v, ok
		}

		_, ok = cp[key].(map[interface{}]interface{})
		if !ok {
			return nil, ok
		}

		cp = cp[key].(map[interface{}]interface{})
	}

	return nil, false
}

func (d *DataBag) Set(name string, value interface{}) {
	keys := namespaceSplitter(name, d.namespaceSep)

	count := len(keys)
	if count == 1 {
		d.data[name] = value
		return
	}

	cp := d.data
	for i := 0; i < count; i++ {
		key := keys[i]
		if i == count-1 {
			cp[key] = value
			break
		}

		_, ok := cp[key]
		if !ok || reflect.TypeOf(cp[key]).Kind() != reflect.Map {
			cp[key] = make(map[interface{}]interface{})
		}

		cp = cp[key].(map[interface{}]interface{})
	}
}

func (d DataBag) All() map[interface{}]interface{} {
	return d.data
}

func (d *DataBag) Merge(b Bag) {
	m := deepMergeMap(d.All(), b.All())
	for name, value := range m {
		d.Set(name.(string), value)
	}
}

func deepCopyMap(dst map[interface{}]interface{}, src map[interface{}]interface{}) {
	for k, v := range src {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			tmp := make(map[interface{}]interface{})
			deepCopyMap(tmp, v.(map[interface{}]interface{}))
			dst[k] = tmp
			continue
		}

		dst[k] = v
	}
}

func deepMergeMap(maps ...map[interface{}]interface{}) map[interface{}]interface{} {
	if len(maps) == 1 {
		return maps[0]
	}

	r := make(map[interface{}]interface{})
	deepCopyMap(r, maps[0])

	for _, m := range maps[1:] {
		for key, value := range m {
			if _, ok := r[key]; !ok {
				r[key] = value
				continue
			}

			isMap := reflect.TypeOf(value).Kind() == reflect.Map &&
				reflect.TypeOf(r[key]).Kind() == reflect.Map
			if isMap {
				r[key] = deepMergeMap(
					r[key].(map[interface{}]interface{}),
					value.(map[interface{}]interface{}),
				)
				continue
			}

			r[key] = value
		}
	}

	return r
}
