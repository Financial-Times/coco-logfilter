package logfilter

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	// 127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] "GET /eom-file/all/e09b49d6-e1fa-11e4-bb7f-00144feab7de HTTP/1.1" 200 53706 919 919
	re1 = regexp.MustCompile(`^([\d.]+) (\S+) (\S+) \[([\w:/]+\s[+-]\d{4})\] \"(.+?)\" (\d{3}) (\d+|-) (\d+) (\d+)`)

	// 172.31.30.229 - - [19/Jun/2015:09:24:24 +0000] "GET /foo/bar/baz HTTP/1.1" 200 1836 "referrer" "user-agent-123 version 2"
	// 172.31.30.229 - - [19/Jun/2015:09:24:24 +0000] "GET /foo/bar/baz HTTP/1.1" 200 1836 "referrer" "user-agent-123 version 2" 1234
	re2 = regexp.MustCompile(`^([\d.]+) +(\S+) +(\S+) +\[([\w:/]+\s[+-]\d{4})\] +\"(.+?)\" +(\d{3}) +(\d+|-) +\"(.+?)\" +\"(.+?)\"( +(\d+|-))?`)

	// ERROR [2015-08-08 00:18:05,872] com.ft.binaryingester.health.BinaryWriterDependencyHealthCheck:  Exception during dependency version check|[dw-18 - GET /__health]! com.sun.jersey.api.client.ClientHandlerException: java.net.SocketTimeoutException: Read timed out|! at com.sun.jersey.client.apache4.ApacheHttpClient4Handler.handle(ApacheHttpClient4Handler.java:187) ~[app.jar:na]|! at com.sun.jersey.api.client.filter.GZIPContentEncodingFilter.handle(GZIPContentEncodingFilter.java:120) ~[app.jar:na]|! at com.sun.jersey.api.client.Client.handle(Client.java:652) ~[app.jar:na]|! at com.ft.jerseyhttpwrapper.ResilientClient.handle(ResilientClient.java:142) ~[app.jar:na]|! at com.sun.jersey.api.client.WebResource.handle(WebResource.java:682) ~[app.jar:na]|! at com.sun.jersey.api.client.WebResource.access$200(WebResource.java:74) ~[app.jar:na]|! at com.sun.jersey.api.client.WebResource$Builder.get(WebResource.java:509) ~[app.jar:na]|! at com.ft.binaryingester.health.BinaryWriterDependencyHealthCheck.checkAdvanced(BinaryWriterDependencyHealthCheck.java:48) ~[app.jar:na]|! at com.ft.platform.dropwizard.AdvancedHealthCheck.executeAdvanced(AdvancedHealthCheck.java:21) [app.jar:na]|! at com.ft.platform.dropwizard.HealthChecks.runAdvancedHealthChecksIn(HealthChecks.java:22) [app.jar:na]|! at com.ft.platform.dropwizard.AdvancedHealthChecksRunner.run(AdvancedHealthChecksRunner.java:36) [app.jar:na]|! at com.ft.platform.dropwizard.AdvancedHealthCheckServlet.doGet(AdvancedHealthCheckServlet.java:40) [app.jar:na]|! at javax.servlet.http.HttpServlet.service(HttpServlet.java:735) [app.jar:na]|! at javax.servlet.http.HttpServlet.service(HttpServlet.java:848) [app.jar:na]|! at io.dropwizard.jetty.NonblockingServletHolder.handle(NonblockingServletHolder.java:49) [app.jar:na]|! at org.eclipse.jetty.servlet.ServletHandler$CachedChain.doFilter(ServletHandler.java:1515) [app.jar:na]|! at org.eclipse.jetty.servlets.UserAgentFilter.doFilter(UserAgentFilter.java:83) [app.jar:na]|! at org.eclipse.jetty.servlets.GzipFilter.doFilter(GzipFilter.java:34
	re4 = regexp.MustCompile(`([A-Z]{4,5})\s{1,2}\[([0-9\-:,\s]*)\] (.*)`)

	// This is for backwards compatability, until we can rip out the old SLA dashboard
	//[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=content
	pamRegexOLD = regexp.MustCompile(`UUID=([\da-f-]*) transaction_id=([\S]+) publishDate=(\d*) publishOk=(\w*) duration=(\d*) endpoint=([\w-]*)`) //[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=content

	//[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 readEnv=prod-uk transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=content
	pamRegex = regexp.MustCompile(`UUID=([\da-f-]*) readEnv=([\w-]*) transaction_id=([\S]+) publishDate=(\d*) publishOk=(\w*) duration=(\d*) endpoint=([\w-]*)`)

	// 172.17.0.1 usr 13/Jun/2016:13:36:23 /test 200 148866 "curl/7.49.1"
	varnishRegex = regexp.MustCompile(`^[\d\.\,\s]+\s+(\S+)\s+[\w:\/]+\s+(\S+)\s+([0-9]{3})\s+([0-9\.]+)\s+\"([\S\s]+)\"\stransaction_id=([\S]+)`)
)

func Extract(message string) (v interface{}, extracted bool, format string) {
	v, extracted = extractAccEntry(message)
	if extracted {
		return v, extracted, "access"
	}
	v, extracted = extractPamEntity(message)
	if extracted {
		return v, extracted, "pam"
	}
	v, extracted = extractOldPamEntity(message)
	if extracted {
		return v, extracted, "oldpam"
	}
	v, extracted = extractVarnishEntity(message)
	if extracted {
		return v, extracted, "varnish"
	}
	v, extracted = extractJsonEntity(message)
	if extracted {
		return v, extracted, "json"
	}
	v, extracted = extractAppEntry(message)
	return v, extracted, "app"
}

func extractJsonEntity(message string) (map[string]interface{}, bool) {
	res := make(map[string]interface{})
	err := json.Unmarshal([]byte(message), &res)
	if err != nil {
		return nil, false
	}
	//the mdc field is added by java json logging library and it is not necessary,
	// so it needs to be removed
	delete(res, "mdc")
	return res, true
}

func extractAccEntry(message string) (ent accessEntry, extracted bool) {
	v, extracted := extractAccEntryRE1(message)
	if extracted {
		return v, extracted
	}
	v, extracted = extractAccEntryRE2(message)
	if extracted {
		return v, extracted
	}
	return
}

func extractAccEntryRE1(msg string) (ent accessEntry, extracted bool) {
	ent = accessEntry{}
	matches := re1.FindStringSubmatch(msg)
	if len(matches) == 10 {
		ent.RemoteServer = matches[1]
		//todo 2 & 3
		ent.Timestamp = matches[4]
		ent.Method, ent.URL, ent.Protocol = methodURLProtocol(matches[5])
		ent.Status = atoi(matches[6])
		ent.LenBytes = atoi(matches[7])
		// todo 8,9,10
		extracted = true
	}
	return
}

func extractAccEntryRE2(msg string) (ent accessEntry, extracted bool) {
	ent = accessEntry{}
	matches := re2.FindStringSubmatch(msg)
	if len(matches) == 12 {
		ent.RemoteServer = matches[1]
		//todo 2 & 3
		ent.Timestamp = matches[4]
		ent.Method, ent.URL, ent.Protocol = methodURLProtocol(matches[5])
		ent.Status = atoi(matches[6])
		ent.LenBytes = atoi(matches[7])
		// todo 8
		ent.UserAgent = matches[9]
		ent.TimeMs = atoi(matches[11])
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

func extractOldPamEntity(msg string) (pam pamEntity, extracted bool) {
	pam = pamEntity{}
	matches := pamRegexOLD.FindStringSubmatch(msg)
	if len(matches) == 7 {
		pam.UUID = matches[1]
		pam.TransactionID = matches[2]
		pam.PublishDate = matches[3]
		pam.PublishOk = matches[4]
		pam.Duration = matches[5]
		pam.Endpoint = matches[6]
		extracted = true
	}
	return
}

func extractPamEntity(msg string) (pam pamEntity, extracted bool) {
	pam = pamEntity{}
	matches := pamRegex.FindStringSubmatch(msg)
	if len(matches) == 8 {
		pam.UUID = matches[1]
		pam.ReadEnv = matches[2]
		pam.TransactionID = matches[3]
		pam.PublishDate = matches[4]
		pam.PublishOk = matches[5]
		pam.Duration = matches[6]
		pam.Endpoint = matches[7]
		extracted = true
	}
	return
}

func extractVarnishEntity(msg string) (varnish varnishEntity, extracted bool) {
	varnish = varnishEntity{}
	matches := varnishRegex.FindStringSubmatch(msg)
	if len(matches) == 7 {
		varnish.AuthUser = matches[1]
		varnish.URI = matches[2]
		varnish.Status = matches[3]
		varnish.Resptime = matches[4]
		varnish.UserAgent = matches[5]
		varnish.TransactionID = matches[6]
		extracted = true
	}
	return
}

func methodURLProtocol(s string) (string, string, string) {
	mup := strings.Split(s, " ")
	if len(mup) != 3 {
		log.Fatalf("failed to split methode, url protocol from  %s.  This is a bug", s)
	}
	return mup[0], mup[1], mup[2]
}

func atoi(s string) int {
	if s == "-" || s == "" {
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
	URL          string `json:"url,omitempty"`
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
	Message   string `json:"-"`
}

type pamEntity struct {
	UUID          string `json:"uuid"`
	ReadEnv       string `json:"readEnv"`
	TransactionID string `json:"transaction_id"`
	PublishDate   string `json:"publishDate"`
	PublishOk     string `json:"publishOk"`
	Duration      string `json:"duration"`
	Endpoint      string `json:"endpoint"`
}

type varnishEntity struct {
	AuthUser      string `json:"authuser,omitempty"`
	URI           string `json:"uri,omitempty"`
	Status        string `json:"status,omitempty"`
	Resptime      string `json:"resptime,omitempty"`
	UserAgent     string `json:"useragent,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
}
