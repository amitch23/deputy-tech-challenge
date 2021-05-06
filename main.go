package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Userroles contains an array of userroles
type UserRoles struct {
	Userroles []UserRole `json:"userroles"`
}

//  UserRole represents a user role
type UserRole struct {
	Id     int
	Name   string
	Parent int
}

// Users represents an array of users
type Users struct {
	Users []User `json:"users"`
}

// User represents a user
type User struct {
	Id   int
	Name string
	Role int
}

// getDirectSubordinates finds the direct subordinates of a role (one layer down)
func getDirectSubordinates(userroles UserRoles, role int) []int {
	var children []int

	for _, userrole := range userroles.Userroles {
		// if a role has as its parent a given role, we know it's a direct child
		if role == userrole.Parent {
			children = append(children, userrole.Id)
		}
	}
	return children
}

// getAllSubordinatesRolesForRole finds all the subordinate roles for a given role
func getAllSubordinatesRolesForRole(userroles UserRoles, userrole int) []int {

	roleToSearch := userrole
	var childrenRoles []int
	for {
		// Get direct child(ren) roles of role
		directChildren := getDirectSubordinates(userroles, roleToSearch)
		if len(directChildren) > 0 {
			// if any children roles are found, search for each role's subsequent children recursively.
			for i, _ := range directChildren {
				roleToSearch = directChildren[i]
				getAllSubordinatesRolesForRole(userroles, roleToSearch)
			}
			// Add children to slice
			childrenRoles = append(childrenRoles, directChildren...)
		} else {
			// if there were no children found, we're done searching
			break
		}
	}
	return childrenRoles
}

// getSubordinates fetches all users with the subordinate roles
func getSubordinates(users Users, subordinateRoles []int) []User {
	// make a map of role: []users for faster lookup
	allUsers := users.Users
	roleUserMap := make(map[int][]User)
	for i := range allUsers {
		var usersToAppend []User
		users, found := roleUserMap[allUsers[i].Role]
		if !found {
			// set role as key to user slice if role not found in map
			usersToAppend = append(usersToAppend, allUsers[i])
			roleUserMap[allUsers[i].Role] = usersToAppend
		} else {
			// if role already in map, add users
			users = append(users, allUsers[i])
			roleUserMap[allUsers[i].Role] = users
		}
	}

	// fetch all subordinate users based on the role key from map above
	var subordinateUsers []User
	for _, v := range subordinateRoles {
		users, found := roleUserMap[v]
		if found {
			subordinateUsers = append(subordinateUsers, users...)
		} else {
			continue
		}
	}
	return subordinateUsers
}

// getRoleFromUserId returns a user's role from a given user Id
func getRoleFromUserId(Id int, users Users) (int, bool) {
	foundUser := false
	for _, u := range users.Users {
		if u.Id == Id {
			foundUser = true
			return u.Role, foundUser
		}
	}
	return 0, foundUser
}

func main() {
	// populate userrole structs from JSON
	userRolesJSON, err := os.Open("userroles.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer userRolesJSON.Close()

	rolesByteValue, _ := ioutil.ReadAll(userRolesJSON)
	var userroles UserRoles
	json.Unmarshal(rolesByteValue, &userroles)

	//  populate users structs from JSON
	usersJSON, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer usersJSON.Close()

	usersByteValue, _ := ioutil.ReadAll(usersJSON)
	var users Users
	json.Unmarshal(usersByteValue, &users)

	// Read in user ID flag from command line
	userId := flag.Int("userid", 1, "User Id must be an integer")
	flag.Parse()

	// get role of user from user id
	role, found := getRoleFromUserId(*userId, users)
	if found {
		// Get all subordinate roles for the role and set to a map of role: [childrole1, childrole2...]
		allChildRoles := getAllSubordinatesRolesForRole(userroles, role)
		roleToChildren := make(map[int][]int)
		if len(allChildRoles) > 0 {
			roleToChildren[role] = allChildRoles
		} else {
			fmt.Printf("User with role %v has no subordinates ", role)
			log.Fatal(err)
		}

		// Get all users who have the subordinate roles
		subordinateUsers := getSubordinates(users, roleToChildren[role])
		if len(subordinateUsers) > 0 {
			fmt.Println(subordinateUsers)
		} else {
			fmt.Printf("No subordinates found for subordinate roles.")
			log.Fatal(err)
		}
	} else {
		fmt.Printf("User not found with Id %v", *userId)
		log.Fatal(err)
	}
}
