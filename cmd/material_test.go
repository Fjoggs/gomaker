package gomaker

import (
	"testing"
)

func TestGetMaterial(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) common/caulk 32 0 0 0.5 0.5 134217728 0 0", ""},
		{"// Entity 0", ""},
		{"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0", "testmap/test_texture"},
		{"( 96 64 192 ) ( 240 64 128 ) ( 96 64 128 ) testmap/test-texture-2 384 256 0 0.25 0.25 134217728 0 0", "testmap/test-texture-2"},
		{"( 216 -64 120 ) ( 200 -192 128 ) ( 216 -192 120 ) common/caulk 0 32 0 0.5 0.5 134217728 0 0", ""},
		{"( 112 -64 192 ) ( 128 -192 184 ) ( 112 -192 192 ) testmap_a1/23-texture 384 0 0 0.25 0.25 134217728 0 0", "testmap_a1/23-texture"},
		{"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap-b5/texture2 32 0 0 0.5 0.5 134217728 0 0", "testmap-b5/texture2"},
		{"}", ""},
	}
	for _, test := range tests {
		value := getMaterial(test.input)
		if value != test.expected {
			t.Errorf("Expected %s got %s for %v", value, test.input, test)
		}
	}
}

func TestIsCustomMaterial(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"common/caulk", false},
		{"common_alphascale/", false},
		{"sfx/", false},
		{"sfx/something", false},
		{"liquids/", false},
		{"effects/", false},
		{"testmap/", true},
		{"testmap2/", true},
	}
	for _, test := range tests {
		value := isCustomMaterial(test.input)
		if value != test.expected {
			t.Errorf("Expected %v got %s", value, test.input)
		}
	}
}
func TestFormatPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"textures/common/caulk", "common/caulk"},
		{"common_alphascale/", "common_alphascale/"},
		{"texture/sfx/lol", "texture/sfx/lol"},
		{"textures/ab_c-12/testy", "ab_c-12/testy"},
	}
	for _, test := range tests {
		value := formatPath(test.input)
		if value != test.expected {
			t.Errorf("Expected %v got %s", value, test.input)
		}
	}
}

func TestIsTexture(t *testing.T) {
	tests := []struct {
		input           string
		expectedBool    bool
		expectedTexture string
	}{
		{"test_texture", true, "resources/test_texture.jpg"},
		{"test_texture_2", false, "test_texture_2"},
		{"test_texture_3", true, "resources/test_texture_3.tga"},
	}
	rootPath := "resources/"
	for _, test := range tests {
		actualBool, actualTexture := isTexture(test.input, rootPath)
		if actualBool != test.expectedBool && actualTexture != test.expectedTexture {
			t.Errorf("Expected %v got %v for %v", test.expectedBool, actualBool, test.input)
		}
	}
}

func TestAddTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"test_texture", "test_texture/"},
		{"test_texture_2", "test_texture_2/"},
		{"test_texture_3/", "test_texture_3/"},
		{"long/path/no/slash", "long/path/no/slash/"},
		{"long/path/yes/slash/", "long/path/yes/slash/"},
	}
	for _, test := range tests {
		actual := addTrailingSlash(test.input)
		if actual != test.expected {
			t.Errorf("Expected %v got %s for %v", test.expected, actual, test.input)
		}
	}
}
