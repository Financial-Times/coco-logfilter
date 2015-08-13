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

	// ERROR [2015-08-08 00:18:05,872] com.ft.binaryingester.health.BinaryWriterDependencyHealthCheck:  Exception during dependency version check|[dw-18 - GET /__health]! com.sun.jersey.api.client.ClientHandlerException: java.net.SocketTimeoutException: Read timed out|! at com.sun.jersey.client.apache4.ApacheHttpClient4Handler.handle(ApacheHttpClient4Handler.java:187) ~[app.jar:na]|! at com.sun.jersey.api.client.filter.GZIPContentEncodingFilter.handle(GZIPContentEncodingFilter.java:120) ~[app.jar:na]|! at com.sun.jersey.api.client.Client.handle(Client.java:652) ~[app.jar:na]|! at com.ft.jerseyhttpwrapper.ResilientClient.handle(ResilientClient.java:142) ~[app.jar:na]|! at com.sun.jersey.api.client.WebResource.handle(WebResource.java:682) ~[app.jar:na]|! at com.sun.jersey.api.client.WebResource.access$200(WebResource.java:74) ~[app.jar:na]|! at com.sun.jersey.api.client.WebResource$Builder.get(WebResource.java:509) ~[app.jar:na]|! at com.ft.binaryingester.health.BinaryWriterDependencyHealthCheck.checkAdvanced(BinaryWriterDependencyHealthCheck.java:48) ~[app.jar:na]|! at com.ft.platform.dropwizard.AdvancedHealthCheck.executeAdvanced(AdvancedHealthCheck.java:21) [app.jar:na]|! at com.ft.platform.dropwizard.HealthChecks.runAdvancedHealthChecksIn(HealthChecks.java:22) [app.jar:na]|! at com.ft.platform.dropwizard.AdvancedHealthChecksRunner.run(AdvancedHealthChecksRunner.java:36) [app.jar:na]|! at com.ft.platform.dropwizard.AdvancedHealthCheckServlet.doGet(AdvancedHealthCheckServlet.java:40) [app.jar:na]|! at javax.servlet.http.HttpServlet.service(HttpServlet.java:735) [app.jar:na]|! at javax.servlet.http.HttpServlet.service(HttpServlet.java:848) [app.jar:na]|! at io.dropwizard.jetty.NonblockingServletHolder.handle(NonblockingServletHolder.java:49) [app.jar:na]|! at org.eclipse.jetty.servlet.ServletHandler$CachedChain.doFilter(ServletHandler.java:1515) [app.jar:na]|! at org.eclipse.jetty.servlets.UserAgentFilter.doFilter(UserAgentFilter.java:83) [app.jar:na]|! at org.eclipse.jetty.servlets.GzipFilter.doFilter(GzipFilter.java:34
	re4 = regexp.MustCompile("\\s*([A-Z]{4,5})\\s*\\[([0-9\\-:,\\s]*)\\] (.*)")
)

func Extract(message string) (v interface{}, extracted bool) {
	v, extracted = extractAccEntry(message)
	if extracted {
		return v, extracted
	}
	return extractAppEntry(message)
}

func extractAccEntry(msg string) (ent accessEntry, extracted bool) {
	ent = accessEntry{}
	matches := re1.FindStringSubmatch(msg)
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

	matches = re2.FindStringSubmatch(msg)
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

func extractAppEntry(msg string) (ent appEntry, extracted bool) {
	ent = appEntry{}
	matches := re4.FindStringSubmatch(msg)
	if len(matches) == 4 {
		ent.Level = matches[1]
		ent.Timestamp = matches[2]
		ent.Message = matches[3]
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
	if s == "-" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("failed to parse %s as integer. this is a bug", s)
	}
	return i
}

type accessEntry struct {
	RemoteServer string `json:"remote-server,omitempty"`
	User         string `json:"user,omitempty"`
	Password     string `json:"password,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	Method       string `json:"method,omitempty"`
	Url          string `json:"url,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
	Status       int    `json:"status,omitempty"`
	LenBytes     int    `json:"byte-length,omitempty"`
	Referrer     string `json:"referrer,omitempty"`
	UserAgent    string `json:"user-agent,omitempty"`
	TimeMs       int    `json:"time-ms,omitempty"`
}

type appEntry struct {
	Level     string `json:"level,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Message   string `json:"msg,omitempty"`
}
