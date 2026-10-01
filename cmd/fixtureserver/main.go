// Command fixtureserver answers GET /x from <dir>/x.json, so a test can stand
// up a fake truth and a fake mirror that disagree in known ways. Each
// "{self}" in a fixture becomes the server's own base URL.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("dir", ".", "fixture directory")
	listen := flag.String("listen", "127.0.0.1:0", "listen address")
	flag.Parse()
	self := "http://" + *listen
	err := http.ListenAndServe(*listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := filepath.Join(*dir, filepath.FromSlash(strings.TrimPrefix(r.URL.Path, "/"))+".json")
		b, err := os.ReadFile(name)
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Link", "<"+self+r.URL.Path+"?page=2>; rel=\"next\"")
		_, _ = w.Write([]byte(strings.ReplaceAll(string(b), "{self}", self)))
	}))
	fmt.Fprintln(os.Stderr, "fixtureserver:", err)
	os.Exit(1)
}
