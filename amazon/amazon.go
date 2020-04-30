package amazon

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MessageType - what is being requested by pressing a button
type MessageType int

/*
ErrNoDetails is returned from the connection if all the environment variables aren't defined
*/
var ErrNoDetails error = errors.New("Need CLIENT_ID, MQTT_HOST, CERT_PEM & KEY_PEM")

const (
	// Unknown we don't know what was pressed
	Unknown MessageType = iota
	// Playlist - I want a store playlist playing
	Playlist
	// Stop - I want to stop playing
	Stop
	// Album - Play a random album
	Album
)

var stringToMessageType map[string]MessageType = map[string]MessageType{
	"SINGLE": Playlist,
	"DOUBLE": Stop,
	"LONG":   Album,
}

type buttonMessage struct {
	SerialNumber   string
	BatteryVoltage string
	ClickType      string
}

/*
ConnectAndSubscribe will connect to the Amazon MQTT and returns a channel
where all received messages will be published
*/
func ConnectAndSubscribe() (chan MessageType, error) {

	clientID, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		return nil, ErrNoDetails
	}

	host, ok := os.LookupEnv("MQTT_HOST")
	if !ok {
		return nil, ErrNoDetails
	}

	tlsconfig, err := newTLSConfig()
	if err != nil {
		return nil, err
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(host)
	opts.SetClientID(clientID).SetTLSConfig(tlsconfig)

	// Start the connection
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	ch := make(chan MessageType)

	var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("MSG: %s ", msg.Payload())
		var buttonData buttonMessage
		if err := json.Unmarshal(msg.Payload(), &buttonData); err != nil {
			panic(err)
		}
		buttonPress, ok := stringToMessageType[buttonData.ClickType]
		if ok {
			ch <- buttonPress
		} else {
			ch <- Unknown
		}

	}
	c.Subscribe("#", 0, f)
	fmt.Println("Subscribed")

	return ch, nil

}

func newTLSConfig() (*tls.Config, error) {

	certPem, ok := os.LookupEnv("CERT_PEM")
	if !ok {
		return nil, ErrNoDetails
	}
	certData, err := base64.StdEncoding.DecodeString(certPem)
	if err != nil {
		return nil, err
	}
	keyPem, ok := os.LookupEnv("KEY_PEM")
	if !ok {
		return nil, ErrNoDetails
	}
	keyData, err := base64.StdEncoding.DecodeString(keyPem)
	if err != nil {
		return nil, err
	}
	// Import client certificate/key pair
	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: nil,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}
