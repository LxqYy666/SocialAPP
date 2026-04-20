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
