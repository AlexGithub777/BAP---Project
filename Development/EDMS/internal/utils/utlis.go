package utils

import (
	"html/template"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer implements echo.Renderer
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer() (*TemplateRenderer, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templateDir := filepath.Join(rootDir, "templates")
	t := template.New("")

	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = t.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{templates: t}, nil
}

// Render implements echo.Renderer
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
