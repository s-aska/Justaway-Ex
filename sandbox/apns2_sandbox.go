package main

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"log"
	"os"
)

func main() {

	certificateFile := os.Getenv("JUSTAWAY_APNS_CERT_PATH")  // cert.pem

	cert, pemErr := certificate.FromPemFile(certificateFile, "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns.Notification{}
	notification.DeviceToken = os.Args[len(os.Args)-1]
	notification.Topic = "pw.aska.Justaway"
	notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

	client := apns.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("APNs ID:", res.ApnsID)
}
