/*
    CDR Importer
    Copyright (C) 2026 Elliot Michael Keavney

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program. If not, see https://github.com/yet-another-pbx/cdr-importer/blob/main/LICENSE
*/

package main

import (
	"fmt"
)

//----------------------------------------------------------------------------------------------------

// Constant for cdrimporter Env file
const fileCDRImporterEnv string = "/etc/yap/cdrimporter.env"

// Clear Screen Function
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// American National Standards Institute (ANSI) reset colour code
const resetColour = "\033[0m"

// American National Standards Institute (ANSI) text colour code
const textBoldBlack = "\033[1;30m"
const textBoldWhite = "\033[1;37m"

// American National Standards Institute (ANSI) background colour codes
const bgRed = "\033[41m"
const bgGreen = "\033[42m"
const bgYellow = "\033[43m"
const bgBlue = "\033[46m"
const bgMagenta = "\033[45m"

func main() {
	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("     " + bgRed + textBoldWhite + " ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃    ██████╗ ██████╗  ██████╗      ██╗ ███╗   ███╗ ██████╗   ██████╗  ██████╗  ████████╗ ███████╗ ██████╗    ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ██╔════╝ ██╔══██╗ ██╔══██╗     ██║ ████╗ ████║ ██╔══██╗ ██╔═══██╗ ██╔══██╗ ╚══██╔══╝ ██╔════╝ ██╔══██╗   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ██║      ██║  ██║ ██████╔╝     ██║ ██╔████╔██║ ██████╔╝ ██║   ██║ ██████╔╝    ██║    █████╗   ██████╔╝   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ██║      ██║  ██║ ██╔══██╗     ██║ ██║╚██╔╝██║ ██╔═══╝  ██║   ██║ ██╔══██╗    ██║    ██╔══╝   ██╔══██╗   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ╚██████╗ ██████╔╝ ██║  ██║     ██║ ██║ ╚═╝ ██║ ██║      ╚██████╔╝ ██║  ██║    ██║    ███████╗ ██║  ██║   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃    ╚═════╝ ╚═════╝  ╚═╝  ╚═╝     ╚═╝ ╚═╝     ╚═╝ ╚═╝       ╚═════╝  ╚═╝  ╚═╝    ╚═╝    ╚══════╝ ╚═╝  ╚═╝   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃         Source code for CDR Importer available at https://github.com/yet-another-pbx/cdr-importer          ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                For use with YAP (Yet Another PBX) - https://github.com/yet-another-pbx/yap                 ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + "                                                     " + textBoldWhite + resetColour + bgRed + textBoldWhite + "                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + textBoldWhite + "   Type \"exit\" or \"quit\" to terminate CDR Importer   " + resetColour + bgRed + textBoldWhite + "                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + textBoldWhite + "                                                     " + resetColour + bgRed + textBoldWhite + "                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                            ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ " + resetColour)

	fmt.Print(textBoldBlack)
	fmt.Println("")
	fmt.Println("          Select an option [1-9]:\n")
	fmt.Println("          [1] List all supplier IDs and names\n")
	fmt.Println("          [2] List a CDR from a particular month for a supplier\n")
	fmt.Println("          [3] List all call rates for a supplier\n")
	fmt.Println("          [4] Edit a call rate for a supplier\n")
	fmt.Println("          [5] Add a new supplier\n")
	fmt.Println("          [6] Delete an existing supplier\n")
	fmt.Println("          [7] Import a new CDR into an existing supplier\n")
	fmt.Println("          [8] Delete a previsouly imported CDR from a supplier\n")
	fmt.Println("          [9] Insert a CDR from a particular month for a supplier into YAP (Yet Another PBX)\n")
	fmt.Println(resetColour)
}

// Contributor(s):
// Elliot Michael Keavney
