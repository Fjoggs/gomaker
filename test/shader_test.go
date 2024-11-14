package test

import (
	"reflect"
	"testing"

	"gomaker/internal/shader"
)

func TestExtractTexturesFromUsedShaders(t *testing.T) {
	input := map[string]int{
		"testmap/test_texture_3": 1,
		"testmap/test_shader":    1,
		"testmap/test_texture":   1,
		"testmap/test_shader_2":  1,
	}
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
	actual, actualShaderNames, actualShaderFiles := shader.ExtractTexturesFromUsedShaders(
		input,
		"data/baseq3/scripts",
	)

	if !reflect.DeepEqual(actual, expectedTextures) {
		t.Errorf("Expected textures %v got %v for %v", expectedTextures, actual, input)
	}

	if !isEqual(actualShaderNames, expectedShaderNames) {
		t.Errorf(
			"Expected shader names %v got %v for %v",
			expectedShaderNames,
			actualShaderNames,
			input,
		)
	}

	if !isEqual(actualShaderFiles, expectedShaderFiles) {
		t.Errorf(
			"Expected shader files %v got %v for %v",
			expectedShaderFiles,
			actualShaderFiles,
			input,
		)
	}
}

func TestCombineTexturesFromShaders(t *testing.T) {
	textures := map[string]int{}
	shaderNames := []string{}
	input := []shader.Shader{
		{
			Name:     "testmap/test_shader_2",
			Lines:    []string{},
			Textures: map[string]int{"testmap/test_shader_4": 1, "testmap/test_shader_5": 1},
		},
		{
			Name:     "testmap/test_shader",
			Lines:    []string{},
			Textures: map[string]int{"testmap/test_shader_2": 1, "testmap/test_shader_3": 1},
		},
	}
	expectedTextures := map[string]int{
		"testmap/test_shader_4": 1,
		"testmap/test_shader_5": 1,
		"testmap/test_shader_2": 1,
		"testmap/test_shader_3": 1,
	}
	expectedShaders := []string{"testmap/test_shader_2", "testmap/test_shader"}
	actualTextures, actualShaders := shader.CombineTexturesFromShaders(input, textures, shaderNames)

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
		expected           []shader.Shader
	}{
		{
			map[string]int{"testmap/test_shader_2": 1},
			"test_shader_2.shader",
			[]shader.Shader{
				{
					Name:  "testmap/test_shader_2",
					Lines: []string{},
					Textures: map[string]int{
						"testmap/test_shader_4": 1,
						"testmap/test_shader_5": 1,
					},
				},
			},
		},
		{
			map[string]int{"testmap/test_shader": 1},
			"testmap.shader",
			[]shader.Shader{
				{
					Name:  "testmap/test_shader",
					Lines: []string{},
					Textures: map[string]int{
						"testmap/test_shader_2": 1,
						"testmap/test_shader_3": 1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		actual := shader.ParseShaderFile(
			test.shadersFromMapFile,
			test.shaderFileName,
			"data/baseq3/scripts",
		)
		if len(actual) != len(test.expected) {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
		}
		for index, actualValue := range actual {
			if actualValue.Name != test.expected[index].Name {
				t.Errorf(
					"Expected %v got %s for %v",
					test.expected[index].Name,
					actualValue.Name,
					test,
				)
			}
			equalTextures := reflect.DeepEqual(actualValue.Textures, test.expected[index].Textures)
			if !equalTextures {
				t.Errorf(
					"Expected %v got %v for %v",
					test.expected[index].Textures,
					actualValue.Textures,
					test,
				)
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
		actual := shader.ShaderIsUsed(test.shadersFromMapFile, test.shaderName)
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
		actual := shader.IsShaderName(test.line)
		if actual != test.expected {
			t.Errorf("Expected %v got %v for %v", test.expected, actual, test)
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
