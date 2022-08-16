package oid

import (
	"errors"
	"fmt"
	"strings"
)

type OID struct {
	Value         string
	Tag           string
	Expl          string
	Name          string
	Attr          string
	CKM           string
	CKK           string
	CertExtension bool
}

func (o *OID) Decode(lines []string) error {
	if o == nil {
		return errors.New("nil oid")
	}
	for _, line := range lines {
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(string(line), " ", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		if k == "" {
			continue
		}
		v := strings.TrimSpace(kv[1])
		switch k {
		case "OID":
			o.Value = v
		case "TAG":
			o.Tag = v
		case "EXPL":
			o.Expl = strings.Trim(v, `"`)
		case "NAME":
			o.Name = v
		case "ATTR":
			o.Attr = v
		case "CERT_EXTENSION":
			if v == "SUPPORTED" {
				o.CertExtension = true
			}
		case "CKK":
			o.CKK = v
		case "CKM":
			o.CKM = v
		default:
			return fmt.Errorf("unknown field '%s' with value '%s' (oid %s)", k, v, o.Value)
		}
	}
	return nil
}
