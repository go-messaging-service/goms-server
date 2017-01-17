package technicalCommon_test

import (
	"goMS/src/technical/common"
	"testing"
)

func TestContainsStringEmpty(t *testing.T) {
	testSlice := make([]string, 0)
	testString := "test"

	if technicalCommon.ContainsString(testSlice, testString) {
		t.FailNow()
	}
}

func TestContainsStringNotContainc(t *testing.T) {
	testSlice := make([]string, 4)
	testSlice[0] = "a"
	testSlice[1] = "b"
	testSlice[2] = "c"
	testSlice[3] = "d"

	testString := "test"

	if technicalCommon.ContainsString(testSlice, testString) {
		t.FailNow()
	}
}

func TestContainsStringCorrectContains(t *testing.T) {
	testSlice := make([]string, 4)
	testSlice[0] = "a"
	testSlice[1] = "b"
	testSlice[2] = "c"
	testSlice[3] = "d"

	for _, v := range testSlice {
		testString := v
		if !technicalCommon.ContainsString(testSlice, testString) {
			t.FailNow()
		}
	}
}
