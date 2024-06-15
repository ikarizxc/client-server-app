package repository

import (
	"dcsa-lab/internal/entities"
	repo "dcsa-lab/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetById(t *testing.T) {
	test_repo := repo.NewRepository()

	users := []*entities.User{
		{
			Username: "test1",
			Email:    "test1@gmail.com",
			Password: "test1",
			IsAdmin:  false,
		},
		{
			Username: "test2",
			Email:    "test2@gmail.com",
			Password: "test2",
			IsAdmin:  false,
		},
		{
			Username: "test3",
			Email:    "test3@gmail.com",
			Password: "test3",
			IsAdmin:  false,
		},
	}

	for _, v := range users {
		test_repo.Add(v)
	}

	for i := len(users) - 1; i >= 0; i-- {
		result, _ := test_repo.GetById(i + 2)
		expected := users[i]

		assert.Equal(t, expected, result, "Incorrect result.")
	}
}

func TestGetByUsername(t *testing.T) {
	test_repo := repo.NewRepository()
	userToCheck := &entities.User{
		Username: "test1",
		Email:    "test1@gmail.com",
		Password: "test1",
		IsAdmin:  false,
	}
	test_repo.Add(userToCheck)

	testCases := []struct {
		Username     string
		ExpectedUser *entities.User
	}{
		{
			Username:     "test1",
			ExpectedUser: userToCheck,
		},
		{
			Username:     "wrong",
			ExpectedUser: nil,
		},
	}

	for _, testCase := range testCases {
		resultUser, _ := test_repo.GetByUsername(testCase.Username)
		assert.Equal(t, testCase.ExpectedUser, resultUser, "Incorrect result.")
	}
}
