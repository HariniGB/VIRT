package controllers

import (
	"fmt"
	"os/exec"
	"log"
)


func createUser(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("invalid request. user name and password can not be empty")
	}

	cmd := exec.Command("openstack", "user", "create", "--password", password, username)
	log.Printf("Running create user command and waiting for it to finish...")
	err := cmd.Run()

	if err != nil {
		log.Println("User creation failed with error: ", err)
		return fmt.Errorf("unable to create user")
	}

	return nil
}

func createProject(username string) error {
	if username == "" {
		return fmt.Errorf("invalid request. user name can not be empty")
	}

	cmd := exec.Command("openstack", "project", "create", "--domain", "default",
		"--description", fmt.Sprintf("%s project", username), username)
	log.Printf("Running create project command and waiting for it to finish...")
	err := cmd.Run()

	if err != nil {
		log.Println("Project creation failed with error: ", err)
		return fmt.Errorf("unable to create project")
	}

	return nil
}

func assignRole(username string) error {
	if username == "" {
		return fmt.Errorf("invalid request. user name can not be empty")
	}

	cmd := exec.Command("openstack", "role", "add", "--project", username,
		"--user", username, "user")
	log.Printf("Running add roll command and waiting for it to finish...")
	err := cmd.Run()

	if err != nil {
		log.Println("Project creation failed with error: ", err)
		return fmt.Errorf("unable to create project")
	}

	return nil
}



func CreateAccount(username, password string) error{
	err := createUser(username, password)
	if err != nil {
		return err
	}

	err = createProject(username)
	if err != nil {
		return err
	}

	err = assignRole(username)
	if err != nil {
		return err
	}

	return nil
}