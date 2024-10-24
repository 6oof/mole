package helpers

import "bytes"

func MoleAsciiArt() string {

	var buffer bytes.Buffer

	buffer.WriteString("• ▌ ▄ ·.       ▄▄▌  ▄▄▄ .\n")
	buffer.WriteString("·██ ▐███▪▪     ██•  ▀▄.▀·\n")
	buffer.WriteString("▐█ ▌▐▌▐█· ▄█▀▄ ██▪  ▐▀▀▪▄\n")
	buffer.WriteString("██ ██▌▐█▌▐█▌.▐▌▐█▌▐▌▐█▄▄▌\n")
	buffer.WriteString("▀▀  █▪▀▀▀ ▀█▄▀▪.▀▀▀  ▀▀▀ \n")

	return buffer.String()
}
