package whodunit

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// StatusTable represents the table that gets logged in the terminal.
type StatusTable struct {
	*tablewriter.Table
	assetType    AssetType
	statusFilter AssetStatus
}

// NewStatusTable returns a new instance of a status table.
func NewStatusTable(assetType AssetType, status AssetStatus) *StatusTable {
	if assetType != AssetTypeRecognition && status == AssetStatusInProcess {
		panic("You can only specify the in-process filter for recognitions")
	}

	return &StatusTable{
		Table:        tablewriter.NewWriter(os.Stdout),
		assetType:    assetType,
		statusFilter: status,
	}
}

// Log is a convenience method for looping through the episodes in all seasons
// and logging their status in the terminal.
func (st *StatusTable) Log() {
	totalCount := 0
	for season := 1; season <= SeasonCount; season++ {
		s := NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			panic("Could not get season episodes")
		}

		for _, ep := range s.AllEpisodes() {
			if st.AddRow(ep) {
				totalCount++
			}
		}
	}

	st.RenderTable(totalCount)
}

// RenderTable shows the table in the terminal.
func (st *StatusTable) RenderTable(totalCount int) {
	st.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	st.SetAlignment(tablewriter.ALIGN_LEFT)
	st.SetHeader([]string{"Season", "Episode", "Title", "Status"})
	st.SetFooter([]string{"", "", "Total", strconv.Itoa(totalCount)})
	st.Render()
}

// AddRow adds a new row to the table associated with the episode.
func (st *StatusTable) AddRow(ep *Episode) bool {
	status := ep.AssetStatus(st.assetType)
	if st.statusFilter != status && st.statusFilter != AssetStatusAny {
		return false
	}

	title := ep.DisplayTitle()
	if strings.Contains(title, "Helle") {
		title = "The Disappearance Of Helle..."
	}

	statusDisplay := st.statusDisplay(status)
	row := []string{
		strconv.Itoa(ep.SeasonNumber),
		strconv.Itoa(ep.EpisodeNumber),
		title,
		statusDisplay,
	}

	fgStyle := tablewriter.Normal
	var fgColor int

	if status == AssetStatusComplete {
		fgColor = tablewriter.FgGreenColor
	} else if status == AssetStatusPending {
		fgColor = tablewriter.FgYellowColor
	} else if status == AssetStatusInProcess {
		fgColor = tablewriter.FgCyanColor
	} else {
		fgStyle = tablewriter.Bold
		fgColor = tablewriter.FgRedColor
	}

	st.Rich(row, []tablewriter.Colors{
		{fgStyle, fgColor},
		{fgStyle, fgColor},
		{fgStyle, fgColor},
		{fgStyle, fgColor},
	})
	return true
}

func (st *StatusTable) statusDisplay(status AssetStatus) string {
	switch status {
	case AssetStatusPending:
		return "Pending"
	case AssetStatusInProcess:
		return "In Process"
	case AssetStatusMissing:
		return "Missing"
	case AssetStatusComplete:
		return "Complete"
	}
	return "Unknown"
}
