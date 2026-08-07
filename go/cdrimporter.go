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
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
	"slices"
	"strconv"
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

//----------------------------------------------------------------------------------------------------

// Function to terminate the running program if a user types "exit" or "quit"
func exit(value string) {
	if value == "exit" || value == "quit" {
		fmt.Print(resetColour)
		clearScreen()
		os.Exit(0)
	}
}

//----------------------------------------------------------------------------------------------------

// Option 1 function
func option1() {

}

//----------------------------------------------------------------------------------------------------

func main() {
	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("     " + bgRed + textBoldWhite + " ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃    ██████╗ ██████╗  ██████╗     ██╗ ███╗   ███╗ ██████╗   ██████╗  ██████╗  ████████╗ ███████╗ ██████╗    ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ██╔════╝ ██╔══██╗ ██╔══██╗    ██║ ████╗ ████║ ██╔══██╗ ██╔═══██╗ ██╔══██╗ ╚══██╔══╝ ██╔════╝ ██╔══██╗   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ██║      ██║  ██║ ██████╔╝    ██║ ██╔████╔██║ ██████╔╝ ██║   ██║ ██████╔╝    ██║    █████╗   ██████╔╝   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ██║      ██║  ██║ ██╔══██╗    ██║ ██║╚██╔╝██║ ██╔═══╝  ██║   ██║ ██╔══██╗    ██║    ██╔══╝   ██╔══██╗   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃   ╚██████╗ ██████╔╝ ██║  ██║    ██║ ██║ ╚═╝ ██║ ██║      ╚██████╔╝ ██║  ██║    ██║    ███████╗ ██║  ██║   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃    ╚═════╝ ╚═════╝  ╚═╝  ╚═╝    ╚═╝ ╚═╝     ╚═╝ ╚═╝       ╚═════╝  ╚═╝  ╚═╝    ╚═╝    ╚══════╝ ╚═╝  ╚═╝   ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃         Source code for CDR Importer available at https://github.com/yet-another-pbx/cdr-importer         ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                For use with YAP (Yet Another PBX) - https://github.com/yet-another-pbx/yap                ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + "                                                     " + textBoldWhite + resetColour + bgRed + textBoldWhite + "                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + textBoldWhite + "   Type \"exit\" or \"quit\" to terminate CDR Importer   " + resetColour + bgRed + textBoldWhite + "                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + textBoldWhite + "                                                     " + resetColour + bgRed + textBoldWhite + "                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ " + resetColour)

	fmt.Print(textBoldBlack)
	fmt.Println("")
	fmt.Println("          [1] List all supplier IDs and names\n")
	fmt.Println("          [2] List a CDR from a particular month for a supplier\n")
	fmt.Println("          [3] List all call rates for a supplier\n")
	fmt.Println("          [4] Edit a call rate for a supplier\n")
	fmt.Println("          [5] Add a new supplier\n")
	fmt.Println("          [6] Delete an existing supplier\n")
	fmt.Println("          [7] Import a new CDR into an existing supplier\n")
	fmt.Println("          [8] Delete a previsouly imported CDR from a supplier\n")
	fmt.Println("          [9] Insert a CDR from a particular month for a supplier into YAP (Yet Another PBX)\n")
	fmt.Println("")
	fmt.Print("          Select an option [1-9]:\n")
	var option string
	fmt.Scan(&option)
	exit(option)
	fmt.Println(resetColour)

	// Conditional statment to determine what happens when an option is input

	// Values allowed for option
	var optionList = []string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "exit", "Exit", "EXIT", "quit", "Quit", "QUIT"}
	validOption := slices.Contains(optionList, option)

	// Get the values from inside the CDR Importer configuration file
	err := godotenv.Load(fileCDRImporterEnv)
	if err != nil {
		panic("Error loading " + fileCDRImporterEnv + " file for database details")
	}

	// Get the database connection details
	dbUsername := os.Getenv("dbUsername")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")
	dbAddress := os.Getenv("dbAddress")
	dbPort := os.Getenv("dbPort")
	dbTransport := os.Getenv("dbTransport")
	dbTLS := os.Getenv("dbTLS")

	// Validate the dbAddress is an IP address
	validateDbAddress := validator.New()
	validateDbAddressErr := validateDbAddress.Var(dbAddress, "required,ip_addr")

	// Validate the dbPortInt is a number
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic("DATABASE PORT MUST BE A NUMBER IN " + fileCDRImporterEnv)
	}

	// Values allowed for dbTransport Variable
	var transportList = []string{"tcp", "udp"}
	validDbTransport := slices.Contains(transportList, dbTransport)

	// Values allowed for dbTls variable
	var dbTLSList = []string{"false", "true"}
	validDbTLS := slices.Contains(dbTLSList, dbTLS)

	if dbUsername == "" {
		panic("DATABASE USERNAME CANNOT BE EMPTY IN " + fileCDRImporterEnv)
	} else if dbPassword == "" {
		panic("DATABASE PASSOWRD CANNOT BE EMPTY IN " + fileCDRImporterEnv)
	} else if dbName == "" {
		panic("DATABASE NAME CANNOT BE EMPTY IN " + fileCDRImporterEnv)
	} else if validateDbAddressErr != nil && dbAddress != "localhost" {
		panic("DATABASE ADDRESS MUST BE A VALID INTERENT PROTOCOL (IP) ADDRESS OR localhost IN " + fileCDRImporterEnv)
	} else if dbPortInt <= 0 || dbPortInt >= 65536 {
		panic("DATABASE PORT MUST BE IN THE NUMBER RANGE 1-65535 IN " + fileCDRImporterEnv)
	} else if dbTransport == "" {
		panic("DATABASE TRANSPORT OPTION CANNOT BE EMPTY IN " + fileCDRImporterEnv)
	} else if validDbTransport == false {
		panic("DATABASE TRANSPORT OPTION MUST BE udp OR tcp IN " + fileCDRImporterEnv)
	} else if dbTLS == "" {
		panic("DATABASE TLS OPTION CANNOT BE EMPTY IN " + fileCDRImporterEnv)
	} else if validDbTLS == false {
		panic("DATABASE TRANSPORT OPTION MUST BE false OR true IN " + fileCDRImporterEnv)
	}

}

// Contributor(s):
// Elliot Michael Keavney
