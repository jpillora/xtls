# xtls

TLS utils prints certificates from file, stdin or hostname, more features to come...

[![GoDev](https://img.shields.io/static/v1?label=godoc&message=reference&color=00add8)](https://pkg.go.dev/github.com/jpillora/xtls)
[![CI](https://github.com/jpillora/xtls/workflows/CI/badge.svg)](https://github.com/jpillora/xtls/actions?workflow=CI)

### Features

* `xtls [file-path-to-cert]` print a PEM certificate file
* `cat [file-path-to-cert] | xtls` print a PEM certificate from stdin
* `xtls [domain]` print a certificate from a TLS connection to `domain`

### Install

**Binaries**

[![Releases](https://img.shields.io/github/release/jpillora/xtls.svg)](https://github.com/jpillora/xtls/releases)
[![Releases](https://img.shields.io/github/downloads/jpillora/xtls/total.svg)](https://github.com/jpillora/xtls/releases)

Find [the latest pre-compiled binaries here](https://github.com/jpillora/xtls/releases/latest)

**Source**

```sh
$ go install github.com/jpillora/xtls@latest
```

### Example

(_terminal colours not shown_)

```
xtls google.com

2022/11/30 14:26:55 connect to: google.com
2022/11/30 14:26:55 dialing google.com:443
2022/11/30 14:26:55 tls server name: google.com
2022/11/30 14:26:55 tls version: 772
2022/11/30 14:26:55 tls ciphersuite: 4865
[
  #1 Certificate{
    Raw: 0x30820e1c30820d04a003020102021100... (3.62KB)
    RawTBSCertificate: 0x30820d04a003020102021100ee642cf8... (3.34KB)
    RawSubjectPublicKeyInfo: 0x3059301306072a8648ce3d020106082a... (91B)
    RawSubject: "0·1·0···U····*.google.com"
    RawIssuer:
      0F1·0⇥··U····US1"0···U·⏎
      ··Google·Trust·Services·LLC1·0···U···⏎
      GTS·CA·1C3
    Signature: 0xe12f0a82387178c543258fa55d301bde... (256B)
    SignatureAlgorithm: "SHA256-RSA" (x509.SignatureAlgorithm)
    PublicKeyAlgorithm: "ECDSA" (x509.PublicKeyAlgorithm)
    PublicKey: 0x140001225c0 (interface {})
    Version: 3 (int)
    SerialNumber: "316876404775094258481874580025735671050" (*big.Int)
    Issuer: "CN=GTS·CA·1C3,O=Google·Trust·Services·LLC,C=US" (pkix.Name)
    Subject: "CN=*.google.com" (pkix.Name)
    NotBefore: "2022-11-02T13:43:09Z"
    NotAfter: "2023-01-25T13:43:08Z"
    KeyUsage: DigitalSignature
    Extensions: [
      #1 Extension{
        Id: Certificate Key Usage <2.5.29.15>
        Critical: true
        Value: 0x03020780 (4B)
      }
      #2 Extension{
        Id: Extended Key Usage <2.5.29.37>
        Critical: false
        Value:
          0⏎
          ··+·······
      }
      #3 Extension{
        Id: Certificate Basic Constraints <2.5.29.19>
        Critical: true
        Value: "0·"
      }
      #4 Extension{
        Id: Certificate Subject Key ID <2.5.29.14>
        Critical: false
        Value: 0x04144e17794eaeac2a1d45701aff5618... (22B)
      }
      #5 Extension{
        Id: Certificate Authority Key Identifier <2.5.29.35>
        Critical: false
        Value: 0x301680148a747faf85cdee95cd3d9cd0... (24B)
      }
      #6 Extension{
        Id: Certificate Authority Information Access <1.3.6.1.5.5.7.1.1>
        Critical: false
        Value: 0x305c302706082b06010505073001861b... (94B)
      }
      #7 Extension{
        Id: Certificate Subject Alternate Name <2.5.29.17>
        Critical: false
        Value: 0x308209a0820c2a2e676f6f676c652e63... (2.47KB)
      }
      #8 Extension{
        Id: Certificate Policies <2.5.29.32>
        Critical: false
        Value: 0x30183008060667810c010201300c060a... (26B)
      }
      #9 Extension{
        Id: CRL Distribution Points <2.5.29.31>
        Critical: false
        Value: 0x30333031a02fa02d862b687474703a2f... (53B)
      }
      #10 Extension{
        Id: 1.3.6.1.4.1.11129.2.4.2
        Critical: false
        Value: 0x0481f200f0007600e83ed0da3ef50635... (245B)
      }
    ]
    ExtKeyUsage: [
      #1 ServerAuth
    ]
    BasicConstraintsValid: true
    IsCA: false
    MaxPathLen: -1 (int)
    MaxPathLenZero: false
    SubjectKeyId: 0x4e17794eaeac2a1d45701aff56189a5a... (20B)
    AuthorityKeyId: 0x8a747faf85cdee95cd3d9cd0e24614f3... (20B)
    OCSPServer: [
      #1 "http://ocsp.pki.goog/gts1c3"
    ]
    IssuingCertificateURL: [
      #1 "http://pki.goog/repo/certs/gts1c3.der"
    ]
    DNSNames: [
      #1 "*.google.com"
      #2 "*.appengine.google.com"
      #3 "*.bdn.dev"
      #4 "*.origin-test.bdn.dev"
      #5 "*.cloud.google.com"
      #6 "*.crowdsource.google.com"
      #7 "*.datacompute.google.com"
      #8 "*.google.ca"
      #9 "*.google.cl"
      #10 "*.google.co.in"
      #11 "*.google.co.jp"
      #12 "*.google.co.uk"
      #13 "*.google.com.ar"
      #14 "*.google.com.au"
      #15 "*.google.com.br"
      #16 "*.google.com.co"
      #17 "*.google.com.mx"
      #18 "*.google.com.tr"
      #19 "*.google.com.vn"
      #20 "*.google.de"
      #21 "*.google.es"
      #22 "*.google.fr"
      #23 "*.google.hu"
      #24 "*.google.it"
      #25 "*.google.nl"
      #26 "*.google.pl"
      #27 "*.google.pt"
      #28 "*.googleadapis.com"
      #29 "*.googleapis.cn"
      #30 "*.googlevideo.com"
      #31 "*.gstatic.cn"
      #32 "*.gstatic-cn.com"
      #33 "googlecnapps.cn"
      #34 "*.googlecnapps.cn"
      #35 "googleapps-cn.com"
      #36 "*.googleapps-cn.com"
      #37 "gkecnapps.cn"
      #38 "*.gkecnapps.cn"
      #39 "googledownloads.cn"
      #40 "*.googledownloads.cn"
      #41 "recaptcha.net.cn"
      #42 "*.recaptcha.net.cn"
      #43 "recaptcha-cn.net"
      #44 "*.recaptcha-cn.net"
      #45 "widevine.cn"
      #46 "*.widevine.cn"
      #47 "ampproject.org.cn"
      #48 "*.ampproject.org.cn"
      #49 "ampproject.net.cn"
      #50 "*.ampproject.net.cn"
      #51 "google-analytics-cn.com"
      #52 "*.google-analytics-cn.com"
      #53 "googleadservices-cn.com"
      #54 "*.googleadservices-cn.com"
      #55 "googlevads-cn.com"
      #56 "*.googlevads-cn.com"
      #57 "googleapis-cn.com"
      #58 "*.googleapis-cn.com"
      #59 "googleoptimize-cn.com"
      #60 "*.googleoptimize-cn.com"
      #61 "doubleclick-cn.net"
      #62 "*.doubleclick-cn.net"
      #63 "*.fls.doubleclick-cn.net"
      #64 "*.g.doubleclick-cn.net"
      #65 "doubleclick.cn"
      #66 "*.doubleclick.cn"
      #67 "*.fls.doubleclick.cn"
      #68 "*.g.doubleclick.cn"
      #69 "dartsearch-cn.net"
      #70 "*.dartsearch-cn.net"
      #71 "googletraveladservices-cn.com"
      #72 "*.googletraveladservices-cn.com"
      #73 "googletagservices-cn.com"
      #74 "*.googletagservices-cn.com"
      #75 "googletagmanager-cn.com"
      #76 "*.googletagmanager-cn.com"
      #77 "googlesyndication-cn.com"
      #78 "*.googlesyndication-cn.com"
      #79 "*.safeframe.googlesyndication-cn.com"
      #80 "app-measurement-cn.com"
      #81 "*.app-measurement-cn.com"
      #82 "gvt1-cn.com"
      #83 "*.gvt1-cn.com"
      #84 "gvt2-cn.com"
      #85 "*.gvt2-cn.com"
      #86 "2mdn-cn.net"
      #87 "*.2mdn-cn.net"
      #88 "googleflights-cn.net"
      #89 "*.googleflights-cn.net"
      #90 "admob-cn.com"
      #91 "*.admob-cn.com"
      #92 "googlesandbox-cn.com"
      #93 "*.googlesandbox-cn.com"
      #94 "*.gstatic.com"
      #95 "*.metric.gstatic.com"
      #96 "*.gvt1.com"
      #97 "*.gcpcdn.gvt1.com"
      #98 "*.gvt2.com"
      #99 "*.gcp.gvt2.com"
      #100 "*.url.google.com"
      #101 "*.youtube-nocookie.com"
      #102 "*.ytimg.com"
      #103 "android.com"
      #104 "*.android.com"
      #105 "*.flash.android.com"
      #106 "g.cn"
      #107 "*.g.cn"
      #108 "g.co"
      #109 "*.g.co"
      #110 "goo.gl"
      #111 "www.goo.gl"
      #112 "google-analytics.com"
      #113 "*.google-analytics.com"
      #114 "google.com"
      #115 "googlecommerce.com"
      #116 "*.googlecommerce.com"
      #117 "ggpht.cn"
      #118 "*.ggpht.cn"
      #119 "urchin.com"
      #120 "*.urchin.com"
      #121 "youtu.be"
      #122 "youtube.com"
      #123 "*.youtube.com"
      #124 "youtubeeducation.com"
      #125 "*.youtubeeducation.com"
      #126 "youtubekids.com"
      #127 "*.youtubekids.com"
      #128 "yt.be"
      #129 "*.yt.be"
      #130 "android.clients.google.com"
      #131 "developer.android.google.cn"
      #132 "developers.android.google.cn"
      #133 "source.android.google.cn"
    ]
    PermittedDNSDomainsCritical: false
    CRLDistributionPoints: [
      #1 "http://crls.pki.goog/gts1c3/QqFxbi9M48c.crl"
    ]
    PolicyIdentifiers: [
      #1 2.23.140.1.2.1
      #2 1.3.6.1.4.1.11129.2.5.3
    ]
  }
  #2 Certificate{
    Raw: 0x308205963082037ea003020102020d02... (1.43KB)
    RawTBSCertificate: 0x3082037ea003020102020d0203bc5359... (898B)
    RawSubjectPublicKeyInfo: 0x30820122300d06092a864886f70d0101... (294B)
    RawSubject:
      0F1·0⇥··U····US1"0···U·⏎
      ··Google·Trust·Services·LLC1·0···U···⏎
      GTS·CA·1C3
    RawIssuer:
      0G1·0⇥··U····US1"0···U·⏎
      ··Google·Trust·Services·LLC1·0···U····GTS·Root·R1
    Signature: 0x897dac205c0c3cbe9aa857951bb4aefa... (512B)
    SignatureAlgorithm: "SHA256-RSA" (x509.SignatureAlgorithm)
    PublicKeyAlgorithm: "RSA" (x509.PublicKeyAlgorithm)
    PublicKey: 0x1400010e180 (interface {})
    Version: 3 (int)
    SerialNumber: "159612451717983579589660725350" (*big.Int)
    Issuer: "CN=GTS·Root·R1,O=Google·Trust·Services·LLC,C=US" (pkix.Name)
    Subject: "CN=GTS·CA·1C3,O=Google·Trust·Services·LLC,C=US" (pkix.Name)
    NotBefore: "2020-08-13T00:00:42Z"
    NotAfter: "2027-09-30T00:00:42Z"
    KeyUsage:
      DigitalSignature
      CertSign
      CRLSign
    Extensions: [
      #1 Extension{
        Id: Certificate Key Usage <2.5.29.15>
        Critical: true
        Value: 0x03020186 (4B)
      }
      #2 Extension{
        Id: Extended Key Usage <2.5.29.37>
        Critical: false
        Value: "0···+·········+·······"
      }
      #3 Extension{
        Id: Certificate Basic Constraints <2.5.29.19>
        Critical: true
        Value: 0x30060101ff020100 (8B)
      }
      #4 Extension{
        Id: Certificate Subject Key ID <2.5.29.14>
        Critical: false
        Value: 0x04148a747faf85cdee95cd3d9cd0e246... (22B)
      }
      #5 Extension{
        Id: Certificate Authority Key Identifier <2.5.29.35>
        Critical: false
        Value: 0x30168014e4af2b26711a2b4827852f52... (24B)
      }
      #6 Extension{
        Id: Certificate Authority Information Access <1.3.6.1.5.5.7.1.1>
        Critical: false
        Value: 0x305a302606082b06010505073001861a... (92B)
      }
      #7 Extension{
        Id: CRL Distribution Points <2.5.29.31>
        Critical: false
        Value: 0x302b3029a027a0258623687474703a2f... (45B)
      }
      #8 Extension{
        Id: Certificate Policies <2.5.29.32>
        Critical: false
        Value: 0x304e3038060a2b06010401d679020503... (80B)
      }
    ]
    ExtKeyUsage: [
      #1 ServerAuth
      #2 ClientAuth
    ]
    BasicConstraintsValid: true
    IsCA: true
    MaxPathLen: 0 (int)
    MaxPathLenZero: true
    SubjectKeyId: 0x8a747faf85cdee95cd3d9cd0e24614f3... (20B)
    AuthorityKeyId: 0xe4af2b26711a2b4827852f52662ceff0... (20B)
    OCSPServer: [
      #1 "http://ocsp.pki.goog/gtsr1"
    ]
    IssuingCertificateURL: [
      #1 "http://pki.goog/repo/certs/gtsr1.der"
    ]
    PermittedDNSDomainsCritical: false
    CRLDistributionPoints: [
      #1 "http://crl.pki.goog/gtsr1/gtsr1.crl"
    ]
    PolicyIdentifiers: [
      #1 1.3.6.1.4.1.11129.2.5.3
      #2 2.23.140.1.2.1
      #3 2.23.140.1.2.2
    ]
  }
  #3 Certificate{
    Raw: 0x308205623082044aa003020102021077... (1.38KB)
    RawTBSCertificate: 0x3082044aa003020102021077bd0d6cdb... (1.1KB)
    RawSubjectPublicKeyInfo: 0x30820222300d06092a864886f70d0101... (550B)
    RawSubject:
      0G1·0⇥··U····US1"0···U·⏎
      ··Google·Trust·Services·LLC1·0···U····GTS·Root·R1
    RawIssuer:
      0W1·0⇥··U····BE1·0···U·⏎
      ··GlobalSign·nv-sa1·0···U····Root·CA1·0···U····GlobalSign·Root·CA
    Signature: 0x34a41eb128a3d0b47617a6317a21e9d1... (256B)
    SignatureAlgorithm: "SHA256-RSA" (x509.SignatureAlgorithm)
    PublicKeyAlgorithm: "RSA" (x509.PublicKeyAlgorithm)
    PublicKey: 0x1400010e290 (interface {})
    Version: 3 (int)
    SerialNumber: "159159747900478145820483398898491642637" (*big.Int)
    Issuer: "CN=GlobalSign·Root·CA,OU=Root·CA,O=GlobalSign·nv-sa,C=BE" (pkix.Name)
    Subject: "CN=GTS·Root·R1,O=Google·Trust·Services·LLC,C=US" (pkix.Name)
    NotBefore: "2020-06-19T00:00:42Z"
    NotAfter: "2028-01-28T00:00:42Z"
    KeyUsage:
      DigitalSignature
      CertSign
      CRLSign
    Extensions: [
      #1 Extension{
        Id: Certificate Key Usage <2.5.29.15>
        Critical: true
        Value: 0x03020186 (4B)
      }
      #2 Extension{
        Id: Certificate Basic Constraints <2.5.29.19>
        Critical: true
        Value: 0x30030101ff (5B)
      }
      #3 Extension{
        Id: Certificate Subject Key ID <2.5.29.14>
        Critical: false
        Value: 0x0414e4af2b26711a2b4827852f52662c... (22B)
      }
      #4 Extension{
        Id: Certificate Authority Key Identifier <2.5.29.35>
        Critical: false
        Value: 0x30168014607b661a450d97ca89502f7d... (24B)
      }
      #5 Extension{
        Id: Certificate Authority Information Access <1.3.6.1.5.5.7.1.1>
        Critical: false
        Value: 0x3052302506082b060105050730018619... (84B)
      }
      #6 Extension{
        Id: CRL Distribution Points <2.5.29.31>
        Critical: false
        Value: 0x30293027a025a0238621687474703a2f... (43B)
      }
      #7 Extension{
        Id: Certificate Policies <2.5.29.32>
        Critical: false
        Value: 0x30323008060667810c01020130080606... (52B)
      }
    ]
    BasicConstraintsValid: true
    IsCA: true
    MaxPathLen: -1 (int)
    MaxPathLenZero: false
    SubjectKeyId: 0xe4af2b26711a2b4827852f52662ceff0... (20B)
    AuthorityKeyId: 0x607b661a450d97ca89502f7d04cd34a8... (20B)
    OCSPServer: [
      #1 "http://ocsp.pki.goog/gsr1"
    ]
    IssuingCertificateURL: [
      #1 "http://pki.goog/gsr1/gsr1.crt"
    ]
    PermittedDNSDomainsCritical: false
    CRLDistributionPoints: [
      #1 "http://crl.pki.goog/gsr1/gsr1.crl"
    ]
    PolicyIdentifiers: [
      #1 2.23.140.1.2.1
      #2 2.23.140.1.2.2
      #3 1.3.6.1.4.1.11129.2.5.3.2
      #4 1.3.6.1.4.1.11129.2.5.3.3
    ]
  }
]
2022/11/30 14:26:55 tls rtt 163.077708ms
```