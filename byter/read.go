package byter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"
)

func NewReader(data []byte) *byter {
	return &byter{
		Buff: bytes.NewBuffer(data),
	}
}

func (b *byter) readInteger(kind reflect.Kind, size int, endian string) uint64 {
	buffer := make([]byte, size)
	_, err := b.Buff.Read(buffer)
	if err != nil {
		log.Printf("Error reading integer: %v\n", err)
		return 0
	}

	if kind == reflect.Uint8 {
		return uint64(buffer[0])
	}

	switch endian {
		case "big":
			switch kind {
			case reflect.Uint16:
				return uint64(binary.BigEndian.Uint16(buffer))
			case reflect.Uint32:
				return uint64(binary.BigEndian.Uint32(buffer))
			case reflect.Uint64:
				return binary.BigEndian.Uint64(buffer)
			}
		case "little":
			switch kind {
			case reflect.Uint16:
				return uint64(binary.LittleEndian.Uint16(buffer))
			case reflect.Uint32:
				return uint64(binary.LittleEndian.Uint32(buffer))
			case reflect.Uint64:
				return binary.LittleEndian.Uint64(buffer)
			}
	}

	log.Printf("Unsupported endianness: %s\n", endian)
	return 0
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
				if endianess == "" {
					endianess = "big"
				}
				uintValue := b.readInteger(field.Kind(), size, endianess)
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
            return 0, errors.New("malformed remaining length")
        }
        multiplier *= 128
        if (encodedByte & 0x80) == 0 {
            break
        }
    }
    return value, nil
}