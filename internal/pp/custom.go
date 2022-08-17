package pp

import (
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"strings"

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
	uses := []string{}
	if ku&x509.KeyUsageDigitalSignature > 0 {
		uses = append(uses, "DigitalSignature")
	}
	if ku&x509.KeyUsageContentCommitment > 0 {
		uses = append(uses, "ContentCommitment")
	}
	if ku&x509.KeyUsageKeyEncipherment > 0 {
		uses = append(uses, "KeyEncipherment")
	}
	if ku&x509.KeyUsageDataEncipherment > 0 {
		uses = append(uses, "DataEncipherment")
	}
	if ku&x509.KeyUsageKeyAgreement > 0 {
		uses = append(uses, "KeyAgreement")
	}
	if ku&x509.KeyUsageCertSign > 0 {
		uses = append(uses, "CertSign")
	}
	if ku&x509.KeyUsageCRLSign > 0 {
		uses = append(uses, "CRLSign")
	}
	if ku&x509.KeyUsageEncipherOnly > 0 {
		uses = append(uses, "EncipherOnly")
	}
	if ku&x509.KeyUsageDecipherOnly > 0 {
		uses = append(uses, "DecipherOnly")
	}
	if len(uses) == 0 {
		return "<usage undefined>"
	}
	return strings.Join(uses, "\n")
}

func customASNObjectID(asnOID asn1.ObjectIdentifier) string {
	s := asnOID.String()
	o, ok := oid.Lookup(s)
	v := o.Expl
	if o.Expl == "" {
		v = o.Name
	}
	if !ok {
		return s
	}
	return fmt.Sprintf("%s <%s>", v, s)
}
