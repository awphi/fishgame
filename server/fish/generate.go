package fish

import (
	"github.com/mroth/weightedrand/v2"
)

type Fish struct {
	Type   *FishDefinition
	Weight float64
}

func createFishGenerator(registry map[string]FishDefinition, temp float64) (*weightedrand.Chooser[FishDefinition, uint], error) {
	fishChoices := make([]weightedrand.Choice[FishDefinition, uint], len(registry))
	i := 0
	for _, fish := range registry {
		weight := uint(float64(fish.Rarity) * fish.WaterTemperatureDistribution.Sample(temp))
		fishChoices[i] = weightedrand.NewChoice(fish, weight)
		i++
	}

	return weightedrand.NewChooser(fishChoices...)

}

func GenerateFish(temp float64) (Fish, error) {
	result := Fish{}
	fishGenerator, err := createFishGenerator(fishRegistry, temp)

	if err != nil {
		return result, err
	}

	def := fishGenerator.Pick()
	result.Weight = def.WeightDistribution.Random()
	result.Type = &def

	return result, nil
}
