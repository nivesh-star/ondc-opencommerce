// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Server handles messages from Pub/Sub topic and send callbacks to Buyer App.
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/benbjohnson/clock"
	log "github.com/golang/glog"
	"google.golang.org/api/option"

	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/clients/keyclient"
	localstackclient "partner-innovation.googlesource.com/googleondcaccelerator.git/shared/clients/localstack-aws-client"
	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/config"
	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/models/model"
	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/signing-authentication/authentication"
)

const (
	EncryptionPrivateKey = "MC4CAQEwBQYDK2VuBCIEIGhebaUMS8k7G3g8gpm/qx4+a8pZglOTO7RXDAcE9ehY"
	request_id           = "smaple_request_id"
	OndcPublicKey        = "MCowBQYDK2VuAyEAduMuZgmtpjdCuxv+Nc49K0cB6tL/Dj3HZetvVN7ZekM="
	SigningPrivateKey    = "7qDUVwqw7Oe13JTa8nAM9ktLj12E4pxBDDxZN8qVwtvGglywynYJUPJo6B/vB5/Rwn2XSAKlKT5snQupvOU4/Q=="
)

type server struct {
	conf         config.RequestActionConfig
	pubsubClient *sns.Client
	httpClient   *http.Client
	keyClient    keyClient
	// transactionClient *transactionclient.Client
	clk clock.Clock

	subs []string
}
type Notification struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SubscribeURL     string `json:"SubscribeURL,omitempty"`
	UnsubscribeURL   string `json:"UnsubscribeURL,omitempty"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
}

type MessageData struct {
	Action string `json:"action,omitempty"`
	Data   string `json:"data,omitempty"`
}

type keyClient interface {
	ServiceSigningPrivateKeyset(context.Context) ([]byte, error)
	AddKey(ctx context.Context, key string, payload []byte) error
}

func main() {
	flag.Set("alsologtostderr", "true")
	ctx := context.Background()

	// configPath, ok := os.LookupEnv("CONFIG")
	// if !ok {
	// 	log.Exit("CONFIG env is not set")
	// }

	conf, err := config.Read[config.RequestActionConfig]("/Users/sandeep.sharma/workspace/nivesh/ondc-opencommerce/shared/config/testdata/callback_action.json")
	if err != nil {
		log.Exit(err)
	}

	keyClient, err := keyclient.NewAwsClient(ctx, conf.ProjectID, conf.SecretID)
	if err != nil {
		log.Exit(err)
	}

	srv, err := initServer(ctx, conf, clock.New(), keyClient, nil, nil)
	if err != nil {
		log.Exit(err)
	}
	//defer srv.close()
	log.Info("Server initialization successs")

	srv.serve(ctx)
}

func initServer(ctx context.Context, conf config.RequestActionConfig, clk clock.Clock, keyClient keyClient, pubsubOpts, transportOpts []option.ClientOption) (*server, error) {
	// validate client
	if keyClient == nil {
		return nil, fmt.Errorf("init server: Key Client is nil")
	}

	// init clients
	pubsubClient, err := localstackclient.NewSNSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("init server: %s", err)
	}

	//TODO: Remove
	keybytes, _ := base64.StdEncoding.DecodeString(SigningPrivateKey)
	err = keyClient.AddKey(ctx, "signingKey", keybytes)
	if err != nil {
		log.Fatal("failed to create signing key in aws secretes manager", err)
	}

	pubsubClient.Subscribe(ctx, &sns.SubscribeInput{
		Protocol: aws.String("http"),
		TopicArn: aws.String(conf.SubscriptionID[0]),
		Endpoint: aws.String("http://cea5-2405-201-4012-867-2d4c-99f6-c9dd-c239.ngrok-free.app/sns"),
	})
	// transactionClient, err := transactionclient.New(ctx, conf.ProjectID, conf.InstanceID, conf.DatabaseID, transportOpts...)
	// if err != nil {
	// 	return nil, fmt.Errorf("init server: %s", err)
	// }

	// validate the subscriptions
	// subs := make([]*pubsub.Subscription, 0, len(conf.SubscriptionID))
	// for _, subID := range conf.SubscriptionID {
	// 	sub := pubsubClient.Subscription(subID)

	// 	ok, err := sub.Exists(ctx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if !ok {
	// 		return nil, fmt.Errorf("init server: subscription %q does not exist", sub.ID())
	// 	}

	// 	subs = append(subs, sub)
	// }
	// _, err = pubsubClient.Subscribe(ctx, &sns.SubscribeInput{
	// 	Protocol: aws.String("http"),
	// 	TopicArn: &conf.SubscriptionID[0],
	// 	Endpoint: aws.String("http://localhost:8081/sns"),
	// })
	if err != nil {
		return nil, fmt.Errorf("init server: failed to subscribe to topic %v", err)
	}

	server := &server{
		conf:         conf,
		pubsubClient: pubsubClient,
		httpClient:   http.DefaultClient,
		keyClient:    keyClient,
		//transactionClient: transactionClient,
		clk:  clk,
		subs: conf.SubscriptionID,
	}
	return server, nil
}

// close closed underlying connections.
// func (s *server) close() {
// 	s.pubsubClient.Close()
// }

// serve handles multiple Pub/Sub subscriptions in parallel.
func (s *server) serve(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/sns", s.handleSubscription)
	if err := http.ListenAndServe("0.0.0.0:8081", mux); err != nil {
		log.Fatalf("failed to start HTTP server, %v", err)
	}
}

// handleSubscription receives and handles messages from the Pub/Sub subscription.
func (s *server) handleSubscription(w http.ResponseWriter, r *http.Request) {
	// defer func() {
	// 	// Ack the msg irrespective of whether the message was successfully processed or not
	// 	// since we do not want the msg to be retried.
	// 	w.WriteHeader(http.StatusOK)
	// 	log.Infof("Handling of message ends")
	// }()

	payload := Notification{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}

	log.Info("Received: ", payload.Type, "MessageId: ", payload.MessageId)

	messageData := MessageData{}
	if payload.Type == "Notification" {
		if err := json.Unmarshal([]byte(payload.Message), &messageData); err != nil {
			log.Fatalf("Failed to unmarshal nested message: %v", err)
		}

		var originalReq model.GenericRequest
		if err := json.Unmarshal([]byte(messageData.Data), &originalReq); err != nil {
			log.Errorf("Unmarshal request failed: %v", err)
			return
		}

		// Determine the request endpoint
		var url string
		if messageData.Action == "search" {
			url = s.conf.GatewayURL
		} else {
			url = originalReq.Context.BppURI
		}

		// Replace BAP data so that the callback is sended to our BAP API Service
		*originalReq.Context.BapID = s.conf.SubscriberID
		*originalReq.Context.BapURI = s.conf.SubscriberURL
		adjustedReqJSON, err := json.Marshal(originalReq)
		if err != nil {
			log.Errorf("Marshal adjusted request failed: %v", err)
			return
		}

		request, err := s.createONDCRequest(r.Context(), messageData.Action, url, adjustedReqJSON)
		if err != nil {
			log.Errorf("Creating request failed: %v", err)
			return
		}

		// bodyBytes, err := io.ReadAll(request.Body)
		// if err != nil {
		// 	return
		// }

		// hhdr := request.Header.Get("Authorization")
		// info, err := authentication.ExtractInfoFromHeader(hhdr)
		// if err != nil {
		// 	log.Errorf("gaye")
		// }
		// kk, _ := base64.StdEncoding.DecodeString("xoJcsMp2CVDyaOgf7wef0cJ9l0gCpSk+bJ0LqbzlOP0=")
		// fmt.Println("sign locl ver: ", authentication.VerifyRequest(info.Signature, bodyBytes, kk, info.Created, info.Expired))
		// send a request to ONDC network
		response, err := s.httpClient.Do(request)
		if err != nil {
			log.Errorf("Sending request to ONDC network failed: %v", err)
			return
		}
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Errorf("Reading response body failed: %v", err)
			return
		}

		// if err := s.storeTransaction(ctx, action, adjustedReqJSON, responseBody); err != nil {
		// 	log.Errorf("Storing transaction failed: %v", err)
		// 	return
		// }

		if response.StatusCode != http.StatusOK {
			log.Infof("Sending request to ONDC network got an error: status code %d, body %s", response.StatusCode, responseBody)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		request, err := http.NewRequest(http.MethodGet, payload.SubscribeURL, nil)
		if err != nil {
			log.Error("failed to confirm topic subscription", err.Error())
			return
		}
		// send a request to ONDC network
		_, err = s.httpClient.Do(request)
		if err != nil {
			log.Errorf("Sending request to ONDC network failed: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// createONDCRequest create a HTTP request for ONDC network with a Authorization header.
func (s *server) createONDCRequest(ctx context.Context, action, url string, body []byte) (*http.Request, error) {
	keyset, err := s.keyClient.ServiceSigningPrivateKeyset(ctx)
	if err != nil {
		return nil, err
	}

	currentTime := s.clk.Now()
	// Use outer bound of request ttl which is 30 seconds.
	expiredTime := currentTime.Add(30 * time.Second)
	authHeader, err := authentication.CreateAuthSignature(body, keyset, currentTime.Unix(), expiredTime.Unix(), s.conf.SubscriberID, s.conf.KeyID)
	if err != nil {
		return nil, err
	}

	requestURL := url + "/" + action
	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authHeader)
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// func (s *server) storeTransaction(ctx context.Context, action string, requestBody []byte, responseBody []byte) error {
// 	switch action {
// 	case "search":
// 		return storeTransaction[model.SearchRequest](ctx, s, action, requestBody, responseBody)
// 	case "select":
// 		return storeTransaction[model.SelectRequest](ctx, s, action, requestBody, responseBody)
// 	case "init":
// 		return storeTransaction[model.InitRequest](ctx, s, action, requestBody, responseBody)
// 	case "confirm":
// 		return storeTransaction[model.ConfirmRequest](ctx, s, action, requestBody, responseBody)
// 	case "track":
// 		return storeTransaction[model.TrackRequest](ctx, s, action, requestBody, responseBody)
// 	case "cancel":
// 		return storeTransaction[model.CancelRequest](ctx, s, action, requestBody, responseBody)
// 	case "update":
// 		return storeTransaction[model.UpdateRequest](ctx, s, action, requestBody, responseBody)
// 	case "status":
// 		return storeTransaction[model.StatusRequest](ctx, s, action, requestBody, responseBody)
// 	case "rating":
// 		return storeTransaction[model.RatingRequest](ctx, s, action, requestBody, responseBody)
// 	case "support":
// 		return storeTransaction[model.SupportRequest](ctx, s, action, requestBody, responseBody)
// 	}
// 	return nil
// }

// func storeTransaction[R model.BPPRequest](ctx context.Context, s *server, action string, requestBody []byte, responseBody []byte) error {
// 	var request R
// 	if err := json.Unmarshal(requestBody, &request); err != nil {
// 		return err
// 	}
// 	msgContext := request.GetContext()

// 	var response model.AckResponse
// 	if err := json.Unmarshal(responseBody, &response); err != nil {
// 		return err
// 	}

// 	data := transactionclient.TransactionData{
// 		ID:              *msgContext.TransactionID,
// 		Type:            "REQUEST-ACTION",
// 		API:             action,
// 		MessageID:       *msgContext.MessageID,
// 		Payload:         request,
// 		ProviderID:      *msgContext.BapID,
// 		MessageStatus:   response.Message.Ack.Status,
// 		ReqReceivedTime: time.Now(),
// 	}

// 	if response.Error != nil {
// 		data.ErrorType = response.Error.Type
// 		data.ErrorCode = *response.Error.Code
// 		data.ErrorMessage = response.Error.Message
// 		data.ErrorPath = response.Error.Path
// 	}

// 	return s.transactionClient.StoreTransaction(ctx, data)
// }
