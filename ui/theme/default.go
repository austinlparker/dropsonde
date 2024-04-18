package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	PrimaryBG   lipgloss.AdaptiveColor
	SecondaryBG lipgloss.AdaptiveColor
	PrimaryFG   lipgloss.AdaptiveColor
	SecondaryFG lipgloss.AdaptiveColor
	Accent      lipgloss.AdaptiveColor
	Primary     lipgloss.AdaptiveColor
	Secondary   lipgloss.AdaptiveColor
}

var (
	OTelBlue       = "#425CC7"
	OTelYellow     = "#F5A800"
	OTelOffsetBlue = "#9393C8"
	OTelLightWhite = "#EFEDFF"
	OTelSpotBrown  = "#BFA975"
)

func DefaultTheme() Theme {
	return Theme{
		PrimaryBG:   lipgloss.AdaptiveColor{OTelBlue, OTelBlue},
		SecondaryBG: lipgloss.AdaptiveColor{OTelOffsetBlue, OTelOffsetBlue},
		PrimaryFG:   lipgloss.AdaptiveColor{OTelLightWhite, OTelLightWhite},
		SecondaryFG: lipgloss.AdaptiveColor{OTelLightWhite, OTelLightWhite},
		Accent:      lipgloss.AdaptiveColor{OTelYellow, OTelYellow},
		Primary:     lipgloss.AdaptiveColor{OTelBlue, OTelBlue},
		Secondary:   lipgloss.AdaptiveColor{OTelYellow, OTelYellow},
	}
}
