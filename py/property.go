// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Property object

package py

// A python Property object
type Property struct {
	Fget func(self Object) (Object, error)
	Fset func(self, value Object) error
	Fdel func(self Object) error
	Doc  string
}

var PropertyType = ObjectType.NewType("property", "property object", PropertyNew, nil)

// Type of this object
func (p *Property) Type() *Type {
	return PropertyType
}

func PropertyNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var callable Object
	err = UnpackTuple(args, kwargs, "property", 1, 1, &callable)
	if err != nil {
		return nil, err
	}
	return &Property{
		Fget: func(self Object) (Object, error) {
			return Call(callable, Tuple{self}, nil)
		},
	}, nil
}

func (p *Property) M__get__(instance, owner Object) (Object, error) {
	if p.Fget == nil {
		return nil, ExceptionNewf(AttributeError, "can't get attribute")
	}
	return p.Fget(instance)
}

func (p *Property) M__set__(instance, value Object) (Object, error) {
	if p.Fset == nil {
		return nil, ExceptionNewf(AttributeError, "can't set attribute")
	}
	return None, p.Fset(instance, value)
}

func (p *Property) M__delete__(instance Object) (Object, error) {
	if p.Fdel == nil {
		return nil, ExceptionNewf(AttributeError, "can't delete attribute")
	}
	return None, p.Fdel(instance)
}

// Interfaces
var _ I__get__ = (*Property)(nil)
var _ I__set__ = (*Property)(nil)
var _ I__delete__ = (*Property)(nil)
