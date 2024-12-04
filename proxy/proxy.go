package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Proxy struct {
	Targets []*url.URL
	Cache   *redis.Client
	Counter uint64
}

func NewProxy(targets []string, redisAddr string) *Proxy {
	urls := make([]*url.URL, len(targets))
	for i, target := range targets {
		parsed, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Error parsing target URL: %v", err)
		}
		urls[i] = parsed
	}

	cache := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &Proxy{
		Targets: urls,
		Cache:   cache,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	sortedQuery := make([]string, 0, len(queryParams))
	for key := range queryParams {
		sortedQuery = append(sortedQuery, key)
	}
	sort.Strings(sortedQuery)

	queryString := ""
	for _, key := range sortedQuery {
		queryString += key + "=" + queryParams.Get(key) + "&"
	}

	cacheKey := r.Method + ":" + r.URL.Path + "?" + queryString

	cached, err := p.Cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		log.Printf("Cache hit: %s, %s", cacheKey, cached)
		w.Write([]byte(cached))
	} else {
		log.Printf("Cache miss: %s", cacheKey)
	}

	target := p.getNextTarget()
	if target == nil {
		log.Println("No available DW node!")
		http.Error(w, "The DW node is unavailable", http.StatusServiceUnavailable)
		return
	}

	log.Printf("The request will be processed by the DW node: %s", target.String())

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading the response: %v", err)
			return err
		}
		resp.Body = io.NopCloser(bytes.NewReader(body))
		resp.Body.Close()

		err = p.Cache.Set(context.Background(), cacheKey, string(body), 0).Err()
		if err != nil {
			log.Printf("Error saving the response to cache: %v", err)
		}

		return nil
	}

	proxy.ServeHTTP(w, r)
}

func (p *Proxy) getNextTarget() *url.URL {
	index := atomic.AddUint64(&p.Counter, 1)
	return p.Targets[index%uint64(len(p.Targets))]
}

func Run(port string, nodePorts []string) {
	if len(nodePorts) == 0 {
		nodePorts = []string{"http://localhost:8081", "http://localhost:8082"}
	} else {
		for idx, nPort := range nodePorts {
			nodePorts[idx] = "http://localhost:" + nPort
		}
	}

	proxy := NewProxy(nodePorts, "localhost:6379")

	fmt.Println("The proxy is running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))
}

