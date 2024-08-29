package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
)

const limitCountErrorsRead = 10

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Не указаны параметры сервера: адрес и порт прослушивания")
	}

	addr := os.Args[1] + ":" + os.Args[2]
	cl, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cl.Close()

	br := bufio.NewReader(cl)
	countErrorsRead := 0
	go sendAlive(cl)
	for {
		s, err := br.ReadString('\n')

		if err != nil {
			countErrorsRead++
			if countErrorsRead >= limitCountErrorsRead {
				log.Fatalf("Превышен лимит количества попыток (%d) читать данные с сервера. Возможно сервер недоступен", limitCountErrorsRead)
			}
		}

		if s != "" {
			log.Printf("Пословица: %s", s)
		}
	}
}

func sendAlive(c net.Conn) {
	t := time.AfterFunc(time.Second*1, func() {
		log.Println("Отправляю alive")
		c.Write([]byte("alive"))
	})
	for {
		<-t.C
	}
}
