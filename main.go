package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/nomad/api"
)

func main() {
	client, _ := api.NewClient(&api.Config{Address: "http://127.0.0.1:4646"})

	eventsClient := client.EventStream()

	topics := map[api.Topic][]string{
		// api.TopicAll: {"*"},
		// api.TopicDeployment: {"*"},
	}

	ctx := context.Background()
	eventCh, err := eventsClient.Stream(ctx, topics, 0, &api.QueryOptions{})
	if err != nil {
		fmt.Printf("[stream] received error %s", err)
		os.Exit(1)
	}

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-eventCh:
			if event.Err != nil {
				fmt.Printf("[event] received error %s", err)
				break
			}

			// Ignore heartbeats used to keep connection alive
			if event.IsHeartbeat() {
				continue
			}

			for _, e := range event.Events {
				fmt.Printf("[%s] %s %d %v\n", e.Topic, e.Type, e.Index, e.Payload)
				// deployment, err := e.Deployment()
				// if err != nil {
				// 	fmt.Printf("received error %s", err)
				// 	continue
				// }

				// fmt.Printf("Deployment update %s", deployment.Status)
			}
		}
	}
}
