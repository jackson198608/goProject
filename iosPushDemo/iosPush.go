package main

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"log"
)

func main() {

	//cert, pemErr := certificate.FromP12File("lingdang.p12", "gouminwang")
	cert, pemErr := certificate.FromPemFile("pro-lingdang.pem", "gouminwang")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns.Notification{}
	notification.DeviceToken = "fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2"
	notification.Topic = "com.goumin.bell"
	//notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below
	notification.Payload = []byte(`{"aps":{"alert":"go push here!!!","sound":"default","badge":1,"type":6,"mark":""}}`) // See Payload section below

	client := apns.NewClient(cert).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("APNs ID:", res.ApnsID)
}
