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

var PropertyType = ObjectType.NewType("property",
	`property(fget=None, fset=None, fdel=None, doc=None) -> property attribute

fget is a function to be used for getting an attribute value, and likewise
fset is a function for setting, and fdel a function for del'ing, an
attribute.  Typical use is to define a managed attribute x:

class C(object):
    def getx(self): return self._x
    def setx(self, value): self._x = value
    def delx(self): del self._x
    x = property(getx, setx, delx, "I'm the 'x' property.")

Decorators make defining new properties or modifying existing ones easy:

class C(object):
    @property
    def x(self):
        "I am the 'x' property."
        return self._x
    @x.setter
    def x(self, value):
        self._x = value
    @x.deleter
    def x(self):
        del self._x`, PropertyNew, nil)

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

// Properties
func init() {
	PropertyType.Dict["getter"] = MustNewMethod("getter", func(self, getter Object) (Object, error) {
		p := self.(*Property)
		p.Fget = func(self Object) (Object, error) {
			return Call(getter, Tuple{self}, nil)
		}
		return p, nil
	}, 0, "Descriptor to change the getter on a property.")
	PropertyType.Dict["setter"] = MustNewMethod("setter", func(self, setter Object) (Object, error) {
		p := self.(*Property)
		p.Fset = func(self, value Object) error {
			_, err := Call(setter, Tuple{self, value}, nil)
			return err
		}
		return p, nil
	}, 0, "Descriptor to change the setter on a property.")
	PropertyType.Dict["deleter"] = MustNewMethod("deleter", func(self, deleter Object) (Object, error) {
		p := self.(*Property)
		p.Fdel = func(self Object) error {
			_, err := Call(deleter, Tuple{self}, nil)
			return err
		}
		return p, nil
	}, 0, "Descriptor to change the deleter on a property.")
}

// Interfaces
var _ I__get__ = (*Property)(nil)
var _ I__set__ = (*Property)(nil)
var _ I__delete__ = (*Property)(nil)
