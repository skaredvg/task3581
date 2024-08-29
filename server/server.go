package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const timeoutSaying = time.Second * 3

var proverbs []string

func init() {
	log.Println("Запущена процедура init()")

	f, err := os.OpenFile("..//resource//proverbs.txt", os.O_RDONLY, 0111)
	if err != nil {
		log.Fatalf("Отсутствует файл %s  с пословицами", "saying.txt")
	}
	defer f.Close()

	proverbs = make([]string, 0)

	br := bufio.NewReader(f)
	for {
		l, err := br.ReadString(byte('\n'))
		if err == io.EOF {
			break
		}
		proverbs = append(proverbs, l)
	}
	if len(proverbs) == 0 {
		log.Fatal("Пословицы не загружены")
	}
}

func main() {
	//1. Открыт файл с пословицами (init())
	//2. Загрузить его в слайс (init())
	//3. Создать объект-сервер
	if len(os.Args) < 3 {
		log.Fatal("Не указаны параметры сервера: адрес и порт прослушивания")
	}

	addr := os.Args[1] + ":" + os.Args[2]
	log.Printf("Адрес и порт прослушивания: %s", addr)

	srv, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	//4. Зарегистрировать функцию - обработчик
	//5. Запустить сервер
	//7. Запустить цикл прослушивания подключений

	for {
		c, err := srv.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		go handleConnect(c)
	}
}

func handleConnect(c net.Conn) {
	defer c.Close()

	log.Printf("Обработка соединения %s\n", c.RemoteAddr().String())

	done := make(chan any)
	defer close(done)

	go aliveSession(c, done)

	for {
		select {
		case <-done:
			log.Printf("Соединение с %s закрыто", c.RemoteAddr().String())
			return
		case <-time.After(timeoutSaying):
			res, _ := getRandSaying()
			log.Printf("Отправляю поговорку %s в %s\n", res, c.RemoteAddr().String())
			if res != "" {
				c.Write([]byte(res))
			}
		}

	}

}

func getRandSaying() (string, error) {
	i := int(rand.Int63n(int64(len(proverbs))))
	if i >= len(proverbs) {
		return "", fmt.Errorf("%s", "Неверный индекс для выбора пословицы")
	}

	return proverbs[i], nil
}

func aliveSession(c net.Conn, done chan<- any) {
	b := make([]byte, 10)
	for {
		<-time.After(time.Second * 3)
		log.Printf("Принимаю alive от %s", c.RemoteAddr().String())
		n, _ := c.Read(b)
		if n == 0 {
			done <- 1
			break
		}
	}

}
