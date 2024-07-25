/**
 * @Author: steven
 * @Description:
 * @File: queue_google_pubsub
 * @Date: 25/12/23 08.21
 */

package queue

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
)

type GooglePubSubManager interface {
	MustConnect() GooglePubSubManager
	Connect() error
	Dispose()

	Publish(ctx context.Context, topic string, packet Packet[any]) error
	PublishWithResult(ctx context.Context, topic string, packet Packet[any]) (*pubsub.PublishResult, error)
	PublishAndWaitAck(ctx context.Context, topic string, packet Packet[any]) (string, error)

	Subscribe(ctx context.Context, topic string, onReceivedMessage func(context.Context, *pubsub.Message)) error

	AckWithResult(ctx context.Context, id string, ackResult *pubsub.AckResult, result func(err error, status pubsub.AcknowledgeStatus))
}

type googlePubSubManager struct {
	projectId string
	scope     string
	client    *pubsub.Client
	opts      []option.ClientOption
	tm        telemetry.Manager
	tc        trace.Tracer
}

func (g *googlePubSubManager) spanName(v string) string {
	return rules.WhenTrueRE1(len(g.scope) > 0, func() string {
		return g.scope + "." + v
	}, func() string {
		return v
	})
}

func publish[T any](ctx context.Context, topic *pubsub.Topic, data T) (*pubsub.PublishResult, error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return topic.Publish(ctx, &pubsub.Message{
		Data: byteData,
	}), nil
}

func (g *googlePubSubManager) PublishWithResult(ctx context.Context, topic string, packet Packet[any]) (*pubsub.PublishResult, error) {
	return wrapTelemetryTuple2[*pubsub.PublishResult, error](
		ctx, g.tc, g.spanName("pubsub.publish_msg_wr"), []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindProducer),
			trace.WithAttributes(attribute.String("key", packet.Key)),
			trace.WithAttributes(attribute.String("topic", topic)),
		}, func(newCtx context.Context) (*pubsub.PublishResult, error) {
			return publish(newCtx, g.client.Topic(topic), packet)
		})
}

func (g *googlePubSubManager) Publish(ctx context.Context, topic string, packet Packet[any]) error {
	return wrapTelemetryTuple1(ctx, g.tc, g.spanName("pubsub.publish_msg"), []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attribute.String("key", packet.Key)),
		trace.WithAttributes(attribute.String("topic", topic)),
	}, func(newCtx context.Context) error {
		_, err := publish(newCtx, g.client.Topic(topic), packet)
		return err
	})
}

func (g *googlePubSubManager) PublishAndWaitAck(ctx context.Context, topic string, packet Packet[any]) (string, error) {
	return wrapTelemetryTuple2[string, error](ctx, g.tc, g.spanName("pubsub.publish_msg_wack"), []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attribute.String("key", packet.Key)),
		trace.WithAttributes(attribute.String("topic", topic)),
	}, func(newCtx context.Context) (string, error) {
		rs, err := publish(newCtx, g.client.Topic(topic), packet)
		if err != nil {
			return "", err
		}
		return rs.Get(newCtx)
	})
}

func (g *googlePubSubManager) Subscribe(ctx context.Context, topic string, onReceivedMessage func(context.Context, *pubsub.Message)) error {
	sub := g.client.Subscription(topic)
	return sub.Receive(ctx, func(innerCtx context.Context, message *pubsub.Message) {
		wrapTelemetry(innerCtx, g.tc, g.spanName("pubsub.receive"), []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(attribute.String("topic", topic)),
		}, func(newCtx context.Context) {
			onReceivedMessage(newCtx, message)
		})
	})
}

func (g *googlePubSubManager) AckWithResult(ctx context.Context, id string, ackResult *pubsub.AckResult, result func(err error, status pubsub.AcknowledgeStatus)) {
	status, err := ackResult.Get(ctx)
	if err != nil {
		result(err, status)
		return
	}
	switch status {
	case pubsub.AcknowledgeStatusSuccess:
	case pubsub.AcknowledgeStatusInvalidAckID:
		err = errors.New("message failed to ack with response of Invalid. ID: " + id)
	case pubsub.AcknowledgeStatusPermissionDenied:
		err = errors.New("message failed to ack with response of Permission Denied. ID: " + id)
	case pubsub.AcknowledgeStatusFailedPrecondition:
		err = errors.New("message failed to ack with response of Failed Precondition. ID: " + id)
	case pubsub.AcknowledgeStatusOther:
		err = errors.New("message failed to ack with response of Other. ID: " + id)
	default:
	}
	result(err, status)
	return
}

func (g *googlePubSubManager) MustConnect() GooglePubSubManager {
	if err := g.Connect(); err != nil {
		panic(err)
	}
	return g
}

func (g *googlePubSubManager) Connect() error {
	c, err := pubsub.NewClient(context.Background(), g.projectId, g.opts...)
	if err != nil {
		return err
	}
	g.client = c
	return nil
}

func (g *googlePubSubManager) Dispose() {
	if g.client == nil {
		return
	}
	_ = g.client.Close()
}

func NewGooglePubSub(projectId, scope string, tm telemetry.Manager, opts ...option.ClientOption) GooglePubSubManager {
	return &googlePubSubManager{
		projectId: projectId,
		scope:     scope,
		tm:        tm,
		tc:        tm.Tracer(),
		opts:      opts,
	}
}
