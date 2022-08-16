package oid

func Lookup(oid string) (OID, bool) {
	parsed.Do(parseRaw)
	o, ok := oids[oid]
	return o, ok
}
