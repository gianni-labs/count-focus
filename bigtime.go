package main

import "strings"

const largeTimeGlyphHeight = 6

var largeTimeGlyphs = map[rune][]string{
	'0': {
		" ██████╗ ",
		"██╔═████╗",
		"██║██╔██║",
		"████╔╝██║",
		"╚██████╔╝",
		" ╚═════╝ ",
	},
	'1': {
		"  ██╗ ",
		" ███║ ",
		" ╚██║ ",
		"  ██║ ",
		"  ██║ ",
		"  ╚═╝ ",
	},
	'2': {
		"██████╗ ",
		"╚════██╗",
		" █████╔╝",
		"██╔═══╝ ",
		"███████╗",
		"╚══════╝",
	},
	'3': {
		"██████╗ ",
		"╚════██╗",
		" █████╔╝",
		" ╚═══██╗",
		"██████╔╝",
		"╚═════╝ ",
	},
	'4': {
		"██╗  ██╗",
		"██║  ██║",
		"███████║",
		"╚════██║",
		"     ██║",
		"     ╚═╝",
	},
	'5': {
		"███████╗",
		"██╔════╝",
		"██████╗ ",
		"╚════██╗",
		"██████╔╝",
		"╚═════╝ ",
	},
	'6': {
		" ██████╗ ",
		"██╔════╝ ",
		"██████╗  ",
		"██╔══██╗ ",
		"╚█████╔╝ ",
		" ╚════╝  ",
	},
	'7': {
		"███████╗",
		"╚════██║",
		"    ██╔╝",
		"   ██╔╝ ",
		"   ██║  ",
		"   ╚═╝  ",
	},
	'8': {
		" █████╗ ",
		"██╔══██╗",
		"╚█████╔╝",
		"██╔══██╗",
		"╚█████╔╝",
		" ╚════╝ ",
	},
	'9': {
		" █████╗ ",
		"██╔══██╗",
		"╚██████║",
		" ╚═══██║",
		" █████╔╝",
		" ╚════╝ ",
	},
	':': {
		"   ",
		"██╗",
		"╚═╝",
		"██╗",
		"╚═╝",
		"   ",
	},
}

func renderLargeTime(value string) string {
	lines := make([]string, largeTimeGlyphHeight)

	for _, ch := range value {
		glyph, ok := largeTimeGlyphs[ch]
		if !ok {
			continue
		}

		for i := range largeTimeGlyphHeight {
			if lines[i] != "" {
				lines[i] += "  "
			}
			lines[i] += glyph[i]
		}
	}

	return strings.Join(lines, "\n")
}

func largeTimeWidth(value string) int {
	lines := strings.Split(renderLargeTime(value), "\n")
	width := 0
	for _, line := range lines {
		if len([]rune(line)) > width {
			width = len([]rune(line))
		}
	}
	return width
}
