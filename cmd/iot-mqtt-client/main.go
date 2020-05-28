package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/hill-daniel/iot-mqtt-client/mqtt"
	pb "github.com/hill-daniel/iot-mqtt-client/proto"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	tcpsPort                = 8883
	tickerIntervalInSeconds = 10
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file", err)
	}

	nano := time.Now().UnixNano()
	rand.Seed(nano)
	log.Printf("seeded with %d", nano)
}

func main() {
	ticker := time.NewTicker(tickerIntervalInSeconds * time.Second)
	done := make(chan bool)
	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)

	host := os.Getenv("IOT_ENDPOINT")
	brokerURL := fmt.Sprintf("tcps://%s:%d%s", host, tcpsPort, "/mqtt")
	tlsConfig, err := createTLSConfig()
	if err != nil {
		log.Fatal("failed to create TLS config", err)
	}
	clientID := os.Getenv("CLIENT_ID")

	publisher := mqtt.New(brokerURL, clientID, tlsConfig)
	topic := os.Getenv("TOPIC")

	go func() {
		for {
			select {
			case <-ticker.C:
				device := createDevice()
				if err := publisher.Publish(topic, device); err != nil {
					log.Printf("failed to publish %v, %v", device, err)
				} else {
					log.Printf("published  %v to topic %s", device, topic)
				}
			case <-sigC:
				log.Printf("got an interrupt, stopping...")
				ticker.Stop()
				close(done)
			}
		}
	}()
	<-done
	fmt.Println("stopped")
}

func createTLSConfig() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(os.Getenv("CA_CERT_PATH"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CA cert file")
	}
	certPool.AppendCertsFromPEM(pemCerts)

	cert, err := tls.LoadX509KeyPair(os.Getenv("DEVICE_CERT_PATH"), os.Getenv("DEVICE_KEY_PATH"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read device cert file or key")
	}

	// just to check
	_, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse certificate")
	}

	return &tls.Config{
		RootCAs:            certPool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}, nil
}

func createDevice() *pb.Device {
	return &pb.Device{
		// does not matter for our use case
		Name: "Test Device",
		// does not matter for our use case
		Id:   4711,
		Lat:  rand.Float64(),
		Long: rand.Float64(),
		// helper function, translate protobuf <-> go time
		LastUpdated: ptypes.TimestampNow(),
	}
}
