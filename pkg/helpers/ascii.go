package helpers

import "bytes"

func MoleAsciiArt() string {

	var buffer bytes.Buffer

	buffer.WriteString("                  ___             \n")
	buffer.WriteString("                 /\\_ \\            \n")
	buffer.WriteString("  ___ ___     ___\\//\\ \\      __   \n")
	buffer.WriteString("/  __  __ \\  / __ \\\\ \\ \\   / __ \\ \n")
	buffer.WriteString("/\\ \\/\\ \\/\\ \\/\\ \\_\\ \\\\_\\ \\_/\\  __/ \n")
	buffer.WriteString("\\ \\_\\ \\_\\ \\_\\ \\____//\\____\\ \\____\\\n")
	buffer.WriteString(" \\/_/\\/_/\\/_/\\/___/ \\/____/\\/____/\n")

	return buffer.String()
}
