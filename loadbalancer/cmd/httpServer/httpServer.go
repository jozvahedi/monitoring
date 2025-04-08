package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	Config "github.com/jozvahedi/loadbalancer/loadbalancer/config"
	auth "github.com/jozvahedi/loadbalancer/loadbalancer/internal/auth"
	middleware "github.com/jozvahedi/loadbalancer/loadbalancer/internal/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received.",
		},
	)
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func init() {
	prometheus.MustRegister(requestsCounter)
}

type Backend struct {
	url   string
	alive bool
}
type ServerPool struct {
	backend []*Backend
	//status  string
	counter int
	mux     sync.Mutex
}

func isBackendAlive(s string) bool {

	r, e := http.Head(s)
	return e == nil && r.StatusCode == 200

}
func (sp *ServerPool) ChangeAliveStatus(i int, alive bool) {
	sp.backend[i].alive = alive
}
func (sp *ServerPool) HealthCheck() {
	for {
		for i, b := range sp.backend {
			alive := isBackendAlive(b.url)
			sp.ChangeAliveStatus(i, alive)

		}
	}

}

func (s *ServerPool) Rotate() Backend {

	s.mux.Lock()
	s.counter = (s.counter) % len(s.backend)
	b := s.counter
	s.mux.Unlock()
	s.counter++
	return *s.backend[b]

}
func (s *ServerPool) GetNextValidPeer() (b Backend) {
	for i := 0; i < len(s.backend); i++ {
		nextPeer := s.Rotate()
		if nextPeer.alive {
			return nextPeer
		}

	}
	return
}

var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

var HttpClient = &http.Client{Transport: tr}

func (sp *ServerPool) handler(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	end := sp.GetNextValidPeer()
	requestsCounter.Inc()
	httpClientReq, err := http.NewRequest(req.Method, fmt.Sprintf("%s%s", end.url, req.URL.String()), req.Body)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	for name, headers := range req.Header {
		for _, h := range headers {
			httpClientReq.Header.Add(name, h)
		}
	}

	resp, err := HttpClient.Do(httpClientReq)

	if err != nil {
		fmt.Printf("error %s", err)
		//falt
		return
	}
	//success
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Fprintf(w, " %s", string(body))

}

func HttpServer(server, port string) {
	Serverpl := ServerPool{backend: []*Backend{
		{url: "http://127.0.0.1:8081"},
		{url: "http://127.0.0.1:8082"},
		{url: "http://127.0.0.1:8083"},
		{url: "http://127.0.0.1:8084"},
		{url: "http://127.0.0.1:8085"},
	}}

	go Serverpl.HealthCheck()

	authService := auth.NewBasicAuthService()
	loggingMiddleware := middleware.LoggingMiddleware{}
	ipWhitelistMiddleware := middleware.IPWhitelistMiddleware{
		Whitelist: []string{"127.0.0.1", "::1"},
	}
	authMiddleware := middleware.BasicAuthMiddleware{
		AuthService: authService,
	}
	mid := map[string]middleware.Middleware{"authService": authMiddleware, "loggingMiddleware": loggingMiddleware, "ipWhitelistMiddleware": ipWhitelistMiddleware}

	mux := http.NewServeMux()

	for _, v := range Config.JsonConfigFile.Middelwarepath {

		middlewareArray := []middleware.Middleware{}
		for _, name := range v.Middelware {
			middlewareArray = append(middlewareArray, mid[name.Name])
		}

		mux.HandleFunc(v.Path, middleware.Chain(
			Serverpl.handler, middlewareArray...,
		))
	}

	mux.Handle("/metrics", promhttp.Handler())

	fmt.Println("Main Server run at " + server + ":" + port)
	http.ListenAndServe(server+":"+port, mux)
}
