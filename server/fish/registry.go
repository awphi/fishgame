package fish

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/awphi/fishgame/util"
)

type ResourceDefinition struct {
	Name string
}

type FishResourceMap = map[*ResourceDefinition]float64

type FishDefinition struct {
	Name                         string
	Value                        uint
	Rarity                       uint
	WaterTemperatureDistribution util.Distribution
	WeightDistribution           util.Distribution
	Resource                     FishResourceMap
}

type rawFishDefinition struct {
	Resource                     map[string]float64
	WaterTemperatureDistribution struct{ Type string }
	WeightDistribution           struct{ Type string }
}

func createResourceRegistry(data []byte) map[string]ResourceDefinition {
	reg := map[string]ResourceDefinition{}
	json.Unmarshal(data, &reg)
	return reg
}

func createFishRegistry(data []byte) map[string]FishDefinition {
	reg := map[string]FishDefinition{}
	objmap := map[string]json.RawMessage{}

	json.Unmarshal(data, &objmap)

	for id, fishData := range objmap {
		fishDef := FishDefinition{Resource: map[*ResourceDefinition]float64{}}
		rawFishDef := rawFishDefinition{}

		json.Unmarshal(objmap[id], &rawFishDef)

		// copy valid resources in, using references to the resource registry
		for resourceId, value := range rawFishDef.Resource {
			resource, has := resourceRegistry[resourceId]
			if !has {
				fmt.Printf("resource '%s' is invalid (in fish '%s'), skipping resource... \n", resourceId, id)
				continue
			}
			fishDef.Resource[&resource] = value
		}

		// construct distributions of proper types based on string "type" prop
		waterTemperatureDist, err1 := util.NewDistribution(rawFishDef.WaterTemperatureDistribution.Type)
		weightDist, err2 := util.NewDistribution(rawFishDef.WeightDistribution.Type)

		if err1 != nil || err2 != nil {
			fmt.Printf("failed to load fish '%s' due to unrecognised weight or temperature distribution, skipping...\n", id)
			continue
		}

		fishDef.WeightDistribution = weightDist
		fishDef.WaterTemperatureDistribution = waterTemperatureDist

		// finally unmarshal everything simple in
		json.Unmarshal(fishData, &fishDef)
		reg[id] = fishDef
	}

	return reg
}

//go:embed resources.json
var resourceJsonData []byte
var resourceRegistry = createResourceRegistry(resourceJsonData)

//go:embed fish.json
var fishJsonData []byte
var fishRegistry = createFishRegistry(fishJsonData)
