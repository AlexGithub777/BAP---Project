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
func NewTemplateRenderer() *TemplateRenderer {
	rootDir, _ := os.Getwd()
	templateDir := filepath.Join(rootDir, "templates", "*.html")
	dashboardComponentDir := filepath.Join(rootDir, "templates", "components", "dashboard", "*.html")
	adminComponentDir := filepath.Join(rootDir, "templates", "components", "admin", "*.html")
	t := template.Must(template.New("").ParseGlob(templateDir))
	t = template.Must(t.ParseGlob(dashboardComponentDir))
	t = template.Must(t.ParseGlob(adminComponentDir))
	return &TemplateRenderer{templates: t}
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
