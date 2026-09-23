package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestConnectWithProxy(t *testing.T) {
	target := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer target.Close()

	request := make(chan *http.Request, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request <- r.Clone(context.Background())
		upstream, err := net.Dial("tcp", r.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer upstream.Close()

		downstream, rw, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Errorf("hijack proxy connection: %v", err)
			return
		}
		defer downstream.Close()
		if _, err := rw.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n"); err != nil {
			t.Errorf("write CONNECT response: %v", err)
			return
		}
		if err := rw.Flush(); err != nil {
			t.Errorf("flush CONNECT response: %v", err)
			return
		}

		go io.Copy(upstream, downstream)
		_, _ = io.Copy(downstream, upstream)
	}))
	defer proxy.Close()

	proxyURL := strings.Replace(proxy.URL, "http://", "http://user:p%40ss@", 1)
	if err := connectWithProxy(target.URL, proxyURL); err != nil {
		t.Fatalf("connect through proxy: %v", err)
	}

	select {
	case got := <-request:
		if got.Method != http.MethodConnect {
			t.Errorf("method = %q, want CONNECT", got.Method)
		}
		wantTarget := strings.TrimPrefix(target.URL, "https://")
		if got.Host != wantTarget {
			t.Errorf("target = %q, want %q", got.Host, wantTarget)
		}
		const wantAuthorization = "Basic dXNlcjpwQHNz"
		if authorization := got.Header.Get("Proxy-Authorization"); authorization != wantAuthorization {
			t.Errorf("proxy authorization = %q, want %q", authorization, wantAuthorization)
		}
	case <-time.After(time.Second):
		t.Fatal("proxy did not receive CONNECT request")
	}
}

func TestDialCONNECTProxyPreservesBufferedData(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		_, _ = conn.Read(buf)
		_, _ = io.WriteString(conn, "HTTP/1.1 200 Connection Established\r\n\r\nready")
	}()

	proxyURL, err := url.Parse("http://" + listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	conn, err := dialCONNECTProxy(context.Background(), &net.Dialer{Timeout: time.Second}, proxyURL, "target.example:443")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	data := make([]byte, len("ready"))
	if _, err := io.ReadFull(conn, data); err != nil {
		t.Fatal(err)
	}
	if string(data) != "ready" {
		t.Fatalf("buffered data = %q, want ready", data)
	}
}

func TestDialCONNECTProxyRejectsFailure(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "authentication required", http.StatusProxyAuthRequired)
	}))
	defer proxy.Close()

	proxyURL, err := url.Parse(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dialCONNECTProxy(context.Background(), &net.Dialer{Timeout: time.Second}, proxyURL, "target.example:443")
	if err == nil || !strings.Contains(err.Error(), "407 Proxy Authentication Required") {
		t.Fatalf("error = %v, want proxy status", err)
	}
}

func TestTargetAddress(t *testing.T) {
	tests := []struct {
		input      string
		address    string
		serverName string
	}{
		{"example.com", "example.com:443", "example.com"},
		{"example.com:8443", "example.com:8443", "example.com"},
		{"https://example.com/path", "example.com:443", "example.com"},
		{"https://example.com:8443/path", "example.com:8443", "example.com"},
		{"::1", "[::1]:443", "::1"},
		{"[::1]:8443", "[::1]:8443", "::1"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			address, serverName, err := targetAddress(test.input)
			if err != nil {
				t.Fatal(err)
			}
			if address != test.address || serverName != test.serverName {
				t.Fatalf("targetAddress() = (%q, %q), want (%q, %q)", address, serverName, test.address, test.serverName)
			}
		})
	}
}

func TestParseProxyURL(t *testing.T) {
	tests := map[string]string{
		"proxy.example":         "http://proxy.example:80",
		"proxy.example:8080":    "http://proxy.example:8080",
		"http://proxy.example":  "http://proxy.example:80",
		"https://proxy.example": "https://proxy.example:443",
	}
	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := parseProxyURL(input)
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != want {
				t.Fatalf("proxy URL = %q, want %q", got, want)
			}
		})
	}

	for _, input := range []string{"socks5://proxy.example", "://bad"} {
		t.Run(fmt.Sprintf("reject_%s", input), func(t *testing.T) {
			if _, err := parseProxyURL(input); err == nil {
				t.Fatalf("parseProxyURL(%q) unexpectedly succeeded", input)
			}
		})
	}
}
