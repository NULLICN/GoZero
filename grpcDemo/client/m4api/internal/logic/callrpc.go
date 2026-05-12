package logic

import (
	"context"
	"fmt"
	"time"

	"client/greeterclient"
	"client/m4api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

func callRpc(ctx context.Context, client greeterclient.Greeter, name string) (*types.SayHelloRes, error) {
	const (
		maxRetries  = 3
		callTimeout = 3 * time.Second
	)

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * 100 * time.Millisecond
			logx.Infof("[gateway] retry %d/%d after %v", attempt, maxRetries-1, backoff)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		callCtx, cancel := context.WithTimeout(ctx, callTimeout)
		resp, err := client.SayHello(callCtx, &greeterclient.HelloReq{Name: name})
		cancel()

		if err != nil {
			lastErr = err
			logx.Errorf("[gateway] call failed (attempt %d/%d): %v", attempt+1, maxRetries, err)
			continue
		}

		logx.Infof("[gateway] rpc response: %s", resp.Message)
		return &types.SayHelloRes{Message: resp.Message}, nil
	}

	return nil, fmt.Errorf("all %d attempts exhausted: %w", maxRetries, lastErr)
}
