/**
 * @Author: steven
 * @Description:
 * @File: quueu_pubsub
 * @Date: 25/12/23 09.15
 */

package queue

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
)

type googlePubsubManager struct {
	projectId string
	scope     string
	client    *pubsub.Client
	opts      []option.ClientOption
	tm        telemetry.Manager
	tc        trace.Tracer
}

type googlePubsubPubOption struct {
	waitAck            bool
	waitActCallback    func(serverId string, err error)
	waitResult         bool
	waitResultCallback func(*pubsub.PublishResult, error)
}

func GooglePubsubPubOptWithWaitAck(v bool, result func(string, error)) PublishOption[googlePubsubPubOption] {
	return publishOptionFunc[googlePubsubPubOption](func(t *googlePubsubPubOption) {
		t.waitAck = v
		t.waitActCallback = result
	})
}

func GooglePubsubPubOptWithWaitResult(v bool, result func(*pubsub.PublishResult, error)) PublishOption[googlePubsubPubOption] {
	return publishOptionFunc[googlePubsubPubOption](func(t *googlePubsubPubOption) {
		t.waitResult = v
		t.waitResultCallback = result
	})
}

func newGooglePubsubPubOption(opts ...PublishOption[googlePubsubPubOption]) *googlePubsubPubOption {
	po := &googlePubsubPubOption{}
	for _, opt := range opts {
		opt.apply(po)
	}
	return po
}

type googlePubsubActOption struct {
	ackResult *pubsub.AckResult
	result    func(err error, status pubsub.AcknowledgeStatus)
}

func (g *googlePubsubManager) spanName(v string) string {
	return rules.WhenTrueRE1(len(g.scope) > 0, func() string {
		return g.scope + "." + v
	}, func() string {
		return v
	})
}

func googlePubsubPublish[T any](ctx context.Context, topic *pubsub.Topic, data T) (*pubsub.PublishResult, error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return topic.Publish(ctx, &pubsub.Message{
		Data: byteData,
	}), nil
}

func (g *googlePubsubManager) MustConnect() Manager[pubsub.Message, googlePubsubPubOption] {
	if err := g.Connect(); err != nil {
		panic(err)
	}
	return g
}

func (g *googlePubsubManager) Connect() error {
	c, err := pubsub.NewClient(context.Background(), g.projectId, g.opts...)
	if err != nil {
		return err
	}
	g.client = c
	return nil
}

func (g *googlePubsubManager) Dispose() {
	if g.client == nil {
		return
	}
	_ = g.client.Close()
}

func (g *googlePubsubManager) Publish(ctx context.Context, topic string, packet Packet[any], opts ...PublishOption[googlePubsubPubOption]) error {
	pubOpts := newGooglePubsubPubOption(opts...)
	return wrapTelemetryTuple1(ctx, g.tc, g.spanName("gpm.publish"), []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attribute.String("key", packet.Key)),
		trace.WithAttributes(attribute.String("topic", topic)),
	}, func(newCtx context.Context) error {
		rs, err := publish(newCtx, g.client.Topic(topic), packet)
		if pubOpts.waitResult {
			pubOpts.waitResultCallback(rs, err)
		}
		if err == nil && rs != nil && pubOpts.waitAck {
			var serverId string
			serverId, err = rs.Get(newCtx)
			pubOpts.waitActCallback(serverId, err)
		}
		return err
	})
}

func (g *googlePubsubManager) Subscribe(ctx context.Context, topic string, onReceivedMessage func(context.Context, *pubsub.Message)) error {
	sub := g.client.Subscription(topic)
	return sub.Receive(ctx, func(innerCtx context.Context, message *pubsub.Message) {
		wrapTelemetry(innerCtx, g.tc, g.spanName("gpm.receive"), []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(attribute.String("topic", topic)),
		}, func(newCtx context.Context) {
			onReceivedMessage(newCtx, message)
		})
	})
}

func (g *googlePubsubManager) Ack(ctx context.Context, id string, ackResult any, opts ...ActOption[any]) {
	//TODO implement me
	panic("implement me")
}

func (g *googlePubsubManager) NAck(ctx context.Context, id string, ackResult any, opts ...ActOption[any]) {
	//TODO implement me
	panic("implement me")
}

func NewPubSub(projectId, scope string, tm telemetry.Manager, opts ...option.ClientOption) Manager[pubsub.Message, googlePubsubPubOption] {
	return &googlePubsubManager{
		projectId: projectId,
		scope:     scope,
		tm:        tm,
		tc:        tm.Tracer(),
		opts:      opts,
	}
}
