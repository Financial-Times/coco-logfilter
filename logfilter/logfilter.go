package main

import (
	"encoding/json"
	"flag"
	"github.com/Financial-Times/coco-logfilter"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	blacklistedProperties = []string{
		"_GID",
		"_CAP_EFFECTIVE",
		"SYSLOG_FACILITY",
		"PRIORITY",
		"SYSLOG_IDENTIFIER",
		"_BOOT_ID",
		"_CMDLINE",
		"_COMM",
		"_EXE",
		"_SYSTEMD_CGROUP",
		"_SYSTEMD_SLICE",
		"_TRANSPORT",
		"_UID",
		"__CURSOR",
		"__MONOTONIC_TIMESTAMP",
	}

	blacklistedUnits = map[string]bool{
		"splunk-forwarder.service": true,
	}
)

var environmentTag *string

func main() {
	environmentTag = flag.String("environment", "", "set the environment tag to use in the outputted json")
	flag.Parse()

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		m := make(map[string]interface{})
		err := dec.Decode(&m)
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		unit := m["_SYSTEMD_UNIT"]
		if unitString, ok := unit.(string); ok {
			if blacklistedUnits[unitString] {
				continue
			}
		}
		munge(m)
		removeBlacklistedProperties(m)
		enc.Encode(m)
	}
}

func munge(m map[string]interface{}) {

	m["platform"] = "up-coco"
	if *environmentTag != "" {
		m["environment"] = *environmentTag
	}

	message, ok := fixBytesToString(m["MESSAGE"]).(string)
	if !ok {
		return
	}

	message = fixNewLines(message)
	m["MESSAGE"] = message

	trans_id := extractTransactionId(message)
	if trans_id != "" {
		m["transaction_id"] = trans_id
	}

	ent, ok := logfilter.Extract(message)
	if !ok {
		return
	}

	// hackity
	j, err := json.Marshal(ent)
	if err != nil {
		panic(err)
	}
	entMap := make(map[string]interface{})
	err = json.Unmarshal(j, &entMap)
	if err != nil {
		panic(err)
	}
	for k, v := range entMap {
		m[k] = v
	}
}

var trans_regex = regexp.MustCompile(`\btransaction_id=([\w-]*)`)

func extractTransactionId(message string) string {
	matches := trans_regex.FindAllStringSubmatch(message, -1)
	if len(matches) != 0 {
		return matches[0][1]
	}

	return ""
}

// workaround for cases where a string has been turned into a
// byte array, or more accurately an array of float64, since
// we've been via json.
// TODO: remove this hack once the underlying cause is found
func fixBytesToString(message interface{}) interface{} {
	intArray, ok := message.([]interface{})
	if !ok {
		return message
	}

	data := make([]byte, len(intArray))
	for i, v := range intArray {
		f64, ok := v.(float64)
		if !ok {
			return message
		}
		data[i] = byte(f64)
	}
	return string(data)
}

func fixNewLines(message string) string {
	return strings.Replace(message, "|", "\n", -1)
}

func removeBlacklistedProperties(m map[string]interface{}) {
	for _, p := range blacklistedProperties {
		delete(m, p)
	}
}
