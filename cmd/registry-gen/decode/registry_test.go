package decode

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeProducts(t *testing.T) {
	jsonInput := `
[
  {
    "vid": 1,
    "name": "LIFX",
    "defaults": {
      "hev": false,
      "color": false,
      "chain": false,
      "matrix": false,
      "relays": false,
      "buttons": false,
      "infrared": false,
      "multizone": false,
      "temperature_range": null,
      "extended_multizone": false
    },
    "products": [
      {
        "pid": 1,
        "name": "LIFX Original 1000",
        "features": {
          "color": true,
          "chain": false,
          "matrix": false,
          "infrared": false,
          "multizone": false,
          "temperature_range": [2500, 9000]
        },
        "upgrades": [
          {
            "major": 2,
            "minor": 80,
            "features": {
              "temperature_range": [1500, 9000]
            }
          }
        ]
      }
    ]
  }
]
`

	ps, err := DecodeProductsRegistry([]byte(jsonInput))
	require.NoError(t, err)
	require.NotNil(t, ps)

	// Validate Enums
	require.Len(t, ps, 1)
	require.Equal(t, 1, ps[0].PID)
	require.Equal(t, "LIFX Original 1000", ps[0].Name)
	require.Equal(t, true, ps[0].Features.Color)
	require.Equal(t, false, ps[0].Features.Chain)
	require.Equal(t, false, ps[0].Features.Matrix)
	require.Equal(t, false, ps[0].Features.Infrared)
	require.Equal(t, false, ps[0].Features.Multizone)
	require.Len(t, ps[0].Features.TemperatureRange, 2)
	require.Equal(t, 2500, ps[0].Features.TemperatureRange[0])
	require.Equal(t, 9000, ps[0].Features.TemperatureRange[1])
	require.Len(t, ps[0].Upgrades, 1)
	require.Equal(t, 2, ps[0].Upgrades[0].Major)
	require.Equal(t, 80, ps[0].Upgrades[0].Minor)
	require.Len(t, ps[0].Upgrades[0].Features.TemperatureRange, 2)
	require.Equal(t, 1500, ps[0].Upgrades[0].Features.TemperatureRange[0])
	require.Equal(t, 9000, ps[0].Upgrades[0].Features.TemperatureRange[1])
}

func TestDecodeUplightCoordinates(t *testing.T) {
	for _, test := range []struct {
		name, fields string
		want         *UplightCoordinates
	}{
		{"missing", ``, nil},
		{"null", `"uplight_coords":null`, nil},
		{"origin", `"uplight_coords":{"x":0,"y":0}`, &UplightCoordinates{X: 0, Y: 0}},
		{"ceiling", `"uplight_coords":{"x":7,"y":7}`, &UplightCoordinates{X: 7, Y: 7}},
		{"capsule", `"uplight_coords":{"x":15,"y":7}`, &UplightCoordinates{X: 15, Y: 7}},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := fmt.Sprintf(`[{"vid":1,"products":[{"pid":176,"features":{%s},"upgrades":[{"major":4,"minor":110,"features":{%s}}]}]}]`, test.fields, test.fields)
			products, err := DecodeProductsRegistry([]byte(input))
			require.NoError(t, err)
			require.Equal(t, test.want, products[0].Features.UplightCoords)
			require.Equal(t, test.want, products[0].Upgrades[0].Features.UplightCoords)
			data, err := json.Marshal(products[0].Features)
			require.NoError(t, err)
			if test.want == nil {
				require.NotContains(t, string(data), "uplight_coords")
			} else {
				require.Contains(t, string(data), "uplight_coords")
			}
			var features FeatureSet
			require.NoError(t, json.Unmarshal(data, &features))
			require.Equal(t, test.want, features.UplightCoords)
		})
	}
}

func TestDecodeCheckedInUplightProducts(t *testing.T) {
	data, err := os.ReadFile("../src/products.json")
	require.NoError(t, err)
	products, err := DecodeProductsRegistry(data)
	require.NoError(t, err)
	want := map[int]UplightCoordinates{
		176: {X: 7, Y: 7}, 177: {X: 7, Y: 7},
		201: {X: 15, Y: 7}, 202: {X: 15, Y: 7},
		265: {X: 7, Y: 7}, 266: {X: 7, Y: 7},
	}
	for _, p := range products {
		if coords, ok := want[p.PID]; ok {
			require.Equal(t, &coords, p.Features.UplightCoords, "pid %d", p.PID)
			delete(want, p.PID)
		}
	}
	require.Empty(t, want, "uplight fixtures missing from checked-in registry")
}
