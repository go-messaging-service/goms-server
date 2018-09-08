package util_test

import (
	"testing"

	"github.com/go-messaging-service/goms-server/src/util"
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

	b := util.ContainsString(testSlice, testString)
	if b {
		t.FailNow()
	}
}

func TestContainsStringNotContainc(t *testing.T) {
	testSlice := getTestSlice()

	testString := "test"

	if util.ContainsString(testSlice, testString) {
		t.FailNow()
	}
}

func TestContainsStringCorrectContains(t *testing.T) {
	testSlice := getTestSlice()

	for _, v := range testSlice {
		testString := v
		if !util.ContainsString(testSlice, testString) {
			t.FailNow()
		}
	}
}

func TestRemoveStringsCorrect(t *testing.T) {
	testSlice := getTestSlice()

	sliceToRemove := make([]string, 4)
	sliceToRemove[0] = "b"
	sliceToRemove[1] = "d"

	resultSlice := util.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || util.ContainsString(resultSlice, "a")
	result = result || !util.ContainsString(resultSlice, "b")
	result = result || util.ContainsString(resultSlice, "c")
	result = result || !util.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringsEmptyToRemove(t *testing.T) {
	testSlice := getTestSlice()

	sliceToRemove := make([]string, 0)

	resultSlice := util.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || util.ContainsString(resultSlice, "a")
	result = result || util.ContainsString(resultSlice, "b")
	result = result || util.ContainsString(resultSlice, "c")
	result = result || util.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringsRemoveNotExisting(t *testing.T) {
	testSlice := getTestSlice()

	sliceToRemove := make([]string, 2)
	sliceToRemove[0] = "e"
	sliceToRemove[1] = "f"

	resultSlice := util.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || util.ContainsString(resultSlice, "a")
	result = result || util.ContainsString(resultSlice, "b")
	result = result || util.ContainsString(resultSlice, "c")
	result = result || util.ContainsString(resultSlice, "d")
	result = result || !util.ContainsString(resultSlice, "e")
	result = result || !util.ContainsString(resultSlice, "f")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringsRemoveFromEmpty(t *testing.T) {
	testSlice := make([]string, 0)

	sliceToRemove := make([]string, 2)
	sliceToRemove[0] = "e"
	sliceToRemove[1] = "f"

	resultSlice := util.RemoveStrings(testSlice, sliceToRemove)

	result := false
	result = result || !util.ContainsString(resultSlice, "e")
	result = result || !util.ContainsString(resultSlice, "f")
	result = result || (len(resultSlice) == 0)

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringCorrect(t *testing.T) {
	testSlice := getTestSlice()

	resultSlice := util.RemoveString(testSlice, "b")
	resultSlice = util.RemoveString(resultSlice, "d")

	result := false
	result = result || util.ContainsString(resultSlice, "a")
	result = result || !util.ContainsString(resultSlice, "b")
	result = result || util.ContainsString(resultSlice, "c")
	result = result || !util.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringEmptyToRemove(t *testing.T) {
	testSlice := getTestSlice()

	resultSlice := util.RemoveString(testSlice, "")

	result := false
	result = result || util.ContainsString(resultSlice, "a")
	result = result || util.ContainsString(resultSlice, "b")
	result = result || util.ContainsString(resultSlice, "c")
	result = result || util.ContainsString(resultSlice, "d")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringRemoveNotExisting(t *testing.T) {
	testSlice := getTestSlice()

	resultSlice := util.RemoveString(testSlice, "e")
	resultSlice = util.RemoveString(resultSlice, "f")

	result := false
	result = result || util.ContainsString(resultSlice, "a")
	result = result || util.ContainsString(resultSlice, "b")
	result = result || util.ContainsString(resultSlice, "c")
	result = result || util.ContainsString(resultSlice, "d")
	result = result || !util.ContainsString(resultSlice, "e")
	result = result || !util.ContainsString(resultSlice, "f")

	if !result {
		t.FailNow()
	}
}

func TestRemoveStringRemoveFromEmpty(t *testing.T) {
	testSlice := make([]string, 0)

	resultSlice := util.RemoveString(testSlice, "e")
	resultSlice = util.RemoveString(resultSlice, "f")

	result := false
	result = result || !util.ContainsString(resultSlice, "e")
	result = result || !util.ContainsString(resultSlice, "f")
	result = result || (len(resultSlice) == 0)

	if !result {
		t.FailNow()
	}
}
