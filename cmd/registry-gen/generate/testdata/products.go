// Code generated. DO NOT EDIT.
// Source: https://github.com/LIFX/public-protocol@0000000000000000000000000000000000000000
// Generated: 2006-01-02T15:04:05Z
package registry

// ProductsByPID maps LIFX Product IDs to products.
var ProductsByPID = map[int]Product{
	1: {
		PID:  1,
		Name: "LIFX Original 1000",
		Features: FeatureSet{
			HEV:               false,
			Color:             true,
			Chain:             false,
			Matrix:            false,
			Relays:            false,
			Buttons:           false,
			Infrared:          false,
			Multizone:         false,
			TemperatureRange:  []int{2500, 9000},
			ExtendedMultizone: false,
		},
		Upgrades: []Upgrade{
			{
				Major: 2,
				Minor: 80,
				Features: FeatureSet{
					HEV:               false,
					Color:             false,
					Chain:             false,
					Matrix:            false,
					Relays:            false,
					Buttons:           false,
					Infrared:          false,
					Multizone:         false,
					TemperatureRange:  []int{1500, 9000},
					ExtendedMultizone: false,
				},
			},
		},
	},
}

// Product represents a LIFX Product and its capabilities.
type Product struct {
	PID      int        `json:"pid"`
	Name     string     `json:"name"`
	Features FeatureSet `json:"features"`
	Upgrades []Upgrade  `json:"upgrades"`
}

// FeatureSet describes the capabilities of a LIFX Product.
type FeatureSet struct {
	HEV               bool  `json:"hev"`
	Color             bool  `json:"color"`
	Chain             bool  `json:"chain"`
	Matrix            bool  `json:"matrix"`
	Relays            bool  `json:"relays"`
	Buttons           bool  `json:"buttons"`
	Infrared          bool  `json:"infrared"`
	Multizone         bool  `json:"multizone"`
	TemperatureRange  []int `json:"temperature_range"`
	ExtendedMultizone bool  `json:"extended_multizone"`
}

// Upgrade describes the features added to an existing LIFX Product
// and available from the given version.
type Upgrade struct {
	Major    int        `json:"major"`
	Minor    int        `json:"minor"`
	Features FeatureSet `json:"features"`
}
