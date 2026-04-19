package realtime

func GetFriends(userId string) <-chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)
		switch userId {
		case "1":
			ch <- []string{"2", "3", "4"}
		case "2":
			ch <- []string{"1", "3", "4"}
		case "3":
			ch <- []string{"1", "2", "4"}
		case "4":
			ch <- []string{"1"}
		default:
			ch <- []string{}
		}

	}()

	return ch
}
