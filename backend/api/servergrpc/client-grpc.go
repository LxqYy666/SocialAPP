package servergrpc

import (
	"Server/models"
	"Server/protos"
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client protos.NotificationGrpcServiceClient
}

func NewClient() (*Client, error) {
	conn, err := grpc.NewClient(":5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := protos.NewNotificationGrpcServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendGrpcNotification(ctx context.Context, req *protos.NotificationGrpcRequest) error {
	_, err := c.client.SendGrpcNotification(ctx, req)
	return err
}

func SendNotification(notification models.Notification) error {
	client, err := NewClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	req := &protos.NotificationGrpcRequest{
		XId:      notification.ID.Hex(),
		Details:  notification.Details,
		Mainuid:  notification.MainUID,
		Targetid: notification.TargetUID,
		Isreaded: notification.IsReaded,
		CreatedAt: &timestamp.Timestamp{
			Seconds: notification.CreatedAt.Unix(),
		},
		User: &protos.Usergrpc{
			Name:   notification.Name,
			Avatar: notification.Avartar,
		},
	}
	return client.SendGrpcNotification(ctx, req)
}
