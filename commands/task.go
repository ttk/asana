package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/ttk/asana/api"
)

func Task(c *cli.Context) {
	taskId := api.FindTaskId(c.Args().First(), true)

	verbose := c.Bool("verbose")
	withComments := c.Bool("with-comments")

	shouldFetchStories := verbose || withComments
	t, stories := api.Task(taskId, shouldFetchStories)
	attachments := api.Attachments(taskId)

	if withComments && !verbose && stories != nil {
		var commentStories []api.Story_t
		for _, s := range stories {
			if s.Type == "comment" {
				commentStories = append(commentStories, s)
			}
		}
		stories = commentStories
	}

	if c.Bool("json") {
		output := map[string]interface{}{
			"task": t,
		}
		if stories != nil {
			output["stories"] = stories
		}
		if len(attachments) > 0 {
			output["attachments"] = attachments
		}
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
		return
	}

	fmt.Println("Task:")
	fmt.Printf("    Name: %s\n", t.Name)
	fmt.Printf("    Due On: %s\n", t.Due_on)
	fmt.Printf("    Gid: %s\n", t.Gid)
	fmt.Printf("    URL: https://app.asana.com/0/0/%s\n", t.Gid)

	showTags(t.Tags)
	showCustomFields(t.CustomFields)
	showAttachments(attachments)

	fmt.Printf("\nNotes:\n%s\n", indentText(t.Notes, "    "))

	if stories != nil {
		showComments(stories)
		showSystem(stories)
	}
}

func showComments(stories []api.Story_t) {
	var comments []api.Story_t
	for _, s := range stories {
		if s.Type == "comment" {
			comments = append(comments, s)
		}
	}

	if len(comments) > 0 {
		fmt.Printf("Comments:\n")
		for _, s := range comments {
			fmt.Printf("  - text: %s\n", indentTextExceptFirst(s.Text, "          "))
			fmt.Printf("    created_at: %s\n", s.Created_at)
			fmt.Printf("    created_by: %s\n", s.Created_by.Name)
		}
	}
}

func showSystem(stories []api.Story_t) {
	var systemStories []api.Story_t
	for _, s := range stories {
		if s.Type != "comment" {
			systemStories = append(systemStories, s)
		}
	}

	if len(systemStories) > 0 {
		fmt.Printf("\nSystem:\n")
		for _, s := range systemStories {
			fmt.Printf("  - text: %s\n", indentTextExceptFirst(s.Text, "          "))
			fmt.Printf("    created_at: %s\n", s.Created_at)
			fmt.Printf("    created_by: %s\n", s.Created_by.Name)
		}
	}
}

func indentText(text string, indent string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}

func indentTextExceptFirst(text string, indent string) string {
	lines := strings.Split(text, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}
	return strings.Join(lines, "\n")
}

func showTags(tags []api.Base) {
	if len(tags) > 0 {
		fmt.Print("Tags: ")
		for i, tag := range tags {
			print(tag.Name)
			if len(tags) != 1 && i != (len(tags)-1) {
				print(", ")
			}
		}
		println("")
	}
}

func showCustomFields(fields []api.CustomField_t) {
	if len(fields) > 0 {
		fmt.Println("\nCustom Fields:")
		for _, field := range fields {
			if field.DisplayValue != "" {
				fmt.Printf("    %s: %s\n", field.Name, field.DisplayValue)
			}
		}
	}
}

func showAttachments(attachments []api.Attachment_t) {
	if len(attachments) > 0 {
		fmt.Printf("\nAttachments (%d):\n", len(attachments))
		for i, att := range attachments {
			fmt.Printf("    [%d] %s (GID: %s)\n", i, att.Name, att.Gid)
		}
	}
}
