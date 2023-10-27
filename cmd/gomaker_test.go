package gomaker

import (
	"reflect"
	"testing"
)

func TestGomaker(t *testing.T) {
	main()
}

func TestReadMap(t *testing.T) {
	filePath := "resources/test.map"
	expected := Materials{map[string]int{"testmap/test_texture_3.tga": 2, "testmap/test_texture.jpg": 2}, map[string]int{"not/a/texture": 1}}
	actual := readMap(filePath, "resources/textures/")

	equalTextures := reflect.DeepEqual(actual.textures, expected.textures)
	if !equalTextures {
		t.Errorf("Expected %v got %v", expected.textures, actual.textures)
	}
	equalShaders := reflect.DeepEqual(actual.shaders, expected.shaders)
	if !equalShaders {
		t.Errorf("Expected %v got %v", expected.shaders, actual.shaders)
	}

}

func TestGetMaterials(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap/texture 32 0 0 0.5 0.5 134217728 0 0", []string{"testmap/texture"}},
		{"// Entity 0", []string{}},
		{"{", []string{}},
		{`"classname" "misc_model"`, []string{}},
		{`"origin" "-924 -4 536"`, []string{}},
		{`"model" "resources/test-model.ase"`, []string{}},
		{`"angles" "-0 0 -180"`, []string{}},
		{`"_remap" "*;textures/testmap/test_texture"`, []string{}},
		{"}", []string{"testmap/test_texture"}},
		{"// Brush 1337", []string{}},
		{"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0", []string{"testmap/test_texture"}},
		{"// entity 1", []string{}},
		{"{", []string{}},
		{"}", []string{}},
		{"// brush 0", []string{}},
		{"{", []string{}},
		{"}", []string{}},
		{"// Entity 2", []string{}},
		{"{", []string{}},
		{`"classname" "misc_model"`, []string{}},
		{`"origin" "-924 -4 536"`, []string{}},
		{`"model" "resources/test-model.ase"`, []string{}},
		{"}", []string{"testmap/test_texture"}},
		{"// Entity 3", []string{}},
		{"{", []string{}},
		{`"classname" "misc_model"`, []string{}},
		{`"origin" "-924 -4 536"`, []string{}},
		{`"model" "resources/test-material.obj"`, []string{}},
		{"}", []string{"testmap/test_texture"}},
	}
	for index, test := range tests {
		actual := getMaterials(test.input)
		if !isEqual(actual, test.expected) {
			t.Errorf("Expected %v got %s for %s, index %d", test.expected, actual, test.input, index)
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
			`"model" "maps/models/test-model.ase"`,
			`"angles" "-0 0 -180"`,
			`"_remap" "*;textures/testmap/test_texture"`,
			"}",
		}, []string{"testmap/test_texture"}},
	}
	for _, test := range tests {
		actual := handleEntity(test.input)
		if !isEqual(actual, test.expected) {
			t.Errorf("Excpted %v got %s for %v", test.expected, actual, test)
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

func TestSortMaterials(t *testing.T) {
	tests := []struct {
		input    []string
		expected Materials
	}{
		{
			[]string{"testmap/test_texture_3", "not/a/texture", "testmap/test_texture"},
			Materials{map[string]int{"testmap/test_texture_3.tga": 1, "testmap/test_texture.jpg": 1}, map[string]int{"not/a/texture": 1}},
		},
	}

	for _, test := range tests {
		actual := sortMaterials(test.input, "resources/textures/")
		equalTextures := reflect.DeepEqual(actual.textures, test.expected.textures)
		if !equalTextures {
			t.Errorf("Expected %v got %v for %s", test.expected.textures, actual.textures, test.input)
		}
		equalShaders := reflect.DeepEqual(actual.shaders, test.expected.shaders)
		if !equalShaders {
			t.Errorf("Expected %v got %v for %s", test.expected.shaders, actual.shaders, test.input)
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
