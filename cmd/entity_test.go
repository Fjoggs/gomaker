package gomaker

import (
	"testing"
)

func TestIsEntity(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) common/caulk 32 0 0 0.5 0.5 134217728 0 0", false},
		{"// Entity 0", true},
		{"// Brush 1337", false},
		{"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0", false},
		{"// entity 1", true},
		{"// brush 0", false},
	}
	for _, test := range tests {
		value := isEntity(test.input)
		if value != test.expected {
			t.Errorf("Expected %v got %v for %v", test.expected, value, test)
		}
	}
}

func TestParseEntity(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"origin" "-924 -4 536"`,
			`"model" "resources/test-model.ase"`,
			`"angles" "-0 0 -180"`,
			"}",
		}, []string{"testmap/test_texture"}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"origin" "-924 -4 536"`,
			`"model" "resources/test-model-2.ase"`,
			`"angles" "-0 0 -180"`,
			"}",
		}, []string{"texture_test/concrete_tile", "texture_test/texture-2"}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"model" "resources/test-material.obj"`,
			"}",
		}, []string{"testmap/test_texture"}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"model" "resources/test-material-2.obj"`,
			"}",
		}, []string{"texture_test/concrete_tile"}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"model" "resources/test-material-2.obj"`,
			`"_remap" ":*;textures/test_texture/texture-2"`,
			"}",
		}, []string{"test_texture/texture-2"}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"origin" "-924 -4 536"`,
			`"model" "maps/models/test-model.ase"`,
			`"angles" "-0 0 -180"`,
			`"_remap" "*;textures/testmap/test_texture"`,
			"}",
		}, []string{"testmap/test_texture"}},
		{[]string{
			"{",
			`"classname" "worldspawn"`,
			`"message" "Test map"`,
			`"ambient" "10"`,
			"}",
		}, []string{}},
		{[]string{"{", "}"}, []string{}},
	}
	for _, test := range tests {
		actual := parseEntity(test.input)
		if !isEqual(actual, test.expected) {
			t.Errorf("Expected %s got %s for %v", test.expected, actual, test.input)
		}
	}

}

func TestModelPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"model" "resources/test-model.ase"`, "resources/test-model.ase"},
		{`"model" "resources/test-model-2.ase"`, "resources/test-model-2.ase"},
		{`"model" "resources/test-model-3.obj"`, "resources/test-model-3.obj"},
		{`"model" "maps/models/test-model.ase"`, "maps/models/test-model.ase"},
	}
	for _, test := range tests {
		actual := modelPath(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %s for %v", test.expected, actual, test)
		}
	}
}
func TestRemapTexture(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{", ""},
		{`"classname" "misc_model"`, ""},
		{`"origin" "-924 -4 536"`, ""},
		{`"angles" "-0 0 -180"`, ""},
		{`"_remap" "*;textures/testmap/test_texture"`, "testmap/test_texture"},
		{"}", ""},
	}
	for _, test := range tests {
		actual := remapTexture(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %s for %v", test.expected, actual, test)
		}
	}
}

func TestParseModel(t *testing.T) {
	tests := []struct {
		path     string
		expected []string
	}{
		{"resources/test-model.ase", []string{"testmap/test_texture"}},
		{"resources/test-model-2.ase", []string{"texture_test/concrete_tile", "texture_test/texture-2"}},
		{"resources/test-material.mtl", []string{"testmap/test_texture"}},
		{"resources/test-material-2.mtl", []string{"texture_test/concrete_tile"}},
	}
	for _, test := range tests {
		actual := parseModel(test.path)
		if !isEqual(actual, test.expected) {
			t.Errorf("Expected %s got %s", test.expected, actual)
		}
	}
}

func TestObjTexture(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{`map_Kd /long/path/for/some/reason/textures/testmap/test_texture.jpg`, "testmap/test_texture"},
		{`map_Kd \slash\wrong\way\textures\texture_test\concrete_tile.jpg`, "texture_test/concrete_tile"},
	}
	for _, test := range tests {
		actual := objTexture(test.path)
		if actual != test.expected {
			t.Errorf("Expected %s got %s for %v", test.expected, actual, test)
		}
	}
}

func TestAseTexture(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{`*BITMAP "//../textures/testmap/test_texture.jpg"`, "testmap/test_texture"},
		{`*BITMAP	"..\textures\texture_test\concrete_tile.tga"`, "texture_test/concrete_tile"},
	}
	for _, test := range tests {
		actual := aseTexture(test.path)
		if actual != test.expected {
			t.Errorf("Expected %s got %s for %v", test.expected, actual, test)
		}
	}
}
