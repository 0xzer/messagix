package byter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

func NewReader(data []byte) *byter {
	return &byter{
		Buff: bytes.NewBuffer(data),
	}
}

func (b *byter) readInteger(kind reflect.Kind, size int, endian string) (uint64, error) {
	if endian == "" {
		endian = "big"
	}
	buffer := make([]byte, size)
	_, err := b.Buff.Read(buffer)
	if err != nil {
		return 0, fmt.Errorf("error reading integer: %w", err)
	}

	if funcs, exists := endianOperations[endian]; exists {
		if ops, ok := funcs[kind]; ok {
			return ops.Read(buffer), nil
		}
	}

	return 0, fmt.Errorf("unsupported kind or endianness")
}

func (b *byter) ReadToStruct(s interface{}) error {
	values := reflect.ValueOf(s)
	if values.Kind() == reflect.Ptr && values.Elem().Kind() == reflect.Struct {
		values = values.Elem()
	} else {
		return fmt.Errorf("expected a pointer to a struct")
	}

	for i := 0; i < values.NumField(); i++ {
		if b.Buff.Len() <= 0 {
			return nil
		}
		field := values.Field(i)
		if !field.CanSet() {
			continue
		}
		switch field.Kind() {
		case reflect.Bool:
			boolByte := b.Buff.Next(1)
			field.SetBool(boolByte[0] != 0)			
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			hasVLQTag := values.Type().Field(i).Tag.Get("vlq")
			if hasVLQTag != "" {
				uintValue, err := b.DecodeVLQ()
				if err != nil {
					return err
				}
				field.SetUint(uint64(uintValue))
			} else {
				size := int(field.Type().Size())
				endianess := values.Type().Field(i).Tag.Get("endian")
				uintValue, err := b.readInteger(field.Kind(), size, endianess)
				if err != nil {
					return err
				}
				field.SetUint(uintValue)
			}
		case reflect.String:
			strLen, _ := binary.ReadUvarint(b.Buff)
			data := make([]byte, strLen)
			b.Buff.Read(data)
			field.SetString(string(data))
		case reflect.Interface:
			if field.IsNil() {
				continue
			}
		default:
			return fmt.Errorf("unsupported type %s for field %s", field.Type(), values.Type().Field(i).Name)
		}
	}

	return nil
}

func (b *byter) DecodeVLQ() (int, error) {
    multiplier := 1
    value := 0
    var encodedByte byte
    var err error

    for {
        encodedByte, err = b.Buff.ReadByte()
        if err != nil {
            return 0, err
        }
        value += int(encodedByte&0x7F) * multiplier
        if multiplier > 2097152 {
            return 0, errors.New("malformed vlq encoded data")
        }
        multiplier *= 128
        if (encodedByte & 0x80) == 0 {
            break
        }
    }
    return value, nil
}