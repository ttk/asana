package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"

	"github.com/ttk/asana/api"
	"github.com/ttk/asana/utils"
)

const (
	CacheDuration = "5m"
)

func Tasks(c *cli.Context) {
	if c.Bool("no-cache") {
		fromAPI(c, false)
	} else {
		if utils.Older(CacheDuration, utils.CacheFile()) || c.Bool("refresh") {
			fromAPI(c, true)
		} else {
			txt, err := ioutil.ReadFile(utils.CacheFile())
			if err == nil {
				var tasks []api.Task_t
				lines := regexp.MustCompile("\n").Split(string(txt), -1)
				for _, line := range lines {
					if len(line) < 1 {
						continue
					}
					// Parse cache line format: index:gid:due_on:name
					parts := regexp.MustCompile(":").Split(line, 4)
					if len(parts) >= 4 {
						tasks = append(tasks, api.Task_t{
							Gid:    parts[1],
							Due_on: parts[2],
							Name:   parts[3],
						})
					}
				}
				outputTasks(c, tasks)
			} else {
				fromAPI(c, true)
			}
		}
	}
}

func fromAPI(c *cli.Context, saveCache bool) {
	tasks := api.Tasks(url.Values{}, false)
	if saveCache {
		cache(tasks)
	}
	outputTasks(c, tasks)
}

func outputTasks(c *cli.Context, tasks []api.Task_t) {
	if c.Bool("json") {
		outputTasksJSON(tasks)
	} else {
		renderTasksTable(tasks)
	}
}

func outputTasksJSON(tasks []api.Task_t) {
	type TaskJSON struct {
		Gid    string `json:"gid"`
		Due_on string `json:"due_on"`
		Name   string `json:"name"`
	}

	var jsonTasks []TaskJSON
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, TaskJSON{
			Gid:    t.Gid,
			Due_on: t.Due_on,
			Name:   t.Name,
		})
	}

	output, err := json.MarshalIndent(jsonTasks, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(output))
}

func renderTasksTable(tasks []api.Task_t) {
	// Get terminal width (default to 80 if detection fails)
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth <= 0 {
		termWidth = 80
	}

	// Calculate width for Task Name column
	// Leave room for: # column (5), Gid column (15), Due Date column (12), borders and padding (~8)
	taskNameWidth := termWidth - 40
	if taskNameWidth < 20 {
		taskNameWidth = 20 // Minimum width for readability
	}

	// Configure per-column widths
	colWidths := tw.NewMapper[int, int]()
	colWidths.Set(0, 5)             // # column
	colWidths.Set(1, 15)            // Gid column
	colWidths.Set(2, 12)            // Due Date column
	colWidths.Set(3, taskNameWidth) // Task Name column

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap: tw.WrapNormal,
				},
				ColMaxWidths: tw.CellWidth{
					PerColumn: colWidths,
				},
			},
		}),
	)

	table.Header("#", "Gid", "Due Date", "Task Name")

	for i, t := range tasks {
		table.Append(strconv.Itoa(i), t.Gid, t.Due_on, t.Name)
	}

	table.Render()
}

func cache(tasks []api.Task_t) {
	f, _ := os.Create(utils.CacheFile())
	defer f.Close()
	for i, t := range tasks {
		f.WriteString(strconv.Itoa(i) + ":")
		f.WriteString(t.Gid + ":")
		f.WriteString(t.Due_on + ":")
		f.WriteString(t.Name + "\n")
	}
}
