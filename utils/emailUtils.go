package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"tmp-invoicely.co/models"
)

// Struct to hold invoice email data
type emailData struct {
	Invoice models.Invoice
	Company models.InvoiceConfiguration
}

// LoadEmailTemplate loads an HTML email template and fills it with invoice data
func LoadEmailTemplate(templatePath string, invoice models.Invoice, company models.InvoiceConfiguration) (string, error) {
	data := emailData{invoice, company}

	log.Println("4")

	// Parse the email template file
	tmpl, err := template.New("invoice").ParseFiles(templatePath)
	if err != nil {
		log.Println(err, templatePath)
		return "", fmt.Errorf("failed to load template from %s: %w", templatePath, err)
	}

	log.Println("5")

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	log.Println("6")

	return tpl.String(), nil
}

func SendWelcomeEmail(name, email string) (*sesv2.SendEmailOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(GetEnv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	client := sesv2.NewFromConfig(awsConfig)
	awsEmail := GetEnv("AWS_AUTH_EMAIL")

	if awsEmail == "" {
		return nil, fmt.Errorf("AWS_AUTH_EMAIL is missing in .env file")
	}

	subject := "Welcome to Invoecly!"
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
	    <meta name="viewport" content="width=device-width, initial-scale=1">
	</head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
	    <div style="max-width: 600px; margin: auto; background: #ffffff; padding: 20px; 
	                box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); border-radius: 10px; text-align: center;">
	        
	        <h1 style="color: #333;">Welcome to Invoecly, %s!</h1>
	        <p style="color: #666; font-size: 16px; line-height: 1.6;">
	            We're excited to help you manage your invoices effortlessly. Get started by setting up your profile, adding your first customer, and creating your first invoice.
	        </p>
	        <a href="https://app-invoicely-co.vercel.app/" 
	           style="display: inline-block; background: #007BFF; color: #ffffff; text-decoration: none; 
	                  padding: 12px 24px; border-radius: 5px; margin-top: 20px; font-size: 16px;">
	           Get Started
	        </a>
	        <p style="margin-top: 20px; font-size: 14px; color: #999;">
	            Need help? Contact us at <a href="mailto:support@invoecly.com" style="color: #007BFF;">support@app-invoecly-co.com</a>
	        </p>
	    </div>
	</body>
	</html>`, name)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &awsEmail,
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: &subject},
				Body: &types.Body{
					Html: &types.Content{Data: &htmlBody},
				},
			},
		},
	}

	result, err := client.SendEmail(ctx, input)
	if err != nil {
		log.Printf("Failed to send welcome email: %v\n", err)
		return nil, err
	}

	log.Println("Welcome email sent successfully")
	return result, err
}

// SendInvoiceEmail sends an invoice email to the customer
func SendInvoiceEmail(invoice models.Invoice, company models.InvoiceConfiguration) (*sesv2.SendEmailOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("1")

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(GetEnv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}
	log.Println("2")

	client := sesv2.NewFromConfig(awsConfig)
	awsEmail := GetEnv("AWS_AUTH_EMAIL")

	if awsEmail == "" {
		return nil, fmt.Errorf("AWS_AUTH_EMAIL is missing in .env file")
	}
	log.Println("3")

	if invoice.Customer.Email == "" {
		return nil, fmt.Errorf("recipient email is missing in the invoice")
	}
	log.Println(invoice, company)
	// Load and parse the invoice email template
	htmlBody, err := LoadEmailTemplate("static/templates/invoice_template.html", invoice, company)

	if err != nil {
		return nil, fmt.Errorf("failed to load invoice email template: %w", err)
	}

	subject := fmt.Sprintf("Invoice #%s from %s", invoice.InvoiceNo, company.Name)

	log.Println(htmlBody)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &awsEmail,
		Destination: &types.Destination{
			ToAddresses: []string{invoice.Customer.Email},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: &subject},
				Body: &types.Body{
					Html: &types.Content{Data: &htmlBody},
				},
			},
		},
	}

	result, err := client.SendEmail(ctx, input)

	if err != nil {
		log.Printf("Failed to send invoice email (InvoiceID: %s, CustomerEmail: %s): %v\n", invoice.InvoiceID, invoice.Customer.Email, err)
		return nil, err
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	log.Println("Invoice email sent successfully:\n", string(resultJSON))

	log.Println("Invoice email sent successfully", invoice.Customer.Email)
	return result, nil
}
