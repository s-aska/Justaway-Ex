package main

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"log"
	"os"
)

func main() {

	pemFile := os.Getenv("JUSTAWAY_APNS_PEM_PATH") // apns.pem

	cert, pemErr := certificate.FromPemFile(pemFile, "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = os.Args[len(os.Args)-1]
	notification.Topic = "pw.aska.Justaway"
	notification.Payload = payload.NewPayload().Alert("hello")

	client := apns2.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("APNs ID:", res.ApnsID)
}

// http://sreecharans.blogspot.jp/2011/08/how-to-build-apple-push-notification.html
// - Launch Keychain Assistant from your local Mac and from the login keychain, filter by the Certificates category.
//   You will see an expandable option called “Apple Push Services”
// - Expand this option then right click on “Apple Push Services” > Export “Apple Push Services ID123″. Save this as apns-cert.p12 file somewhere you can access it.
// - Do the same again for the “Private Key” that was revealed when you expanded “Apple Push Services” ensuring you save it as apns-key.p12 file.
// openssl pkcs12 -clcerts -nokeys -out apns-cert.pem -in apns-cert.p12
// openssl pkcs12 -nocerts -out apns-key.pem -in apns-key.p12
// openssl rsa -in apns-key.pem -out apns-key-noenc.pem
// cat apns-cert.pem apns-key-noenc.pem > apns.pem
