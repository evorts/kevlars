/**
 * @Author: steven
 * @Description:
 * @File: queue
 * @Date: 25/12/23 08.21
 */

package queue

import (
	"context"
)

type Manager[T any, POPT PublishOptionType] interface {
	MustConnect() Manager[T, POPT]
	Connect() error
	Dispose()

	Publish(ctx context.Context, topic string, packet Packet[any], opts ...PublishOption[POPT]) error
	Subscribe(ctx context.Context, topic string, onReceivedMessage func(context.Context, *T)) error

	Ack(ctx context.Context, id string, ackResult any, opts ...ActOption[any])
	NAck(ctx context.Context, id string, ackResult any, opts ...ActOption[any])
}

type PublishOptionType interface {
	googlePubsubPubOption | natsPubsubPubOption
}

type ActOptionType interface {
	googlePubsubPubOption | natsPubsubPubOption
}
