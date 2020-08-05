package models

//Group comment
type Group struct {
	GroupName string `json:"name"`
}

//CreateGroup comment
func CreateGroup(groupName string) (int64, error) {
	stmt := "INSERT INTO groups (GroupName) VALUES ($1) ON CONFLICT DO NOTHING"
	result, err := db.Exec(stmt, groupName)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// CreateGroups comment
func CreateGroups(groups []string) (int64, error) {
	var totalAffected int64
	for _, group := range groups {
		affected, err := CreateGroup(group)
		if err != nil {
			return 0, err
		}
		totalAffected += affected
	}
	return totalAffected, nil
}

// GroupExists comment
func GroupExists(groupName string) (bool, error) {
	stmt := "SELECT COUNT(*) FROM groups WHERE GroupName = $1"
	rows, err := db.Query(stmt, groupName)
	if err != nil {
		return false, err
	}

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}
	if count == 0 {
		return false, err
	}

	return true, nil
}

// DeleteGroup comment
func DeleteGroup(groupName string) (int64, error) {
	stmt := "DELETE FROM Groups WHERE GroupName = $1"
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
