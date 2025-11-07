package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type FileCommand struct {
	Action           string `json:"action"`
	FilePattern      string `json:"file_pattern"`
	ReplaceExtension string `json:"replace_extension"`
	Destination      string `json:"destination"`
}

func executeCommand(cmd FileCommand) {
	switch cmd.Action {
	case "rename":
		files, _ := filepath.Glob(cmd.FilePattern)
		for _, file := range files {
			newName := file[:len(file)-len(filepath.Ext(file))] + "." + cmd.ReplaceExtension
			fmt.Println("Renaming", file, "→", newName)
			os.Rename(file, newName)
		}
	case "move":
		files, _ := filepath.Glob(cmd.FilePattern)
		for _, file := range files {
			newPath := filepath.Join(cmd.Destination, filepath.Base(file))
			fmt.Println("Moving", file, "→", newPath)
			os.Rename(file, newPath)
		}
	default:
		fmt.Println("Unknown action:", cmd.Action)
	}
}

func main() {
	r := gin.Default()
	client := resty.New()

	r.POST("/ai-command", func(c *gin.Context) {
		var userInput struct {
			Prompt string `json:"prompt"`
		}
		if err := c.BindJSON(&userInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{"prompt": userInput.Prompt}).
			Post("http://localhost:5001/interpret")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
			return
		}

		var cmd FileCommand
		if err := json.Unmarshal(resp.Body(), &cmd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse AI output"})
			return
		}

		executeCommand(cmd)
		c.JSON(http.StatusOK, cmd)
	})

	r.Run(":8080")
}
