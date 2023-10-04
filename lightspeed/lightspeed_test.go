package lightspeed_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/0xzer/messagix"
	"github.com/0xzer/messagix/lightspeed"
	"github.com/0xzer/messagix/table"
)

func TestDecode(t *testing.T) {
	data, _ := os.ReadFile("test_data.json")
	var res *messagix.PublishResponseData
	err := json.Unmarshal(data, &res)
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
}