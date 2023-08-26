package lightspeed

import (
	"encoding/json"
	"log"
)

type StepType int
const (
	BLOCK                       StepType = 1
	LOAD                        StepType = 2
	STORE                       StepType = 3
	STORE_ARRAY                 StepType = 4
	CALL_STORED_PROCEDURE       StepType = 5
	CALL_NATIVE_TYPE_OPERATION  StepType = 6
	CALL_NATIVE_OPERATION       StepType = 7
	LIST                        StepType = 8
	UNDEFINED                   StepType = 9
	INFINITY                    StepType = 10
	NAN                         StepType = 11
	RETURN                      StepType = 12
	BOOL_TO_STR                 StepType = 13
	BLOBS_TO_STRING             StepType = 14
	BLOBS_OF_STRING             StepType = 15
	TO_BLOB                     StepType = 16
	I64_OF_FLOAT                StepType = 17
	I64_TO_FLOAT                StepType = 18
	I64_FROM_STRING             StepType = 19
	I64_TO_STRING               StepType = 20
	READ_GK                     StepType = 21
	READ_QE                     StepType = 22
	IF                          StepType = 23
	OR                          StepType = 24
	AND                         StepType = 25
	NOT                         StepType = 26
	IS_NULL                     StepType = 27
	ENFORCE_NOT_NULL            StepType = 28
	GENERIC_EQUAL               StepType = 29
	I64_EQUAL                   StepType = 30
	BLOB_EQUAL                  StepType = 31
	GENERIC_NOT_EQUAL           StepType = 32
	I64_NOT_EQUAL               StepType = 33
	BLOB_NOT_EQUAL              StepType = 34
	GENERIC_GREATER_THAN        StepType = 35
	I64_GREATER_THAN            StepType = 36
	BLOB_GREATER_THAN           StepType = 37
	GENERIC_GREATER_THAN_OR_EQUAL StepType = 38
	I64_GREATER_THAN_OR_EQUAL   StepType = 39
	BLOB_GREATER_THAN_OR_EQUAL  StepType = 40
	GENERIC_LESS_THAN           StepType = 41
	I64_LESS_THAN               StepType = 42
	BLOB_LESS_THAN              StepType = 43
	GENERIC_LESS_THAN_OR_EQUAL  StepType = 44
	I64_LESS_THAN_OR_EQUAL      StepType = 45
	BLOB_LESS_THAN_OR_EQUAL     StepType = 46
	THROW                       StepType = 47
	LOG_CONSOLE                 StepType = 48
	LOGGER_LOG                  StepType = 49
	NATIVE_OP_ARRAY_CREATE      StepType = 50
	NATIVE_OP_ARRAY_APPEND      StepType = 51
	NATIVE_OP_ARRAY_GET_SIZE    StepType = 52
	NATIVE_OP_MAP_CREATE        StepType = 53
	NATIVE_OP_MAP_GET           StepType = 54
	NATIVE_OP_MAP_SET           StepType = 55
	NATIVE_OP_MAP_KEYS          StepType = 56
	NATIVE_OP_MAP_DELETE        StepType = 57
	NATIVE_OP_MAP_HAS           StepType = 58
	NATIVE_OP_STR_JOIN          StepType = 59
	NATIVE_OP_CURRENT_TIME      StepType = 60
	NATIVE_OP_JSON_STRINGIFY    StepType = 61
	NATIVE_OP_RNG_NUM           StepType = 62
	NATIVE_OP_LOCALIZATION_SUPPORTED StepType = 63
	NATIVE_OP_LOCALIZATION_SUPPORTED_V2 StepType = 64
	NATIVE_OP_RESOLVE_LOCALIZED  StepType = 65
	NATIVE_OP_RESOLVE_LOCALIZED_V2 StepType = 66
	ADD                         StepType = 68
	I64_ADD                     StepType = 69
	I64_CAST                    StepType = 70
	READ_JUSTKNOB               StepType = 71
	READ_IGGK                   StepType = 72
)

type LightSpeedData struct {
	Name string `json:"name"`
	Steps interface{} `json:"step"`
}

type Dependency struct {
	Name string `json:"name,omitempty"`
	Value DependencyValue `json:"value,omitempty"`
}

type DependencyValue struct {
	ReferenceName string `json:"__dr,omitempty"`
}

func DependenciesToMap(dep interface{}) map[string]string {
	var converted []Dependency
	b, err := json.Marshal(dep)
	if err != nil {
		log.Fatal(err)
	}

	err  = json.Unmarshal(b, &converted)
	if err != nil {
		log.Fatal(err)
	}

	depMap := make(map[string]string, 0)
	for _, d := range converted {
		depMap[d.Name] = d.Value.ReferenceName
	}

	return depMap
}