# Word of wisdom + POW
The app designed to demonstrate the usage of proof-of-work algorithm for DOS protection.
There are two main components:

1. The TCP server that generates the challenge and verifies the solution.
2. The client that solves the challenge and sends the solution to the server.
 
For successful solution the server sends a random word of wisdom to the client.
In order to imitate some work on the server before sending the response the server also solves proof-of-work challenge.
The difficulty of the challenge is dynamically adjusted based on the number of active connections to the server.

## The choice of the algorithm
The algorithm is based on the [Hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm.
It's a CPU-bound algorithm that requires a lot of CPU time to solve.
We didn't use Memory-bound algorithms because out of memory situation can potentially cause app failure instead of just to require more time to compute.
Our benchmark showed stable results for the algorithm based on difficulty.

## Usage
Build and run the server:
```
docker build -f .\server.Dockerfile -t wow .
docker run --rm -p 8080:8080 --cpus=1 -e POW_SERVER_DIFFICULTY=15 -e POW_CLIENT_MIN_DIFFICULTY=15 -e POW_CLIENT_MAX_DIFFICULTY=25 wow
```
Parameters:
- `POW_SERVER_DIFFICULTY` - the difficulty of the challenge for the server (to simulate some work).
- `POW_CLIENT_MIN_DIFFICULTY` - the minimum difficulty of the challenge for the client.
- `POW_CLIENT_MAX_DIFFICULTY` - the maximum difficulty of the challenge for the client.

Build and run the client:
```
docker build -f .\client.Dockerfile -t wow-client .
docker run --rm -p 8082:8081 --cpus=1 -e TEST_REQUESTS_NUM=5000 -e TEST_REQUESTS_INTERVAL=1 wow-client
```
Parameters:
- `TEST_REQUESTS_NUM` - the number of requests to send to the server.
- `TEST_REQUESTS_INTERVAL` - the interval between requests in milliseconds.

## Results
Benchmark for challenge solving on my machine:
```
goos: windows
goarch: amd64
pkg: github.com/websmee/word-of-wisdom/pow
cpu: 12th Gen Intel(R) Core(TM) i7-12700F
BenchmarkSolve5
BenchmarkSolve5-20         81865             14371 ns/op
BenchmarkSolve10
BenchmarkSolve10-20         5996            188798 ns/op
BenchmarkSolve15
BenchmarkSolve15-20          895           1323513 ns/op
BenchmarkSolve20
BenchmarkSolve20-20           73          16204967 ns/op
BenchmarkSolve25
BenchmarkSolve25-20            1        3620536800 ns/op
PASS
```

Based on the benchmark results we can take the range 15 to 25 difficulty as reasonable for our tests.
The difficulty of 15 is taking 0.0013 seconds to solve.
The difficulty of 25 is taking 3.5 seconds to solve.

For our test we had server parameters as follows:
- `POW_SERVER_DIFFICULTY=15`
- `POW_CLIENT_MIN_DIFFICULTY=15`
- `POW_CLIENT_MAX_DIFFICULTY=25`

First we run the client with 5000 requests and 1 millisecond interval between requests
and got 0.6 seconds average client solution time and 0.08 seconds average server response time.
Then we run 5 clients simultaneously with 5000 requests and 1 millisecond interval between requests for each client
and got 2 seconds average client solution time and 0.09 seconds average server response time.
The difficulty was automatically adjusted to handle the load,
the more requests were spawned the harder it was for each client to solve new challenges and for the server the load stayed around the same.

## Conclusion
The algorithm is working as expected.
It's possible to adjust the difficulty of the challenge to handle the load on the fly.
The difficulty range for the server should be chosen based on the work server does,
required response times, CPU resources and expected load highs and lows.

## P.S.
There's a DOS attack called **SYN flood**, but we didn't touch it in this exercise,
because it requires low level control of the TCP handshake and much more time
for research and development.