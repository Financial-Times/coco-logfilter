package main

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Financial-Times/coco-logfilter"
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
		"_SELINUX_CONTEXT",
		"__REALTIME_TIMESTAMP",
		"_PID",
	}

	blacklistedUnits = map[string]bool{
		"splunk-forwarder.service":   true,
		//"docker.service":             true,
		"diamond.service":            true,
		"logstash-forwarder.service": true,
		"kubelet.service": true,
		"flanneld.service": true,

	}

	blacklistedStrings = []string{
		"transaction_id=SYNTHETIC-REQ",
	}

	blacklistedSyslogIds = map[string]bool{
		"dockerd": true,
	}

	blacklistedContainerTags = []string{
		"gcr.io/google_containers/heapster",
		"gcr.io/google_containers/kubedns-amd64",
		"gcr.io/google_containers/addon-resizer",
	}

	propertyMapping = map[string]string{
		"_SYSTEMD_UNIT": "SYSTEMD_UNIT",
		"_MACHINE_ID":   "MACHINE_ID",
		"_HOSTNAME":     "HOSTNAME",
	}
)

var environmentTag *string = new(string)

func main() {
	*environmentTag = os.Getenv("ENV")

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

		syslogId := m["SYSLOG_IDENTIFIER"]
		if syslogIdString, ok := syslogId.(string); ok {
			if blacklistedSyslogIds[syslogIdString] {
				continue
			}
		}

		containerTag := m["CONTAINER_TAG"]
		if containerTagString, ok := containerTag.(string); ok {
			if (containsBlacklistedString(containerTagString, blacklistedContainerTags)) {
				continue
			}
		}

		message := fixBytesToString(m["MESSAGE"]).(string)

		if containsBlacklistedString(message, blacklistedStrings) {
			continue
		}

		munge(m, message)
		removeBlacklistedProperties(m)
		renameProperties(m)
		enc.Encode(m)
	}
}

func containsBlacklistedString(message string, blacklistedStrings []string) bool {
	for _, blacklistedString := range blacklistedStrings {
		if strings.Contains(message, blacklistedString) {
			return true
		}
	}
	return false
}

func munge(m map[string]interface{}, message string) {

	m["platform"] = "up-coco"
	if *environmentTag != "" {
		m["environment"] = *environmentTag
	}

	podName := extractPodName(m["CONTAINER_NAME"])
	if podName != "" {
		m["POD_NAME"] = podName
	}

	serviceName := extractServiceName(m["CONTAINER_NAME"])
	if serviceName != "" && serviceName != "POD" {
		m["SERVICE_NAME"] = serviceName
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

func extractServiceName(containerTag interface{}) string {
	containerNameSplitByUnderscores := splitByUnderscores(containerTag)

	if len(containerNameSplitByUnderscores) >= 1 {
		stringArray := strings.Split(containerNameSplitByUnderscores[1], ".")
		return stringArray[0]
	}

	return ""
}

func extractPodName(containerTag interface{}) string {
	containerNameSplitByUnderscores := splitByUnderscores(containerTag)

	if len(containerNameSplitByUnderscores) > 2 {
		return containerNameSplitByUnderscores[2]
	}

	return ""
}

func splitByUnderscores(i interface{}) []string {
	if s, ok := i.(string); ok {
		items := strings.Split(s, "_")
		return items
	}

	return []string{}
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

func renameProperties(m map[string]interface{}) {
	for p, r := range propertyMapping {
		value := m[p]
		if value != nil {
			delete(m, p)
			m[r] = value
		}
	}

}
