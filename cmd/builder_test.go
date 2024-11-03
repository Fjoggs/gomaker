package gomaker

import (
	"os"
	"testing"
)

func TestCreatePk3(t *testing.T) {
	createPk3("output", "testmap", true)

	_, err := os.Stat("output/testmap.pk3")

	if err != nil {
		t.Errorf("PK3 does not exist: %s", err)
	}
	deleteFolderAndSubFolders("output")
}

func TestCreateDirectory(t *testing.T) {
	expected := true
	actual := createDirectory("testcreate", "")
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
	deleteFolderAndSubFolders("testcreate")
}

func TestZipOutputFolder(t *testing.T) {
	createDirectory("output", "")
	createDirectory("testmap", "output")
	pk3DirPath := "output/testmap/"
	createDirectory("env", pk3DirPath)
	createDirectory("maps", pk3DirPath)
	createDirectory("textures", pk3DirPath)
	createDirectory("randomdir", pk3DirPath+"textures")
	createDirectory("scripts", pk3DirPath)
	createDirectory("sounds", pk3DirPath)
	createDirectory("levelshots", pk3DirPath)

	err := zipOutputFolder("output", "testmap")

	if err != nil {
		t.Errorf("Error while creating pk3: %s", err)
	}

	_, statErr := os.Stat("output/testmap.pk3")

	if statErr != nil {
		t.Errorf("ZIP does not exist: %s", statErr)
	}

	deleteFolderAndSubFolders("output")
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

		actual := addResourceIfExists("resources", test.input, "output")
		if actual != test.expected {
			t.Errorf("Expected %v got %v", test.expected, actual)
		}
	}
	deleteFolderAndSubFolders("output")
}

func TestDeleteFolderAndSubFolders(t *testing.T) {
	createDirectory("testdelete", "output/")
	deleteFolderAndSubFolders("output/testdelete")
}

func TestAddArenaFile(t *testing.T) {
	expected := "scripts/testmap.arena"
	actual := getArenaFile("resources", "testmap")
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
		actual := getLevelshot("resources", test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}
func TestExtractFolderPaths(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", "test"},
		{"this/is/a/test", "this/is/a/test"},
		{"this/is/also/a/test.txt", "this/is/also/a"},
		{"resources/scripts/testmap.arena", "resources/scripts"},
	}

	for _, test := range tests {
		actual := extractFolderPaths(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}
