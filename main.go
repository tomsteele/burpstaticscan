package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/franela/goreq"
)

var dir *string
var root *string
var host *string
var port *string
var iport int
var burpBuddyURL *string
var localURL string

var count int

func init() {
	host = flag.String("host", "localhost", "host for the file server to listen on")
	port = flag.String("port", "9999", "port for the file server to listen on")
	dir = flag.String("dir", "", "directory with code to scan")
	root = flag.String("root", "", "root to serve from")
	burpBuddyURL = flag.String("burpbuddy", "http://localhost:8001", "HTTP API URL for burpbuddy")
	flag.Parse()

	if *dir == "" {
		log.Fatal("--dir is required")
	}

	d, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("Error getting absolute path of dir: %s", err.Error())
	}
	dir = &d

	iport, err = strconv.Atoi(*port)
	if err != nil {
		log.Fatalf("Error converting port to int: %s", err.Error())
	}
	count = 0
	localURL = fmt.Sprintf("http://%s:%s%s", *host, *port, *root)
}

func main() {

	go func() {
		log.Fatal(http.ListenAndServe(*host+":"+*port, http.StripPrefix(*root+"/", http.FileServer(http.Dir(*dir)))))
	}()

	log.Printf("Static file server listening on %s, serving %s", localURL, *dir)
	log.Printf("Adding %s to scope", localURL)

	s := ScopeRequest{URL: localURL + "/"}
	resp, err := goreq.Request{
		Method: "POST",
		Uri:    *burpBuddyURL + "/scope",
		Body:   s,
	}.Do()

	if err != nil {
		log.Fatalf("Error adding site to scope: %s", err.Error())
	}
	if resp.StatusCode != 201 {
		log.Fatal("Non 201 status code when adding to scope")
	}
	resp.Body.Close()

	log.Println("Walking directory, each file will be sent to burp's passive scan")

	if err := filepath.Walk(*dir, walker); err != nil {
		log.Printf("Error during directory walk: %s", err.Error())
	}

	log.Printf("%d files sent to burp", count)
}

func walker(path string, info os.FileInfo, err error) error {
	parts := strings.Split(path, *dir)
	if parts[1] == "" {
		return nil
	}
	fname := strings.Replace(parts[1], "\\", "/", -1)
	resp, err := http.Get(localURL + fname)

	if err != nil {
		log.Printf("Error getting %s from go http server: %s", fname, err.Error())
		return nil
	}
	if resp.StatusCode != 200 {
		log.Printf("Non 200 status code from go for %s", fname)
		return nil
	}

	rawResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("Error dumping response: %s", err.Error())
		return nil
	}
	rawRequest, err := httputil.DumpRequest(resp.Request, true)
	if err != nil {
		log.Printf("Error dumping request: %s", err.Error())
		return nil
	}

	r := PassivScanRequest{
		Host:     *host,
		Port:     iport,
		UseHTTPS: false,
		Request:  base64.StdEncoding.EncodeToString(rawRequest),
		Response: base64.StdEncoding.EncodeToString(rawResponse),
	}

	presp, err := goreq.Request{
		Method: "POST",
		Uri:    *burpBuddyURL + "/scan/passive",
		Body:   r,
	}.Do()

	if err != nil {
		log.Printf("Error during passive scan request to burpbuddy: %s", err.Error())
		return nil
	}

	if presp.StatusCode != 201 {
		log.Println("Non 201 status code during passive scan request to burpbuddy")
		return nil
	}

	defer resp.Body.Close()
	defer presp.Body.Close()
	count += 1

	return nil
}
