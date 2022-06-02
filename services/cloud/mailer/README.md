# Mailer Service

The purpose of `mailer` service is to send out emails that was scheduled on the AMQP queue. 

In order to send email the message in following format shoule be posted to the 'mailer' queue:

```json
{
    "to": "receiver@example.com",  // receiver address
    "templateName": "test-template", // name of the email template  
    "values":{     // set of values used by template
        "Name": "John Doe",
        "Message": "Hello mailer"
     }
}
```

## Templates 
Mailer uses go templates to generate  email content. Templates should be placed in `templates` folder that could cofigured 
via `TemplatesPath` option. Refer to [config.go](/pkg/config.go)

File name without extension is used as a template name. After processign the template it's entire content will be sent over SMTP 
so template should contain fileds like `Subject` and `Content-Type:`. Refer to and example in [templates](/templates/test-tempalte.tmpl) folder.

`Values` portion of a mailer queue message is used as input data for processing the template. 

## Deployment 

Mailer service requires SMTP server to be configured in order to send out emails.  Refer to SMTP section [config.go](/pkg/config.go)

## Failed Messages

Messages that are not sent out successfully are stored in `dead-letter-mailer` queue for future investigation and processing. 
