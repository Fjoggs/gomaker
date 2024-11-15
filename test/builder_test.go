package test

import (
	"archive/zip"
	"fmt"
	"os"
	"slices"
	"testing"

	"gomaker/internal/builder"
)

func TestBuildPk3(t *testing.T) {
	expected := []string{
		"/",
		"testmap.txt",
		"cfg-maps/",
		"cfg-maps/testmap.cfg",
		"levelshots/",
		"levelshots/testmap.jpg",
		"maps/",
		"maps/testmap.bsp",
		"maps/testmap.map",
		"maps/testmap/lm_0000.tga",
		"maps/testmap/lm_0001.tga",
		"maps/testmap/lm_0002.tga",
		"scripts/",
		"scripts/testmap.arena",
		"scripts/testmap.shader",
		"scripts/test_shader_2.shader",
		"sound/",
		"sound/testmap/sound-file.wav",
		"textures/",
		"textures/testmap/test_model_texture_1.jpg",
		"textures/testmap/test_model_texture_2.tga",
		"textures/testmap/test_shader_2.tga",
		"textures/testmap/test_shader_3.jpg",
		"textures/testmap/test_texture.jpg",
		"textures/testmap/test_texture_3.tga",
	}

	pk3Path := builder.BuildPk3("testmap", "data/baseq3")

	_, err := os.Stat(pk3Path)
	if err != nil {
		t.Fatalf("PK3 does not exist: %s", err)
	}

	readCloser, err := zip.OpenReader(pk3Path)
	if err != nil {
		t.Fatalf("Open reader blew up: %s", err)
	}
	defer readCloser.Close()

	numOfPaths := len(readCloser.File)
	for _, f := range readCloser.File {
		path := f.Name
		if !slices.Contains(expected, path) {
			t.Fatalf("Expected %s to be in %v", path, expected)
		}
	}

	expectedNumOfPaths := len(expected)
	if numOfPaths != len(expected) {
		t.Fatalf("Expected number of paths to be %v but got %v", expectedNumOfPaths, numOfPaths)
	}
}

func TestCreatePk3(t *testing.T) {
	resources := []string{"scripts/testmap.arena", "levelshots/testmap.jpg", "maps/testmap.map"}
	pk3Path := builder.CreatePk3("data/baseq3", resources, "testmap")

	expected := []string{
		"/",
		"levelshots/",
		"levelshots/testmap.jpg",
		"maps/",
		"maps/testmap.map",
		"scripts/",
		"scripts/testmap.arena",
	}

	_, err := os.Stat(pk3Path)
	if err != nil {
		t.Errorf("PK3 does not exist: %s", err)
	}

	readCloser, err := zip.OpenReader(pk3Path)
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

	builder.DeleteFolderAndSubFolders("output")
}

func TestCreateDirectory(t *testing.T) {
	expected := true
	actual := builder.CreateDirectory("testcreate")
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
	builder.DeleteFolderAndSubFolders("testcreate")
}

func TestZipOutputFolderAsPk3(t *testing.T) {
	rootFolder := "output"
	builder.CreateDirectory(rootFolder)
	createSubFolder(rootFolder, "env")
	createSubFolder(rootFolder, "maps")
	createSubFolder(rootFolder, "textures")
	createSubFolder(rootFolder+"/textures", "randomdir")
	createSubFolder(rootFolder, "scripts")
	createSubFolder(rootFolder, "sounds")
	createSubFolder(rootFolder, "levelshots")

	pk3Path, err := builder.ZipOutputFolderAsPk3("output", "testmap")
	if err != nil {
		t.Errorf("Error while creating pk3: %s", err)
	}

	_, err = os.Stat(pk3Path)
	if err != nil {
		t.Errorf("ZIP does not exist: %s", err)
	}

	_, err = os.Stat("output")
	if err == nil {
		t.Errorf("Output folder still exists %s", err)
	}
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

		actual := builder.AddResourceIfExists("data/baseq3", test.input, "output")
		if actual != test.expected {
			t.Errorf("Expected %v got %v", test.expected, actual)
		}
	}
	builder.DeleteFolderAndSubFolders("output")
}

func TestDeleteFolderAndSubFolders(t *testing.T) {
	builder.CreateDirectory("output")
	createSubFolder("output", "testdelete")
	builder.DeleteFolderAndSubFolders("output/testdelete")
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
		actual := builder.GetCfgFile("data/baseq3", test.input)
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
		actual := builder.GetReadme("data/baseq3", test.input)
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
		actual := builder.GetBspFile("data/baseq3", test.input)
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
		actual := builder.GetMapFile("data/baseq3", test.input)
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
		actualLightmaps := builder.GetExternalLightmaps("data/baseq3", test.input)
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
	actual := builder.GetArenaFile("data/baseq3", "testmap")
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
		actual := builder.GetLevelshot("data/baseq3", test.input)
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
		{"data/baseq3/scripts/testmap.arena", "data/baseq3/scripts"},
	}

	for _, test := range tests {
		actual := builder.ExtractFolderPaths(test.input)
		if actual != test.expected {
			t.Errorf("Expected %s got %v", test.expected, actual)
		}
	}
}

func createSubFolder(rootFolder string, folderName string) {
	path := fmt.Sprintf("%s/%s", rootFolder, folderName)
	err := os.Mkdir(path, 0777)
	if err != nil {
		fmt.Printf("createSubFolder for path %s failed with err: %s", path, err)
	}
}
