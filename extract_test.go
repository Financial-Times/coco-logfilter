package logfilter

import (
	"testing"
)

func TestRepoMirrorLogExample(t *testing.T) {
	in := `172.31.30.229 - - [19/Jun/2015:09:24:24 +0000] "GET /v1/images/f467d023d63178a6686daab33049b7fec024f88e5b64898e9c89dafaaa4e1d8a/ancestry HTTP/1.1" 200 1836 "-" "docker/1.5.0 go/go1.3.3 git-commit/a8a31ef-dirty kernel/3.19.3 os/linux arch/amd64"`

	out, ok := extractAccEntry(in)

	if !ok {
		t.Fatal("failed to extract values")
	}

	if out.Status != 200 {
		t.Errorf("expected status %d but got %d\n", 200, out.Status)
	}
	// TODO:
}

func TestMethodeAPIExample(t *testing.T) {
	in := `127.0.0.1 - - [21/Apr/2015:12:15:34 +0000] "GET /eom-file/all/e09b49d6-e1fa-11e4-bb7f-00144feab7de HTTP/1.1" 200 53706 919 919`

	out, ok := extractAccEntry(in)
	if !ok {
		t.Fatal("failed to extract values")
	}
	if out.Status != 200 {
		t.Errorf("expected status %d but got %d\n", 200, out.Status)
	}
	// TODO:
}

func TestCmsNotifierPostExample(t *testing.T) {
	in := `172.17.42.1 -  -  [24/Jun/2015:11:09:36 +0000] "POST /notify HTTP/1.1" 500 - "-" "curl/7.42.0" 2197`
	out, ok := extractAccEntry(in)
	t.Logf("out status value %v", out.Status)
	if !ok {
		t.Fatal("failed to extract values")
	}
	if out.Status != 500 {
		t.Errorf("expected status %d but got %d\n", 500, out.Status)
	}
	// TODO:
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
		message        string
		UUID           string
		Transaction_id string
		PublishDate    string
		PublishOk      string
		Duration       string
		Endpoint       string
	}{
		{`[splunkMetrics] 2015/12/21 10:01:37.336610 UUID=08d30fb4-a7b3-11e5-955c-1e1d6de94879 transaction_id=tid_28pbiavoqs publishDate=1450692093737000000 publishOk=true duration=6 endpoint=content`,
			"08d30fb4-a7b3-11e5-955c-1e1d6de94879",
			"tid_28pbiavoqs",
			"1450692093737000000",
			"true",
			"6",
			"content"},
	}

	for _, test := range tests {
		pamEntity, ok := extractPamEntity(test.message)
		if !ok {
			t.Fatalf("failed to extract values '%s'", test.message)
		}
		if pamEntity.UUID != test.UUID {
			t.Errorf("message: %s\nexpected UUID %s, actual UUID %s", test.message, test.UUID, pamEntity.UUID)
		}
		if pamEntity.TransactionId != test.Transaction_id {
			t.Errorf("message: %s\nexpected transaction_id %s, actual transaction_id %s", test.message, test.Transaction_id, pamEntity.TransactionId)
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
