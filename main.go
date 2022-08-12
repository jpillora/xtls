package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/jpillora/opts"
	"github.com/jpillora/xtls/internal/pp"
)

var version = "0.0.0-src"

type pemCmd struct {
	File string `help:"an optional file, stdin will be used by default"`
}

func (p pemCmd) Run() error {
	r := os.Stdin
	if p.File != "" || p.File == "-" {
		f, err := os.Open(p.File)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	var block *pem.Block
	data := b
	for len(data) > 0 {
		block, data = pem.Decode(data)
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return err
			}
			pp.Print(cert)
		default:
			fmt.Printf("unsupported PEM type: %s\n", block.Type)
		}
	}
	return nil
}

type hostCmd struct {
	Host string `opts:"mode=arg" help:"a required <host> containing both hostname and port"`
}

func (c hostCmd) Run() error {
	if h, p, err := net.SplitHostPort(c.Host); err != nil {
		return err
	} else if h == "" {
		return errors.New("empty hostname")
	} else if p == "" {
		return errors.New("empty port")
	}
	// ctx := context.Background()
	d := net.Dialer{
		Timeout: 10 * time.Second,
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		VerifyConnection: func(cs tls.ConnectionState) error {
			pp.Print(cs.PeerCertificates)
			return nil
		},
	}
	t0 := time.Now()
	log.Printf("dialing %s", c.Host)
	conn, err := tls.DialWithDialer(&d, "tcp", c.Host, config)
	if err != nil {
		return err
	}
	log.Printf("tls rtt %s", time.Since(t0))
	conn.Close()
	return nil
}

func main() {
	type config struct{}
	c := config{}
	opts.New(&c).
		Version(version).
		AddCommand(opts.New(&pemCmd{}).Name("pem")).
		AddCommand(opts.New(&hostCmd{}).Name("host")).
		Parse().
		RunFatal()
}
