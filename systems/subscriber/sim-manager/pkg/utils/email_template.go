package utils

import (
	"bytes"
	"text/template"
)

func GenerateEmailBody(network string, name string, qrCode string) (string, error) {
	bodyTemplate := `
	Subject:[Ukama] {{ .Values.NETWORK}} invited you to use their network
MIME-version: 1.0;
Content-Type: multipart/mixed;
        boundary="XXXXboundary text"

This is a multipart message in MIME format.

--XXXXboundary text
Content-Type: text/html; charset="UTF-8";
<!DOCTYPE html>
<html>
<head>
<link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
<style>
.button {
  background-color: #2190F6;
font-family: arial narrow;
font-style: normal;
font-weight: 500;
letter-spacing: 0.4px;
  color: #ffffff;
font-size: 14px;
line-height: 24px;
padding-top:5px;
  box-shadow: 0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12);
border-radius: 4px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 14px;
  font-weight: 500;
  width: 234px;
height: 30px;
  cursor: pointer;
}
.button:hover {
  background-color: #21a1f6;
}

</style>
</head>
<body>
<div style="width:600px">
<h3>Download your eSIM </h3>
<p style="margin-bottom:5px,margin-top:3px,font-size:10px">Hi {{ .Values.NAME}},</p>
<div style="maring-bottom:3px,font-size:10px">{{ .Values.NETWORK}} has invited you to join their network. To get started, scan the QR code below to download your eSIM, and start using the Ukama network. You have X GBs of free data, and will be charged at $5/GB while roaming.
</div><br>
<div style="maring-bottom:3px,font-size:10px">
This invitation <strong>expires in 48 hours </strong>for security purposes, but another invitation will be sent if you click on the link past the time period.
</div><br>
<div style="maring-bottom:3px,font-size:10px"> 
Note: you don’t have to do anything with your current SIM in your phone./other instructions regarding how to navigate this
</div><br>
<div>

<img src="cid:qrcid"  alt={{ .Values.QRCODE}}  width="213" height="211" alt="qrcode"/>

</div><br>
<div style="margin-bottom:3px,font-size:16px">Having trouble with the QR code? Use the following link instead: <a href="{{ .Values.QRCODE}}" >https://{{ .Values.QRCODE}}</a></div> 

<div style="white-space: pre">
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

--XXXXboundary text
Content-Type: image/png;
Content-Disposition: inline; filename="qr.png"
Content-Transfer-Encoding: base64
Content-ID: <qrcid>
Content-Location: qr.png

{{ .Values.QRCODE}}

--XXXXboundary text--
`

	tmpl, err := template.New("email").Parse(bodyTemplate)
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
