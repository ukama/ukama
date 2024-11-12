/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package engine

import (
	"fmt"
	"io"

	wk "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type WkGenerator struct {
	pdfGen *wk.PDFGenerator
}

func NewWkGenerator() (*WkGenerator, error) {
	pdfgen, err := wk.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to get a new WK PDF generator: %w", err)
	}

	return &WkGenerator{
		pdfGen: pdfgen,
	}, nil
}

func (w *WkGenerator) Configure(tmpFile io.ReadWriteCloser, dpi uint) error {
	w.pdfGen.AddPage(wk.NewPageReader(tmpFile))
	w.pdfGen.PageSize.Set(wk.PageSizeA4)
	w.pdfGen.Dpi.Set(dpi)

	return nil
}

func (w *WkGenerator) Generate(path string) error {
	err := w.pdfGen.Create()
	if err != nil {
		return fmt.Errorf("failed to convert html file to PDF from generator: %w", err)
	}

	err = w.pdfGen.WriteFile(path)
	if err != nil {
		return fmt.Errorf("failed to write PDF file to disk: %v", err)
	}

	return nil
}
