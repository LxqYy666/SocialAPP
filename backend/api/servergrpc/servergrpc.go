package servergrpc

import (
	"Server/database"
	"Server/models"
	pb "Server/protos"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Server struct {
	pb.UnimplementedRealtimeChatServiceServer
}

func (s *Server) GetUserFollowingFollowers(ctx context.Context, req *pb.UserID) (*pb.UserIDsListResponse, error) {
	userSchema := database.DB.Collection("users")

	userID := req.GetUserid()
	if userID == "" {
		return nil, fmt.Errorf("userid is empty")
	}

	var user models.User
	uid, _ := bson.ObjectIDFromHex(userID)

	err := userSchema.FindOne(ctx, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	friendsMap := make(map[string]bool)
	for _, id := range user.Following {
		friendsMap[id] = true
	}
	for _, id := range user.Followers {
		friendsMap[id] = true
	}

	var friends []string
	for key := range friendsMap {
		friends = append(friends, key)
	}

	return &pb.UserIDsListResponse{
		UserIDsLists: []*pb.UserIDsList{
			{
				UserIdsList: friends,
			},
		},
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	if req.GetSender() == "" || req.GetReceiver() == "" {
		return nil, fmt.Errorf("sender and receiver is needed")
	}

	sendID, err := bson.ObjectIDFromHex(req.GetSender())
	receID, err := bson.ObjectIDFromHex(req.GetReceiver())
	if err != nil {
		return nil, fmt.Errorf("error when trans")
	}

	//检查用户是否存在
	userSchame := database.DB.Collection("users")
	var sender, receiver models.User

	err = userSchame.FindOne(ctx, bson.M{"_id": sendID}).Decode(&sender)
	err = userSchame.FindOne(ctx, bson.M{"_id": receID}).Decode(&receiver)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	//保存消息
	message := models.Message{
		Sender:   req.GetSender(),
		Receiver: req.GetReceiver(),
		Content:  req.GetContent(),
	}

	_, err = database.DB.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("fail to save message")
	}

	//更新未读消息
	unReadMsgSchema := database.DB.Collection("unReadedMsgs")
	existingRecord := bson.M{}
	err = unReadMsgSchema.FindOneAndUpdate(
		ctx,
		bson.M{"mainUserId": req.GetReceiver(), "otherUserId": req.GetSender()},
		bson.M{"$inc": bson.M{"numOfUnReadedMsg": 1}, "$set": bson.M{"isReaded": false}},
	).Decode(&existingRecord)
	if err != nil {
		_, err := unReadMsgSchema.InsertOne(ctx, bson.M{"mainUserId": req.GetReceiver(), "otherUserId": req.GetSender(), "numOfUnReadedMsg": 1, "isReaded": false})
		if err != nil {
			return nil, fmt.Errorf("fail to update unread message")
		}
	}

	return &pb.MessageResponse{
		Message: req.GetContent(),
	}, nil
}
