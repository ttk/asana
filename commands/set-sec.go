package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/ttk/asana/api"
)

func SetSection(c *cli.Context) {
	if c.NArg() < 2 {
		log.Fatal("fatal: set-sec requires 2 arguments: [task index] [section name]")
	}

	taskIndex := c.Args().Get(0)
	sectionName := c.Args().Get(1)
	projectKeyword := c.String("project")

	taskId := api.FindTaskId(taskIndex, false)
	task, _ := api.Task(taskId, false)

	var projectGid string

	if len(task.Projects) == 0 {
		log.Fatal("fatal: Task is not part of any project")
	}

	if len(task.Projects) > 1 && projectKeyword == "" {
		log.Fatal("fatal: Task is part of multiple projects. Please specify a project using -p or --project")
	}

	if projectKeyword == "" {
		projectGid = task.Projects[0].Gid
	} else {
		projectGid = matchProject(task.Projects, projectKeyword)
	}

	sections := api.GetSectionsForProject(projectGid)

	if len(sections) == 0 {
		log.Fatal("fatal: Project has no sections")
	}

	sectionGid := matchSection(sections, sectionName)

	api.AddTaskToSection(sectionGid, taskId)
	fmt.Printf("Moved task '%s' to section '%s'\n", task.Name, sectionName)
}

func matchProject(projects []api.Base, keyword string) string {
	keyword = strings.ToLower(keyword)
	var matches []api.Base

	for _, project := range projects {
		if strings.Contains(strings.ToLower(project.Name), keyword) {
			matches = append(matches, project)
		}
	}

	if len(matches) == 0 {
		log.Fatalf("fatal: No project found matching '%s'", keyword)
	}

	if len(matches) > 1 {
		log.Fatalf("fatal: Multiple projects found matching '%s'. Please be more specific. Matches: %v",
			keyword, getProjectNames(matches))
	}

	return matches[0].Gid
}

func matchSection(sections []api.Section_t, keyword string) string {
	keyword = strings.ToLower(keyword)
	var matches []api.Section_t

	for _, section := range sections {
		if strings.Contains(strings.ToLower(section.Name), keyword) {
			matches = append(matches, section)
		}
	}

	if len(matches) == 0 {
		log.Fatalf("fatal: No section found matching '%s'", keyword)
	}

	if len(matches) > 1 {
		log.Fatalf("fatal: Multiple sections found matching '%s'. Please be more specific. Matches: %v",
			keyword, getSectionNames(matches))
	}

	return matches[0].Gid
}

func getProjectNames(projects []api.Base) []string {
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	return names
}

func getSectionNames(sections []api.Section_t) []string {
	names := make([]string, len(sections))
	for i, s := range sections {
		names[i] = s.Name
	}
	return names
}