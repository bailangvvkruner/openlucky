package web

import (
	"embed"
	"mime"
	"path"
	"strings"
)

//go:embed index.html manifest.webmanifest static/css/app.css static/js/*.js static/js/views/*.js
var files embed.FS

func Read(requestPath string) ([]byte, string, error) {
	clean := cleanPath(requestPath)
	if clean == "" {
		clean = "index.html"
	}

	data, err := files.ReadFile(clean)
	if err != nil && path.Ext(clean) == "" {
		clean = "index.html"
		data, err = files.ReadFile(clean)
	}
	if err != nil {
		return nil, "", err
	}

	return data, contentType(clean), nil
}

func cleanPath(requestPath string) string {
	clean := strings.TrimSpace(requestPath)
	clean = strings.TrimPrefix(clean, "/lucky/")
	clean = strings.TrimPrefix(clean, "/")
	clean = path.Clean("/" + clean)
	clean = strings.TrimPrefix(clean, "/")
	if clean == "." {
		return ""
	}
	return clean
}

func contentType(assetPath string) string {
	if strings.HasSuffix(assetPath, ".webmanifest") {
		return "application/manifest+json; charset=utf-8"
	}
	if assetPath == "index.html" {
		return "text/html; charset=utf-8"
	}
	if contentType := mime.TypeByExtension(path.Ext(assetPath)); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}
