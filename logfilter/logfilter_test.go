package main

import (
	"reflect"
	"testing"
)

func TestFixBytesToString(t *testing.T) {
	// happy path
	input := []interface{}{float64('A'), float64('B')}
	output := fixBytesToString(input)
	expected := "AB"
	if output != expected {
		t.Errorf("expected output %v but got %v\n", expected, output)
	}
}

const expectedNewlines = `A
B
C`

func TestFixNewlines(t *testing.T) {
	input := "A|B|C"
	output := fixNewLines(input)
	if output != expectedNewlines {
		t.Errorf("expected %v but got %v\n", expectedNewlines, output)
	}
}

var rawJson = map[string]interface{}{
	"MESSAGE":               "message",
	"_HOSTNAME":             "hostname",
	"_MACHINE_ID":           "machine",
	"_SYSTEMD_UNIT":         "system",
	"_GID":                  "gid",
	"_COMM":                 "comm",
	"_EXE":                  "exe",
	"_CAP_EFFECTIVE":        "cap",
	"SYSLOG_FACILITY":       "syslog",
	"PRIORITY":              "priority",
	"SYSLOG_IDENTIFIER":     "syslogi",
	"_BOOT_ID":              "boot",
	"_CMDLINE":              "cmd",
	"_SYSTEMD_CGROUP":       "cgroup",
	"_SYSTEMD_SLICE":        "slice",
	"_TRANSPORT":            "transport",
	"_UID":                  "uid",
	"__CURSOR":              "cursor",
	"__MONOTONIC_TIMESTAMP": "monotonic",
}

var blacklistFilteredJson = map[string]interface{}{
	"MESSAGE":       "message",
	"_HOSTNAME":     "hostname",
	"_MACHINE_ID":   "machine",
	"_SYSTEMD_UNIT": "system",
}

func TestApplyPropertyBlacklist(t *testing.T) {
	removeBlacklistedProperties(rawJson)
	if !reflect.DeepEqual(rawJson, blacklistFilteredJson) {
		t.Errorf("expected %v but got %v\n", blacklistFilteredJson, rawJson)
	}
}

func TestEnvTag(t *testing.T) {
	s := ""
	environmentTag = &s

	m := make(map[string]interface{})
	munge(m)

	if m["environment"] != nil {
		t.Errorf("didn't expect to find environment %v", m["environment"])
	}

	s = "foo"
	environmentTag = &s

	munge(m)

	if m["environment"] != "foo" {
		t.Errorf("expected foo but got  %v", m["environment"])
	}

}
