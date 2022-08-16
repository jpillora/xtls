package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jpillora/opts"
	"github.com/jpillora/xtls/internal/pp"
)

var version = "0.0.0-src"

func main() {
	type config struct {
		Args []string `opts:"mode=arg, help=url or hostname or file path, max=1"`
	}
	c := config{}
	opts.New(&c).
		Version(version).
		Parse()

	arg := "-"
	if len(c.Args) > 0 {
		arg = c.Args[0]
	}
	if err := run(arg); err != nil {
		log.Fatalf("errored: %s", err)
	}

}

func run(arg string) error {
	if arg == "" || arg == "-" {
		return stdin()
	}
	if len(arg) < 1024 {
		if s, err := os.Stat(arg); err == nil && !s.IsDir() {
			return file(arg)
		}
		return connect(arg)
	}
	return errors.New("unknown input")
}

func stdin() error {
	return reader(os.Stdin)
}

func file(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return reader(f)
}

func reader(r io.Reader) error {
	r = io.LimitReader(r, 1024*1024) //cert bigger than 1MB? dont think so
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	// optionally decode base64s
	if b64, err := base64.StdEncoding.DecodeString(string(b)); err == nil && len(b64) > 0 {
		b = b64
	}
	var block *pem.Block
	data := b
	for len(data) > 0 {
		block, data = pem.Decode(data)
		switch block.Type {
		case "CERTIFICATE":
			c, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return err
			}
			if err := cert(c); err != nil {
				return err
			}
		default:
			fmt.Printf("unsupported PEM type: %s\n", block.Type)
		}
	}
	return nil
}

func connect(host string) error {
	log.Printf("connect to: %s", host)
	// parse input
	if u, err := url.Parse(host); err == nil && strings.HasPrefix(u.Scheme, "http") {
		host = u.Host
	}
	if _, _, err := net.SplitHostPort(host); err != nil {
		host += ":443"
	}
	// setup tcp and tls
	d := net.Dialer{
		Timeout: 10 * time.Second,
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		VerifyConnection: func(cs tls.ConnectionState) error {
			log.Printf("tls server name: %s", cs.ServerName)
			log.Printf("tls version: %d", cs.Version)
			log.Printf("tls ciphersuite: %d", cs.CipherSuite)
			for _, c := range cs.PeerCertificates {
				cert(c)
			}
			return nil
		},
	}
	t0 := time.Now()
	log.Printf("dialing %s", host)
	conn, err := tls.DialWithDialer(&d, "tcp", host, config)
	if err != nil {
		return err
	}
	log.Printf("tls rtt %s", time.Since(t0))
	conn.Close()
	return nil
}

func cert(c *x509.Certificate) error {
	pp.Print(c)
	return nil
}
