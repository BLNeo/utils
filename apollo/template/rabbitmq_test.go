package template

import (
	"fmt"
	"testing"
	"time"
)




func TestRabbitChannel_Publish(t *testing.T) {
	tem:= &RabbitMq{
		User:     "guest",
		Password: "guest",
		Host:     "127.0.0.1:5672",
	}
	engin ,err := tem.Engine()
	if err !=nil {t.Fatal(err)}

	for {
		time.Sleep(3000 *time.Millisecond)
		if err:= engin.Publish("test.shenchenxi.testexchange","你好吊毛");err!=nil{
			fmt.Println(err)
			continue
		}
		fmt.Println("success: send message")
	}
}


func TestRabbitChannel_Recieve(t *testing.T) {
	tem:= &RabbitMq{
		User:     "guest",
		Password: "guest",
		Host:     "127.0.0.1:5672",
	}
	engin ,err := tem.Engine()
	if err !=nil {t.Fatal(err)}
	ch ,close,err := engin.ReceiveSub("test.shenchenxi.testexchange","test.name")
	if err !=nil {t.Fatal(err)}
	for {
		select {
		case data:=<- ch:
			fmt.Println(string(data))
		case <- close:
			fmt.Println("收到结束信号")
		}

	}
}


func TestRabbitChannel_RecieveAndPublish(t *testing.T) {
	tem:= &RabbitMq{
		User:     "guest",
		Password: "guest",
		Host:     "127.0.0.1:5672",
	}
	engin ,err := tem.Engine()
	if err !=nil {t.Fatal(err)}
	go func() {
		for  {
			time.Sleep(100 *time.Millisecond)
			engin.Publish("test.shenchenxi.testexchange","你好吊毛")
		}
	}()
	ch ,_,err := engin.ReceiveSub("test.shenchenxi.testexchange","test.name")
	if err !=nil {t.Fatal(err)}
	for {
		data := <-ch
		fmt.Println(string(data))
	}
}