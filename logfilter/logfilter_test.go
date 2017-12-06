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
	"_PID":                  "pid",
	"_SELINUX_CONTEXT":      "selinux context",
	"__REALTIME_TIMESTAMP":  "realtime timestamp",
}

var blacklistFilteredJson = map[string]interface{}{
	"MESSAGE":       "message",
	"_HOSTNAME":     "hostname",
	"_MACHINE_ID":   "machine",
	"_SYSTEMD_UNIT": "system",
}

var blacklistFilteredAndPropertiesRenamedJson = map[string]interface{}{
	"MESSAGE":      "message",
	"HOSTNAME":     "hostname",
	"MACHINE_ID":   "machine",
	"SYSTEMD_UNIT": "system",
}

func TestApplyPropertyBlacklist(t *testing.T) {
	removeBlacklistedProperties(rawJson)
	if !reflect.DeepEqual(rawJson, blacklistFilteredJson) {
		t.Errorf("expected %v but got %v\n", blacklistFilteredJson, rawJson)
	}
}

func TestShouldRenameProperties(t *testing.T) {
	renameProperties(blacklistFilteredJson)
	if !reflect.DeepEqual(blacklistFilteredJson, blacklistFilteredAndPropertiesRenamedJson) {
		t.Errorf("expected %v but got %v\n", blacklistFilteredAndPropertiesRenamedJson, blacklistFilteredJson)
	}
}

func TestEnvTag(t *testing.T) {
	s := ""
	environmentTag = &s

	m := make(map[string]interface{})
	munge(m, "")

	if m["environment"] != nil {
		t.Errorf("didn't expect to find environment %v", m["environment"])
	}

	s = "foo"
	environmentTag = &s

	munge(m, "")

	if m["environment"] != "foo" {
		t.Errorf("expected foo but got  %v", m["environment"])
	}

}

func TestTransactionId(t *testing.T) {
	message := "foo baz baz transaction_id=transid_a-b banana"
	m := map[string]interface{}{
		"MESSAGE": message,
	}
	munge(m, message)

	expected := "transid_a-b"
	actual := m["transaction_id"]
	if actual != expected {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestNoTransactionId(t *testing.T) {
	message := "foo baz baz transazzzction_id=transid_a-b banana"
	m := map[string]interface{}{
		"MESSAGE": message,
	}
	munge(m, message)

	actual := m["transaction_id"]
	if actual != nil {
		t.Errorf("expected nil but got %v", actual)
	}
}

func TestContainsBlacklistedStringWithBlacklistedString(t *testing.T) {
	message := "foo baz baz " + blacklistedStrings[0] + " foo "

	if !containsBlacklistedString(message, blacklistedStrings) {
		t.Error("Expected to detect blacklisted string in test")
	}
}

func TestContainsBlacklistedStringWithoutBlacklistedString(t *testing.T) {
	message := "foo baz baz transazzzction_id=transid_a-b banana"

	if containsBlacklistedString(message, blacklistedStrings) {
		t.Error("Detected black listed string when there was none")
	}
}

func TestExtractPodNameWithEmptyContainerTag(t *testing.T) {
	if podName := extractPodName(""); podName != "" {
		t.Error("Expected empty string as pod name when empty container tag is provided")
	}
}

func TestExtractPodNameWithNonStringContainerTag(t *testing.T) {
	nonStringContainerTag := 1
	if podName := extractPodName(nonStringContainerTag); podName != "" {
		t.Error("Expected empty string as pod name when non string container tag is provided")
	}
}

func TestExtractPodNameWithContainerTagWithoutUnderscores(t *testing.T) {
	if podName := extractPodName("test"); podName != "" {
		t.Error("Expected empty string as pod name when container tag without underscores is provided")
	}
}

func TestExtractPodNameWithContainerTagWithOneUnderscore(t *testing.T) {
	if podName := extractPodName("test_a"); podName != "" {
		t.Error("Expected empty string as pod name when container tag with one underscore is provided")
	}
}

func TestExtractPodNameWithValidContainerTagContainingTwoUnderscores(t *testing.T) {
	if podName := extractPodName("test_a_b"); podName != "b" {
		t.Error("Expected non empty string as pod name when container tag with two underscores is provided")
	}
}

func TestExtractPodNameWithValidContainerTagContainingMoreThanTwoUnderscores(t *testing.T) {
	if podName := extractPodName("test_a_b_c"); podName != "b" {
		t.Error("Expected third substring from container tag as pod name when container tag with more two underscores is provided")
	}
}