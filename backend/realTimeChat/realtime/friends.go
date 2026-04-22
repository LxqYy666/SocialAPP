package realtime

import "realTimeChat/servergrpc"

func GetFriends(userId string) ([]string, error) {
	userFriends, err := servergrpc.GetFollowingFollowersClient(userId)
	if err != nil {
		return nil, err
	}

	var friends []string
	for _, userIDsList := range userFriends {
		friends = append(friends, userIDsList.UserIdsList...)
	}

	return friends, nil
}
