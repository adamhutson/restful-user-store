package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adamhutson/mux-postgres-api/models"
	"github.com/gorilla/mux"
)

// UserWithGroups comment
type UserWithGroups struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	UserID    string   `json:"userid"`
	Groups    []string `json:"groups"`
}

//GroupWithUsers comment
type GroupWithUsers struct {
	GroupName string   `json:"name"`
	UserIDs   []string `json:"userids"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", handleUsers)
	r.HandleFunc("/users/{userid}", handleUser)

	r.HandleFunc("/groups", handleGroups)
	r.HandleFunc("/groups/{groupname}", handleGroup)

	http.ListenAndServe(":8080", r)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsersWithGroups(w, r)
	case http.MethodPost:
		createUserWithGroups(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUserWithGroups(w, r)
	case http.MethodPut:
		updateUserWithGroups(w, r)
	case http.MethodDelete:
		deleteUserWithGroups(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGroups(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createGroup(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getGroupWithUsers(w, r)
	case http.MethodPut:
		updateGroupWithUsers(w, r)
	case http.MethodDelete:
		deleteGroupWithUsers(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func getUsersWithGroups(w http.ResponseWriter, r *http.Request) {
	usersWithGroups := make([]UserWithGroups, 0)

	users, err := models.AllUsers()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, user := range users {
		userWithGroups := UserWithGroups{}
		userWithGroups.FirstName = user.FirstName
		userWithGroups.LastName = user.LastName
		userWithGroups.UserID = user.UserID

		groups, err := models.AllGroupsPerUserID(user.UserID)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		userWithGroups.Groups = groups

		usersWithGroups = append(usersWithGroups, userWithGroups)
	}
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, usersWithGroups)
}

func createUserWithGroups(w http.ResponseWriter, r *http.Request) {
	userWithGroups := UserWithGroups{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.Unmarshal(b, &userWithGroups); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid user data")
		return
	}
	defer r.Body.Close()

	fmt.Println(userWithGroups)

	user := models.User{}
	user.FirstName = userWithGroups.FirstName
	user.LastName = userWithGroups.LastName
	user.UserID = userWithGroups.UserID

	createdUsers, err := models.CreateUser(user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if createdUsers != 1 {
		errorResponse(w, http.StatusConflict, "User already exists")
		return
	}

	_, err = models.CreateGroups(userWithGroups.Groups)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = models.CreateGroupsPerUserID(user.UserID, userWithGroups.Groups)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusCreated, nil)
}

func getUserWithGroups(w http.ResponseWriter, r *http.Request) {
	userWithGroups := UserWithGroups{}

	vars := mux.Vars(r)
	userid := vars["userid"]

	user, err := models.OneUser(userid)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user.UserID == "" {
		errorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	groups, err := models.AllGroupsPerUserID(userid)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	userWithGroups.FirstName = user.FirstName
	userWithGroups.LastName = user.LastName
	userWithGroups.UserID = user.UserID
	userWithGroups.Groups = groups

	jsonResponse(w, http.StatusOK, userWithGroups)
}

func updateUserWithGroups(w http.ResponseWriter, r *http.Request) {
	userWithGroups := UserWithGroups{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.Unmarshal(b, &userWithGroups); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid user data")
		return
	}
	defer r.Body.Close()

	user := models.User{}
	user.FirstName = userWithGroups.FirstName
	user.LastName = userWithGroups.LastName
	user.UserID = userWithGroups.UserID

	updatedUsers, err := models.UpdateUser(user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if updatedUsers != 1 {
		errorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	_, err = models.CreateGroups(userWithGroups.Groups)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = models.CreateGroupsPerUserID(userWithGroups.UserID, userWithGroups.Groups)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusAccepted, nil)
}

func deleteUserWithGroups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]

	deleted, err := models.DeleteUser(userid)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if deleted != 1 {
		errorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	_, err = models.DeleteUserGroupsPerUserID(userid)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusAccepted, nil)
}

func getGroupWithUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupname := vars["groupname"]

	groupExists, err := models.GroupExists(groupname)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if groupExists == false {
		errorResponse(w, http.StatusNotFound, "Group not found")
		return
	}

	userIDs, err := models.AllUserIDsPerGroup(groupname)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	groupWithUsers := GroupWithUsers{}
	groupWithUsers.GroupName = groupname
	groupWithUsers.UserIDs = userIDs

	jsonResponse(w, http.StatusOK, groupWithUsers)
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var groupName string
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.Unmarshal(b, &groupName); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid group data")
		return
	}
	defer r.Body.Close()

	created, err := models.CreateGroup(groupName)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if created != 1 {
		errorResponse(w, http.StatusConflict, "Group already exists")
		return
	}

	jsonResponse(w, http.StatusCreated, nil)
}

func updateGroupWithUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupname := vars["groupname"]

	groupExists, err := models.GroupExists(groupname)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if groupExists == false {
		errorResponse(w, http.StatusNotFound, "Group not found")
		return
	}

	type users struct {
		UserIDs []string `json:"userids"`
	}
	userIDs := users{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.Unmarshal(b, &userIDs); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid group membership data")
		return
	}
	defer r.Body.Close()

	_, err = models.UpdateUserIDsPerGroup(groupname, userIDs.UserIDs)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusAccepted, nil)
}

func deleteGroupWithUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupname := vars["groupname"]

	_, err := models.DeleteUserGroupsPerGroup(groupname)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = models.DeleteGroup(groupname)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusAccepted, nil)
}

func errorResponse(w http.ResponseWriter, status int, err string) {
	jsonResponse(w, status, map[string]string{"error": err})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
