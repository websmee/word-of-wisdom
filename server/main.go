package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/websmee/word-of-wisdom/pow"
)

const salt = "89475ht48907fgh4werg8oj"

type Specification struct {
	PowClientMinDifficulty int64 `split_words:"true" default:"15"`
	PowClientMaxDifficulty int64 `split_words:"true" default:"25"`
	PowServerDifficulty    int64 `split_words:"true" default:"15"`
}

func main() {
	var s Specification
	envconfig.MustProcess("", &s)

	fmt.Println("Server started")
	fmt.Println("pow client min difficulty:", s.PowClientMinDifficulty)
	fmt.Println("pow client max difficulty:", s.PowClientMaxDifficulty)
	fmt.Println("pow server difficulty:", s.PowServerDifficulty)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net listen error: %v\n", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("close listener error: %v\n", err)
		}
	}()

	dm := NewPowDifficultyManager(
		s.PowClientMinDifficulty,
		s.PowClientMaxDifficulty,
		s.PowServerDifficulty,
		time.Millisecond*10,
	)
	go dm.Run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept connection error: %v\n", err)
		}
		go handleConnection(conn, dm)
	}
}

func handleConnection(conn net.Conn, dm *PowDifficultyManager) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close connection error: %v\n", err)
		}
	}()

	difficulty := int(dm.GetDifficulty())
	challenge := pow.Make(conn.RemoteAddr().String()+time.Now().String()+salt, difficulty)
	if _, err := conn.Write(challenge); err != nil {
		log.Printf("write challenge error: %v\n", err)
		return
	}

	solutionBytes := make([]byte, pow.SolutionSize)
	if _, err := conn.Read(solutionBytes); err != nil {
		log.Printf("read solution error: %v\n", err)
		return
	}
	solution := binary.BigEndian.Uint64(solutionBytes)

	if !pow.Verify(challenge, int64(solution)) {
		if _, err := conn.Write([]byte("proof of work verification failed")); err != nil {
			log.Printf("write pow failed error: %v\n", err)
		}
		return
	}

	dm.IncConnections()
	defer dm.DecConnections()

	dm.DoWork() // simulate work
	if _, err := conn.Write([]byte(getWordOfWisdom())); err != nil {
		log.Printf("write word of wisdom error: %v\n", err)
		return
	}
}

func getWordOfWisdom() string {
	quotes := []string{
		"If you want to achieve greatness stop asking for permission. ~Anonymous",
		"Things work out best for those who make the best of how things work out. ~John Wooden",
		"To live a creative life, we must lose our fear of being wrong. ~Anonymous",
		"If you are not willing to risk the usual you will have to settle for the ordinary. ~Jim Rohn",
		"Trust because you are willing to accept the risk, not because it's safe or certain. ~Anonymous",
		"Take up one idea. Make that one idea your life - think of it, dream of it, live on that idea. Let the brain, muscles, nerves, every part of your body, be full of that idea, and just leave every other idea alone. This is the way to success. ~Swami Vivekananda",
		"All our dreams can come true if we have the courage to pursue them. ~Walt Disney",
		"Good things come to people who wait, but better things come to those who go out and get them. ~Anonymous",
		"If you do what you always did, you will get what you always got. ~Anonymous",
		"Success is walking from failure to failure with no loss of enthusiasm. ~Winston Churchill",
	}

	return quotes[rand.Intn(len(quotes))]
}

type PowDifficultyManager struct {
	difficulty        atomic.Int64
	connectionNum     atomic.Int64
	minDifficulty     int64
	maxDifficulty     int64
	workDifficulty    int64
	checkFrequency    time.Duration
	prevConnectionNum int64
}

func NewPowDifficultyManager(
	minDifficulty int64,
	maxDifficulty int64,
	workDifficulty int64,
	checkFrequency time.Duration,
) *PowDifficultyManager {
	manager := &PowDifficultyManager{
		minDifficulty:  minDifficulty,
		maxDifficulty:  maxDifficulty,
		workDifficulty: workDifficulty,
		checkFrequency: checkFrequency,
	}
	manager.difficulty.Store(minDifficulty)

	return manager
}

func (m *PowDifficultyManager) Run() {
	for {
		time.Sleep(m.checkFrequency)
		if m.prevConnectionNum < m.connectionNum.Load() && m.difficulty.Load() < m.maxDifficulty {
			m.difficulty.Add(2)
		} else if m.prevConnectionNum >= m.connectionNum.Load() && m.difficulty.Load() > m.minDifficulty {
			m.difficulty.Add(-2)
		}
		m.prevConnectionNum = m.connectionNum.Load()
	}
}

func (m *PowDifficultyManager) GetDifficulty() int64 {
	return m.difficulty.Load()
}

func (m *PowDifficultyManager) IncConnections() {
	m.connectionNum.Add(1)
}

func (m *PowDifficultyManager) DecConnections() {
	m.connectionNum.Add(-1)
}

func (m *PowDifficultyManager) DoWork() {
	pow.Solve(pow.Make(salt, int(m.workDifficulty)))
}
