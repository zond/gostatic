package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type mockResponseWriter struct {
	w      http.ResponseWriter
	status int
}

func (self *mockResponseWriter) Header() http.Header {
	return self.w.Header()
}

func (self *mockResponseWriter) Write(b []byte) (int, error) {
	return self.w.Write(b)
}

func (self *mockResponseWriter) WriteHeader(i int) {
	self.status = i
	self.w.WriteHeader(i)
}

func logger(h http.Handler) (result http.Handler) {
	result = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock := &mockResponseWriter{w: w}
		h.ServeHTTP(mock, r)
		log.Printf("%v\t%v", r.URL, mock.status)
	})
	return
}

func main() {
	host := flag.String("host", "0.0.0.0:8888", "Where to listen")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := flag.String("dir", wd, "Where to look for files")

	http.ListenAndServe(*host, logger(http.FileServer(http.Dir(*dir))))
}
