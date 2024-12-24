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

// Server handle messages from seller and send it to Buyer App.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	log "github.com/golang/glog"

	localstackclient "partner-innovation.googlesource.com/googleondcaccelerator.git/shared/clients/localstack-aws-client"
	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/config"
	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/models/model"
)

type server struct {
	pubsubClient *sns.Client
	httpClient   *http.Client
	config       config.BuyerAdapterConfig
	subs         []string
}

func main() {
	flag.Set("alsologtostderr", "true")
	ctx := context.Background()

	// configPath, ok := os.LookupEnv("CONFIG")
	// if !ok {
	// 	log.Exit("CONFIG env is not set")
	// }

	conf, err := config.Read[config.BuyerAdapterConfig]("./bap_adapter_config.json")
	if err != nil {
		log.Exit(err)
	}

	snsClient, err := localstackclient.NewSNSClient(ctx)
	if err != nil {
		log.Exit(err)
	}

	srv, err := initServer(ctx, http.DefaultClient, snsClient, conf)
	if err != nil {
		log.Exit(err)
	}
	log.Info("Server initialization successs")

	srv.serve(ctx)
}

func subscriptionExists(client *sns.Client, topicArn string) (bool, error) {
	input := &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	}

	for {
		output, err := client.ListSubscriptionsByTopic(context.TODO(), input)
		if err != nil {
			return false, err
		}

		// for _, sub := range output.Subscriptions {
		// 	if aws.ToString(sub.Endpoint) == endpoint {
		// 		return true, nil
		// 	}
		// }

		if output.NextToken == nil {
			break
		}

		input.NextToken = output.NextToken
	}

	return false, nil
}

func initServer(ctx context.Context, httpClient *http.Client, pubsubClient *sns.Client, conf config.BuyerAdapterConfig) (*server, error) {
	// validate clients
	if httpClient == nil {
		return nil, errors.New("init server: HTTP client is nil")
	}
	if pubsubClient == nil {
		return nil, errors.New("init server: Pub/Sub client is nil")
	}

	//validate the subscriptions
	subs := make([]string, 0, len(conf.SubscriptionID))
	// for _, subID := range conf.SubscriptionID {
	// 	exist, err := subscriptionExists(pubsubClient, subID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("init server: failed in checking if the subscription %q exists: %v", subID, err)
	// 	}
	// 	if !exist {
	// 		return nil, fmt.Errorf("init server: subscription %q does not exist", subID)
	// 	}

	// 	subs = append(subs, subID)
	// }

	// pubsubClient.Subscribe(ctx, &sns.SubscribeInput{
	// 	Protocol: aws.String("http"),
	// 	TopicArn: aws.String(conf.SubscriptionID[0]),
	// 	Endpoint: aws.String("http://localhost:8091"),
	// })

	server := &server{
		pubsubClient: pubsubClient,
		httpClient:   httpClient,
		config:       conf,
		subs:         subs,
	}
	return server, nil
}

// serve handles multiple Pub/Sub subscriptions in parallel.
// func (s *server) serve(ctx context.Context) error {
// 	g, ctx := errgroup.WithContext(ctx)

// 	for _, sub := range s.subs {
// 		// create a subscription as a local variable
// 		// so that it can be passed to handleSubscription safely.
// 		sub := sub
// 		g.Go(func() error {
// 			return s.handleSubscription(ctx, sub)
// 		})
// 	}

// 	log.Info("Ready to receive messages")
// 	return g.Wait()
// }

func (s *server) serve(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/sns", s.handleSubscription)
	if err := http.ListenAndServe("0.0.0.0:8091", mux); err != nil {
		log.Fatalf("failed to start HTTP server, %v", err)
	}
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

// handleSubscription receives and handles messages from the Pub/Sub subscription.
func (s *server) handleSubscription(w http.ResponseWriter, r *http.Request) {

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

		//buyerEndpoint := s.config.BuyerAppURL + "/" + messageData.Action
		buyerEndpoint := os.Getenv("BUYER_APP_URL") + "/" + messageData.Action //for the time being take app url from docker compose.
		response, err := s.httpClient.Post(buyerEndpoint, "application/json", bytes.NewReader([]byte(payload.Message)))
		if err != nil {
			log.Errorf("Calling Buyer App failed: %v", err)
			return
		}
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Errorf("Reading response body failed: %v", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			log.Errorf("Calling Buyer App got an error: status code %d, body %s", response.StatusCode, responseBody)
			return
		}

		log.Info("Handle the message successfully")
		w.WriteHeader(http.StatusOK)
	} else {
		request, err := http.NewRequest(http.MethodGet, payload.SubscribeURL, nil)
		if err != nil {
			log.Error("failed to confirm topic subscription", err.Error())
			return
		}
		// send a request to ONDC network
		_, err = s.httpClient.Do(request)
		w.WriteHeader(http.StatusOK)
	}
}
