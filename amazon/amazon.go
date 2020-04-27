package amazon

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

/*
PressType a constent of how the IOT button was pressed.
*/
type PressType int

/*
ErrNoDetails is returned from the connection if all the environment variables aren't defined
*/
var ErrNoDetails error = errors.New("Need CLIENT_ID, MQTT_HOST, CERT_PEM & KEY_PEM")

const (
	// Unknown we don't know what was pressed
	Unknown PressType = iota
	// Single a single press
	Single
	// Double a double press
	Double
	// Long a long press
	Long
)

var stringToPressType map[string]PressType = map[string]PressType{
	"SINGLE": Single,
	"DOUBLE": Double,
	"LONG":   Long,
}

/*
ConnectAndSubscribe will connect to the Amazon MQTT and returns a channel
where all received messages will be published
*/
func ConnectAndSubscribe() (chan PressType, error) {

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

	ch := make(chan PressType)

	var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		// fmt.Printf("MSG: %s\n", msg.Payload())
		var buttonData buttonMessage
		if err := json.Unmarshal(msg.Payload(), &buttonData); err != nil {
			panic(err)
		}

		buttonPress, ok := stringToPressType[buttonData.ClickType]
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

type buttonMessage struct {
	SerialNumber   string
	BatteryVoltage string
	ClickType      string
}

func newTLSConfig() (*tls.Config, error) {

	certPem, ok := os.LookupEnv("CERT_PEM")
	if !ok {
		return nil, ErrNoDetails
	}
	keyPem, ok := os.LookupEnv("KEY_PEM")
	if !ok {
		return nil, ErrNoDetails
	}

	// Import client certificate/key pair
	cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
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
