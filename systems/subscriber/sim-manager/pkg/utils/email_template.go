package utils

import (
	"bytes"
	"html/template"
	"io/ioutil"
)
func GenerateEmailBody(network string, name string, qrCode string) (string, error) {


	// Read the template content from the file
	templateContent, err := ioutil.ReadFile("email_template.gotmpl")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("email").Parse(string(templateContent))
	if err != nil {
		return "", err
	}

	var bodyBuffer bytes.Buffer
	err = tmpl.Execute(&bodyBuffer, map[string]interface{}{
		"Values": map[string]interface{}{
			"NETWORK": network,
			"NAME":    name,
			"QRCODE":  qrCode,
			
		},
	})
	if err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}
