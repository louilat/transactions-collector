package tps

type User struct {
	Address string `json:"user"`
}

type UsersList []User

func (usersList UsersList) Contains(address string) bool {
	for _, user := range usersList {
		if user.Address == address {
			return true
		}
	}
	return false
}
