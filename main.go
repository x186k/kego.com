package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"
)

var _ = checkFatal

func checkFatal(err error) {
	if err != nil {
		_, fileName, fileLine, _ := runtime.Caller(1)
		log.Fatalf("FATAL %s:%d %v", filepath.Base(fileName), fileLine, err)
	}
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func main() {

	mux := http.NewServeMux()
	f := os.DirFS("html")
	mux.Handle("/", http.FileServer(neuteredFileSystem{http.FS(f)}))



	mux.HandleFunc("/search", func(rw http.ResponseWriter, r *http.Request) {
		//a := time.Now().String()
		time.Sleep(time.Second * 2)
		a := strings.ToLower(r.PostFormValue("search"))
		b := []rune{}
		for _, v := range a {
			if unicode.IsNumber(v) || unicode.IsLetter(v) {
				b = append(b, v)
			}
		}
		final := string(b)

		if final == "" {
			_, _ = rw.Write([]byte("not valid"))
			return
		}

		c := fmt.Sprintf("Good news! Your domain is available: %s.kego.com", final)
		_, _ = rw.Write([]byte(c))
	})

	println("serving http")
	panic(http.ListenAndServe(":80", mux))

}
