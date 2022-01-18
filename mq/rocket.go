package mq

import (
	"context"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketMQOptions struct {
	NameserverAddrs []string
	GroupName       string
	ProduceRetries  int
	ConsumeMode     consumer.MessageModel
}
type MQSending struct {
	Topic   string
	Content []byte
}
type MQSubscribe struct {
	Topic    string
	Tag      string
	Callback func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}
type MQProducer struct {
	producer rocketmq.Producer
}

func (p *MQProducer) Start() error {
	if p == nil || p.producer == nil {
		return nil
	}

	return p.producer.Start()
}
func (p *MQProducer) Shutdown() error {
	if p == nil || p.producer == nil {
		return nil
	}

	return p.producer.Shutdown()
}
func (p *MQProducer) Sync(ctx context.Context, msg MQSending) (bool, error) {
	if p == nil || p.producer == nil {
		return false, nil
	}

	result, err := p.producer.SendSync(context.Background(), &primitive.Message{
		Topic: msg.Topic,
		Body:  msg.Content,
	})
	if err != nil {
		return false, err
	}
	if result.Status == primitive.SendOK {
		return true, nil
	}
	log.Printf("producer.SendSync erro %+v\n", result.Status)
	return false, nil
}

type MQConsumer struct {
	consumer rocketmq.PushConsumer
}

func (p *MQConsumer) Subscribe(subscribe MQSubscribe) error {
	if p == nil || p.consumer == nil {
		return nil
	}

	return p.consumer.Subscribe(subscribe.Topic, consumer.MessageSelector{}, subscribe.Callback)
}
func (p *MQConsumer) Unsubscribe(topic string) error {
	if p == nil || p.consumer == nil {
		return nil
	}

	return p.consumer.Unsubscribe(topic)
}
func (p *MQConsumer) Start() error {
	if p == nil || p.consumer == nil {
		return nil
	}
	return p.consumer.Start()
}
func (p *MQConsumer) Shutdown() error {
	if p == nil || p.consumer == nil {
		return nil
	}

	return p.consumer.Shutdown()
}

func NewRocketMQProducer(opt RocketMQOptions) (*MQProducer, error) {
	nsAddr, err := primitive.NewNamesrvAddr(opt.NameserverAddrs...)
	if err != nil {
		return nil, err
	}
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nsAddr),
		//producer.WithNsResolver(primitive.NewPassthroughResolver(endPoint)),
		producer.WithRetry(opt.ProduceRetries),
		producer.WithGroupName(opt.GroupName),
	)
	if err != nil {
		return nil, err
	}
	return &MQProducer{producer: p}, nil
}

func NewRocketMQConsumer(opt RocketMQOptions) (*MQConsumer, error) {
	nsAddr, err := primitive.NewNamesrvAddr(opt.NameserverAddrs...)
	if err != nil {
		return nil, err
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(nsAddr),
		consumer.WithConsumerModel(opt.ConsumeMode),
		consumer.WithGroupName(opt.GroupName),
	)
	if err != nil {
		return nil, err
	}
	return &MQConsumer{consumer: c}, nil
}
