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
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------

// Constant for cdrimporter Env file
const fileCDRImporterEnv string = "/etc/cdr-importer/cdr-importer.env"

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

// Database struct
type databaseFunctionParameter struct {
	connection *sql.DB
	database   string
}

//----------------------------------------------------------------------------------------------------

// Function to terminate the running program if a user types "exit" or "quit"
func exit(value string) {
	if value == "exit" || value == "Exit" || value == "EXIT" || value == "quit" || value == "Quit" || value == "QUIT" {
		messageBox("Program exited successfully", bgGreen)
		os.Exit(0)
	}
}

// Function to draw box
func messageBox(message string, bgColour string) {
	clearScreen()
	fmt.Print(resetColour)
	topBottomSymbol := strings.Repeat(" ⊛", (len(message)/2)+6)
	inbetweenSpace := strings.Repeat(" ", len(message)+8)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("     " + bgColour + textBoldWhite + topBottomSymbol + " " + resetColour)
	fmt.Println("     " + bgColour + textBoldWhite + " ⊛" + inbetweenSpace + "⊛ " + resetColour)
	fmt.Println("     " + bgColour + textBoldWhite + " ⊛    " + message + "    ⊛ " + resetColour)
	fmt.Println("     " + bgColour + textBoldWhite + " ⊛" + inbetweenSpace + "⊛ " + resetColour)
	fmt.Println("     " + bgColour + textBoldWhite + topBottomSymbol + " " + resetColour)
	fmt.Println("")
	fmt.Println("")
}

func returnToMainMenu() {

	fmt.Print(textBoldBlack)
	fmt.Print("     Press the enter/return key to continue ")
	fmt.Print(resetColour)
	var enter string
	fmt.Scanln(&enter)
	if enter == "" || enter != "" {
		main()
	}
}

// Function to format ISO datetime format
func formatDateTime(dateTime string) string {

	const (
		iso    = "2006-01-02 15:04:05"
		layout = "02/01/2006 15:04:05 PM"
	)

	parse, _ := time.Parse(iso, dateTime)
	return string(parse.Format(layout))
}

//----------------------------------------------------------------------------------------------------

// Option 1 function
// List all supplier IDs and names
func option0(dbDetail databaseFunctionParameter) {

	var (
		id            string
		name          string
		dateTimeAdded string
	)

	option0SQL, err := dbDetail.connection.Query(`SELECT
						        id,
	                                                name,
	                                                date_time_added
	                                              FROM
	                                                cdr_importer.voip_carrier;`)

	// Error
	if err != nil {
		panic(err)
	}

	clearScreen()
	fmt.Println("")
	fmt.Println(textBoldBlack)
	fmt.Println("          ╔═════════════════════════════════════════╗")
	fmt.Println("          ║ [0] List all VoIP carrier IDs and names ║")
	fmt.Println("     ╔════╩════════════╦════════════════════════════╩════════════════════════╦═════════════════════════╗")
	fmt.Println("     ║ VoIP Carrier ID ║                  VoIP Carrier Name                  ║    Date & Time Added    ║")

	for option0SQL.Next() {

		err = option0SQL.Scan(
			&id,
			&name,
			&dateTimeAdded,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Println("     ╠═════════════════╬═════════════════════════════════════════════════════╬═════════════════════════╣")
		fmt.Println("     ║ " + id + strings.Repeat(" ", 16-len(id)) + "║ " + name + strings.Repeat(" ", 52-len(name)) + "║ " + formatDateTime(dateTimeAdded) + "  ║")

	}
	fmt.Println("     ╚═════════════════╩═════════════════════════════════════════════════════╩═════════════════════════╝")
	fmt.Println(resetColour)
	returnToMainMenu()
}

// Option 1 function
func option1() {

}

// Option 2 function
func option2() {

}

// Option 3 function
// List CDRs previously inserted into (Yet Another PBX)
func option3(dbDetail databaseFunctionParameter) {

	var (
		voipCarrierID                string
		voipCarrierName              string
		cdrMonthYear                 string
		yapCDRInsertLogDateTimeAdded string
	)

	option3SQL, err := dbDetail.connection.Query(`SELECT
                                                        voip_carrier_id,
                                                        voip_carrier_name,
                                                        cdr_month_year,
                                                        yap_cdr_insert_log_date_time_added
                                                      FROM
                                                        cdr_importer.view___yap_cdr_insert_log_detail;`)

	// Error
	if err != nil {
		panic(err)
	}

	clearScreen()
	fmt.Println("")
	fmt.Println(textBoldBlack)
	fmt.Println("          ╔══════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ [3] List CDRs previously inserted into (Yet Another PBX) ║")
	fmt.Println("     ╔════╩════════════╦═════════════════════════════════════════════╩═══════╦═════════════════════╦══════════════════════════════════════╗")
	fmt.Println("     ║ VoIP Carrier ID ║                  VoIP Carrier Name                  ║ Month & Year of CDR ║ Date & Time of CDR Inserted into YAP ║")

	for option3SQL.Next() {

		err = option3SQL.Scan(
			&voipCarrierID,
			&voipCarrierName,
			&cdrMonthYear,
			&yapCDRInsertLogDateTimeAdded,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Println("     ╠═════════════════╬═════════════════════════════════════════════════════╬═════════════════════╬══════════════════════════════════════╣")
		fmt.Println("     ║ " + voipCarrierID + strings.Repeat(" ", 16-len(voipCarrierID)) + "║ " + voipCarrierName + strings.Repeat(" ", 52-len(voipCarrierName)) + "║ " + cdrMonthYear + "             ║ " + formatDateTime(yapCDRInsertLogDateTimeAdded) + "               ║")

	}
	fmt.Println("     ╚═════════════════╩═════════════════════════════════════════════════════╩═════════════════════╩══════════════════════════════════════╝")
	fmt.Println(resetColour)
	returnToMainMenu()
}

// Option 4 function
func option4() {

}

// Option 5 function
func option5() {

}

// Option 6 function
func option6() {

}

// Option 7 function
func option7() {

}

// Option 8 function
func option8() {

}

// Option 9 function
func option9() {

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
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + "                                                     " + textBoldWhite + resetColour + bgRed + textBoldWhite + "                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + textBoldWhite + "   Type \"exit\" or \"quit\" to terminate CDR Importer   " + resetColour + bgRed + textBoldWhite + "                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                           " + resetColour + bgMagenta + textBoldWhite + "                                                     " + resetColour + bgRed + textBoldWhite + "                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┃                                                                                                           ┃ " + resetColour)
	fmt.Println("     " + bgRed + textBoldWhite + " ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ " + resetColour)

	fmt.Print(textBoldBlack)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          [0] List all VoIP carrier IDs and names\n")
	fmt.Println("          [1] List a CDR from a particular month for a VoIP carrier ID\n")
	fmt.Println("          [2] List all call rates for a VoIP carrier\n")
	fmt.Println("          [3] List CDRs previously inserted into (Yet Another PBX)\n")
	fmt.Println("          [4] Edit a call rate for a VoIP carrier\n")
	fmt.Println("          [5] Add a new VoIP carrier\n")
	fmt.Println("          [6] Delete an existing VoIP carrier\n")
	fmt.Println("          [7] Import a new CDR into an existing VoIP carrier\n")
	fmt.Println("          [8] Delete a previsouly imported CDR from a VoIP carrier\n")
	fmt.Println("          [9] Insert a CDR from a particular month for a VoIP carrier into YAP (Yet Another PBX)\n")
	fmt.Println("")
	fmt.Print("          Select an option [0-9]: ")
	var option string
	fmt.Scanln(&option)
	fmt.Println(resetColour)

	// Values allowed for option
	var optionList = []string{"", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "exit", "Exit", "EXIT", "quit", "Quit", "QUIT"}
	validOption := slices.Contains(optionList, option)

	// Conditional statment to determine what happens when an option is input
	if validOption == false {
		messageBox("Invalid option - enter option [1-9] or exit", bgYellow)
		returnToMainMenu()
	}

	// If user typed exit or quit then stop program
	exit(option)

	// If the user pressed the enter/return key then re-run the main function
	if option == "" {
		main()
	}

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

	// Open database connection
	dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
	defer dbConnection.Close()

	// Error
	if err != nil {
		panic(err)
	}

	var dbDetail databaseFunctionParameter

	if option == "0" {
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		option0(dbDetail)
	} else if option == "1" {
		option1()
	} else if option == "2" {
		option2()
	} else if option == "3" {
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		option3(dbDetail)
	} else if option == "4" {
		option4()
	} else if option == "5" {
		option5()
	} else if option == "6" {
		option6()
	} else if option == "7" {
		option7()
	} else if option == "8" {
		option8()
	} else if option == "9" {
		option9()
	}
}

// Contributor(s):
// Elliot Michael Keavney
