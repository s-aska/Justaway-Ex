package main

import (
	"fmt"
	"github.com/anachronistic/apns"
	"os"
)

func main() {
	payload := apns.NewPayload()
	payload.Alert = "Ok!"
	payload.Badge = 1
	payload.Sound = "bingbong.aiff"

	pn := apns.NewPushNotification()
	pn.DeviceToken = os.Args[len(os.Args)-1]
	pn.AddPayload(payload)

	certificateFile := os.Getenv("JUSTAWAY_APNS_SANDBOX_CERT_PATH")  // apns-dev-cert.pem
	keyFile := os.Getenv("JUSTAWAY_APNS_SANDBOX_KEY_NOENC_PEM_PATH") // apns-dev-key-noenc.pem

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", certificateFile, keyFile)
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	fmt.Println("  Token:", pn.DeviceToken)
	fmt.Println("  Alert:", alert)
	fmt.Println("Success:", resp.Success)
	fmt.Println("  Error:", resp.Error)
}

// http://sreecharans.blogspot.jp/2011/08/how-to-build-apple-push-notification.html
// - Launch Keychain Assistant from your local Mac and from the login keychain, filter by the Certificates category. You will see an expandable option called “Apple Development Push Services”
// - Expand this option then right click on “Apple Development Push Services” > Export “Apple Development Push Services ID123″. Save this as apns-dev-cert.p12 file somewhere you can access it.
// - Do the same again for the “Private Key” that was revealed when you expanded “Apple Development Push Services” ensuring you save it as apns-dev-key.p12 file.
// openssl pkcs12 -clcerts -nokeys -out apns-dev-cert.pem -in apns-dev-cert.p12
// openssl pkcs12 -nocerts -out apns-dev-key.pem -in apns-dev-key.p12
// openssl rsa -in apns-dev-key.pem -out apns-dev-key-noenc.pem
