# Notification System

The Notification System with Multi-Channel Support is a flexible system that allows ukama system of systems to send notifications through various channels, including email, SMS, Slack, and WhatsApp. It provides a unified interface for sending notifications and supports customizable templates for each channel.
currently only supporting email.

## Features

Send notifications through multiple channels: email, SMS, Slack, and WhatsApp.
Customize notification templates for each channel.
Support for HTML content in email notifications.
Specify recipients, subject, and body content for each notification.

## Sending Emails

To send an email using the Notification System, follow these steps:

Import the mailer package into your code: import "github.com/ukama/ukama/systems/notification/mailer"
Create a new instance of the mailer service: mailer := mailer.NewMailer(smtpServer)
Prepare the necessary email data:
Recipients: Provide a list of email addresses as recipients.
Subject: Specify the subject line for the email.
Body: Create the email body content, which can include HTML templates.
Use the mailer service to send the email:
`
err := mailer.SendEmail(recipients, subject, body)
if err != nil {
// Handle error
}

`
