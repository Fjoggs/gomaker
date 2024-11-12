package builder

import (
	"archive/zip"
	"os"
	"slices"
	"testing"
)

func TestCreatePk3(t *testing.T) {
	resources := []string{"scripts/testmap.arena", "levelshots/testmap.jpg", "maps/testmap.map"}
	CreatePk3("resources", resources, "testmap", true)

	expected := []string{
		"/",
		"levelshots/",
		"levelshots/testmap.jpg",
		"maps/",
		"maps/testmap.map",
		"scripts/",
		"scripts/testmap.arena",
	}

	_, err := os.Stat("output/testmap.pk3")
	if err != nil {
		t.Errorf("PK3 does not exist: %s", err)
	}

	readCloser, err := zip.OpenReader("output/testmap.pk3")
	if err != nil {
		t.Errorf("Open reader blew up: %s", err)
	}
	defer readCloser.Close()

	numOfPaths := len(readCloser.File)
	for _, f := range readCloser.File {
		path := f.Name
		if !slices.Contains(expected, path) {
			t.Errorf("Expected %s to be in %v", path, expected)
		}
	}

	expectedNumOfPaths := len(expected)
	if numOfPaths != len(expected) {
		t.Errorf("Expected number of paths to be %v but got %v", expectedNumOfPaths, numOfPaths)
	}

	DeleteFolderAndSubFolders("output")
}

func TestCreateDirectory(t *testing.T) {
	expected := true
	actual := createDirectory("testcreate", "")
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
	DeleteFolderAndSubFolders("testcreate")
}

func TestZipOutputFolder(t *testing.T) {
	createDirectory("output", "")
	createDirectory("env", "output")
	createDirectory("maps", "output")
	createDirectory("textures", "output")
	createDirectory("randomdir", "output/textures")
	createDirectory("scripts", "output")
	createDirectory("sounds", "output")
	createDirectory("levelshots", "output")

	err := ZipOutputFolderAsPk3("output", "testmap")
	if err != nil {
		t.Errorf("Error while creating pk3: %s", err)
	}

	_, statErr := os.Stat("output/testmap.pk3")

	if statErr != nil {
		t.Errorf("ZIP does not exist: %s", statErr)
	}

	DeleteFolderAndSubFolders("output")
}

func TestAddResourceIfExists(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"scripts/testmap.arena", "output/scripts/testmap.arena"},
		{"levelshots/testmap.jpg", "output/levelshots/testmap.jpg"},
		{"env/something/test.jpg", ""},
		{"maps/testmap/lm_0000.tga", "output/maps/testmap/lm_0000.tga"},
	}

	for _, test := range tests {

		actual := AddResourceIfExists("resources", test.input, "output")
		if actual != test.expected {
			t.Errorf("Expected %v got %v", test.expected, actual)
		}
	}
	DeleteFolderAndSubFolders("output")
}

func TestDeleteFolderAndSubFolders(t *testing.T) {
	createDirectory("testdelete", "output/")
	DeleteFolderAndSubFolders("output/testdelete")
}

func TestGetCfgFile(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", ""},
		{"testmap", "cfg-maps/testmap.cfg"},
	}

	for _, test := range tests {
		actual := GetCfgFile("resources", test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}

func TestGetReadme(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", ""},
		{"testmap", "testmap.txt"},
	}

	for _, test := range tests {
		actual := GetReadme("resources", test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}

func TestGetBspFile(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", ""},
		{"testmap", "maps/testmap.bsp"},
	}

	for _, test := range tests {
		actual := GetBspFile("resources", test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}

func TestGetMapFile(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", ""},
		{"testmap", "maps/testmap.map"},
	}

	for _, test := range tests {
		actual := GetMapFile("resources", test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}

func TestGetExternalLightmaps(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"test", []string{}},
		{
			"testmap",
			[]string{
				"maps/testmap/lm_0000.tga",
				"maps/testmap/lm_0001.tga",
				"maps/testmap/lm_0002.tga",
			},
		},
	}

	for _, test := range tests {
		actualLightmaps := GetExternalLightmaps("resources", test.input)
		actualLength := len(actualLightmaps)
		expectedLength := len(test.expected)
		if actualLength != expectedLength {
			t.Errorf(
				"Expected lightmap slice length to be %v but was %v",
				expectedLength,
				actualLength,
			)
		}
		for index, actual := range actualLightmaps {
			if actual != test.expected[index] {
				t.Errorf("Expected %s got %v", test.expected[index], actual)
			}
		}
	}
}

func TestGetArenaFile(t *testing.T) {
	expected := "scripts/testmap.arena"
	actual := GetArenaFile("resources", "testmap")
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
		actual := GetLevelshot("resources", test.input)
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
		actual := ExtractFolderPaths(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}
