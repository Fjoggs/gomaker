package gomaker

import (
	"reflect"
	"testing"
)

func TestExtractTexturesFromUsedShaders(t *testing.T) {
	input := map[string]int{"testmap/test_texture_3": 1, "testmap/test_shader": 1, "testmap/test_texture": 1, "testmap/test_shader_2": 1}
	expectedTextures := map[string]int{
		"testmap/test_texture_3": 1,
		"testmap/test_shader_2":  1,
		"testmap/test_shader_3":  1,
		"testmap/test_texture":   1,
		"testmap/test_shader_4":  1,
		"testmap/test_shader_5":  1,
	}
	expectedShaderNames := []string{"testmap/test_shader_2", "testmap/test_shader"}
	expectedShaderFiles := []string{"test_shader_2.shader", "testmap.shader"}
	actual, actualShaderNames, actualShaderFiles := extractTexturesFromUsedShaders(input, "resources/scripts")

	if !reflect.DeepEqual(actual, expectedTextures) {
		t.Errorf("Expected textures %v got %v for %v", expectedTextures, actual, input)
	}

	if !isEqual(actualShaderNames, expectedShaderNames) {
		t.Errorf("Expected shader names %v got %v for %v", expectedShaderNames, actualShaderNames, input)
	}

	if !isEqual(actualShaderFiles, expectedShaderFiles) {
		t.Errorf("Expected shader files %v got %v for %v", expectedShaderFiles, actualShaderFiles, input)
	}
}

func TestCombineTexturesFromShaders(t *testing.T) {
	textures := map[string]int{}
	shaderNames := []string{}
	input := []Shader{{"testmap/test_shader_2", []string{}, map[string]int{"testmap/test_shader_4": 1, "testmap/test_shader_5": 1}}, {"testmap/test_shader", []string{}, map[string]int{"testmap/test_shader_2": 1, "testmap/test_shader_3": 1}}}
	expectedTextures := map[string]int{"testmap/test_shader_4": 1, "testmap/test_shader_5": 1, "testmap/test_shader_2": 1, "testmap/test_shader_3": 1}
	expectedShaders := []string{"testmap/test_shader_2", "testmap/test_shader"}
	actualTextures, actualShaders := combineTexturesFromShaders(input, textures, shaderNames)

	if !reflect.DeepEqual(actualTextures, expectedTextures) {
		t.Errorf("Expected %v got %v for %v", expectedTextures, actualTextures, input)
	}
	if !isEqual(actualShaders, expectedShaders) {
		t.Errorf("Expected %v got %v for %v", expectedShaders, actualShaders, input)
	}
}

func TestParseShaderFile(t *testing.T) {
	tests := []struct {
		shadersFromMapFile map[string]int
		shaderFileName     string
		expected           []Shader
	}{
		{map[string]int{"testmap/test_shader_2": 1}, "test_shader_2.shader", []Shader{{"testmap/test_shader_2", []string{}, map[string]int{"testmap/test_shader_4": 1, "testmap/test_shader_5": 1}}}},
		{map[string]int{"testmap/test_shader": 1}, "testmap.shader", []Shader{{"testmap/test_shader", []string{}, map[string]int{"testmap/test_shader_2": 1, "testmap/test_shader_3": 1}}}},
	}

	for _, test := range tests {
		actual := parseShaderFile(test.shadersFromMapFile, test.shaderFileName, "resources/scripts")
		if len(actual) != len(test.expected) {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
		}
		for index, actualValue := range actual {
			if actualValue.name != test.expected[index].name {
				t.Errorf("Expected %v got %s for %v", test.expected[index].name, actualValue.name, test)
			}
			// if !isEqual(actualValue.lines, test.expected[index].lines) {
			// 	t.Errorf("Expected %v got %s for %v", test.expected[index].lines, actualValue.lines, test)
			// }
			equalTextures := reflect.DeepEqual(actualValue.textures, test.expected[index].textures)
			if !equalTextures {
				t.Errorf("Expected %v got %v for %v", test.expected[index].textures, actualValue.textures, test)
			}
		}
	}
}

func TestShaderIsUsed(t *testing.T) {
	tests := []struct {
		shadersFromMapFile map[string]int
		shaderName         string
		expected           bool
	}{
		{map[string]int{"testmap/test_shader": 1}, "testmap/test_shader", true},
		{map[string]int{"textures/testmap/test_shader_2": 1}, "testmap/test_shader", false},
		{map[string]int{"testmap/test_shader_3": 1}, "testmap/test_shader", false},
	}

	for _, test := range tests {
		actual := shaderIsUsed(test.shadersFromMapFile, test.shaderName)
		if actual != test.expected {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
		}
	}
}

func TestIsShaderName(t *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"textures/testmap/test_shader_2", true},
		{"textures/testmap/test_shader_2 {", true},
		{"map textures/testmap/test_shader_5.tga", false},
	}

	for _, test := range tests {
		actual := isShaderName(test.line)
		if actual != test.expected {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
		}
	}
}
