package main

import (
	"bufio"
	"dcsa-lab/internal/models"
	"dcsa-lab/internal/requests"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	reader = bufio.NewReader(os.Stdin)
)

func signUp() {
	signUpData := models.UserData{}

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	signUpData.Username = strings.TrimSpace(username)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	signUpData.Email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	signUpData.Password = strings.TrimSpace(password)

	status, mes, err := requests.SignUp(signUpData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nResponce:\n\t- status: %#v\n\t- body:\n%s\n", status, mes)
}

func login() {
	loginData := models.LoginData{}

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	loginData.Email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	loginData.Password = strings.TrimSpace(password)

	status, mes, err := requests.Login(loginData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nResponce:\n\t- status: %#v\n\t- body:\n%s\n", status, mes)
}

func getAllUsers() {
	status, mes, err := requests.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nResponce:\n\t- status: %#v\n\t- body:\n%s\n", status, mes)
}

func getUserById() {
	fmt.Print("Enter user id: ")
	input, _ := reader.ReadString('\n')
	idstr := strings.TrimSpace(input)

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println(err)
	}

	status, mes, err := requests.GetUserById(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nResponce:\n\t- status: %#v\n\t- body:\n%s\n", status, mes)
}

func deleteUser() {
	fmt.Print("Enter user id: ")
	input, _ := reader.ReadString('\n')
	idstr := strings.TrimSpace(input)

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println(err)
	}

	status, mes, err := requests.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nResponce:\n\t- status: %#v\n\t- body:\n%s\n", status, mes)
}

func updateUser() {
	fmt.Print("Enter user id: ")
	input, _ := reader.ReadString('\n')
	idstr := strings.TrimSpace(input)

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println(err)
	}

	updateData := models.UserData{}

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	updateData.Username = strings.TrimSpace(username)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	updateData.Email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	updateData.Password = strings.TrimSpace(password)

	status, mes, err := requests.UpdateUser(id, updateData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nResponce:\n\t- status: %#v\n\t- body:\n%s\n", status, mes)
}

func main() {
	variants := []string{
		"\n--- --- MENU --- ---\n",
		"1. Sign up\n",
		"2. Login\n",
		"3. Get all users\n",
		"4. Get user by id\n",
		"5. Delete user\n",
		"6. Update user\n",
		"7. Exit\n",
		"\nChoose option: ",
	}

	for {
		for _, l := range variants {
			fmt.Print(l)
		}

		var input int = 0
		fmt.Scanf("%d", &input)

		switch input {
		case 1:
			signUp()
		case 2:
			login()
		case 3:
			getAllUsers()
		case 4:
			getUserById()
		case 5:
			deleteUser()
		case 6:
			updateUser()
		case 7:
			os.Exit(1)
		}
	}
}
