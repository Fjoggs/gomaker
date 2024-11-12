package builder

import (
	"reflect"
	"testing"
)

func TestGetSound(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) testmap/texture 32 0 0 0.5 0.5 134217728 0 0",
			"",
		},
		{"// Entity 0", ""},
		{"{", ""},
		{"// entity 1", ""},
		{`"classname" "target_speaker"`, ""},
		{`"origin" "296 1032 488"`, ""},
		{`"spawnflags" "1"`, ""},
		{`"noise" "sound/testmap/sound-file.wav"`, "sound/testmap/sound-file.wav"},
		{"// entity 2", ""},
		{`"classname" "target_speaker"`, ""},
		{`"origin" "296 1032 488"`, ""},
		{`"spawnflags" "1"`, ""},
		{`"noise" "sound/world/base-file.wav"`, ""},
		{"}", ""},
	}
	for _, test := range tests {
		actual := getSound(test.input)

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v got %v for %s", test.expected, actual, test.input)
		}
	}
}
