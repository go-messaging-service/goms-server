package technicalCommon_test

import (
	"goms-server/src/technical/common"
	"testing"
)

func getTestSlice() []string {
	testSlice := make([]string, 4)

	testSlice[0] = "a"
	testSlice[1] = "b"
	testSlice[2] = "c"
	testSlice[3] = "d"

	return testSlice
}

func TestContainsStringEmpty(t *testing.T) {
	testSlice := make([]string, 0)
	testString := "test"

	b := technicalCommon.ContainsString(testSlice, testString)
	if b {
		t.FailNow()
	}
}

func TestContainsStringNotContainc(t *testing.T) {
	testSlice := getTestSlice()

	testString := "test"

	if technicalCommon.ContainsString(testSlice, testString) {
		t.FailNow()
	}
}

func TestContainsStringCorrectContains(t *testing.T) {
	testSlice := getTestSlice()

	for _, v := range testSlice {
		testString := v
		if !technicalCommon.ContainsString(testSlice, testString) {
			t.FailNow()
		}
	}
}

func TestRemoveStringsCorrect(t *testing.T) {
	testSlice := getTestSlice()

	sliceToRemove := make([]string, 4)
	sliceToRemove[0] = "b"
	sliceToRemove[1] = "d"

	resultSlice := technicalCommon.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || technicalCommon.ContainsString(resultSlice, "a")
	result = result || !technicalCommon.ContainsString(resultSlice, "b")
	result = result || technicalCommon.ContainsString(resultSlice, "c")
	result = result || !technicalCommon.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringsEmptyToRemove(t *testing.T) {
	testSlice := getTestSlice()

	sliceToRemove := make([]string, 0)

	resultSlice := technicalCommon.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || technicalCommon.ContainsString(resultSlice, "a")
	result = result || technicalCommon.ContainsString(resultSlice, "b")
	result = result || technicalCommon.ContainsString(resultSlice, "c")
	result = result || technicalCommon.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringsRemoveNotExisting(t *testing.T) {
	testSlice := getTestSlice()

	sliceToRemove := make([]string, 2)
	sliceToRemove[0] = "e"
	sliceToRemove[1] = "f"

	resultSlice := technicalCommon.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || technicalCommon.ContainsString(resultSlice, "a")
	result = result || technicalCommon.ContainsString(resultSlice, "b")
	result = result || technicalCommon.ContainsString(resultSlice, "c")
	result = result || technicalCommon.ContainsString(resultSlice, "d")
	result = result || !technicalCommon.ContainsString(resultSlice, "e")
	result = result || !technicalCommon.ContainsString(resultSlice, "f")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringsRemoveFromEmpty(t *testing.T) {
	testSlice := make([]string, 0)

	sliceToRemove := make([]string, 2)
	sliceToRemove[0] = "e"
	sliceToRemove[1] = "f"

	resultSlice := technicalCommon.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || !technicalCommon.ContainsString(resultSlice, "e")
	result = result || !technicalCommon.ContainsString(resultSlice, "f")
	result = result || (len(resultSlice) == 0)

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringCorrect(t *testing.T) {
	testSlice := getTestSlice()

	resultSlice := technicalCommon.RemoveString(testSlice, "b")
	resultSlice = technicalCommon.RemoveString(resultSlice, "d")

	result := false
	result = result || technicalCommon.ContainsString(resultSlice, "a")
	result = result || !technicalCommon.ContainsString(resultSlice, "b")
	result = result || technicalCommon.ContainsString(resultSlice, "c")
	result = result || !technicalCommon.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringEmptyToRemove(t *testing.T) {
	testSlice := getTestSlice()

	resultSlice := technicalCommon.RemoveString(testSlice, "")

	result := false
	result = result || technicalCommon.ContainsString(resultSlice, "a")
	result = result || technicalCommon.ContainsString(resultSlice, "b")
	result = result || technicalCommon.ContainsString(resultSlice, "c")
	result = result || technicalCommon.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringRemoveNotExisting(t *testing.T) {
	testSlice := getTestSlice()

	resultSlice := technicalCommon.RemoveString(testSlice, "e")
	resultSlice = technicalCommon.RemoveString(resultSlice, "f")

	result := false
	result = result || technicalCommon.ContainsString(resultSlice, "a")
	result = result || technicalCommon.ContainsString(resultSlice, "b")
	result = result || technicalCommon.ContainsString(resultSlice, "c")
	result = result || technicalCommon.ContainsString(resultSlice, "d")
	result = result || !technicalCommon.ContainsString(resultSlice, "e")
	result = result || !technicalCommon.ContainsString(resultSlice, "f")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringRemoveFromEmpty(t *testing.T) {
	testSlice := make([]string, 0)

	resultSlice := technicalCommon.RemoveString(testSlice, "e")
	resultSlice = technicalCommon.RemoveString(resultSlice, "f")

	result := false
	result = result || !technicalCommon.ContainsString(resultSlice, "e")
	result = result || !technicalCommon.ContainsString(resultSlice, "f")
	result = result || (len(resultSlice) == 0)

	if !result {
		t.FailNow()
	}
}
