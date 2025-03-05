package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai/aimodelservice"
)

func TestSimulateAutoOrder(t *testing.T) {
	InitClient("ai", client.WithHostPorts("192.168.31.223:8889"))
	type args struct {
		ctx         context.Context
		Req         *ai.SimulateAutoOrderRequest
		callOptions []callopt.Option
	}
	tests := []struct {
		name       string
		args       args
		wantStream aimodelservice.AiModelService_SimulateAutoOrderClient
		wantErr    bool
	}{
		{
			name: "echo",
			args: args{
				ctx: context.Background(),
				Req: &ai.SimulateAutoOrderRequest{
					UserMessage: "hello",
				},
			},
			wantStream: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStream, _ := SimulateAutoOrder(tt.args.ctx, tt.args.Req, tt.args.callOptions...)
			for {
				res, err := gotStream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println(res)
			}
		})
	}
}
