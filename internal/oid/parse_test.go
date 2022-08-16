package oid

import (
	"testing"
)

func TestParse(t *testing.T) {
	block := []string{
		`OID 1.3.6.1.5.5.7.1.1`,
		`TAG id-pe-authorityInfoAccess`,
		`EXPL "Certificate Authority Information Access"`,
		`NAME NSS_OID_X509_AUTH_INFO_ACCESS`,
		`CERT_EXTENSION SUPPORTED`,
	}
	o := OID{}
	if err := o.Decode(block); err != nil {
		t.Fatal(err)
	}
}

func TestLookup1(t *testing.T) {
	_, ok := Lookup("1.3.6.1.4.1.11129.2.4.2")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestLookup2(t *testing.T) {
	o, ok := Lookup("2.5.8.1.1")
	if !ok {
		t.Fatal("not found")
	}
	if o.Expl != "RSA Encryption Algorithm" {
		t.Fatal("mismatch expl")
	}
}
