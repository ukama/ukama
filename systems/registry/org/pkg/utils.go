package pkg

import (
	"bytes"
	"html/template"
	"io/ioutil"
)

func GenerateEmailBody(invitationID string, link string, owner string, org string, role string, name string) (string, error) {
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
			"EmailID": invitationID,
			"LINK":    link,
			"OWNER":   owner,
			"ORG":     org,
			"ROLE":    role,
			"NAME":    name,
		},
	})
	if err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}
