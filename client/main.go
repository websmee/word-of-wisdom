package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/websmee/word-of-wisdom/pow"
)

type Specification struct {
	TestRequestsNum      int    `split_words:"true" default:"0"`
	TestRequestsInterval int    `split_words:"true" default:"1"`
	WOWServiceAddr       string `split_words:"true" default:"host.docker.internal:8080"`
}

func main() {
	var s Specification
	envconfig.MustProcess("", &s)

	if s.TestRequestsNum > 0 {
		fmt.Println("requests", s.TestRequestsNum)
		fmt.Println("interval (ms)", s.TestRequestsInterval)
		fmt.Println("running test...")
		st, rt := runTest(s.WOWServiceAddr, s.TestRequestsNum, time.Millisecond*time.Duration(s.TestRequestsInterval))
		fmt.Println("average client solution time (nano):", st)
		fmt.Println("average server response time (nano):", rt)
		return
	}

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("net listen error: %v\n", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("close listener error: %v\n", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept connection error: %v\n", err)
		}
		go handleConnection(conn, s.WOWServiceAddr)
	}
}

func handleConnection(conn net.Conn, addr string) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close connection error: %v\n", err)
		}
	}()

	request := make([]byte, 255)
	if _, err := conn.Read(request); err != nil && err != io.EOF {
		log.Printf("read request error: %v\n", err)
		return
	}

	response, err, _, _ := getResponse(addr)
	if err != nil {
		log.Printf("get response error: %v\n", err)
		r := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       http.NoBody,
		}
		if err := r.Write(conn); err != nil {
			log.Printf("write error response failed: %v\n", err)
		}
		return
	}

	r := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(bytes.Trim([]byte(response), "\x00"))),
	}
	if err := r.Write(conn); err != nil {
		log.Printf("write response error: %v\n", err)
		return
	}
}

func getResponse(addr string) (string, error, int64, int64) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("unable to open tcp connection: %w", err), 0, 0
	}
	defer c.Close()

	challenge := make([]byte, pow.ChallengeSize)
	if _, err := c.Read(challenge); err != nil {
		return "", fmt.Errorf("read challenge error: %w", err), 0, 0
	}

	solutionStart := time.Now().UnixNano()
	solution := pow.Solve(challenge)
	solutionTime := time.Now().UnixNano() - solutionStart

	solutionBytes := make([]byte, pow.SolutionSize)
	binary.BigEndian.PutUint64(solutionBytes, uint64(solution))
	if _, err = c.Write(solutionBytes); err != nil {
		return "", fmt.Errorf("write solution error: %w", err), 0, 0
	}

	response := make([]byte, 255)
	responseStart := time.Now().UnixNano()
	if _, err := c.Read(response); err != nil {
		return "", fmt.Errorf("read response error: %w", err), 0, 0
	}
	responseTime := time.Now().UnixNano() - responseStart

	return string(response), nil, solutionTime, responseTime
}

func runTest(addr string, requestsNum int, interval time.Duration) (int64, int64) {
	if requestsNum == 0 {
		return 0, 0
	}

	var wg sync.WaitGroup
	solutionTimings := make([]int64, requestsNum)
	responseTimings := make([]int64, requestsNum)
	for i := 0; i < requestsNum; i++ {
		i := i
		time.Sleep(interval)
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err, st, rt := getResponse(addr)
			if err != nil {
				fmt.Println(err)
			}
			solutionTimings[i] = st
			responseTimings[i] = rt
		}()
	}
	wg.Wait()

	sumSolutionTimings := int64(0)
	sumResponseTimings := int64(0)
	for i := 0; i < requestsNum; i++ {
		sumSolutionTimings += solutionTimings[i]
		sumResponseTimings += responseTimings[i]
	}

	return sumSolutionTimings / int64(requestsNum), sumResponseTimings / int64(requestsNum)
}
