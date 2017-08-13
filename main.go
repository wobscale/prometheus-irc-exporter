package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	gircclient "github.com/goshuirc/irc-go/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	http.Handle("/probe", http.HandlerFunc(probeHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		w.WriteHeader(400)
		w.Write([]byte("no target provided"))
		return
	}

	tgt, err := url.Parse(target)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("invalid target; not valid uri"))
		return
	}
	if tgt.Scheme != "ircs" && tgt.Scheme != "irc" {
		w.WriteHeader(400)
		w.Write([]byte("target must have ircs or irc scheme"))
		return
	}
	// TODO: prometheus timeout header

	registry := prometheus.NewRegistry()

	up := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "target irc server is up",
	})
	registry.MustRegister(up)

	reactor := gircclient.NewReactor()
	defer reactor.Shutdown("probe done")

	server := reactor.CreateServer("probe")
	server.InitialNick = fmt.Sprintf("promirc_%s", rand.Int31())
	server.InitialUser = "promirc"

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	err = server.Connect(
		tgt.Host,
		tgt.Scheme == "ircs",
		tlsConfig,
	)
	if err != nil {
		log.Printf("[ERROR] Could not connect to target: %v", err)
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
		return
	}

	// TODO timeout; implement this ourselves
	server.WaitForConnection()
	up.Inc()

	if tgt.Scheme == "ircs" {
		state := server.RawConnection.(*tls.Conn).ConnectionState()
		tlsExpiryGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "ssl_expiry",
			Help: "ssl expiry in unixtime, or zero for error",
		})
		registry.MustRegister(tlsExpiryGauge)
		earliest := time.Time{}
		if len(state.PeerCertificates) != 0 {
			earliest = state.PeerCertificates[0].NotAfter
		}

		for _, cert := range state.PeerCertificates {
			if cert.NotAfter.Before(earliest) {
				earliest = cert.NotAfter
			}
		}
		tlsExpiryGauge.Set(float64(earliest.Unix()))
	}

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
