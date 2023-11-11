import libtest as self


doc="test_single"
class C(object):
    @staticmethod
    def foo(): return 42

    def bar(self): return 30

print(C.foo())
print(C().foo())

self.assertEqual(C.foo(), 42)
self.assertEqual(C().foo(), 42)

print('***********',type(C().foo))
print('***********',type(C().bar))

doc="test_staticmethod_function"
@staticmethod
def notamethod(x):
    return x

self.assertRaises(TypeError, notamethod, 1)

doc="finished"
