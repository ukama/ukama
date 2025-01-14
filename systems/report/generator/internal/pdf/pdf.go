/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pdf

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/util/payments"
)

//go:embed templates
var templatesFs embed.FS

const pdfDPI = 300

type PdfEngine interface {
	Configure(io.ReadWriteCloser, uint) error
	Generate(string) error
}

type PdfObject struct {
	body   string
	engine PdfEngine
}

func NewPDFObject(body string, engine PdfEngine) *PdfObject {
	return &PdfObject{
		body:   body,
		engine: engine,
	}
}

func (p *PdfObject) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.New(path.Base(templateFileName)).Funcs(template.FuncMap{
		"formatAmount": payments.ToAmount,
	}).ParseFS(templatesFs, templateFileName)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	p.body = buf.String()

	return nil
}

func (p *PdfObject) GenerateFile(pdfPath string) error {
	tmpFile, err := os.CreateTemp("", "invoice*.html")
	if err != nil {
		return fmt.Errorf("failed to create temporary html file: %w", err)
	}

	defer func() {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			log.Errorf("failed to remove temporary html file: %v", err)
		}
	}()

	defer func() {
		err = tmpFile.Close()
		if err != nil {
			log.Errorf("failed to close html file: %v", err)
		}
	}()

	err = os.WriteFile(tmpFile.Name(), []byte(p.body), 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to temporary html file: %w", err)
	}

	err = p.engine.Configure(tmpFile, pdfDPI)
	if err != nil {
		return fmt.Errorf("failed to configure PDF generator engine: %w", err)
	}

	err = p.engine.Generate(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to generate PDF file to disk: %v", err)
	}

	return nil
}
