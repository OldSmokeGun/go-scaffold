package repository

import (
	"fmt"
	"strconv"
	"strings"
)

func GetPolicyUser(userID int64) string {
	return fmt.Sprintf("user_%d", userID)
}

func FromPolicyUser(user string) (int64, error) {
	us := strings.Trim(user, "user_")
	userID, err := strconv.ParseInt(us, 10, 64)
	if err != nil {
		return 0, handleError(err)
	}
	return userID, nil
}

func GetPolicyRole(roleID int64) string {
	return fmt.Sprintf("role_%d", roleID)
}

func FromPolicyRole(role string) (int64, error) {
	rs := strings.Trim(role, "role_")
	roleID, err := strconv.ParseInt(rs, 10, 64)
	if err != nil {
		return 0, handleError(err)
	}
	return roleID, nil
}
