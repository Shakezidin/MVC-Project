package observability

import (
	"context"
	"encoding/json"

	cloudpubsub "cloud.google.com/go/pubsub"
)

type PubSubPublisher struct {
	client *cloudpubsub.Client
	topic  *cloudpubsub.Topic
	ctx    context.Context
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

	return &PubSubPublisher{
		client: client,
		topic:  client.Topic(topicID),
		ctx:    ctx,
	}, nil
}

func (p *PubSubPublisher) Publish(
	entry LogEntry,
) error {

	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	result := p.topic.Publish(
		p.ctx,
		&cloudpubsub.Message{
			Data: bytes,
		},
	)

	_, err = result.Get(p.ctx)

	return err
}
