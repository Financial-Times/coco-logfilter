package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Financial-Times/coco-logfilter"
	"github.com/stretchr/testify/assert"
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

func TestBlacklistedServices(t *testing.T) {
	testCases := []struct {
		jsonString string
		expected   bool
	}{
		{
			jsonString: `{"CONTAINER_ID":"03d1f4078733","CONTAINER_ID_FULL":"03d1f4078733f75f4505b07d1f8a3e8287ed497d9d54e0e785440cb969378ca3","CONTAINER_NAME":"k8s_cluster-autoscaler_cluster-autoscaler-79d574774-2rxrj_kube-system_a093cbca-fb5a-11e7-a6b6-06263dd4a414_6","CONTAINER_TAG":"gcr.io/google_containers/cluster-autoscaler@sha256:6ceb111a36020dc2124c0d7e3746088c20c7e3806a1075dd9e5fe1c42f744fff","HOSTNAME":"ip-10-172-40-164.eu-west-1.compute.internal","MACHINE_ID":"8d1225f40ee64cc7bcce2f549a41657c","MESSAGE":"I0119 15:38:05.932385 1 leaderelection.go:199] successfully renewed lease kube-system/cluster-autoscaler","POD_NAME":"cluster-autoscaler-79d574774-2rxrj","SERVICE_NAME":"cluster-autoscaler","SYSTEMD_UNIT":"docker.service","_SOURCE_REALTIME_TIMESTAMP":"1516376285932645","_SYSTEMD_INVOCATION_ID":"e3b2703c430f45e8a7075dbcf6b3a588","environment":"upp-prod-publish-eu","platform":"up-coco"}`,
			expected:   false,
		},
	}

	for _, c := range testCases {
		m := make(map[string]interface{})
		json.Unmarshal([]byte(c.jsonString), &m)
		ok := processMessage(m)
		assert.True(t, c.expected == ok)
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

func TestHideSingleAPIKeysInURLQueryParam(t *testing.T) {
	msgWithAPYKey := `10.2.26.0 ops-17-01-2018 30/Jan/2018:08:35:04 /content/notifications-push?apiKey=vhs2aazf3gyywm3wk2sv44wb&type=ALL 200 -2147483648 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36" transaction_id=- miss`
	expectedMsg := `10.2.26.0 ops-17-01-2018 30/Jan/2018:08:35:04 /content/notifications-push?apiKey=vhs2aazf3g**************&type=ALL 200 -2147483648 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36" transaction_id=- miss`
	actualMsg := hideAPIKeysInURLQueryParams(msgWithAPYKey)
	assert.Equal(t, expectedMsg, actualMsg)
}

func TestHideMultipleAPIKeysInURLQueryParams(t *testing.T) {
	msgWithAPYKey := `10.2.26.0 ops-17-01-2018 30/Jan/2018:08:35:04 /content/notifications-push?apiKey=vhs2aazf3gyywm3wk2sv44wb&type=ALL /content/notifications-push?apiKey=wm3wk2sv44wbvhs2aazf3gyy`
	expectedMsg := `10.2.26.0 ops-17-01-2018 30/Jan/2018:08:35:04 /content/notifications-push?apiKey=vhs2aazf3g**************&type=ALL /content/notifications-push?apiKey=wm3wk2sv44**************`
	actualMsg := hideAPIKeysInURLQueryParams(msgWithAPYKey)
	assert.Equal(t, expectedMsg, actualMsg)
}

func TestBypassWithoutAPIKeysInURLQueryParams(t *testing.T) {
	msgWithAPYKey := `10.2.26.0 ops-17-01-2018 30/Jan/2018:08:35:04 /content/notifications-push?type=ALL 200 -2147483648 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36" transaction_id=- miss`
	expectedMsg := `10.2.26.0 ops-17-01-2018 30/Jan/2018:08:35:04 /content/notifications-push?type=ALL 200 -2147483648 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36" transaction_id=- miss`
	actualMsg := hideAPIKeysInURLQueryParams(msgWithAPYKey)
	assert.Equal(t, expectedMsg, actualMsg)
}
