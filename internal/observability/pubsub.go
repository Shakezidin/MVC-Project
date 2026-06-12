package observability

import (
	"context"
	"encoding/json"
	"fmt"

	cloudpubsub "cloud.google.com/go/pubsub/v2"
)

type PubSubPublisher struct {
	client *cloudpubsub.Client
	pubsub *cloudpubsub.Publisher
	ctx    context.Context
	topic  string
}

func NewPubSubPublisher(
	projectID string,
	topicID string,
) (*PubSubPublisher, error) {

	ctx := context.Background()

	client, err := cloudpubsub.NewClient(
		ctx,
		projectID,
	)

	if err != nil {
		return nil, err
	}

	topic := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
	publisher := client.Publisher(topic)

	return &PubSubPublisher{
		client: client,
		pubsub: publisher,
		ctx:    ctx,
		topic:  topic,
	}, nil
}

func (p *PubSubPublisher) Publish(
	entry LogEntry,
) error {

	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	result := p.pubsub.Publish(
		p.ctx,
		&cloudpubsub.Message{
			Data: bytes,
		},
	)

	_, err = result.Get(p.ctx)

	return err
}
