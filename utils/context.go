/**
 * @Author: steven
 * @Description:
 * @File: context
 * @Date: 23/12/23 08.32
 */

package utils

import "context"

func MergeContextWithCancel(c1, c2 context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c1)
	go func() {
		select {
		case <-ctx.Done(): // don't leak go-routine on clean gRPC run
		case <-c2.Done():
			cancel() // c2 canceled, so cancel ctx
		}
	}()
	return ctx, cancel
}
