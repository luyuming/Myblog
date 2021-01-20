package controller

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"net/smtp"
	"sync"
)

type User struct {
	ID       uint
	Name     string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Tel      string `form:"phone" json:"phone"`
	Email    string `form:"email" json:"email"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ProductMsg(msg []byte) {
	var once sync.Once
	once.Do(func() {
		go ConsumeMsg()
	})

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	failOnError(err, "Failed to publish a message")
}

func ConsumeMsg() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
			sendtomail(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func sendtomail(email string) {
	user := "849592709@qq.com"
	password := "xxxxyyyyzzzzwwww"

	// Set up authentication information.
	auth := smtp.PlainAuth("", user, password, "smtp.qq.com")
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	var u User
	json.Unmarshal([]byte(email), &u)
	to := []string{u.Email}
	body := "From: " + "Myblog管理团队<849592709@qq.com>" + "\r\nSubject: " + "欢迎注册Myblog" + "\r\n" + u.Name + "您好,\n" + "\t\t" + "欢迎注册Myblog,这是您的密码：" + u.Password + ", 请妥善保管\n\n" + "马上开始您的博客之旅吧.\n"
	msg := []byte(body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("mail sent sucess.")
	}
}
