/*
MIT License

Copyright (c) 2020 David Bulkow

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package subversion

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSubversion(t *testing.T) {
	ts := httptest.NewServer(http.StripPrefix("/subversion/", NewSvnHandler("")))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/subversion/uuid/notifyCommit"
	q := u.Query()
	q.Set("rev", "1234")
	u.RawQuery = q.Encode()

	body := []byte("A file1\nD file2\nM file3\n")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("empty response")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	exp := http.StatusOK

	if resp.StatusCode != exp {
		t.Fatalf("expected status \"%s\" got \"%s\"", http.StatusText(exp), resp.Status)
	}
}

func TestSubversionShortPath(t *testing.T) {
	ts := httptest.NewServer(http.StripPrefix("/subversion/", NewSvnHandler("")))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/subversion/"

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("empty response")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	exp := http.StatusBadRequest

	if resp.StatusCode != exp {
		t.Fatalf("expected status \"%s\" got \"%s\"", http.StatusText(exp), resp.Status)
	}
}

func TestSubversionBadCommand(t *testing.T) {
	ts := httptest.NewServer(http.StripPrefix("/subversion/", NewSvnHandler("")))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/subversion/uuid/somecommand"

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("empty response")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	exp := http.StatusNotFound

	if resp.StatusCode != exp {
		t.Fatalf("expected status \"%s\" got \"%s\"", http.StatusText(exp), resp.Status)
	}
}

func TestSubversionNoRev(t *testing.T) {
	ts := httptest.NewServer(http.StripPrefix("/subversion/", NewSvnHandler("")))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/subversion/uuid/notifyCommit"

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("empty response")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	exp := http.StatusNotFound

	if resp.StatusCode != exp {
		t.Fatalf("expected status \"%s\" got \"%s\"", http.StatusText(exp), resp.Status)
	}
}

func TestSubversionBadContentType(t *testing.T) {
	ts := httptest.NewServer(http.StripPrefix("/subversion/", NewSvnHandler("")))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/subversion/uuid/notifyCommit"
	q := u.Query()
	q.Set("rev", "1234")
	u.RawQuery = q.Encode()

	body := []byte("A file1\nD file2\nM file3\n")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("empty response")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	exp := http.StatusBadRequest

	if resp.StatusCode != exp {
		t.Fatalf("expected status \"%s\" got \"%s\"", http.StatusText(exp), resp.Status)
	}
}

func TestSubversionBadBody(t *testing.T) {
	ts := httptest.NewServer(http.StripPrefix("/subversion/", NewSvnHandler("")))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/subversion/uuid/notifyCommit"
	q := u.Query()
	q.Set("rev", "1234")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("empty response")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	exp := http.StatusOK

	if resp.StatusCode != exp {
		t.Fatalf("expected status \"%s\" got \"%s\"", http.StatusText(exp), resp.Status)
	}
}
