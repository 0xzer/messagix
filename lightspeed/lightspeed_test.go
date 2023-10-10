package lightspeed_test

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/0xzer/messagix"
	"github.com/0xzer/messagix/lightspeed"
	"github.com/0xzer/messagix/table"
)

func TestDecode(t *testing.T) {
	data, err := os.ReadFile("test_data.json")
	if err != nil {
		log.Fatal(err)
	}

	var res *messagix.PublishResponseData
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Fatal(err)
	}

	deps := table.SPToDepMap(res.Sp)
	var lsData *lightspeed.LightSpeedData
	err = json.Unmarshal([]byte(res.Payload), &lsData)
	if err != nil {
		log.Fatal(err)
	}
	
	lsTable := &table.LSTable{}
	lsDecoder := lightspeed.NewLightSpeedDecoder(deps, lsTable)
	lsDecoder.Decode(lsData.Steps)

	tableReflectionTest(lsTable)
}

func tableReflectionTest(loadedTable *table.LSTable) {
	values := reflect.ValueOf(loadedTable).Elem()

	for i := 0; i < values.NumField(); i++ {
		fieldValue := values.Field(i)
		fieldKind := fieldValue.Kind()
		
		if fieldKind == reflect.Slice && fieldValue.Len() > 0 {
			switch data := fieldValue.Interface().(type) {
			case []table.LSUpsertMessage:
				log.Println(data)
			case []table.LSDeleteThenInsertIGContactInfo:
				log.Println(data, fieldValue.Type().Elem().String())
			default:
			}
		}
	}
}