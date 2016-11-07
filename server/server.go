package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "../pb"
	"regexp"
	"net/http"
	"io/ioutil"
	"sync/atomic"
	"os"
	"sync"
	"time"
)

type server struct{}

var (
	port string = ":50051"
	errorCnt uint64 = 0
	success uint64 = 0
	l sync.Mutex
)

func (s *server) VerifySite(ctx context.Context, in *pb.VerifySiteRequest) (*pb.VerifySiteResponse, error) {
	l.Lock()
	resp, err := http.Get(in.SiteUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	time.Sleep(time.Second)
	l.Unlock()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	re := regexp.MustCompile("<meta."+in.Key+"=\"(.*)\".*content=\"(.*)\".*\\/>")
	if re.MatchString(string(body)) {
		atomic.AddUint64(&success, 1)
		return &pb.VerifySiteResponse{Found: true}, nil

	} else {
		atomic.AddUint64(&errorCnt, 1)
		return &pb.VerifySiteResponse{Found: false}, nil
	}
}

func main() {
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	log.Println(port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Ошибка при прослушке порта: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCheckSiteServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка при обслуживании: %v", err)
	}
}