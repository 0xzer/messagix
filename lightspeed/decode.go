package lightspeed

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type LightSpeedDecoder struct {
	Table interface{} // struct that contains pointers to all the dependencies/stores
	Dependencies map[string]string
	StatementReferences map[int]int64
}

func NewLightSpeedDecoder(dependencies map[string]string, table interface{}) *LightSpeedDecoder {
	return &LightSpeedDecoder{
		Table: table,
		Dependencies: dependencies,
		StatementReferences: make(map[int]int64),
	}
}

func (ls *LightSpeedDecoder) Decode(data interface{}) interface{} {
	s, ok := data.([]interface{})
	if !ok {
		return data
	}
	
	stepType := StepType(int(s[0].(float64)))
	stepData := s[1:]
	switch stepType {
	case BLOCK:
		for _, blockData := range stepData {
			stepDataArr := blockData.([]interface{})
			ls.Decode(stepDataArr)
		}
	case LOAD:
		key, ok := stepData[0].(float64)
		if !ok {
			log.Println("[LOAD] failed to store key to float64")
			return false
		}
		
		shouldLoad, ok := ls.StatementReferences[int(key)]
		if !ok {
			log.Println("[LOAD] failed to fetch statement reference for key:", key)
			return false
		}
		return shouldLoad
	case STORE:
		retVal := ls.Decode(stepData[1])
		ls.StatementReferences[int(stepData[0].(float64))] = retVal.(int64)
		//log.Println(ls.StatementReferences)
	case STORE_ARRAY:
		key, ok := stepData[0].(float64)
		if !ok {
			log.Println(stepData...)
			os.Exit(1)
		}

		shouldStore, ok := stepData[1].(float64)
		if !ok {
			log.Println(stepData...)
			os.Exit(1)
		}

		ls.StatementReferences[int(key)] = int64(shouldStore)
		ls.Decode(s[2:])
	case CALL_STORED_PROCEDURE:
		referenceName := stepData[0].(string)
		ls.handleStoredProcedure(referenceName, stepData[1:])
	case UNDEFINED:
		return nil
	case I64_FROM_STRING:
		i64, err := strconv.ParseInt(stepData[0].(string), 10, 64)
		if err != nil {
			log.Println("[I64_FROM_STRING] failed to convert string to int64:", err.Error())
			return 0
		}
		return i64
	case IF:
		statement := stepData[0]
		result := ls.Decode(statement).(int64)
		if result > 0 {
			ls.Decode(stepData[1])
		} else if len(stepData) >= 3 {
			if stepData[2] != nil {
				ls.Decode(stepData[2])
			}
		}
	default:
		log.Println("got unknown step type:", stepType)
		os.Exit(1)
	}

	return nil
}

func (ls *LightSpeedDecoder) handleStoredProcedure(referenceName string, data []interface{}) {
	depReference, ok := ls.Dependencies[referenceName]
	if !ok {
		log.Println("Skipping dependency with reference name:",referenceName, data)
		return
	}

	reflectedMs := reflect.ValueOf(ls.Table).Elem()
	//log.Println(depReference)
	depField := reflectedMs.FieldByName(depReference)
	
	if !depField.IsValid() {
		log.Println("Skipping dependency with reference name:",referenceName, data)
		return
	}
	var err error
	// get the Type of the elements of the slice
	depFieldsType := depField.Type().Elem()
	
	// create a new instance of the underlying type
	newDepInstance := reflect.New(depFieldsType).Elem()
	for i := 0; i < depFieldsType.NumField(); i++ {
		fieldInfo := depFieldsType.Field(i)
		var index int
		conditionField := fieldInfo.Tag.Get("conditionField")
		if conditionField != "" {
			indexChoices := fieldInfo.Tag.Get("indexes")
			conditionVal := newDepInstance.FieldByName(conditionField)
			index, err = ls.parseConditionIndex(conditionVal.Bool(), indexChoices)
			if err != nil {
				log.Println(fmt.Sprintf("failed to parse condition index in dependency %v for field %v", depFieldsType.Name(), fieldInfo.Name))
				continue
			}
		} else {
			index, _ = strconv.Atoi(fieldInfo.Tag.Get("index"))
		}
		
		kind := fieldInfo.Type.Kind()
		val := ls.Decode(data[index])
		if val == nil { // skip setting field, because the index in the array was [9] which is undefined.
			continue
		}
		
		switch kind {
		case reflect.Int64:
			i64, ok := val.(int64)
			if !ok {
				log.Println(fmt.Sprintf("failed to set int64 to %v in dependency %v for field %v", val, depFieldsType.Name(), fieldInfo.Name))
				continue
			}
			newDepInstance.Field(i).SetInt(i64)
		case reflect.String:
			str, ok := val.(string)
			if !ok {
				log.Println(fmt.Sprintf("failed to set string to %v in dependency %v for field %v", val, depFieldsType.Name(), fieldInfo.Name))
				continue
			}
			newDepInstance.Field(i).SetString(str)
		case reflect.Interface:
			if val == nil {
				continue
			}
			log.Println(val)
			newDepInstance.Field(i).Set(reflect.ValueOf(val))
		case reflect.Bool:
			boolean, ok := val.(bool)
			if !ok {
				log.Println(fmt.Sprintf("failed to set bool to %v in dependency %v for field %v", val, depFieldsType.Name(), fieldInfo.Name))
				continue
			}
			newDepInstance.Field(i).SetBool(boolean)
		default:
			log.Println("invalid kind:", kind, val)
			os.Exit(1)
		}
	}
	newSlice := reflect.Append(depField, newDepInstance)
	depField.Set(newSlice)
}

// conditionVal ? trueIndex : falseIndex
func (ls *LightSpeedDecoder) parseConditionIndex(val bool, choices string) (int, error) {
	indexes := strings.Split(choices, ",")
	var index int
	var err error
	if val {
		index, err = strconv.Atoi(indexes[0])
	} else {
		index, err = strconv.Atoi(indexes[1])
	}
	return index, err
}