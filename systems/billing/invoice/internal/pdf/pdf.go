package pdf

import (
	"embed"
	_ "embed"
	"os"

	"bytes"
	"fmt"
	"html/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/mattn/go-slim"
	log "github.com/sirupsen/logrus"
)

//go:embed templates
var templates embed.FS

const pdfDPI = 300

type InvoicePDF struct {
	body string
}

func NewInvoicePDF(body string) *InvoicePDF {
	return &InvoicePDF{
		body: body,
	}
}

func (i *InvoicePDF) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFS(templates, templateFileName)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	i.body = buf.String()

	return nil
}

func (i *InvoicePDF) ParseSlimTemplate(templateFileName string, data interface{}) error {
	templateFile, err := templates.Open(templateFileName)
	if err != nil {
		return err
	}

	defer templateFile.Close()

	t, err := slim.Parse(templateFile)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, slim.Values{"Invoice": data}); err != nil {
		return err
	}

	i.body = buf.String()

	return nil
}

func (i *InvoicePDF) GeneratePDF(pdfPath string) error {
	tmpFile, err := os.CreateTemp("", "invoice*.html")
	if err != nil {
		return fmt.Errorf("failed to create temporary html invoice file: %w", err)
	}

	defer func() {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			log.Errorf("failed to remove temporary html invoice file: %v", err)
		}
	}()

	defer func() {
		err = tmpFile.Close()
		if err != nil {
			log.Errorf("failed to close html invoice file: %v", err)
		}
	}()

	err = os.WriteFile(tmpFile.Name(), []byte(i.body), 0644)
	if err != nil {
		return fmt.Errorf("failed to write invoice data to temporary html file: %w", err)
	}

	pdfgen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to get a new PDF generator: %w", err)
	}

	pdfgen.AddPage(wkhtmltopdf.NewPageReader(tmpFile))

	pdfgen.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfgen.Dpi.Set(pdfDPI)

	err = pdfgen.Create()
	if err != nil {
		return fmt.Errorf("failed to convert html file to PDF from generator: %w", err)
	}

	err = pdfgen.WriteFile(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to write PDF invoice file to disk: %v", err)
	}

	return nil
}
