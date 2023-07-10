package pkg

import (
	"bytes"
	"html/template"
)

func GenerateEmailBody(invitationID string, link string, owner string, org string, role string, name string) (string, error) {
	bodyTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
	  <meta charset="UTF-8">
	  <link rel="preconnect" href="https://fonts.googleapis.com">
	  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	  <link href="https://fonts.googleapis.com/css2?family=Arial&display=swap" rel="stylesheet">
	  <title>
		{{ .Values.OWNER}} has invited you to join {{ .Values.ORG}} as a {{ .Values.ROLE}}
	  </title>
	  <style>
	  h1 {
		font-family: 'Google Font', sans-serif;
		font-size: 1rem;
	  }
  
	  body {
		font-family: 'Arial', sans-serif;
		text-align: center;
	  }
  
	  .button {
		display: inline-block;
		padding: 10px 20px;
		text-decoration: none;
		background-color: #4285F4;
		border-radius: 4px;
		transition: background-color 0.3s ease;
		with: 200px;
	  }
  
	  .button:hover {
		background-color: #3367D6;
	  }
  
	  </style>
	</head>
	<body>
	  <h1>{{ .Values.OWNER}} has invited you to join {{ .Values.ORG}} as a {{ .Values.ROLE}}</h1>
	  
	  <p>Hi {{ .Values.NAME }},</p>
	
	  <p>You have been invited to join {{ .Values.ORG}} as a {{ .Values.ROLE}}. This invitation will expire after 5 minutes. To accept the invitation and get started with Ukama Console, click the button below.</p>
	  <a href="{{ .Values.LINK }}" class="button" style="color: #ffffff">ACCEPT TEAM INVITATION </a>


	  <p>Having trouble with the button? Use the following link instead: <a href="{{ .Values.LINK}}">{{ .Values.LINK}}</a></p>
	  Email ID: {{ .Values.EmailID }} (You can use this email to follow up in case of any issues)
	  <div>
	  Thanks,
	  The Ukama Team
	  </div><br>
	  <span style="white-space: pre-line">
	  <div style="font-size:8px,font-style:normal">© Ukama Inc. 1233 Quarry Lane #115, Pleasanton, CA 94566</div>
	  <div><img width=325 height=25 style="margin-top:3px" src="https://i.ibb.co/7yHp3jV/Screen-Shot-2022-02-28-at-18-45-11.png" alt="Ukama-verification-Email-footer" border="0"></div>
	  </span>
	  </div>
	</body>
	</html>
	`

	tmpl, err := template.New("email").Parse(bodyTemplate)
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
