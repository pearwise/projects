package box

import (
	"fmt"
	"reflect"
)

type ApplicationContext interface {
	Inject(isDeeped bool, name string, val any) error
	Autowise(isDeeped bool, name string, val any) error
}

type context struct {
	namedStore map[string]reflect.Value
	typedStore map[reflect.Type]reflect.Value
	deepNamedStore map[string]func()reflect.Value
	deepTypedStore map[reflect.Type]func()reflect.Value
}

var instance ApplicationContext

func init() {
	instance = &context{
		namedStore: make(map[string]reflect.Value),
        typedStore: make(map[reflect.Type]reflect.Value),
		deepNamedStore: make(map[string]func()reflect.Value),
		deepTypedStore: make(map[reflect.Type]func()reflect.Value),
	}
}

func (c *context) Inject(isDeeped bool, name string, val any) error {
	if name == "" {
		// type
		if isDeeped {
			v, ok := val.(func() reflect.Value)
			if !ok {
				return fmt.Errorf("inject: deep storage kind is not a func() reflect.Value")
			}
			ty := v().Type()
			if _, ok := c.deepTypedStore[ty];ok {
				return fmt.Errorf("inject: %v is ambiguous", ty)
			} else {
				c.deepTypedStore[ty] = v
			}
		} else {
			v := reflect.ValueOf(val)
			if _, ok := c.typedStore[v.Type()];ok {
				return fmt.Errorf("inject: %v is ambiguous", v.Type())
			} else {
				c.typedStore[v.Type()] = v
			}
		}
	} else {
		// name
		if isDeeped {
			v, ok := val.(func() reflect.Value)
			if !ok {
				return fmt.Errorf("inject: deep storage kind is not a func() reflect.Value")
			}
			if _, ok := c.deepNamedStore[name];ok {
				return fmt.Errorf("inject: %v is ambiguous", name)
			} else {
				c.deepNamedStore[name] = v
			}
		} else {
			v := reflect.ValueOf(val)
			if _, ok := c.namedStore[name];ok {
				return fmt.Errorf("inject: %v is ambiguous", v.Type())
			} else {
				c.namedStore[name] = v
			}
		}
	}
	return nil
}

func (c *context) Autowise(isDeeped bool, name string, val any) error {
	if val == nil {
		return fmt.Errorf("inject: nil value")
	}
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("inject: %v is not a pointer", rv)
	}
	ri := reflect.Indirect(rv)
	rt := ri.Type()
	if name == "" {
		// type
		if isDeeped {
			if fn, ok := c.deepTypedStore[rt];ok {
				v := fn()
				if v.CanConvert(rt) {
					ri.Set(v.Convert(rt))
					return nil
				}
			}
		} else {
			if v, ok := c.typedStore[rt];ok&&v.CanConvert(rt) {
				ri.Set(v.Convert(rt))
                return nil
			}
		}
	} else {
		// name
		if isDeeped {
			if fn, ok := c.deepNamedStore[name];ok {
                v := fn()
                if v.CanConvert(rt) {
                    ri.Set(v.Convert(rt))
                    return nil
                }
			}
		} else {
			if v, ok := c.namedStore[name];ok&&v.CanConvert(rt) {
                ri.Set(v.Convert(rt))
                return nil
            }
		}
	}

	return fmt.Errorf("inject: value can not convert to %s", ri.Type())
}

func Get[T any]() (T,error) {
	var v T
	return v, instance.Autowise(false, "", &v)
}

func GetByName[T any](name string) (T, error) {
	var v T
	if name == "" {
		return v, fmt.Errorf("inject: name can not be empty")
	}
	return v, instance.Autowise(false, name, &v)
}

func DeepGet[T any]() (T, error) {
	var v T
	return v, instance.Autowise(true, "", &v)
}

func DeepGetByName[T any](name string) (T, error) {
	var v T
	if name == "" {
		return v, fmt.Errorf("inject: name can not be empty")
	}
	return v, instance.Autowise(false, name, &v)
}

func Put(v any) error {
	return instance.Inject(false, "", v)
}

func PutByName(name string, v any) error {
	if name == "" {
		return fmt.Errorf("inject: name can not be empty")
	}
	return instance.Inject(false, name, v)
}

func DeepPut(v func() reflect.Value) error {
	return instance.Inject(true, "", v)
}

func DeepPutByName(name string, v any) error {
	if name == "" {
        return fmt.Errorf("inject: name can not be empty")
    }
	return instance.Autowise(true, name, v)
}