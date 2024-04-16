package model

import (
	"math"
	"strings"
)

func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func GetListWidth(w int) int {
	if w <= 0 {
		return 50
	}
	return int(math.Ceil(0.25 * float64(w)))
}

// stateViewNames maps stateView values to their string names.
var stateViewNames = map[stateView]string{
	ListView:  "List",
	ValueView: "Value",
	HelpView:  "Help",
	OpAmpView: "OpAMP",
}

// method that returns the name of the stateView
func (v stateView) Name() string {
	// check if the name exists in the map
	if name, ok := stateViewNames[v]; ok {
		return name
	}

	// return a default value
	return "Unknown"
}
