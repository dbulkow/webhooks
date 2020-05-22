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
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Subversion struct {
	token string // access token
}

func NewSvnHandler(token string) http.Handler {
	return &Subversion{}
}

func (s *Subversion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	p := strings.Split(r.URL.Path, "/")
	fmt.Println(p, len(p))
	if len(p) < 2 {
		http.Error(w, "malformed request", http.StatusBadRequest)
		return
	}
	uuid := p[0]
	cmd := p[1]
	fmt.Printf("uuid=%s\ncommand=%s\n", uuid, cmd)

	rev := r.URL.Query().Get("rev")
	if rev == "" {
		http.Error(w, "revision not specified", http.StatusNotFound)
		return
	}
	fmt.Printf("rev=%s\n", rev)

	ct := r.Header.Get("Content-Type")
	switch ct {
	case "text/plain":
	// case "application/json":
	default:
		http.Error(w, fmt.Sprintf("content type \"%s\" not accepted", ct), http.StatusBadRequest)
		return
	}
	fmt.Println("Content-Type:", ct)

	for k, v := range r.Header {
		fmt.Printf("%s\t", k)
		for _, x := range v {
			fmt.Printf("%s ", x)
		}
		fmt.Println()
	}

	// if s.token != "" {
	// 	tok := r.Header.Get("")
	// 	if tok == "" {
	// 		http.Error(w, "unauthorized access", http.StatusUnauthorized)
	// 		return
	// 	}
	// 	if tok != s.token {
	// 		http.Error(w, "unauthorized access", http.StatusUnauthorized)
	// 		return
	// 	}
	// }

	if r.Body != nil {
		io.Copy(os.Stdout, io.LimitReader(r.Body, 4096))
	}

	w.WriteHeader(http.StatusOK)
}
