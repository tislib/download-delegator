package service

import (
	"context"
	"download-delegator/model"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

import (
	"testing"
)

type MockServer struct {
	srv  *http.Server
	addr string
}

func (s *MockServer) ListenAndServeAsync() error {
	srv := s.srv
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.addr = ln.Addr().String()

	go func() {
		srv.Serve(ln)
	}()

	return nil
}

const (
	SMALL  string = "small"
	MEDIUM        = "medium"
	LARGE         = "large"
)

func (s *MockServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/" + SMALL:
		writer.Write([]byte("hello world"))
		break
	case "/" + MEDIUM:
		for i := 0; i < 1000; i++ {
			writer.Write([]byte("hello world"))
		}
		break
	case "/" + LARGE:
		for i := 0; i < 50000; i++ {
			writer.Write([]byte("hello world"))
		}
		break
	case "/" + SMALL + "-timeout":
		time.Sleep(500)
		writer.Write([]byte("hello world"))
		break
	case "/" + MEDIUM + "-timeout":
		time.Sleep(500)
		for i := 0; i < 1000; i++ {
			writer.Write([]byte("hello world"))
		}
		break
	case "/" + LARGE + "-timeout":
		time.Sleep(500)
		for i := 0; i < 50000; i++ {
			writer.Write([]byte("hello world"))
		}
		break
	}
}

func (s *MockServer) Start() {
	s.srv = &http.Server{Addr: ":0", Handler: s}

	s.ListenAndServeAsync()
}

func (s *MockServer) GetBaseUrl() string {
	return "http://" + s.addr + "/"
}

func BenchmarkSimpleDownload(b *testing.B) {
	b.StopTimer()
	mockServer := new(MockServer)
	mockServer.Start()
	b.StartTimer()

	b.Run("small", func(b *testing.B) {
		simpleDownloadTest(b, mockServer, SMALL)
	})

	b.Run("medium", func(b *testing.B) {
		simpleDownloadTest(b, mockServer, MEDIUM)
	})

	b.Run("large", func(b *testing.B) {
		simpleDownloadTest(b, mockServer, LARGE)
	})
}

func BenchmarkDownloadTimeout(b *testing.B) {
	b.StopTimer()
	mockServer := new(MockServer)
	mockServer.Start()
	b.StartTimer()

	b.Run("small", func(b *testing.B) {
		simpleDownloadTest(b, mockServer, SMALL+"-timeout")
	})

	b.Run("medium", func(b *testing.B) {
		simpleDownloadTest(b, mockServer, MEDIUM+"-timeout")
	})

	b.Run("large", func(b *testing.B) {
		simpleDownloadTest(b, mockServer, LARGE+"-timeout")
	})
}

func simpleDownloadTest(b *testing.B, mockServer *MockServer, mockName string) {
	b.RunParallel(func(pb *testing.PB) {
		config := model.DownloadConfig{
			Url: mockServer.GetBaseUrl() + mockName,
		}
		for pb.Next() {
			err := DownloaderServiceInstance.Get(ioutil.Discard, context.TODO(), config)

			if err != nil {
				log.Panic(err)
			}
		}
	})
}
