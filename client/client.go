package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "../pb"


	"sync"
)

const (
	address     = "localhost:50051"
)
var (
	SiteArray []string = []string{
		"https://ya.ru/",
		"https://www.google.ru/",
		"https://mail.ru/",
		"https://www.avito.ru/",
		"https://youtube.com/",
		"https://ok.ru/",
		"https://www.avito.ru/",
		"https://news.yandex.ru/",
		"http://www.rambler.ru/",
		"http://www.rambler.ru/",
		"http://echo.msk.ru",
		"https://www.lenta.ru/",
		"http://www.fontanka.ru/",
		"https://gmail.com/",
	}
	tag string = "charset"
	elem uint64 = 0
	l sync.Mutex
	goRoutineCount = 5
	c pb.CheckSiteClient
)
func getSite() string {
	l.Lock()
	site := ""
	if int(elem) < len(SiteArray) {
		site = SiteArray[elem]
		elem++
	}
	l.Unlock()
	return site
}
func checkSite() {
	site_url := getSite()
	if site_url != "" {
		r, err := c.VerifySite(context.Background(), &pb.VerifySiteRequest{SiteUrl:site_url,Key:tag})
		if err != nil {
			log.Fatalf("Ошибка проверки[%s]: %v", err,site_url)
		}
		if r.Found {
			log.Printf("Ключ найден[%s]\n",site_url)
		} else {
			log.Printf("Ключ ненайден[%s]\n",site_url)
		}
		checkSite()
	} else {
		log.Println("Закончил работу")
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не смогли приконнектиться: %v", err)
	}
	defer conn.Close()
	c = pb.NewCheckSiteClient(conn)
	for goRoutineCount > 0 {
		goRoutineCount--
		checkSite()
	}
	select{}
}