package py

import (
	"reflect"
)

const (
	goObjectDefaultDoc = ``
)

type GoObject struct {
	Info       *StructInfo
	Value      *reflect.Value // 存储结构体实例的 reflect.Value
	ObjectType *Type
}

func (s *GoObject) Type() *Type {
	return ObjectType
}

func NewGoObject(data interface{}) *GoObject {
	if data == nil {
		return nil
	}
	info, value := GetStructInfo(data)
	if info == nil {
		return nil
	}
	gs := &GoObject{
		Info:  info,
		Value: value,
	}
	docStr := goObjectDefaultDoc
	if docField, err := gs.GetField("__doc__"); err == nil {
		if str, ok := docField.(string); ok {
			docStr = str
		}
	}

	gs.ObjectType = ObjectType.NewType(info.Name, docStr, nil, nil)
	return gs
}

func (s *GoObject) GetField(name string) (interface{}, error) {
	fieldInfo, ok := s.Info.Fields[name]
	if !ok {
		return nil, ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", s.Type().Name, name)
	}
	field := s.Value.FieldByName(fieldInfo.Name)
	if !field.IsValid() {
		return nil, ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", s.Type().Name, name)
	}
	return field.Interface(), nil
}

func (s *GoObject) SetField(name string, value interface{}) error {
	fieldInfo, ok := s.Info.Fields[name]
	if !ok {
		return ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", s.Type().Name, name)
	}
	field := s.Value.FieldByName(fieldInfo.Name)
	if !field.IsValid() {
		return ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", s.Type().Name, name)
	}

	if !field.CanSet() {
		return ExceptionNewf(AttributeError, "'%s' attribute is not settable", name)
	}

	fieldValue := reflect.ValueOf(value)
	valueType := fieldValue.Type()
	if !fieldInfo.Type.AssignableTo(valueType) {
		return ExceptionNewf(TypeError, "'%s' type is not assignable to '%s' type", valueType, fieldInfo.Type)
	}
	field.Set(fieldValue)
	return nil
}

func (s *GoObject) CallMethod(name string, args ...interface{}) ([]interface{}, error) {
	methodInfo, ok := s.Info.Methods[name]
	if !ok {
		return nil, ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", s.Type().Name, name)
	}
	// 参数数量检查
	numArgs := len(args)
	numParams := len(methodInfo.Params)
	if numArgs != numParams {
		return nil, checkNumberOfArgs(name, numArgs, numParams, numParams, numParams)
	}
	// 构建参数列表
	inputArgs := make([]reflect.Value, numArgs)
	for i, arg := range args {
		inputArgs[i] = reflect.ValueOf(arg)
	}

	method := s.Value.MethodByName(methodInfo.Name)
	if !method.IsValid() {
		return nil, ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", s.Type().Name, name)
	}

	// 调用方法
	results := method.Call(inputArgs)

	// 转换结果为 interface{} 切片
	resultInterfaces := make([]interface{}, len(results))
	for i, result := range results {
		resultInterfaces[i] = result.Interface()
	}
	return resultInterfaces, nil
}
