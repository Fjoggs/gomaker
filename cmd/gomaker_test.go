package gomaker

import (
	"reflect"
	"testing"
)

func TestGomaker(t *testing.T) {
	main()
}

func TestReadMap(t *testing.T) {
	mapName := "test.map"
	expected := Materials{map[string]int{"testmap/test_texture_3.tga": 1, "testmap/test_texture.jpg": 1, "testmap/test_shader_2.tga": 1, "testmap/test_shader_3.jpg": 1, "testmap/test_model_texture_1.jpg": 1, "testmap/test_model_texture_2.tga": 1}, map[string]int{"testmap/test_texture_3": 2, "testmap/test_texture": 2}}
	expectedSounds := map[string]int{"sound/testmap/sound-file.wav": 1}
	expectedShaderNames := []string{"testmap/test_shader_2", "testmap/test_shader"}
	actual, actualSounds, actualShaderNames := readMap(mapName, "resources")

	if !reflect.DeepEqual(actual, expected.textures) {
		t.Errorf("Expected %v got %v", expected.textures, actual)
	}

	if !reflect.DeepEqual(actualSounds, expectedSounds) {
		t.Errorf("Expected %v got %v", expectedSounds, actual)
	}

	if !isEqual(actualShaderNames, expectedShaderNames) {
		t.Errorf("Expected %v got %v", expectedShaderNames, actualShaderNames)
	}
}

func TestAddMaterials(t *testing.T) {
	line := "( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap/test_shader_3 32 0 0 0.5 0.5 134217728 0 0"
	materials := map[string]int{"testmap/test_texture_3": 1, "testmap/test_texture": 1}
	expected := map[string]int{"testmap/test_texture_3": 1, "testmap/test_texture": 1, "testmap/test_shader_3": 1}

	addMaterials(line, materials)
	if !reflect.DeepEqual(materials, expected) {
		t.Errorf("Expected %v got %v", expected, materials)
	}
}

func TestAddSounds(t *testing.T) {
	line := `"noise" "sound/testmap/sound-file.wav"`
	sounds := map[string]int{"sound/testmap/sound-file-2.wav": 1}
	expected := map[string]int{"sound/testmap/sound-file-2.wav": 1, "sound/testmap/sound-file.wav": 1}

	addSounds(line, sounds)
	if !reflect.DeepEqual(sounds, expected) {
		t.Errorf("Expected %v got %v", expected, sounds)
	}
}

func TestGetMaterials(t *testing.T) {
	emptyMap := map[string]int{}
	tests := []struct {
		input    string
		expected map[string]int
	}{
		{"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap/texture 32 0 0 0.5 0.5 134217728 0 0", map[string]int{"testmap/texture": 1}},
		{"// Entity 0", emptyMap},
		{"{", emptyMap},
		{`"classname" "misc_model"`, emptyMap},
		{`"origin" "-924 -4 536"`, emptyMap},
		{`"model" "resources/models/test-model.ase"`, emptyMap},
		{`"angles" "-0 0 -180"`, emptyMap},
		{`"_remap" "*;textures/testmap/test_texture"`, emptyMap},
		{"}", map[string]int{"testmap/test_texture": 1}},
		{"// Brush 1337", emptyMap},
		{"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0", map[string]int{"testmap/test_texture": 1}},
		{"// entity 1", emptyMap},
		{"{", emptyMap},
		{"}", emptyMap},
		{"// brush 0", emptyMap},
		{"{", emptyMap},
		{"}", emptyMap},
		{"// Entity 2", emptyMap},
		{"{", emptyMap},
		{`"classname" "misc_model"`, emptyMap},
		{`"origin" "-924 -4 536"`, emptyMap},
		{`"model" "resources/models/test-model.ase"`, emptyMap},
		{"}", map[string]int{"testmap/test_model_texture_1": 1}},
		{"// Entity 3", emptyMap},
		{"{", emptyMap},
		{`"classname" "misc_model"`, emptyMap},
		{`"origin" "-924 -4 536"`, emptyMap},
		{`"model" "resources/models/test-material.obj"`, emptyMap},
		{"}", map[string]int{"testmap/test_model_texture_2": 1}},
	}
	for index, test := range tests {
		actual := getMaterials(test.input)

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v got %v for %s, index %d", test.expected, actual, test.input, index)
		}
	}
}

func TestHandleBrush(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap/texture 32 0 0 0.5 0.5 134217728 0 0", "testmap/texture"},
		{"// Entity 0", ""},
		{"// Brush 1337", ""},
		{"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0", "testmap/test_texture"},
		{"// entity 1", ""},
		{"// brush 0", ""},
	}
	for _, test := range tests {
		actual := handleBrush(test.input)
		if actual != test.expected {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
		}
	}
}

func TestHandleEntity(t *testing.T) {
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
			`"model" "maps/models/test-model.ase"`,
			`"angles" "-0 0 -180"`,
			`"_remap" "*;textures/testmap/test_texture"`,
			"}",
		}, map[string]int{"testmap/test_texture": 1}},
	}
	for _, test := range tests {
		actual := handleEntity(test.input)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
		}
	}
}

func TestIsClosingBracket(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"{", false},
		{"}", true},
		{")", false},
		{"// Entity 0", false},
		{"", false},
	}
	for _, test := range tests {
		actual := isClosingBracket(test.input)
		if actual != test.expected {
			t.Errorf("Expected %v got %v for %s", test.expected, actual, test.input)
		}
	}
}

func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMergeMaps(t *testing.T) {
	tests := []struct {
		source      map[string]int
		destination map[string]int
		expected    map[string]int
	}{
		{
			map[string]int{"testmap/test_model_texture_3": 1},
			map[string]int{"testmap/test_model_texture_1": 1, "testmap/test_model_texture_2": 1},
			map[string]int{"testmap/test_model_texture_1": 1, "testmap/test_model_texture_2": 1, "testmap/test_model_texture_3": 1},
		},
		{
			map[string]int{"testmap/test_model_texture_1": 1, "testmap/test_model_texture_2": 1},
			map[string]int{"testmap/test_model_texture_3": 1},
			map[string]int{"testmap/test_model_texture_1": 1, "testmap/test_model_texture_2": 1, "testmap/test_model_texture_3": 1},
		},
	}

	for _, test := range tests {
		mergeMaps(test.source, test.destination)
		if !reflect.DeepEqual(test.destination, test.expected) {
			t.Errorf("Expected %v got %v for %v", test.expected, test.destination, test)
		}
	}
}
