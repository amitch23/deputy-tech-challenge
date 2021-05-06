package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var userroles = []UserRole{
	{
		Id:     1,
		Name:   "System Administrator",
		Parent: 0,
	},
	{
		Id:     2,
		Name:   "Location Manager",
		Parent: 1,
	},
	{
		Id:     3,
		Name:   "Supervisor",
		Parent: 2,
	},
	{
		Id:     4,
		Name:   "Employee",
		Parent: 3,
	},
	{
		Id:     5,
		Name:   "Trainer",
		Parent: 3,
	},
}

var users = []User{
	{
		Id:   1,
		Name: "Adam Admin",
		Role: 1,
	},
	{
		Id:   5,
		Name: "Mary Manager",
		Role: 2,
	},
	{
		Id:   4,
		Name: "Sam Supervisor",
		Role: 3,
	},
	{
		Id:   2,
		Name: "Emily Employee",
		Role: 4,
	},
	{
		Id:   3,
		Name: "Maya Employee",
		Role: 4,
	},
	{
		Id:   6,
		Name: "Steve Trainer",
		Role: 5,
	},
}

var userrolesSlice = UserRoles{
	Userroles: userroles,
}

var usersSlice = Users{
	Users: users,
}

type roleFromIdTestCase struct {
	input  int
	output int
	found  bool
}

type childRoleTestCase struct {
	input  int
	output []int
}

type userTestCase struct {
	input  []int
	output []User
}

var directChildrenTestCases = []childRoleTestCase{
	{0, []int{1}},
	{1, []int{2}},
	{2, []int{3}},
	{3, []int{4, 5}},
	// role 4 has no subordinates so should return an empty slice
	{4, nil},
}

var allChildrenTestCases = []childRoleTestCase{
	{0, []int{1, 2, 3, 4, 5}},
	{1, []int{2, 3, 4, 5}},
	{2, []int{3, 4, 5}},
	{3, []int{4, 5}},
	// role 4 has no subordinates, so should return an empty slice
	{4, nil},
}

var getSubordinatesTestCase = []userTestCase{
	{[]int{4, 5, 6}, []User{
		{
			Id:   2,
			Name: "Emily Employee",
			Role: 4,
		},
		{
			Id:   3,
			Name: "Maya Employee",
			Role: 4,
		},
		{
			Id:   6,
			Name: "Steve Trainer",
			Role: 5,
		},
	}},
	{[]int{1, 2, 3, 4, 5}, users},
	// Role 7 does not exist so should return an empty slice
	{[]int{7}, nil},
}

var getUserRoleTestCase = []roleFromIdTestCase{
	{1, 1, true},
	{2, 4, true},
	{3, 4, true},
	{4, 3, true},
	{5, 2, true},
	{6, 5, true},
	// User ID 117 doesn't exist, so we expect no role returned
	{117, 0, false},
}

func TestGetDirectSubordinates(t *testing.T) {
	// Check that the input results in the expected output as defined in directChildrenTestCases
	for _, test := range directChildrenTestCases {
		directSubordinates := getDirectSubordinates(userrolesSlice, test.input)
		assert.Equal(t, directSubordinates, test.output, "Given role did not return the expected direct subordinate roles.")
	}
}

func TestGetAllSubordinatesForRole(t *testing.T) {
	// Check that the input results in the expected output as defined in allChildrenTestCases
	for _, test := range allChildrenTestCases {
		allSubordinates := getAllSubordinatesRolesForRole(userrolesSlice, test.input)
		assert.Equal(t, allSubordinates, test.output, "Given role did not return the expected subordinate roles.")
	}
}

func TestGetSubordinates(t *testing.T) {
	// Check that the input results in the expected output as defined in getSubordinatesTestCase
	for _, test := range getSubordinatesTestCase {
		subordinates := getSubordinates(usersSlice, test.input)
		assert.Equal(t, subordinates, test.output, "Subordinate roles did not match expected users returned.")
	}
}

func TestGetRoleFromUserId(t *testing.T) {
	// Check that the input results in the expected output as defined in getUserRoleTestCase
	for _, test := range getUserRoleTestCase {
		role, found := getRoleFromUserId(test.input, usersSlice)
		assert.Equal(t, role, test.output, "User Id did not return the expected user role")
		assert.Equal(t, found, test.found, "User Id did not return expected boolean for whether user was found.")
	}
}
