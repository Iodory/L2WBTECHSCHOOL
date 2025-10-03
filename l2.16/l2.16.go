package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <URL>")
		return
	}
	startURL := os.Args[1]

	// Парсим URL
	parsed, err := url.Parse(startURL)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return
	}

	// Скачиваем HTML
	htmlData, err := downloadFile(startURL, "index.html")
	if err != nil {
		fmt.Println("Download failed:", err)
		return
	}
	fmt.Println("✓ Скачан HTML:", "index.html")

	// Парсим ресурсы
	resources := extractResources(strings.NewReader(string(htmlData)), parsed)

	// Скачиваем все найденные ресурсы
	for _, resURL := range resources {
		saveResource(resURL, parsed)
	}
	fmt.Println("✅ Загрузка завершена!")
}

// downloadFile — скачивает файл по URL и сохраняет в файл localName
func downloadFile(resourceURL, localName string) ([]byte, error) {
	resp, err := http.Get(resourceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(localName, data, 0644)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// extractResources — извлекает ресурсы из HTML (img, script, link)
func extractResources(r io.Reader, base *url.URL) []string {
	resources := []string{}
	doc, err := html.Parse(r)
	if err != nil {
		return resources
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "img":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resources = append(resources, resolveURL(base, attr.Val))
					}
				}
			case "script":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resources = append(resources, resolveURL(base, attr.Val))
					}
				}
			case "link":
				var href string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if href != "" {
					resources = append(resources, resolveURL(base, href))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return resources
}

// resolveURL — делает абсолютный URL
func resolveURL(base *url.URL, link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	return base.ResolveReference(u).String()
}

// saveResource — скачивает ресурс и сохраняет с правильным именем
func saveResource(resourceURL string, base *url.URL) {
	parsed, err := url.Parse(resourceURL)
	if err != nil {
		fmt.Println("⚠️ Пропуск ресурса:", resourceURL)
		return
	}

	// Скачиваем только если тот же домен
	if parsed.Host != base.Host {
		return
	}

	filename := filepath.Base(parsed.Path)
	if filename == "" || filename == "/" {
		filename = "index"
	}

	_, err = downloadFile(resourceURL, filename)
	if err != nil {
		fmt.Println("⚠️ Ошибка скачивания:", resourceURL, "-", err)
	} else {
		fmt.Println("✓ Скачан ресурс:", filename)
	}
}
