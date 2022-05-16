package main

import (
	"log"
	"time"

	"github.com/Philip-21/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false ///not keeping the cooection always active unless needed
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	//used in production for a live server
	//Server.Username
	//Server.Password
	//Server.Encryption

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject) //from an email, To an email
	//sending the email in html format
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Panicln(err)
	} else {
		log.Println("Email sent")
	}
}

func listenForMail() {
	//go routine runs indefinitely in the background to make things asyncronasly(when a message is sent the application,it still runs without it stopping )
	go func() {
		//listens all the time for incoming data
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}
