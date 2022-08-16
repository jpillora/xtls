package oid

import (
	"bytes"
	_ "embed"
	"log"
	"strings"
	"sync"
)

var (
	//go:embed oids.txt
	oidsRaw []byte
	parsed  sync.Once
	oids    map[string]OID
)

func parseRaw() {
	oids = map[string]OID{}
	blocks := bytes.Split(oidsRaw, []byte{'\n', '\n'})
	for i, block := range blocks {
		if i == 0 {
			// TODO expose metadata as package vars
			// if strings.HasPrefix(line, "CVS_ID") {
			// 	continue
			// }
			continue
		}
		strs := []string{}
		for _, b := range bytes.Split(block, []byte{'\n'}) {
			s := string(b)
			if s == "" || strings.HasPrefix(s, "#") {
				continue
			}
			strs = append(strs, s)
		}
		if len(strs) == 0 {
			continue
		}
		o := OID{}
		if err := o.Decode(strs); err != nil {
			log.Fatalf("parse raw failed: %s", err)
		}
		if o.Value == "" {
			log.Fatalf("oid block with no value: %s", string(block))
		}
		oids[o.Value] = o
	}
}
