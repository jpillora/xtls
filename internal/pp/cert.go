package pp

import (
	"crypto/x509"
	"encoding/asn1"
	"fmt"

	"github.com/jpillora/xtls/internal/oid"
)

func customPrinter(val any) string {
	switch v := val.(type) {
	case x509.ExtKeyUsage:
		return customExtKeyUsage(v)
	case x509.KeyUsage:
		return customKeyUsage(v)
	case asn1.ObjectIdentifier:
		return customASNObjectID(v)
	}
	return ""
}

func customExtKeyUsage(ku x509.ExtKeyUsage) string {
	switch ku {
	case x509.ExtKeyUsageAny:
		return "Any"
	case x509.ExtKeyUsageServerAuth:
		return "ServerAuth"
	case x509.ExtKeyUsageClientAuth:
		return "ClientAuth"
	case x509.ExtKeyUsageCodeSigning:
		return "CodeSigning"
	case x509.ExtKeyUsageEmailProtection:
		return "EmailProtection"
	case x509.ExtKeyUsageIPSECEndSystem:
		return "IPSECEndSystem"
	case x509.ExtKeyUsageIPSECTunnel:
		return "IPSECTunnel"
	case x509.ExtKeyUsageIPSECUser:
		return "IPSECUser"
	case x509.ExtKeyUsageTimeStamping:
		return "TimeStamping"
	case x509.ExtKeyUsageOCSPSigning:
		return "OCSPSigning"
	case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
		return "MicrosoftServerGatedCrypto"
	case x509.ExtKeyUsageNetscapeServerGatedCrypto:
		return "NetscapeServerGatedCrypto"
	case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
		return "MicrosoftCommercialCodeSigning"
	case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
		return "MicrosoftKernelCodeSigning"
	}
	return fmt.Sprintf("<unknown usage %d>", ku)
}

func customKeyUsage(ku x509.KeyUsage) string {
	switch ku {
	case 0:
		return "<usage undefined>"
	case x509.KeyUsageDigitalSignature:
		return "DigitalSignature"
	case x509.KeyUsageContentCommitment:
		return "ContentCommitment"
	case x509.KeyUsageKeyEncipherment:
		return "KeyEncipherment"
	case x509.KeyUsageDataEncipherment:
		return "DataEncipherment"
	case x509.KeyUsageKeyAgreement:
		return "KeyAgreement"
	case x509.KeyUsageCertSign:
		return "CertSign"
	case x509.KeyUsageCRLSign:
		return "CRLSign"
	case x509.KeyUsageEncipherOnly:
		return "EncipherOnly"
	case x509.KeyUsageDecipherOnly:
		return "DecipherOnly"
	}
	return fmt.Sprintf("<unknown usage %d>", ku)
}

func customASNObjectID(asnOID asn1.ObjectIdentifier) string {
	s := asnOID.String()
	o, ok := oid.Lookup(s)
	v := o.Name
	if !ok {
		v = "<unknown oid>"
	}
	return v + " " + s
}
