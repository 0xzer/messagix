package byter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
)

func NewWriter() *byter {
	return &byter{
		Buff: bytes.NewBuffer(make([]byte, 0)),
	}
}

func (b *byter) writeInteger(value uint64, kind reflect.Kind, endian string) error {
	if kind == reflect.Uint8 {
		return b.Buff.WriteByte(byte(value))
	}

	if endian == "" {
		endian = "big"
	}
	switch endian {
		case "big":
			switch kind {
			case reflect.Uint16:
				return binary.Write(b.Buff, binary.BigEndian, uint16(value))
			case reflect.Uint32:
				return binary.Write(b.Buff, binary.BigEndian, uint32(value))
			case reflect.Uint64:
				return binary.Write(b.Buff, binary.BigEndian, value)
			}
		case "little":
			switch kind {
			case reflect.Uint16:
				return binary.Write(b.Buff, binary.LittleEndian, uint16(value))
			case reflect.Uint32:
				return binary.Write(b.Buff, binary.LittleEndian, uint32(value))
			case reflect.Uint64:
				return binary.Write(b.Buff, binary.LittleEndian, value)
			}
		default:
			return fmt.Errorf("received unsupported endianness while trying to write %v: %s", kind, endian)
	}

	return nil
}

func (b *byter) WriteFromStruct(s interface{}) ([]byte, error) {
	values := reflect.ValueOf(s)
	if values.Kind() == reflect.Ptr && values.Elem().Kind() == reflect.Struct {
		values = values.Elem()
	} else {
		return nil, fmt.Errorf("expected a struct")
	}

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		switch field.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if b.isEnum(field) {
				enumValue := uint8(field.Uint()) // TO-DO: support all int vals
				b.Buff.WriteByte(enumValue)
				continue
			}
			hasVLQTag := values.Type().Field(i).Tag.Get("vlq")
			if hasVLQTag != "" {
				err := b.EncodeVLQ(int(field.Uint())) // Convert to int because our VLQ function takes an int
				if err != nil {
					return nil, err
				}
				continue
			} else {
				endianess := values.Type().Field(i).Tag.Get("endian")
				err := b.writeInteger(field.Uint(), field.Kind(), endianess)
				if err != nil {
					return nil, err
				}
			}
		case reflect.String:
			str := field.String()
			f := values.Type().Field(i)
			lengthType := f.Tag.Get("lengthType")
			endianess := f.Tag.Get("endian")
			err := b.writeString(str, lengthType, endianess)
			if err != nil {
				return nil, err
			}
		/*
		case reflect.Struct:
			sBytes, err := b.WriteFromStruct(field)
			if err != nil {
				return nil, err
			}
			b.Buff.Write(sBytes)
		*/
		default:
			log.Printf("Unsupported type %s for field %s\n", field.Type(), values.Type().Field(i).Name)
		}
	}

	return b.Buff.Bytes(), nil
}

func (b *byter) writeString(s string, lengthType string, endianess string) error {
	if endianess == "" {
		endianess = "big"
	}

	b.writeInteger(uint64(len(s)), stringLengthTags[lengthType], endianess)
	_, err := b.Buff.Write([]byte(s))
	return err
}

func (b *byter) EncodeVLQ(value int) error {
    var encodedByte byte
    for {
        encodedByte = byte(value & 0x7F)
        value >>= 7
        if value > 0 {
            encodedByte |= 0x80
        }

        err := b.Buff.WriteByte(encodedByte)
        if err != nil {
            return err
        }

        if value == 0 {
            break
        }
    }

    return nil
}