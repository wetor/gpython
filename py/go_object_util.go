package py

import (
	"reflect"
	"strings"
)

const (
	GoStructFieldRenameTag = "gpython"
)

const (
	MethodTypeNone     uint8 = 1 << iota
	MethodTypeProperty       // property
)

var (
	MethodTypePrefixes = map[string]uint8{
		"Method_":   MethodTypeNone,
		"Property_": MethodTypeProperty,
	}
)

type StructInfo struct {
	Name    string
	Fields  map[string]FieldInfo
	Methods map[string]MethodInfo
}

type FieldInfo struct {
	Name string
	Type reflect.Type
}

type MethodInfo struct {
	Name           string
	Type           uint8  // 方法类型
	ParamsPyFormat string // 参数校验字符串，如OOnn
	Params         []reflect.Type
	Returns        []reflect.Type
}

func GetStructInfo(data interface{}) (*StructInfo, *reflect.Value) {
	tempSrcValue := reflect.ValueOf(data)
	if tempSrcValue.Kind() == reflect.Pointer {
		if tempSrcValue.IsNil() {
			return nil, nil
		}
	}
	structValue := reflect.Indirect(tempSrcValue)

	structType := structValue.Type()
	structInfo := &StructInfo{
		Name:    structType.Name(),
		Fields:  make(map[string]FieldInfo, structType.NumField()),
		Methods: make(map[string]MethodInfo, structType.NumMethod()),
	}

	// 获取字段信息
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		name, ok := field.Tag.Lookup(GoStructFieldRenameTag)
		if !ok {
			name = field.Name
		}
		structInfo.Fields[name] = FieldInfo{
			Name: field.Name,
			Type: field.Type,
		}
	}

	// 获取方法信息
	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i)
		methodType := method.Type
		params := make([]reflect.Type, 0, methodType.NumIn()-1)
		returns := make([]reflect.Type, 0, methodType.NumOut())

		// 跳过self
		for j := 1; j < methodType.NumIn(); j++ {
			params = append(params, methodType.In(j))
		}

		for j := 0; j < methodType.NumOut(); j++ {
			returns = append(returns, methodType.Out(j))
		}

		name := method.Name
		mType := MethodTypeNone
		for prefix, value := range MethodTypePrefixes {
			if strings.HasPrefix(name, prefix) {
				name = name[len(prefix):]
				mType |= value
			}
		}

		structInfo.Methods[name] = MethodInfo{
			Name:           method.Name,
			Type:           mType,
			ParamsPyFormat: "",
			Params:         params,
			Returns:        returns,
		}
	}

	return structInfo, &structValue
}

func StructToObject(src any) (Object, error) {
	tempSrcValue := reflect.ValueOf(src)
	if tempSrcValue.Kind() == reflect.Pointer {
		if tempSrcValue.IsNil() {
			return None, nil
		}
	}
	srcValue := reflect.Indirect(tempSrcValue)

	srcType := srcValue.Type()
	dst := NewStringDictSized(srcType.NumField())
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		value := srcValue.Field(i).Interface()

		keyName := field.Tag.Get("json")
		if len(keyName) == 0 {
			keyName = field.Name
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Pointer:
			obj, err := StructToObject(value)
			if err != nil {
				return nil, err
			}
			dst[keyName] = obj
		default:
			obj, err := ToObject(value)
			if err != nil {
				return nil, err
			}
			dst[keyName] = obj
		}
	}
	return dst, nil
}

func ToObject(goVal any) (Object, error) {
	var pyObj Object
	switch val := goVal.(type) {
	case nil:
		pyObj = None
	case bool:
		pyObj = NewBool(val)
	case int8:
		pyObj = Int(val)
	case uint8:
		pyObj = Int(val)
	case int16:
		pyObj = Int(val)
	case uint16:
		pyObj = Int(val)
	case int:
		pyObj = Int(val)
	case uint:
		pyObj = Int(val)
	case int32:
		pyObj = Int(val)
	case uint32:
		pyObj = Int(val)
	case int64:
		pyObj = Int(val)
	case uint64:
		pyObj = Int(val)
	case float32:
		pyObj = Float(val)
	case float64:
		pyObj = Float(val)
	case string:
		pyObj = String(val)
	case []byte:
		pyObj = Bytes(val)
	case map[string]any:
		pyValDict := NewStringDictSized(len(val))
		for key, value := range val {
			obj, err := ToObject(value)
			if err != nil {
				return nil, err
			}
			pyValDict[key] = obj
		}
		pyObj = pyValDict
	default:
		refVal := reflect.ValueOf(goVal)
		switch refVal.Kind() {
		case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
			reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Int, reflect.Uint:
			// 派生类型转换
			pyObj = Int(refVal.Int())
		case reflect.Float32, reflect.Float64:
			// 派生类型转换
			pyObj = Float(refVal.Float())
		case reflect.String:
			// 派生类型转换
			pyObj = String(refVal.String())
		case reflect.Array, reflect.Slice:
			l := refVal.Len()
			pyValList := NewListWithCapacity(l)
			for i := 0; i < l; i++ {
				obj, err := ToObject(refVal.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				pyValList.Append(obj)
			}
			pyObj = pyValList
		case reflect.Struct, reflect.Pointer:
			obj, err := StructToObject(goVal)
			if err != nil {
				return nil, err
			}
			pyObj = obj
		default:
			return nil, ExceptionNewf(TypeError, "'%s' type is unsupported", reflect.ValueOf(goVal).Type())
		}
	}
	return pyObj, nil
}
