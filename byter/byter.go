package byter

import (
	"bytes"
	"reflect"
)

var stringLengthTags = map[string]reflect.Kind{
	"uint8": reflect.Uint8,
	"uint16": reflect.Uint16,
	"uint32": reflect.Uint32,
	"uint64": reflect.Uint64,
	"int8": reflect.Uint8,
	"byte": reflect.Uint8,
	"int16": reflect.Int16,
	"int32": reflect.Int32,
	"int64": reflect.Int64,
}

type byter struct {
	Buff *bytes.Buffer
}

type EnumMarker interface {
	IsEnum()
}

func (b *byter) isEnum(field reflect.Value) bool {
	return field.CanInterface() && reflect.PtrTo(field.Type()).Implements(reflect.TypeOf((*EnumMarker)(nil)).Elem())
}