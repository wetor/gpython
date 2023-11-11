# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc="property"

class A:
    _value = 123
    @property
    def value(cls):
        return cls._value

a = A()
assert a.value == 123

a._value = 456
assert a.value == 456

try:
    a.value = 666
except AttributeError:
    pass
else:
    assert False, "AttributeError not raised"

assertRaises(TypeError, a.value)

doc="property2"
class C(object):
    @property
    def x(self):
        return self._x
    @x.setter
    def x(self, value):
        self._x = value
    @x.deleter
    def x(self):
        del self._x

c = C()
c.x = 123
assert c.x == 123
c.x = 456
assert c.x == 456

del c.x
try:
    _ = c.x
except AttributeError:
    pass
else:
    assert False, "AttributeError not raised"


doc="finished"
