package servegrpc

import (
	"context"
	"net"
	pb "realTimeNotification/protos"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type notificationServer struct {
	pb.UnimplementedNotificationGrpcServiceServer
	wsmutex       *sync.Mutex
	wsConnections map[string]*websocket.Conn
}

type Notification struct {
	ID        string    `json:"id"`
	Details   string    `json:"details"`
	MainUID   string    `json:"mainuid"`
	TargetID  string    `json:"targetid"`
	IsReaded  bool      `json:"isreaded"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (s *notificationServer) SendGrpcNotification(ctx context.Context, req *pb.NotificationGrpcRequest) (*empty.Empty, error) {
	s.wsmutex.Lock()
	defer s.wsmutex.Unlock()

	if conn, ok := s.wsConnections[req.Mainuid]; ok {
		notification := Notification{
			ID:        req.XId,
			Details:   req.Details,
			TargetID:  req.Targetid,
			IsReaded:  req.Isreaded,
			CreatedAt: time.Unix(req.CreatedAt.Seconds, 0),
			User: User{
				Name:   req.User.Name,
				Avatar: req.User.Avatar,
			},
		}
		err := conn.WriteJSON(notification)
		if err != nil {
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

func StartGrpcServer(wsmutex *sync.Mutex, wsConnections map[string]*websocket.Conn) error {
	listener, err := net.Listen("tcp", ":5002")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationGrpcServiceServer(grpcServer, &notificationServer{
		wsmutex:       wsmutex,
		wsConnections: wsConnections,
	})
	return grpcServer.Serve(listener)
}
