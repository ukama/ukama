package pdf

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	log "github.com/sirupsen/logrus"
)

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
	t, err := template.ParseFiles(templateFileName)
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

func (i *InvoicePDF) GeneratePDF(pdfPath string) error {
	t := time.Now().Unix()

	err := os.WriteFile("storage/"+strconv.FormatInt(int64(t), 10)+".html", []byte(i.body), 0644)
	if err != nil {
		return fmt.Errorf("failed to write html invoice file: %w", err)
	}

	defer func() {
		err = os.Remove("storage/" + strconv.FormatInt(int64(t), 10) + ".html")
		if err != nil {
			log.Errorf("failed to remove html invoice file: %v", err)
		}
	}()

	f, err := os.Open("storage/" + strconv.FormatInt(int64(t), 10) + ".html")
	if err != nil {
		return fmt.Errorf("failed to open html invoice file: %w", err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Errorf("failed to close html invoice file: %v", err)
		}
	}()

	pdfgen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to get a new PDF generator: %w", err)
	}

	pdfgen.AddPage(wkhtmltopdf.NewPageReader(f))

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
