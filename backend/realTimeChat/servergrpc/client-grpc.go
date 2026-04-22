package servergrpc

import (
	"context"
	"fmt"
	"realTimeChat/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetFollowingFollowersClient(id string) ([]*protos.UserIDsList, error) {
	conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("fail to connect to grpc server")
	}
	defer conn.Close()

	client := protos.NewRealtimeChatServiceClient(conn)
	ctx := context.Background()
	req := &protos.UserID{Userid: id}
	resp, err := client.GetUserFollowingFollowers(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetUserIDsLists(), nil
}

func SendMessageClient(sender, receiver, content string) error {

	conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("fail to connect to grpc server")
	}
	defer conn.Close()

	client := protos.NewRealtimeChatServiceClient(conn)
	ctx := context.Background()
	req := &protos.MessageRequest{
		Sender:   sender,
		Receiver: receiver,
		Content:  content,
	}
	_, err = client.SendMessage(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
