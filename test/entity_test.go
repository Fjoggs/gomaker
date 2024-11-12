package builder

import (
	"reflect"
	"testing"
)

func TestIsEntity(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) common/caulk 32 0 0 0.5 0.5 134217728 0 0",
			false,
		},
		{"// Entity 0", true},
		{"// Brush 1337", false},
		{
			"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0",
			false,
		},
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
		expected map[string]int
	}{
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"origin" "-924 -4 536"`,
			`"model" "resources/models/test-model.ase"`,
			`"angles" "-0 0 -180"`,
			"}",
		}, map[string]int{"testmap/test_model_texture_1": 1}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"origin" "-924 -4 536"`,
			`"model" "resources/models/test-model-2.ase"`,
			`"angles" "-0 0 -180"`,
			"}",
		}, map[string]int{"texture_test/concrete_tile": 1, "texture_test/texture-2": 1}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"model" "resources/models/test-material.obj"`,
			"}",
		}, map[string]int{"testmap/test_model_texture_2": 1}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"model" "resources/models/test-material-2.obj"`,
			"}",
		}, map[string]int{"texture_test/concrete_tile": 1}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"model" "resources/models/test-material-2.obj"`,
			`"_remap" ":*;textures/test_texture/texture-2"`,
			"}",
		}, map[string]int{"test_texture/texture-2": 1}},
		{[]string{
			"{",
			`"classname" "misc_model"`,
			`"origin" "-924 -4 536"`,
			`"model" "maps/models/test-model.ase"`,
			`"angles" "-0 0 -180"`,
			`"_remap" "*;textures/testmap/test_texture"`,
			"}",
		}, map[string]int{"testmap/test_texture": 1}},
		{[]string{
			"{",
			`"classname" "worldspawn"`,
			`"message" "Test map"`,
			`"ambient" "10"`,
			"}",
		}, map[string]int{}},
		{[]string{"{", "}"}, map[string]int{}},
	}
	for _, test := range tests {
		actual := parseEntity(test.input)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test.input)
		}
	}
}

func TestModelPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"model" "resources/models/test-model.ase"`, "resources/models/test-model.ase"},
		{`"model" "resources/models/test-model-2.ase"`, "resources/models/test-model-2.ase"},
		{`"model" "resources/models/test-model-3.obj"`, "resources/models/test-model-3.obj"},
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
		expected map[string]int
	}{
		{"resources/models/test-model.ase", map[string]int{"testmap/test_model_texture_1": 1}},
		{
			"resources/models/test-model-2.ase",
			map[string]int{"texture_test/concrete_tile": 1, "texture_test/texture-2": 1},
		},
		{"resources/models/test-material.mtl", map[string]int{"testmap/test_model_texture_2": 1}},
		{"resources/models/test-material-2.mtl", map[string]int{"texture_test/concrete_tile": 1}},
	}
	for _, test := range tests {
		actual := parseModel(test.path)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v got %v", test.expected, actual)
		}
	}
}

func TestObjTexture(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			`map_Kd /long/path/for/some/reason/textures/testmap/test_texture.jpg`,
			"testmap/test_texture",
		},
		{
			`map_Kd \slash\wrong\way\textures\texture_test\concrete_tile.jpg`,
			"texture_test/concrete_tile",
		},
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
