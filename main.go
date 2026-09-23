package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
		Proxy string   `opts:"help=HTTP CONNECT proxy URL (http:// or https://), short=p"`
		Args  []string `opts:"mode=arg, help=url or hostname or file path, max=1"`
	}
	c := config{}
	opts.New(&c).
		Version(version).
		Parse()

	arg := "-"
	if len(c.Args) > 0 {
		arg = c.Args[0]
	}
	if err := runWithProxy(arg, c.Proxy); err != nil {
		log.Fatalf("errored: %s", err)
	}

}

func run(arg string) error {
	return runWithProxy(arg, "")
}

func runWithProxy(arg, proxy string) error {
	if arg == "" || arg == "-" {
		return stdin()
	}
	if len(arg) < 1024 {
		if s, err := os.Stat(arg); err == nil && !s.IsDir() {
			return file(arg)
		}
		return connectWithProxy(arg, proxy)
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
			pp.Print(c)
		default:
			fmt.Printf("unsupported PEM type: %s\n", block.Type)
		}
	}
	return nil
}

func connect(host string) error {
	return connectWithProxy(host, "")
}

func connectWithProxy(host, proxy string) error {
	log.Printf("connect to: %s", host)
	host, serverName, err := targetAddress(host)
	if err != nil {
		return err
	}

	d := net.Dialer{
		Timeout: 10 * time.Second,
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         serverName,
		VerifyConnection: func(cs tls.ConnectionState) error {
			log.Printf("tls server name: %s", cs.ServerName)
			log.Printf("tls version: %d", cs.Version)
			log.Printf("tls ciphersuite: %d", cs.CipherSuite)
			pp.Print(cs.PeerCertificates)
			return nil
		},
	}
	t0 := time.Now()
	var raw net.Conn
	if proxy == "" {
		log.Printf("dialing %s", host)
		raw, err = d.Dial("tcp", host)
	} else {
		proxyURL, parseErr := parseProxyURL(proxy)
		if parseErr != nil {
			return parseErr
		}
		log.Printf("dialing %s via %s", host, proxyURL.Redacted())
		raw, err = dialCONNECTProxy(context.Background(), &d, proxyURL, host)
	}
	if err != nil {
		return err
	}
	defer raw.Close()

	deadline := time.Now().Add(d.Timeout)
	if err := raw.SetDeadline(deadline); err != nil {
		return err
	}
	conn := tls.Client(raw, config)
	if err := conn.Handshake(); err != nil {
		return err
	}
	if err := raw.SetDeadline(time.Time{}); err != nil {
		return err
	}
	log.Printf("tls rtt %s", time.Since(t0))
	return conn.Close()
}

func targetAddress(target string) (address, serverName string, err error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", "", errors.New("empty target")
	}
	if strings.Contains(target, "://") {
		u, parseErr := url.Parse(target)
		if parseErr != nil {
			return "", "", fmt.Errorf("invalid target URL: %w", parseErr)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return "", "", fmt.Errorf("unsupported target URL scheme %q", u.Scheme)
		}
		if u.Host == "" {
			return "", "", errors.New("target URL has no host")
		}
		target = u.Host
	}

	if host, port, splitErr := net.SplitHostPort(target); splitErr == nil {
		if host == "" || port == "" {
			return "", "", fmt.Errorf("invalid target address %q", target)
		}
		return target, strings.Trim(host, "[]"), nil
	}

	host := strings.Trim(target, "[]")
	if strings.Contains(host, ":") && net.ParseIP(host) == nil {
		return "", "", fmt.Errorf("invalid target address %q", target)
	}
	return net.JoinHostPort(host, "443"), host, nil
}

func parseProxyURL(rawURL string) (*url.URL, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported proxy URL scheme %q", u.Scheme)
	}
	if u.Hostname() == "" {
		return nil, errors.New("proxy URL has no host")
	}
	if u.Port() == "" {
		port := "80"
		if u.Scheme == "https" {
			port = "443"
		}
		u.Host = net.JoinHostPort(u.Hostname(), port)
	}
	return u, nil
}

func dialCONNECTProxy(ctx context.Context, d *net.Dialer, proxyURL *url.URL, target string) (net.Conn, error) {
	conn, err := d.DialContext(ctx, "tcp", proxyURL.Host)
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}
	ok := false
	defer func() {
		if !ok {
			conn.Close()
		}
	}()

	deadline := time.Now().Add(d.Timeout)
	if err := conn.SetDeadline(deadline); err != nil {
		return nil, err
	}
	if proxyURL.Scheme == "https" {
		proxyTLS := tls.Client(conn, &tls.Config{ServerName: proxyURL.Hostname()})
		if err := proxyTLS.HandshakeContext(ctx); err != nil {
			return nil, fmt.Errorf("TLS handshake with proxy: %w", err)
		}
		conn = proxyTLS
	}

	header := make(http.Header)
	if proxyURL.User != nil {
		password, _ := proxyURL.User.Password()
		credentials := proxyURL.User.Username() + ":" + password
		header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(credentials)))
	}
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: target},
		Host:   target,
		Header: header,
	}
	if err := req.Write(conn); err != nil {
		return nil, fmt.Errorf("write proxy CONNECT request: %w", err)
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, req)
	if err != nil {
		return nil, fmt.Errorf("read proxy CONNECT response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("proxy CONNECT failed: %s", resp.Status)
	}
	if err := conn.SetDeadline(time.Time{}); err != nil {
		return nil, err
	}

	ok = true
	if reader.Buffered() > 0 {
		return &bufferedConn{Conn: conn, reader: reader}, nil
	}
	return conn, nil
}

type bufferedConn struct {
	net.Conn
	reader *bufio.Reader
}

func (c *bufferedConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}
