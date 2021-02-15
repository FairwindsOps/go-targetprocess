package targetprocess

import (
	"encoding/json"
	"fmt"
	"os"
)

// Example of client.Get() for godoc
func ExampleClient_GetPriority() {
	tpClient, err := NewClient("exampleaccount", "superSecretToken")
	if err != nil {
		fmt.Println("Failed to create tp client:", err)
		os.Exit(1)
	}
	priority, err := tpClient.GetPriority("Must",
		"UserStory")
	if err != nil {
		fmt.Println("Failed to get users:", err)
		os.Exit(1)
	}
	jsonBytes, _ := json.Marshal(priority)
	fmt.Print(string(jsonBytes))
}

// Example of client.Get() for godoc
func ExampleUserStory_SetPriority() {
	tpClient, err := NewClient("exampleaccount", "superSecretToken")
	if err != nil {
		fmt.Println("Failed to create tp client:", err)
		os.Exit(1)
	}
	us, err := NewUserStory(tpClient, "Testing", "Description", "Project")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	err = us.SetPriority("Must")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	_, link, err := us.Create()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Printf("Card created here: %s\n", link)
}
