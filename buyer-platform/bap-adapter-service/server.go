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
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	log "github.com/golang/glog"
	"golang.org/x/sync/errgroup"

	localstackclient "partner-innovation.googlesource.com/googleondcaccelerator.git/shared/clients/localstack-aws-client"
	"partner-innovation.googlesource.com/googleondcaccelerator.git/shared/config"
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

	configPath, ok := os.LookupEnv("CONFIG")
	if !ok {
		log.Exit("CONFIG env is not set")
	}

	conf, err := config.Read[config.BuyerAdapterConfig](configPath)
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

	if err := srv.serve(ctx); err != nil {
		log.Exitf("Serving failed: %v", err)
	}
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
	for _, subID := range conf.SubscriptionID {
		exist, err := subscriptionExists(pubsubClient, subID)
		if err != nil {
			return nil, fmt.Errorf("init server: failed in checking if the subscription %q exists: %v", subID, err)
		}
		if !exist {
			return nil, fmt.Errorf("init server: subscription %q does not exist", subID)
		}

		subs = append(subs, subID)
	}

	server := &server{
		pubsubClient: pubsubClient,
		httpClient:   httpClient,
		config:       conf,
		subs:         subs,
	}
	return server, nil
}

// serve handles multiple Pub/Sub subscriptions in parallel.
func (s *server) serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, sub := range s.subs {
		// create a subscription as a local variable
		// so that it can be passed to handleSubscription safely.
		sub := sub
		g.Go(func() error {
			return s.handleSubscription(ctx, sub)
		})
	}

	log.Info("Ready to receive messages")
	return g.Wait()
}

// handleSubscription receives and handles messages from the Pub/Sub subscription.
func (s *server) handleSubscription(ctx context.Context, sub pubsub.Message) error {
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		defer func() {
			// Ack the msg irrespective of whether the message was successfully processed or not
			// since we do not want the msg to be retried.
			msg.Ack()
			log.Infof("Handling of message %q ends", msg.ID)
		}()

		log.Infof("Receiving a message from %q, message ID: %q", sub.ID(), msg.ID)

		// example actions: `on_search`, `on_select`
		action, ok := msg.Attributes["action"]
		if !ok {
			log.Error(`"action" attribute is not present in the message`)
			return
		}

		buyerEndpoint := s.config.BuyerAppURL + "/" + action
		response, err := s.httpClient.Post(buyerEndpoint, "application/json", bytes.NewReader(msg.Data))
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
		msg.Ack()
	})

	return err
}
