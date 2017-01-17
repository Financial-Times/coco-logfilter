package logfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoMirrorLogExample(t *testing.T) {
	assert := assert.New(t)

	in := `172.31.30.229 - - [19/Jun/2015:09:24:24 +0000] "GET /v1/images/f467d023d63178a6686daab33049b7fec024f88e5b64898e9c89dafaaa4e1d8a/ancestry HTTP/1.1" 200 1836 "-" "docker/1.5.0 go/go1.3.3 git-commit/a8a31ef-dirty kernel/3.19.3 os/linux arch/amd64"`

	out, ok := extractAccEntry(in)

	if !ok {
		t.Fatal("failed to extract values")
	}

	assert.Equal("172.31.30.229", out.RemoteServer)
	assert.Equal("19/Jun/2015:09:24:24 +0000", out.Timestamp)
	assert.Equal("GET", out.Method)
	assert.Equal("/v1/images/f467d023d63178a6686daab33049b7fec024f88e5b64898e9c89dafaaa4e1d8a/ancestry", out.URL)
	assert.Equal("HTTP/1.1", out.Protocol)
	assert.Equal(200, out.Status)
	assert.Equal(1836, out.LenBytes)
	assert.Equal(`docker/1.5.0 go/go1.3.3 git-commit/a8a31ef-dirty kernel/3.19.3 os/linux arch/amd64`, out.UserAgent)

	// TODO:
}

func TestMethodeAPIExample(t *testing.T) {
	assert := assert.New(t)

	in := `127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] "GET /eom-file/all/e09b49d6-e1fa-11e4-bb7f-00144feab7de HTTP/1.1" 200 53706 919 919`

	out, ok := extractAccEntry(in)
	if !ok {
		t.Fatal("failed to extract values")
	}

	assert.Equal("127.0.0.1", out.RemoteServer)
	assert.Equal("21/Apr/2015:12:15:34 +0000", out.Timestamp)
	assert.Equal("GET", out.Method)
	assert.Equal("/eom-file/all/e09b49d6-e1fa-11e4-bb7f-00144feab7de", out.URL)
	assert.Equal("HTTP/1.1", out.Protocol)
	assert.Equal(200, out.Status)
	assert.Equal(53706, out.LenBytes)
	assert.Equal("", out.UserAgent)
	// TODO:
}

func TestCmsNotifierPostExample(t *testing.T) {
	assert := assert.New(t)

	in := `172.17.42.1 -  -  [24/Jun/2015:11:09:36 +0000] "POST /notify HTTP/1.1" 500 - "-" "curl/7.42.0" 2197`
	out, ok := extractAccEntry(in)

	if !ok {
		t.Fatal("failed to extract values")
	}

	assert.Equal("172.17.42.1", out.RemoteServer)
	assert.Equal("24/Jun/2015:11:09:36 +0000", out.Timestamp)
	assert.Equal("POST", out.Method)
	assert.Equal("/notify", out.URL)
	assert.Equal("HTTP/1.1", out.Protocol)
	assert.Equal(500, out.Status)
	assert.Equal(0, out.LenBytes)
	assert.Equal("curl/7.42.0", out.UserAgent)
	// TODO:
}

func TestExtractNeoExample(t *testing.T) {
	assert := assert.New(t)

	in := `172.24.3.248 - - [18/Aug/2016:09:51:35 +0000] "POST /db/data/cypher HTTP/1.1" 200 51 "-" "neoism" 77`
	out, ok := extractAccEntry(in)

	if !ok {
		t.Fatal("failed to extract values")
	}

	assert.Equal("172.24.3.248", out.RemoteServer)
	assert.Equal("18/Aug/2016:09:51:35 +0000", out.Timestamp)
	assert.Equal("POST", out.Method)
	assert.Equal("/db/data/cypher", out.URL)
	assert.Equal("HTTP/1.1", out.Protocol)
	assert.Equal(200, out.Status)
	assert.Equal(51, out.LenBytes)
	assert.Equal("neoism", out.UserAgent)
	assert.Equal(77, out.TimeMs)

	t.Logf("%+v\n", out)
}

func TestExtractAppEntry(t *testing.T) {
	var tests = []struct {
		message string
		level   string
	}{
		{`ERROR [2015-08-07 09:03:45,581] kafka.utils.Utils$: transaction_id=tid_lYpxZctRHb_kafka_bridge fetching topic metadata for topics [Set(NativeCmsPublicationEvents)] from broker [ArrayBuffer(id:0,host:172.23.219.136,port:9092)] failed|[dw-3910 - POST /notify]! kafka.common.KafkaException: fetching topic metadata for topics [Set(NativeCmsPublicationEvents)] from broker [ArrayBuffer(id:0,host:172.23.219.136,port:9092)] failed|! at kafka.client.ClientUtils$.fetchTopicMetadata(ClientUtils.scala:67) ~[app.jar:0.0.1-SNAPSHOT]|! at kafka.producer.BrokerPartitionInfo.updateInfo(BrokerPartitionInfo.scala:82) ~[app.jar:0.0.1-SNAPSHOT]|! at kafka.producer.async.DefaultEventHandler$$anonfun$handle$2.apply$mcV$sp(DefaultEventHandler.scala:78) ~[app.jar:0.0.1-SNAPSHOT]|! at kafka.utils.Utils$.swallow(Utils.scala:167) [app.jar:0.0.1-SNAPSHOT]|! at kafka.utils.Logging$class.swallowError(Logging.scala:106) [app.jar:0.0.1-SNAPSHOT]|! at kafka.utils.Utils$.swallowError(Utils.scala:46) [app.jar:0.0.1-SNAPSHOT]|! at kafka.producer.async.DefaultEventHandler.handle(DefaultEventHandler.scala:78) [app.jar:0.0.1-SNAPSHOT]|! at kafka.producer.Producer.send(Producer.scala:76) [app.jar:0.0.1-SNAPSHOT]|! at kafka.javaapi.producer.Producer.send(Producer.scala:33) [app.jar:0.0.1-SNAPSHOT]|! at com.ft.cmsnotifier.service.KafkaMessageProducer.produceNotifyEvent(KafkaMessageProducer.java:30) [app.jar:0.0.1-SNAPSHOT]|! at com.ft.cmsnotifier.resources.CmsNotifierResource.produceEventForContent(CmsNotifierResource.java:59) [app.jar:0.0.1-SNAPSHOT]|! at com.ft.cmsnotifier.resources.CmsNotifierResource.importContent(CmsNotifierResource.java:43) [app.jar:0.0.1-SNAPSHOT]|! at sun.reflect.GeneratedMethodAccessor15.invoke(Unknown Source) ~[na:na]|! at sun.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43) ~[na:1.8.0_51]|! at java.lang.reflect.Method.invoke(Method.java:497) ~[na:1.8.0_51]|! at com.sun.jersey.spi.container.JavaMethodInvokerFactory$1.invoke(JavaMethodInvokerFactory.java:60) [app.jar:0.0.1-SNAPSHOT]|! at com.sun.je`, `ERROR`},
		{`INFO  [2015-08-13 07:59:00,019] com.ft.platform.dropwizard.HealthCheckPageData:  event="advancedHealthCheck", action="detail", name="BinaryIngesterService", checkName="CanConnectToZooKeeper", ok="true", checkOutput="ZooKeeper instance responded, at least 1 are up."|[dw-19 - GET /__health]`, `INFO`},
		{`DEBUG [2015-08-12 14:14:34,421] io.dropwizard.setup.AdminEnvironment:  health checks = [CanConnectToKafka, CanConnectToZooKeeper, KafkaIsNotLagging, binary-writer ping, binary-writer version, deadlocks]|[main]`, `DEBUG`},
		{`WARN  [2015-08-12 14:14:35,004] kafka.consumer.ZookeeperConsumerConnector:  [content_d414a887c7ca-1439388871813-b0f5b10f], No broker partitions consumed by consumer thread content_d414a887c7ca-1439388871813-b0f5b10f-1 for topic CmsPublicationEvents|[pool-5-thread-1]`, `WARN`},
		{`INFO  [2015-08-13 08:01:28,393] com.ft.platform.dropwizard.HealthCheckPageData:  event="advancedHealthCheck", action="detail", name="CMSNotifierApplication", checkName="Can connect to Kafka on topic: NativeCmsPublicationEvents", ok="true", checkOutput="Kafka is connected and topic is present"|[dw-19 - GET /__health]`, `INFO`},
	}

	for _, test := range tests {
		appEntry, ok := extractAppEntry(test.message)
		if !ok {
			t.Fatalf("failed to extract values '%s'", test.message)
		}
		if appEntry.Level != test.level {
			t.Errorf("message: %s\nexpected level %s, actual level %s", test.message, test.level, appEntry.Level)
		}
	}
}

func TestExtractPamEntity(t *testing.T) {
	var tests = []struct {
		message       string
		UUID          string
		ReadEnv       string
		TransactionID string
		PublishDate   string
		PublishOk     string
		Duration      string
		Endpoint      string
	}{
		{`[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 readEnv=prod-uk transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=content`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"prod-uk",
			"tid_28pbiavoqs",
			"1450692093737000000",
			"true",
			"6",
			"content"},
		{`[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 readEnv=prod-uk transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=notifications-push`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"prod-uk",
			"tid_28pbiavoqs",
			"1450692093737000000",
			"true",
			"6",
			"notifications-push"},
	}

	for _, test := range tests {
		pamEntity, ok := extractPamEntity(test.message)
		if !ok {
			t.Fatalf("failed to extract values '%s'", test.message)
		}
		if pamEntity.UUID != test.UUID {
			t.Errorf("message: %s\nexpected UUID %s, actual UUID %s", test.message, test.UUID, pamEntity.UUID)
		}
		if pamEntity.ReadEnv != test.ReadEnv {
			t.Errorf("message: %s\nexpected ReadEnv %s, actual ReadEnv %s", test.message, test.ReadEnv, pamEntity.ReadEnv)
		}
		if pamEntity.TransactionID != test.TransactionID {
			t.Errorf("message: %s\nexpected transaction_id %s, actual transaction_id %s", test.message, test.TransactionID, pamEntity.TransactionID)
		}
		if pamEntity.PublishDate != test.PublishDate {
			t.Errorf("message: %s\nexpected publishDate %s, actual publishDate %s", test.message, test.PublishDate, pamEntity.PublishDate)
		}
		if pamEntity.PublishOk != test.PublishOk {
			t.Errorf("message: %s\nexpected publishOk %s, actual publishOk %s", test.message, test.PublishOk, pamEntity.PublishOk)
		}
		if pamEntity.Duration != test.Duration {
			t.Errorf("message: %s\nexpected duration %s, actual duration %s", test.message, test.Duration, pamEntity.Duration)
		}
		if pamEntity.Endpoint != test.Endpoint {
			t.Errorf("message: %s\nexpected endpoint %s, actual endpoint %s", test.message, test.Endpoint, pamEntity.Endpoint)
		}
	}
}

func TestExtractOldPamEntity(t *testing.T) {
	var tests = []struct {
		message       string
		UUID          string
		TransactionID string
		PublishDate   string
		PublishOk     string
		Duration      string
		Endpoint      string
	}{
		{`[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=content`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"tid_28pbiavoqs",
			"1450692093737000000",
			"true",
			"6",
			"content"},
		{`[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=notifications-push`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"tid_28pbiavoqs",
			"1450692093737000000",
			"true",
			"6",
			"notifications-push"},
	}

	for _, test := range tests {
		pamEntity, ok := extractOldPamEntity(test.message)
		if !ok {
			t.Fatalf("failed to extract values '%s'", test.message)
		}
		if pamEntity.UUID != test.UUID {
			t.Errorf("message: %s\nexpected UUID %s, actual UUID %s", test.message, test.UUID, pamEntity.UUID)
		}
		if pamEntity.TransactionID != test.TransactionID {
			t.Errorf("message: %s\nexpected transaction_id %s, actual transaction_id %s", test.message, test.TransactionID, pamEntity.TransactionID)
		}
		if pamEntity.PublishDate != test.PublishDate {
			t.Errorf("message: %s\nexpected publishDate %s, actual publishDate %s", test.message, test.PublishDate, pamEntity.PublishDate)
		}
		if pamEntity.PublishOk != test.PublishOk {
			t.Errorf("message: %s\nexpected publishOk %s, actual publishOk %s", test.message, test.PublishOk, pamEntity.PublishOk)
		}
		if pamEntity.Duration != test.Duration {
			t.Errorf("message: %s\nexpected duration %s, actual duration %s", test.message, test.Duration, pamEntity.Duration)
		}
		if pamEntity.Endpoint != test.Endpoint {
			t.Errorf("message: %s\nexpected endpoint %s, actual endpoint %s", test.message, test.Endpoint, pamEntity.Endpoint)
		}
	}
}

func TestExtractSlaPamEntity(t *testing.T) {
	var tests = []struct {
		message        string
		UUID           string
		MetPublishSLA  string
		OkEnvironments string
		TransactionID  string
		PublishDate    string
	}{
		{`[slaMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 metPublishSLA=true okEnvironments=[prod-uk,prod-us] transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 duration=6 endpoint=content`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"true",
			"[prod-uk,prod-us]",
			"tid_28pbiavoqs",
			"1450692093737000000"},
		{`[slaMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 metPublishSLA=false okEnvironments=[prod-uk] transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 duration=6 endpoint=notifications-push`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"false",
			"[prod-uk]",
			"tid_28pbiavoqs",
			"1450692093737000000"},
	}

	for _, test := range tests {
		pamEntity, ok := extractSlaPamEntity(test.message)
		if !ok {
			t.Fatalf("failed to extract values '%s'", test.message)
		}
		if pamEntity.UUID != test.UUID {
			t.Errorf("message: %s\nexpected UUID %s, actual UUID %s", test.message, test.UUID, pamEntity.UUID)
		}
		if pamEntity.MetPublishSLA != test.MetPublishSLA {
			t.Errorf("message: %s\nexpected metPublishSLA %s, actual metPublishSLA %s", test.message, test.MetPublishSLA, pamEntity.MetPublishSLA)
		}
		if pamEntity.OkEnvironments != test.OkEnvironments {
			t.Errorf("message: %s\nexpected okEnvironments %s, actual okEnvironments %s", test.message, test.OkEnvironments, pamEntity.OkEnvironments)
		}
		if pamEntity.TransactionID != test.TransactionID {
			t.Errorf("message: %s\nexpected transaction_id %s, actual transaction_id %s", test.message, test.TransactionID, pamEntity.TransactionID)
		}
		if pamEntity.PublishDate != test.PublishDate {
			t.Errorf("message: %s\nexpected publishDate %s, actual publishDate %s", test.message, test.PublishDate, pamEntity.PublishDate)
		}
	}
}

func TestExtractVarnishEntity(t *testing.T) {
	var tests = []struct {
		message       string
		AuthUser      string
		URI           string
		Status        string
		Resptime      string
		UserAgent     string
		TransactionID string
	}{
		{`82.136.1.214, 172.24.88.199 next 14/Jun/2016:08:24:42 /__enriched-content-read-api/enrichedcontent/409ba29e-b7f8-417f-9847-f3332aa064a6 200 65526 "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36" transaction_id=tid_pete-doing-some-tests hit`,
			"next",
			"/__enriched-content-read-api/enrichedcontent/409ba29e-b7f8-417f-9847-f3332aa064a6",
			"200",
			"65526",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
			"tid_pete-doing-some-tests"},
	}

	for _, test := range tests {
		varnishEntity, ok := extractVarnishEntity(test.message)
		if !ok {
			t.Fatalf("failed to extract values '%s'", test.message)
		}
		if varnishEntity.URI != test.URI {
			t.Errorf("message: %s\nexpected uri %s, actual uri %s", test.message, test.URI, varnishEntity.URI)
		}
		if varnishEntity.Status != test.Status {
			t.Errorf("message: %s\nexpected status %s, actual status %s", test.message, test.Status, varnishEntity.Status)
		}
		if varnishEntity.Resptime != test.Resptime {
			t.Errorf("message: %s\nexpected respTime %s, actual respTime %s", test.message, test.Resptime, varnishEntity.Resptime)
		}
	}
}
