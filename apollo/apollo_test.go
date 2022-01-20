package pkg

import (
	"fmt"
	"testing"
)


var apollo *Apollo

func init (){
	var err error
	apollo ,err = NewApollo("http://apollo.siwei.com:8080","abroad_api","DEV")

	if err !=nil {
		panic(err)
	}
}

func TestApollo_Mysql(t *testing.T) {
	my, err := apollo.Mysql()
	if err !=nil {
		t.Fatal(err)
	}
	for k,v := range my{
		fmt.Println(k,v)
	}
}

func TestApollo_Redis(t *testing.T) {
	my, err := apollo.Redis()
	if err !=nil {
		t.Fatal(err)
	}
	for k,v := range my{
		fmt.Println(k,v)
	}
}

func TestApollo_Mongo(t *testing.T) {
	my, err := apollo.Mongo()
	if err !=nil {
		t.Fatal(err)
	}
	for k,v := range my{
		fmt.Println(k,v)
	}
}

func TestApollo_RabbitMQ(t *testing.T) {
	my, err := apollo.RabbitMq()
	if err !=nil {
		t.Fatal(err)
	}
	for k,v := range my{
		fmt.Println(k,v)
	}
}


