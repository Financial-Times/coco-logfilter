package logfilter

import (
	"testing"
)

func TestRepoMirrorLogExample(t *testing.T) {
	in := `172.31.30.229 - - [19/Jun/2015:09:24:24 +0000] "GET /v1/images/f467d023d63178a6686daab33049b7fec024f88e5b64898e9c89dafaaa4e1d8a/ancestry HTTP/1.1" 200 1836 "-" "docker/1.5.0 go/go1.3.3 git-commit/a8a31ef-dirty kernel/3.19.3 os/linux arch/amd64"`

	out, ok := Extract(in)
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

	out, ok := Extract(in)
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
	out, ok := Extract(in)
	log.Printf("out status value %v", out.Status)
	if !ok {
		t.Fatal("failed to extract values")
	}
	if out.Status != 500 {
		t.Errorf("expected status %d but got %d\n", 500, out.Status)
	}
	// TODO:
}
