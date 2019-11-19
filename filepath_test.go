package main

import "testing"

func TestPathLevel0(t *testing.T) {
	expect := ""
	given := "a/b/c/d"
	recieved := PathLevel(given, 0)

	if expect != recieved {
		t.Errorf("Expected '%s', returns '%s'", expect, recieved)
	}
}

func TestPathLevel1(t *testing.T) {
	expect := "a"
	given := "a/b/c/d"
	recieved := PathLevel(given, 1)

	if expect != recieved {
		t.Errorf("Expected '%s', returns '%s'", expect, recieved)
	}
}

func TestPathLevel2(t *testing.T) {
	expect := "a/b"
	given := "a/b/c/d"
	recieved := PathLevel(given, 2)

	if expect != recieved {
		t.Errorf("Expected '%s', returns '%s'", expect, recieved)
	}
}
