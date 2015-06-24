package logfilter

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	// 127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] "GET /eom-file/all/e09b49d6-e1fa-11e4-bb7f-00144feab7de HTTP/1.1" 200 53706 919 919
	re1 = regexp.MustCompile("^([\\d.]+) (\\S+) (\\S+) \\[([\\w:/]+\\s[+-]\\d{4})\\] \"(.+?)\" (\\d{3}) (\\d+|-) (\\d+) (\\d+)")

	// 172.31.30.229 - - [19/Jun/2015:09:24:24 +0000] "GET /foo/bar/baz HTTP/1.1" 200 1836 "referrer" "user-agent-123 version 2"
	re2 = regexp.MustCompile("^([\\d.]+) +(\\S+) +(\\S+) +\\[([\\w:/]+\\s[+-]\\d{4})\\] +\"(.+?)\" +(\\d{3}) +(\\d+|-) +\"(.+?)\" +\"(.+?)\"")

	// 172.17.42.1 -  -  [24/Jun/2015:11:09:36 +0000] "POST /notify HTTP/1.1" 500 - "-" "curl/7.42.0" 2197
	//re3 = regexp.MustCompile("^([\\d.]+) (\\S+) (\\S+) \\[([\\w:/]+\\s[+-]\\d{4})\\] \"(.+?)\" (\\d{3}) (\\d+|-) \"(.+?)\" \"(.+?)\"")
)

func Extract(message string) (ent accessEntry, extracted bool) {
	matches := re1.FindStringSubmatch(message)
	if len(matches) == 10 {
		ent.RemoteServer = matches[1]
		//todo 2 & 3
		ent.Timestamp = matches[4]
		ent.Method, ent.Url, ent.Protocol = methodUrlProtocol(matches[5])
		ent.Status = atoi(matches[6])
		ent.LenBytes = atoi(matches[7])
		// todo 8,9,10
		extracted = true
	}

	matches = re2.FindStringSubmatch(message)
	if len(matches) == 10 {
		ent.RemoteServer = matches[1]
		//todo 2 & 3
		ent.Timestamp = matches[4]
		ent.Method, ent.Url, ent.Protocol = methodUrlProtocol(matches[5])
		ent.Status = atoi(matches[6])
		ent.LenBytes = atoi(matches[7])
		// todo 8,9
		extracted = true
	}

	return
}

/*
func parseTime(s string) time.Time {
	format := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(format, s)
	if err != nil {
		log.Fatalf("failed to parse date %s. This is a bug\n%v\n", s, err)
	}
	return t
}
*/

func methodUrlProtocol(s string) (string, string, string) {
	mup := strings.Split(s, " ")
	if len(mup) != 3 {
		log.Fatalf("failed to split methode, url protocol from  %s.  This is a bug", s)
	}
	return mup[0], mup[1], mup[2]
}

func atoi(s string) int {
	if s=="-" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("failed to parse %s as integer. this is a bug", s)
	}
	return i
}

type accessEntry struct {
	RemoteServer string `json:",omitempty"`
	User         string `json:",omitempty"`
	Password     string `json:",omitempty"`
	Timestamp    string `json:",omitempty"`
	Method       string `json:",omitempty"`
	Url          string `json:",omitempty"`
	Protocol     string `json:",omitempty"`
	Status       int    `json:",omitempty"`
	LenBytes     int    `json:",omitempty"`
	Referrer     string `json:",omitempty"`
	UserAgent    string `json:",omitempty"`
	TimeMs       int    `json:",omitempty"`
}
