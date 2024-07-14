package optional

import "go.opentelemetry.io/collector/confmap"

// Optional is a type that can be used to represent a value that may or may not be present.
// It supports three flavors: Some(value), None(), and WithDefault(defaultValue).
type Optional[T any] struct {
	hasValue bool
	value    T

	defaultFn DefaultFunc[T]
}

type DefaultFunc[T any] func() T

var _ confmap.Unmarshaler = (*Optional[any])(nil)

// Some creates an Optional with a value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, hasValue: true}
}

// None creates an Optional with no value.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// WithDefault creates an Optional which has no value
// unless user config provides some, in which case
// the defaultValue is used as a starting point,
// which may be overridden by the user provided values.
func WithDefault[T any](defaultFn DefaultFunc[T]) Optional[T] {
	return Optional[T]{defaultFn: defaultFn}
}

func (o Optional[T]) HasValue() bool {
	return o.hasValue
}

func (o Optional[T]) Value() T {
	return o.value
}

func (o *Optional[T]) Unmarshal(conf *confmap.Conf) error {
	// we assume that Unmarshal will not be called if conf has no value.
	if o.defaultFn != nil {
		o.value = o.defaultFn()
		o.hasValue = true
	}
	if err := conf.Unmarshal(&o.value); err != nil {
		return err
	}
	o.hasValue = true
	return nil
}
