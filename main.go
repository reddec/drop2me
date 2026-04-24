package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

//go:embed static/form.html static/success.html
var staticFS embed.FS

var (
	tplSuccess *template.Template
	bindAddr   string
	uploadDir  string
	maxSize    int64
)

func init() {
	flag.StringVar(&bindAddr, "bind", envOr("DROP2ME_BIND", ":8080"), "binding address")
	flag.StringVar(&uploadDir, "dir", envOr("DROP2ME_DIR", "."), "upload directory")
	flag.Int64Var(&maxSize, "max-size", envInt64("DROP2ME_MAX_SIZE", 0), "max upload size in bytes (0 = unlimited)")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt64(key string, fallback int64) int64 {
	if v, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
		return v
	}
	return fallback
}

func main() {
	flag.Parse()

	var err error
	tplSuccess, err = template.ParseFS(staticFS, "static/success.html")
	if err != nil {
		log.Fatalf("parse success template: %v", err)
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatalf("failed to create upload directory: %v", err)
	}

	_, port, err := net.SplitHostPort(bindAddr)
	if err != nil {
		log.Fatalf("invalid bind address %q: %v", bindAddr, err)
	}

	urls := listURLs(port)
	if len(urls) > 0 {
		fmt.Println("Open one of these URLs to upload files:")
		for _, u := range urls {
			fmt.Printf("  %s\n", u)
		}
		fmt.Println()
		printQR(urls[0])
		fmt.Println()
	}

	http.HandleFunc("/", handler)

	fmt.Printf("Listening on %s\n\n", bindAddr)
	server := &http.Server{
		Addr:              bindAddr,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func listURLs(port string) []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var urls []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := addrIP(addr)
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue
			}
			if v4 := ip.To4(); v4 != nil {
				urls = append(urls, fmt.Sprintf("http://%s:%s/", v4, port))
			} else {
				urls = append(urls, fmt.Sprintf("http://[%s]:%s/", ip, port))
			}
		}
	}
	return urls
}

func addrIP(a net.Addr) net.IP {
	switch v := a.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	}
	return nil
}

func printQR(s string) {
	qr, err := qrcode.New(s, qrcode.Medium)
	if err != nil {
		log.Printf("QR generation failed: %v", err)
		return
	}
	fmt.Println(qr.ToString(false))
}

func createUniqueFile(dir, name string) (*os.File, string, error) {
	name = filepath.Base(name)
	dst := filepath.Join(dir, name)
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err == nil {
		return f, dst, nil
	}
	if !os.IsExist(err) {
		return nil, "", err
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	for i := 1; ; i++ {
		dst = filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
		f, err = os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
		if err == nil {
			return f, dst, nil
		}
		if !os.IsExist(err) {
			return nil, "", err
		}
	}
}

// ---------------------------------------------------------------------------
// HTTP
// ---------------------------------------------------------------------------

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data, _ := staticFS.ReadFile("static/form.html")
		w.Write(data)
	case http.MethodPost:
		handleUpload(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if maxSize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var names []string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name := part.FileName()
		if name == "" {
			continue
		}
		f, dst, err := createUniqueFile(uploadDir, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name = filepath.Base(dst)

		_, err = io.Copy(f, part)
		f.Close()
		if err != nil {
			os.Remove(dst)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		names = append(names, name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tplSuccess.Execute(w, names)
}
