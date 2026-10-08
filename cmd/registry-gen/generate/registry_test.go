package generate

import (
	"embed"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alessio-palumbo/lifxregistry-go/cmd/registry-gen/decode"
)

//go:embed testdata
var testdataFS embed.FS
var testNow = time.Date(2006, time.January, 2, 15, 04, 05, 0, time.UTC)

func TestGenerateProductsRegistry(t *testing.T) {
	products := []decode.Product{
		{
			PID:  1,
			Name: "LIFX Original 1000",
			Features: decode.FeatureSet{
				Color:            true,
				Chain:            false,
				Matrix:           false,
				Infrared:         false,
				Multizone:        false,
				TemperatureRange: []int{2500, 9000},
			},
			Upgrades: []decode.Upgrade{
				{
					Major: 2,
					Minor: 80,
					Features: decode.FeatureSet{
						TemperatureRange: []int{1500, 9000},
					},
				},
			},
		},
	}

	want, err := fs.ReadFile(testdataFS, "testdata/products.go")
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	// Create temp directory for output
	tmpDir, err := os.MkdirTemp("", "testproducts_gen")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sourceCommit := "0000000000000000000000000000000000000000"
	now = func() time.Time { return testNow }
	t.Cleanup(func() { now = time.Now })
	if err := GenerateProductsRegistry(products, tmpDir, sourceCommit); err != nil {
		t.Fatalf("TestGenerateProductsRegistry failed: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(tmpDir, "products.go"))
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("generated output does not match golden file\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGenerateUplightCoordinates(t *testing.T) {
	products := []decode.Product{
		{PID: 176, Name: "LIFX Ceiling", Features: decode.FeatureSet{Matrix: true, UplightCoords: &decode.UplightCoordinates{X: 7, Y: 7}}},
		{PID: 201, Name: "LIFX Ceiling 13x26", Features: decode.FeatureSet{Matrix: true, UplightCoords: &decode.UplightCoordinates{X: 15, Y: 7}}},
		{PID: 999, Name: "Origin fixture", Features: decode.FeatureSet{UplightCoords: &decode.UplightCoordinates{}},
			Upgrades: []decode.Upgrade{{Major: 4, Minor: 110, Features: decode.FeatureSet{UplightCoords: &decode.UplightCoordinates{X: 1, Y: 2}}}}},
		{PID: 1, Name: "No uplight"},
	}
	dir := t.TempDir()
	if err := GenerateProductsRegistry(products, dir, "fixture"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "products.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"&UplightCoordinates{X: 7, Y: 7}", "&UplightCoordinates{X: 15, Y: 7}", "&UplightCoordinates{X: 0, Y: 0}", "&UplightCoordinates{X: 1, Y: 2}", `json:"uplight_coords,omitempty"`} {
		if !strings.Contains(string(data), text) {
			t.Fatalf("missing generated value: %s", text)
		}
	}
	if strings.Count(string(data), "UplightCoords:") != 4 {
		t.Fatal("generated literals must omit absent uplight coordinates")
	}
	// Compile generated output as a standalone package. The generator's decode
	// types must not leak into the public registry code.
	if output, err := exec.Command("go", "test", path).CombinedOutput(); err != nil {
		t.Fatalf("generated registry does not compile: %v\n%s", err, output)
	}
}
