package gomaker

import (
	"os"
	"testing"
)

func TestCreatePk3(t *testing.T) {
	createPk3("output", "resources", true)

	_, err := os.Stat("output/resources.pk3")

	if err != nil {
		t.Errorf("PK3 does not exist: %s", err)
	}

}

func TestCreateDirectory(t *testing.T) {
	expected := true
	actual := createDirectory("testcreate", "output/")
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
	deleteFolderAndSubFolders("output/testcreate")
}

func TestAddResourceIfExists(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"scripts/testmap.arena", true},
		{"levelshots/testmap.jpg", true},
		{"env/something/test.jpg", false},
	}

	for _, test := range tests {

		actual := addResourceIfExists(test.input, "resources", "output")
		if actual != test.expected {
			t.Errorf("Expected %v got %v", test.expected, actual)
		}
	}
}

func TestDeleteFolderAndSubFolders(t *testing.T) {
	createDirectory("testdelete", "output/")
	deleteFolderAndSubFolders("output/testdelete")
}

func TestAddArenaFile(t *testing.T) {
	expected := "scripts/testmap.arena"
	actual := getArenaFile("testmap")
	if actual != expected {
		t.Errorf("Expected %s got %v", expected, actual)
	}
}

func TestGetLevelshot(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"testmap", "levelshots/testmap.jpg"},
		{"testmap2", "levelshots/testmap2.tga"},
		{"testmap3", ""},
	}

	for _, test := range tests {
		actual := getLevelshot(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}
