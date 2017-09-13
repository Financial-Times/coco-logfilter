package main

import (
	"fmt"
	"reflect"
	"testing"

	"encoding/json"
	"github.com/Financial-Times/coco-logfilter"
	"github.com/stretchr/testify/assert"
	"strings"
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
	testCases := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "standard API call",
			message:  "foo baz baz transaction_id=transid_a-b banana",
			expected: "transid_a-b",
		},
		{
			name:     "message without transaction id",
			message:  "foo baz baz transzzzaction_id=transid_a-b banana",
			expected: "",
		},
		{
			name:     "PAM notifications feed transaction id may contain colon character",
			message:  "INFO  [2017-01-19 12:05:13,478] com.ft.api.util.transactionid.TransactionIdFilter: transaction_id=tid_pam_notifications_pull_2017-01-19T12:05:13Z [REQUEST HANDLED] uri=/content/notifications time_ms=2 status=200 exception_was_thrown=false [dw-1968]",
			expected: "tid_pam_notifications_pull_2017-01-19T12:05:13Z",
		},
		{
			name:     "transaction_id should not include parenthesis or quotes",
			message:  "foo baz baz \"My User Agent (transaction_id=transid_a-b)\" banana",
			expected: "transid_a-b",
		},
	}

	for _, c := range testCases {
		m := map[string]interface{}{
			"MESSAGE": c.message,
		}
		munge(m, c.message)

		actual, found := m["transaction_id"]
		if len(c.expected) == 0 {
			assert.False(t, found, fmt.Sprintf("expected no transaction_id for %s", c.name))
		} else {
			assert.Equal(t, c.expected, actual, fmt.Sprintf("transaction_id for %s", c.name))
		}
	}
}

func TestContainsBlacklistedStringWithBlacklistedString(t *testing.T) {
	message := "foo baz baz " + blacklistedStrings[0] + " foo "

	if !containsBlacklistedString(message) {
		t.Error("Expected to detect blacklisted string in test")
	}

}

func TestContainsBlacklistedStringWithoutBlacklistedString(t *testing.T) {
	message := "foo baz baz transazzzction_id=transid_a-b banana"

	if containsBlacklistedString(message) {
		t.Error("Detected black listed string when there was none")
	}

}

func TestClusterStatus(t *testing.T) {
	trueVar := true
	falseVar := false

	testCases := []struct {
		jsonString string
		dnsAddress string
		tag        string
		expected   *bool
	}{
		{
			jsonString: `{"@time":"2017-09-12T14:19:28.199162596Z","HOSTNAME":"ip-172-24-159-194.eu-west-1.compute.internal","MACHINE_ID":"1234","MESSAGE":"{\"@time\":\"2017-09-12T14:19:28.199162596Z\",\"content_type\":\"Suggestions\",\"event\":\"SaveNeo4j\",\"level\":\"info\",\"monitoring_event\":\"true\",\"msg\":\"%s successfully written in Neo4jSuggestions\",\"service_name\":\"suggestions-rw-neo4j\",\"transaction_id\":\"tid_u7pkkludzd\",\"uuid\":\"0ec3c76b-9be4-4d76-b1f9-5414460a8bc1\"}","SYSTEMD_UNIT":"suggestions-rw-neo4j@1.service","_SYSTEMD_INVOCATION_ID":"1234","content_type":"Suggestions","environment":"xp","event":"SaveNeo4j","level":"info","monitoring_event":"true","msg":"%s successfully written in Neo4jSuggestions","platform":"up-coco","service_name":"suggestions-rw-neo4j","transaction_id":"tid_test","uuid":"a3f63cda-97af-11e7-b83c-9588e51488a0"}`,
			dnsAddress: "google.com",
			tag:        "ns",
			expected:   &trueVar,
		},
		{
			jsonString: `{"@time":"2017-09-12T14:19:28.199162596Z","HOSTNAME":"ip-172-24-159-194.eu-west-1.compute.internal","MACHINE_ID":"1234","MESSAGE":"{\"@time\":\"2017-09-12T14:19:28.199162596Z\",\"content_type\":\"Suggestions\",\"event\":\"SaveNeo4j\",\"level\":\"info\",\"monitoring_event\":\"true\",\"msg\":\"%s successfully written in Neo4jSuggestions\",\"service_name\":\"suggestions-rw-neo4j\",\"transaction_id\":\"tid_u7pkkludzd\",\"uuid\":\"0ec3c76b-9be4-4d76-b1f9-5414460a8bc1\"}","SYSTEMD_UNIT":"suggestions-rw-neo4j@1.service","_SYSTEMD_INVOCATION_ID":"1234","content_type":"Suggestions","environment":"xp","event":"SaveNeo4j","level":"info","monitoring_event":"true","msg":"%s successfully written in Neo4jSuggestions","platform":"up-coco","service_name":"suggestions-rw-neo4j","transaction_id":"tid_test","uuid":"a3f63cda-97af-11e7-b83c-9588e51488a0"}`,
			dnsAddress: "google.com",
			tag:        "invalid",
			expected:   &falseVar,
		},
		{
			jsonString: `{"@time":"2017-09-12T14:19:28.199162596Z","HOSTNAME":"ip-172-24-159-194.eu-west-1.compute.internal","MACHINE_ID":"1234","MESSAGE":"{\"@time\":\"2017-09-12T14:19:28.199162596Z\",\"content_type\":\"Suggestions\",\"event\":\"SaveNeo4j\",\"level\":\"info\",\"msg\":\"%s successfully written in Neo4jSuggestions\",\"service_name\":\"suggestions-rw-neo4j\",\"transaction_id\":\"tid_u7pkkludzd\",\"uuid\":\"a0ec3c76b-9be4-4d76-b1f9-5414460a8bc1\"}","SYSTEMD_UNIT":"suggestions-rw-neo4j@1.service","_SYSTEMD_INVOCATION_ID":"1234","content_type":"Suggestions","environment":"xp","event":"SaveNeo4j","level":"info","msg":"%s successfully written in Neo4jSuggestions","platform":"up-coco","service_name":"suggestions-rw-neo4j","transaction_id":"tid_test","uuid":"a3f63cda-97af-11e7-b83c-9588e51488a0"}`,
			dnsAddress: "google.com",
			tag:        "ns",
			expected:   nil,
		},
	}

	for _, c := range testCases {
		mc = logfilter.NewMonitoredClusterService(c.dnsAddress, c.tag)
		m := make(map[string]interface{})
		json.NewDecoder(strings.NewReader(c.jsonString)).Decode(&m)
		processMessage(m)
		if c.expected == nil {
			assert.Nil(t, m["active_cluster"])
		} else {
			assert.Equal(t, *c.expected, m["active_cluster"])
		}
	}
}
