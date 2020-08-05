package models

import (
	"fmt"
	"strings"
)

// UserGroup comment
type UserGroup struct {
	UserID    string `json:"userid"`
	GroupName string `json:"groupname"`
}

// AllGroupsPerUserID comment
func AllGroupsPerUserID(userID string) ([]string, error) {
	groupNames := make([]string, 0)

	stmt := "SELECT GroupName FROM UserGroups WHERE UserID = $1"
	rows, err := db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupName string
		if err := rows.Scan(&groupName); err != nil {
			return nil, err
		}

		groupNames = append(groupNames, groupName)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groupNames, nil
}

// AllUserIDsPerGroup comment
func AllUserIDsPerGroup(groupName string) ([]string, error) {
	userIDs := make([]string, 0)

	stmt := "SELECT UserID FROM UserGroups WHERE GroupName = $1"
	rows, err := db.Query(stmt, groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		userIDs = append(userIDs, userID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

// CreateGroupsPerUserID comment
func CreateGroupsPerUserID(userID string, groupNames []string) (int64, error) {
	stmt := "INSERT INTO UserGroups (UserID, GroupName) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	var totalAffected int64
	for _, groupName := range groupNames {
		result, err := db.Exec(stmt, userID, groupName)
		if err != nil {
			return 0, err
		}
		inserted, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		totalAffected += inserted
	}

	inClause := strings.Join(groupNames, "','")
	stmt = fmt.Sprintf("DELETE FROM UserGroups WHERE UserID = $1 AND GroupName NOT IN ('%s')", inClause)
	result, err := db.Exec(stmt, userID)
	if err != nil {
		return 0, err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	totalAffected += deleted

	return totalAffected, nil
}

// UpdateUserIDsPerGroup comment
func UpdateUserIDsPerGroup(groupName string, userIDs []string) (int64, error) {
	stmt := "INSERT INTO UserGroups (UserID, GroupName) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	var totalAffected int64
	for _, userID := range userIDs {
		result, err := db.Exec(stmt, userID, groupName)
		if err != nil {
			return 0, err
		}
		inserted, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		totalAffected += inserted
	}

	return totalAffected, nil
}

// DeleteUserGroupsPerUserID comment
func DeleteUserGroupsPerUserID(userID string) (int64, error) {
	stmt := "DELETE FROM UserGroups WHERE UserID = $1"
	result, err := db.Exec(stmt, userID)
	if err != nil {
		return 0, err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}

// DeleteUserGroupsPerGroup comment
func DeleteUserGroupsPerGroup(groupName string) (int64, error) {
	stmt := "DELETE FROM UserGroups WHERE GroupName = $1"
	result, err := db.Exec(stmt, groupName)
	if err != nil {
		return 0, err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
