package gomaker

import (
	"testing"
)

func TestReadMap(t *testing.T) {
	filePath := "resources/test.map"
	expected := []string{"testmap/test_texture", "testmap/test-texture-2", "testmap_a1/23-texture"}
	actual := readMap(filePath)
	if !isEqual(actual, expected) {
		t.Errorf("Excpted %s got %s", expected, actual)
	}
}
func TestReadLine(t *testing.T) {
	with := "( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap/texture 32 0 0 0.5 0.5 134217728 0 0"
	actual := readLine(with)
	expected := "testmap/texture"
	if actual != expected {
		t.Errorf("Expected %s got %s", expected, actual)
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
