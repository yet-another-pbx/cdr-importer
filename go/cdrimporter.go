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
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
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
const bgBlue = "\033[44m"
const bgPurple = "\033[45m"
const bgCyan = "\033[46m"

// Validation Messages
const validationMessageAlphaNum string = "alphanumeric and up to 50 characters long"
const validationMessageNumeric string = "numeric and up to 12 numbers long"

//----------------------------------------------------------------------------------------------------

// Database struct
type databaseFunctionParameter struct {
	connection       *sql.DB
	username         string
	password         string
	database         string
	address          string
	port             string
	transport        string
	tls              string
	column           string
	table            string
	columnWhere      string
	columnWhereValue string
}

// Make table struct
type sqlFunctionParameter struct {
	tableType                          string
	callDirection                      string
	voipCarrierID                      string
	filePath                           string
	ignoreFirstCSVLine                 bool
	cdrTagColumnNumber                 string
	cdrNumberDialledColumnNumber       string
	cdrDescriptionColumnNumber         string
	cdrChargeCodeColumnNumber          string
	cdrDurationColumnNumber            string
	cdrDateTimeColumnNumber            string
	cdrTimeColumn                      string
	cdrMonthYearColumnNumber           string
	rateCardDescriptionColumnNumber    string
	rateCardChargeCodeColumnNumber     string
	rateCardPricePerMinuteColumnNumber string
	rateCardPricePerCallColumnNumber   string
}

// Function to return slice of tableType
func callDirectionSlice() []string {
	callDirectionList := []string{"inbound", "outbound"}
	return callDirectionList
}

// Function to return slice of yesNoList
func yesNoSlice() []string {
	yesNoList := []string{"yes", "Yes", "YES", "y", "Y", "no", "No", "NO", "n", "N"}
	return yesNoList
}

//----------------------------------------------------------------------------------------------------

// Function to validate user input utlising the Go validator version 10 package
func validateInput(value string, valueType string) (validation bool) {

	validateInput := validator.New()
	// Conditional statments are used for each type of validation needed
	if valueType == "alphaNum" {
		validateInputAlphaNumErr := validateInput.Var(value, "alphanumspace,min=1,max=50")
		validateInputSymbolErr := validateInput.Var(value, "excludes=`!\"£$%^&*()-_=+{}[];:@'#~\\.<>/?")
		if validateInputAlphaNumErr != nil || validateInputSymbolErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "filePath" {
		validateInputDirErr := validateInput.Var(value, "file")
		if validateInputDirErr != nil {
			validation = false
			fmt.Println(value)
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "monthYear" {
		validateInputMonthYearErr := validateInput.Var(value, "datetime=01/2006,min=7,max=7")
		if validateInputMonthYearErr != nil {
			validation = false
			fmt.Println(value)
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "rate" {
		validateInputNumberErr := validateInput.Var(value, "numeric,max=12")
		if validateInputNumberErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else {
		panic("The validateInput function can only take the following arguments: alphaNum, filePath, monthYear or rate")
	}
}

// Function to validate column numbers input
func validateColumnNumber(value string) (validation bool) {

	// Validate the value is an int
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return false
	} else {
		if valueInt > 50 {
			return false
		} else if valueInt <= 0 {
			return false
		} else {
			return true
		}
	}
}

//----------------------------------------------------------------------------------------------------

// Function to terminate the running program if a user types "exit" or "quit"
func exit(value string) {

	if value == "exit" || value == "Exit" || value == "EXIT" || value == "quit" || value == "Quit" || value == "QUIT" {
		messageBox("Program exited successfully", bgGreen)
		fmt.Println("")
		os.Exit(0)
	}
}

// Function to return to the main menu if a user types "menu"
func mainMenu(value string) {

	if value == "menu" || value == "Menu" || value == "MENU" {
		cdrimporter()
	}
}

// Function to return to main menu when the enter key is pressed
func returnToMainMenu() {

	fmt.Print(textBoldBlack)
	fmt.Print("     Press the enter/return key to continue ")
	fmt.Print(resetColour)
	var enterMainMenu string
	fmt.Scanln(&enterMainMenu)
	if enterMainMenu == "" || enterMainMenu != "" {
		cdrimporter()
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
}

//----------------------------------------------------------------------------------------------------

// Function to select value from a table with the WHERE clause and return the value
func selectWhere(dbDetail databaseFunctionParameter) string {

	var selectWhere string
	selectWhereQuery, err := dbDetail.connection.Query(`SELECT
				  			      `+dbDetail.column+`
						            FROM
							      `+dbDetail.database+`.`+dbDetail.table+`
							    WHERE
							      `+dbDetail.columnWhere+` = ?;`, dbDetail.columnWhereValue)

	if err != nil {
		panic(err)
	}

	for selectWhereQuery.Next() {
		err := selectWhereQuery.Scan(&selectWhere)
		if err != nil {
			panic(err)
		}
	}
	return selectWhere
}

// Function to retrive VoIP carrier name(s) and ID(s) from the voip_carrier table
func voipCarrierSlice(dbDetail databaseFunctionParameter) ([][]string, []string) {

	// Get VoIP carrier name and ID from the database and append to slice
	var voipCarrierIDNameList [][]string
	var voipCarrierIDList []string

	var voipCarrierIDName string
	var voipCarrierID string

	voipCarrierIDNameSQL, err := dbDetail.connection.Query(`SELECT
						                  id,
                                                                  name
                                                                FROM
                                                                  cdr_importer.voip_carrier;`)

	// Error
	if err != nil {
		panic(err)
	}

	for voipCarrierIDNameSQL.Next() {

		err = voipCarrierIDNameSQL.Scan(
			&voipCarrierID,
			&voipCarrierIDName,
		)

		// Error
		if err != nil {
			panic(err)
		}

		var voipCarrierIDAndName []string

		voipCarrierIDAndName = append([]string{voipCarrierID}, []string{voipCarrierIDName}...)
		voipCarrierIDNameList = append(voipCarrierIDNameList, voipCarrierIDAndName)
		voipCarrierIDList = append(voipCarrierIDList, voipCarrierID)

	}
	return voipCarrierIDNameList, voipCarrierIDList
}

// Function to create table
func makeTable(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	var lastColumn string

	lastColumn = "column_50 varchar(255)"

	if sqlDetail.cdrTimeColumn == "n/a" || sqlDetail.cdrTimeColumn == "N/A" {
		if sqlDetail.tableType == "cdr" {
			lastColumn = "column_50 varchar(255), column_no_time varchar(0)"
		}
	}

	dbDetail.connection.Exec(`CREATE TABLE cdr_importer.` + sqlDetail.callDirection + `_` + sqlDetail.tableType + `_` + sqlDetail.voipCarrierID + ` (
		column_1 varchar(255), column_2 varchar(255), column_3 varchar(255), column_4 varchar(255), column_5 varchar(255),
		column_6 varchar(255), column_7 varchar(255), column_8 varchar(255), column_9 varchar(255), column_10 varchar(255),
                column_11 varchar(255), column_12 varchar(255), column_13 varchar(255), column_14 varchar(255), column_15 varchar(255),
                column_16 varchar(255), column_17 varchar(255), column_18 varchar(255), column_19 varchar(255), column_20 varchar(255),
                column_21 varchar(255), column_22 varchar(255), column_23 varchar(255), column_24 varchar(255), column_25 varchar(255),
                column_26 varchar(255), column_27 varchar(255), column_28 varchar(255), column_29 varchar(255), column_30 varchar(255),
                column_31 varchar(255), column_32 varchar(255), column_33 varchar(255), column_34 varchar(255), column_35 varchar(255),
                column_36 varchar(255), column_37 varchar(255), column_38 varchar(255), column_39 varchar(255), column_40 varchar(255),
                column_41 varchar(255), column_42 varchar(255), column_43 varchar(255), column_44 varchar(255), column_45 varchar(255),
                column_46 varchar(255), column_47 varchar(255), column_48 varchar(255), column_49 varchar(255), ` + lastColumn + `
	);`)
}

// Function to import a CSV file into a database table
func importCSV(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	var ignoreFirstLine string

	if sqlDetail.ignoreFirstCSVLine == true {
		ignoreFirstLine = "1"
	} else if sqlDetail.ignoreFirstCSVLine == false {
		ignoreFirstLine = "0"
	}

	filePath := sqlDetail.filePath
	mysql.RegisterLocalFile(filePath)
	_, err := dbDetail.connection.Exec("LOAD DATA LOCAL INFILE '" + filePath + "' INTO TABLE `" + dbDetail.table + "` COLUMNS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"' ESCAPED BY '\"' LINES TERMINATED BY '\n' IGNORE " + ignoreFirstLine + " LINES;")
	if err != nil {
		messageBox("File cannot be accessed or does not exist", bgRed)
		returnToMainMenu()
	}
}

// Function to create a rate card view
func makeRateCardView(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	dbDetail.connection.Exec(`CREATE VIEW view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + ` AS
	SELECT
		IFNULL(` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.rateCardDescriptionColumnNumber + `, '') AS 'description',
		IFNULL(` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.rateCardChargeCodeColumnNumber + `, '') AS 'charge_code',
                IFNULL(` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.rateCardPricePerMinuteColumnNumber + `, '') AS 'price_per_minute',
                IFNULL(` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.rateCardPricePerCallColumnNumber + `, '') AS 'price_per_call'
	FROM ` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `;`)
}

// Function to create a CDR view
func makeCDRView(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	var timeColumn string

	timeColumn = sqlDetail.cdrTimeColumn

	if sqlDetail.cdrTimeColumn == "n/a" || sqlDetail.cdrTimeColumn == "N/A" {
		timeColumn = "no_time"
	}

	dbDetail.connection.Exec(`CREATE VIEW view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + ` AS
        SELECT
                IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrTagColumnNumber + `, '') AS 'tag',
		IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrNumberDialledColumnNumber + `, '') AS 'number_dialled',
		IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrDescriptionColumnNumber + `, '') AS 'description',
                IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrChargeCodeColumnNumber + `, '') AS 'charge_code',
                IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrDurationColumnNumber + `, '') AS 'duration',
                IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrDateTimeColumnNumber + `, '') AS 'date_time',
                IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + timeColumn + `, '') AS 'time',                
                IFNULL(` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.column_` + sqlDetail.cdrMonthYearColumnNumber + `, '') AS 'month_year'
        FROM ` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `;`)
}

// Function to create a rate card and CDR joined view
func makeCDRRateCardView(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	dbDetail.connection.Exec(`CREATE VIEW view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + ` AS
	SELECT
		view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.tag AS 'cdr_tag',
		view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.number_dialled AS 'cdr_number_dialled',
		view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.description AS 'cdr_description',
		view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.charge_code AS 'cdr_charge_code',
		view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.duration AS 'cdr_duration',
		ROUND(TIME_TO_SEC(view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.duration), 0) AS 'cdr_duration_in_seconds',
		CASE
		  WHEN ROUND(TIME_TO_SEC(view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.duration), 0) < 60 THEN 60
		  ELSE ROUND(TIME_TO_SEC(view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.duration), 0)
		END AS 'cdr_duration_in_seconds_with_minimum_duration_60_seconds',
		DATE_FORMAT(view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.date_time, '%d/%m/%Y %H:%i:%s') AS 'cdr_date_time',
		TIME_FORMAT(view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.time, '%H:%i:%s') AS 'cdr_time',
		view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.month_year AS 'cdr_month_year',
		view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.description AS 'rate_card_description',
		view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.price_per_minute AS 'rate_card_price_per_minute',
		(view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.price_per_minute)/60 AS 'rate_card_price_per_second',
		view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.price_per_call AS 'rate_card_price_per_call'
	FROM view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + ` ` +
		`INNER JOIN view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `
	ON view___` + sqlDetail.callDirection + `_cdr_` + sqlDetail.voipCarrierID + `.charge_code = view___` + sqlDetail.callDirection + `_rate_card_` + sqlDetail.voipCarrierID + `.charge_code;`)
}

// Function to create itemised CDR view
func makeItemisedCDRView(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	dbDetail.connection.Exec(`CREATE VIEW view___` + sqlDetail.callDirection + `_itemised_cdr_` + sqlDetail.voipCarrierID + ` AS
	SELECT
		IFNULL(yap.view___invoice_item.customer_id, 'NO CUSTOMER') AS yap_customer_id,
		IFNULL(yap.view___invoice_item.customer_name, 'NO CUSTOMER') AS yap_customer_name,
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_tag,
		IFNULL(yap.view___invoice_item.service_product_name, '') AS yap_service_product_name,		
		IFNULL(yap.view___invoice_item.invoice_item_sales_tax_rate, '') AS yap_invoice_item_sales_tax_rate,
		IFNULL(yap.view___invoice_item.invoice_item_sales_tax_status, '') AS yap_invoice_item_sales_tax_status,
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_number_dialled,
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_description,
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_duration,
		IFNULL(cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_date_time, '') AS 'cdr_date_time',
		IFNULL(cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_time, '') AS 'cdr_time',
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_month_year,
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_price_per_minute,
		cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_price_per_call,
		ROUND((cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_duration_in_seconds * cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_price_per_second) + cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_price_per_call, 10) AS 'no_minimum_charge_for_call_duration',
		ROUND((cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_duration_in_seconds_with_minimum_duration_60_seconds * cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_price_per_second) + cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.rate_card_price_per_call, 10) AS 'minimum_charge_for_call_duration_below_60_seconds'
	FROM cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + ` ` +
		`LEFT JOIN yap.view___invoice_item
	ON cdr_importer.view___` + sqlDetail.callDirection + `_cdr_rate_card_` + sqlDetail.voipCarrierID + `.cdr_tag = yap.view___invoice_item.invoice_item_tag;`)
}

// Function to create total CDR view
func makeTotalCDRView(dbDetail databaseFunctionParameter, sqlDetail sqlFunctionParameter) {

	dbDetail.connection.Exec(`CREATE VIEW view___` + sqlDetail.callDirection + `_total_cdr_` + sqlDetail.voipCarrierID + ` AS
	SELECT
  		yap_customer_id,
  		yap_customer_name,
  		cdr_tag,
  		yap_service_product_name,
  		yap_invoice_item_sales_tax_rate,
  		yap_invoice_item_sales_tax_status,
  		cdr_month_year,
  		ROUND(SUM(no_minimum_charge_for_call_duration),2) AS total_no_minimum_charge_for_call_duration,
		ROUND(SUM(minimum_charge_for_call_duration_below_60_seconds),2) AS total_minimum_charge_for_call_duration_below_60_seconds
	FROM cdr_importer.view___` + sqlDetail.callDirection + `_itemised_cdr_` + sqlDetail.voipCarrierID + ` ` +
		`GROUP BY cdr_tag;`)
}

// Function to input carrierID
func voipCarrierIDInput(dbDetail databaseFunctionParameter, optionNumber string) string {

	var voipCarrierID string

	voipCarrierIDNameList, _ := voipCarrierSlice(dbDetail)

	fmt.Println(textBoldBlack)
	fmt.Println("     ╔═════════════════════════════════════════════════════╦═════════════════╗")
	fmt.Println("     ║                  VoIP Carrier Name                  ║ VoIP Carrier ID ║")

	for _, voipCarrierIDNameValue := range voipCarrierIDNameList {
		var id string = voipCarrierIDNameValue[0:][0]
		var name string = voipCarrierIDNameValue[0:][1]
		fmt.Println("     ╠═════════════════════════════════════════════════════╬═════════════════╣")
		fmt.Println("     ║ " + name + strings.Repeat(" ", 52-len(name)) + "║ " + id + strings.Repeat(" ", 16-len(id)) + "║")
	}
	fmt.Println("     ╚═════════════════════════════════════════════════════╩═════════════════╝")
	fmt.Println("")

	_, voipCarrierIDList := voipCarrierSlice(dbDetail)

	fmt.Print(textBoldBlack)
	fmt.Print("     Enter the VoIP carrier ID [Valid input - numeric]: ")
	fmt.Scanln(&voipCarrierID)

	// If the user pressed the enter/return key then re-run the main function
	if voipCarrierID == "" {
		cdrimporter()
	}

	// Return to main menu if menu is input
	mainMenu(voipCarrierID)

	// Check rateCardIgnoreFirstLine is contained in the slice
	validateVoIPCarrierID := slices.Contains(voipCarrierIDList, voipCarrierID)

	if validateVoIPCarrierID == false {
		// Invalid input message displays to the user
		messageBox("The VoIP carrier ID does not exist ", bgYellow)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			if optionNumber == "option1" {
				option1(dbDetail)
			} else if optionNumber == "option2" {
				option2(dbDetail)
			} else if optionNumber == "option3" {
				option3(dbDetail)
			} else if optionNumber == "option4" {
				option4(dbDetail)
			} else if optionNumber == "option5" {
				option5(dbDetail)
			} else if optionNumber == "option6" {
				option6(dbDetail)
			} else if optionNumber == "option7" {
				option7(dbDetail)
			} else if optionNumber == "option8" {
				option8(dbDetail)
			} else if optionNumber == "option9" {
				option9(dbDetail)
			} else {
				panic("Only option1 to option9 is a valid argument to the voipCarrierIDInput function")
			}
		}
	}
	return voipCarrierID
}

// Function to input the call direction
func callDirectionInput(dbDetail databaseFunctionParameter, optionNumber string) string {

	var callDirection string

	callDirectionList := callDirectionSlice()
	fmt.Println("")
	fmt.Print("     Enter the VoIP carrier CDR direction [Valid options - " + strings.Join(callDirectionList, ", ") + "]: ")
	fmt.Scan(&callDirection)
	// Return to main menu if menu is input
	mainMenu(callDirection)

	// Check callDirection is contained in the slice
	validateCallDirection := slices.Contains(callDirectionList, callDirection)

	if validateCallDirection == false {
		// Invalid input message displays to the user
		messageBox("Invalid option, please re-enter either "+(strings.Join(callDirectionList, ", ")+" "), bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			if optionNumber == "option1" {
				option1(dbDetail)
			} else if optionNumber == "option2" {
				option2(dbDetail)
			} else if optionNumber == "option3" {
				option3(dbDetail)
			} else if optionNumber == "option4" {
				option4(dbDetail)
			} else if optionNumber == "option5" {
				option5(dbDetail)
			} else if optionNumber == "option6" {
				option6(dbDetail)
			} else if optionNumber == "option7" {
				option7(dbDetail)
			} else if optionNumber == "option8" {
				option8(dbDetail)
			} else if optionNumber == "option9" {
				option9(dbDetail)
			} else {
				panic("Only option1 to option9 is a valid argument to the callDirectionInput function")
			}
		}
	}
	return callDirection
}

// Function to input the Month and Year for the CDR
func enterMonthYearInput(dbDetail databaseFunctionParameter, optionNumber string) string {

	var enterMonthYear string

	fmt.Println("")
	fmt.Print("     Enter the VoIP carrier CDR month and year [Valid format - MM/YYYY]: ")
	fmt.Scan(&enterMonthYear)
	// Return to main menu if menu is input
	mainMenu(enterMonthYear)

	// Validate monthYear is a date
	validateCDRMonthYear := validateInput(enterMonthYear, "monthYear")

	if validateCDRMonthYear == false {
		// Invalid input message displays to the user
		messageBox("Invalid month and year, please re-enter", bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			if optionNumber == "option1" {
				option1(dbDetail)
			} else if optionNumber == "option2" {
				option2(dbDetail)
			} else if optionNumber == "option3" {
				option3(dbDetail)
			} else if optionNumber == "option4" {
				option4(dbDetail)
			} else if optionNumber == "option5" {
				option5(dbDetail)
			} else if optionNumber == "option6" {
				option6(dbDetail)
			} else if optionNumber == "option7" {
				option7(dbDetail)
			} else if optionNumber == "option8" {
				option8(dbDetail)
			} else if optionNumber == "option9" {
				option9(dbDetail)
			} else {
				panic("Only option1 to option9 is a valid argument to the enterMonthYearInput function")
			}
		}
	}
	return enterMonthYear
}

// Function to input a path for a CSV file, can be a absolute or relative path
func csvImportPathInput(dbDetail databaseFunctionParameter, optionNumber string) string {

	var csvImportPath string

	fmt.Println("")
	fmt.Print("     Enter the path for the CSV file [Example: /root/inbound-rate-card.csv or ./outbound-cdr.csv]: ")
	fmt.Scan(&csvImportPath)
	// Return to main menu if menu is input
	mainMenu(csvImportPath)

	// Check csvImportPath is a file path
	validateCSVImportPath := validateInput(csvImportPath, "filePath")

	if validateCSVImportPath == false {
		// Invalid input message displays to the user
		messageBox("Invalid path for the CSV file, please check the files location and re-enter the path including the file name ", bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			if optionNumber == "option1" {
				option1(dbDetail)
			} else if optionNumber == "option2" {
				option2(dbDetail)
			} else if optionNumber == "option3" {
				option3(dbDetail)
			} else if optionNumber == "option4" {
				option4(dbDetail)
			} else if optionNumber == "option5" {
				option5(dbDetail)
			} else if optionNumber == "option6" {
				option6(dbDetail)
			} else if optionNumber == "option7" {
				option7(dbDetail)
			} else if optionNumber == "option8" {
				option8(dbDetail)
			} else if optionNumber == "option9" {
				option9(dbDetail)
			} else {
				panic("Only option1 to option9 is a valid argument to the csvImportPathInput function")
			}
		}
	}
	return csvImportPath
}

// Function to provide an option to ignore the first line of CSV file when importing
func csvIgnoreFirstLineOption(dbDetail databaseFunctionParameter, optionNumber string) string {

	var csvIgnoreFirstLine string

	yesNoList := yesNoSlice()
	fmt.Println("")
	fmt.Print("     Ignore the first row/line of the CSV file when importing into the database table? [Valid options - " + strings.Join(yesNoList, ", ") + "]: ")
	fmt.Scan(&csvIgnoreFirstLine)
	// Return to main menu if menu is input
	mainMenu(csvIgnoreFirstLine)

	// Check rateCardIgnoreFirstLine is contained in the slice
	validateCSVIgnoreFirstLine := slices.Contains(yesNoList, csvIgnoreFirstLine)

	if validateCSVIgnoreFirstLine == false {
		// Invalid input message displays to the user
		messageBox("Invalid option for ignore first row/line, please re-enter with either "+strings.Join(yesNoList, ", "), bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			if optionNumber == "option1" {
				option1(dbDetail)
			} else if optionNumber == "option2" {
				option2(dbDetail)
			} else if optionNumber == "option3" {
				option3(dbDetail)
			} else if optionNumber == "option4" {
				option4(dbDetail)
			} else if optionNumber == "option5" {
				option5(dbDetail)
			} else if optionNumber == "option6" {
				option6(dbDetail)
			} else if optionNumber == "option7" {
				option7(dbDetail)
			} else if optionNumber == "option8" {
				option8(dbDetail)
			} else if optionNumber == "option9" {
				option9(dbDetail)
			} else {
				panic("Only option1 to option9 is a valid argument to the csvIgnoreFirstLineOption function")
			}

		}
	}
	return csvIgnoreFirstLine
}

//----------------------------------------------------------------------------------------------------

// Option 0 function
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
	                                                cdr_importer.view___voip_carrier;`)

	// Error
	if err != nil {
		panic(err)
	}

	clearScreen()
	fmt.Println("")
	fmt.Println(textBoldBlack)
	fmt.Println("          ╔═════╦═════════════════════════════════════╗")
	fmt.Println("          ║ " + bgBlue + textBoldWhite + "[0]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List all VoIP carrier IDs and names" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩══════╦══════════════════════════════╩══════════════════════╦═════════════════════╗")
	fmt.Println("     ║ VoIP Carrier ID ║                  VoIP Carrier Name                  ║  Date & Time Added  ║")

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

		fmt.Println("     ╠═════════════════╬═════════════════════════════════════════════════════╬═════════════════════╣")
		fmt.Println("     ║ " + id + strings.Repeat(" ", 16-len(id)) + "║ " + name + strings.Repeat(" ", 52-len(name)) + "║ " + dateTimeAdded + " ║")

	}
	fmt.Println("     ╚═════════════════╩═════════════════════════════════════════════════════╩═════════════════════╝")
	fmt.Println(resetColour)
	returnToMainMenu()
}

// Option 1 function
// List a CDR from a particular month and year for a VoIP carrier
func option1(dbDetail databaseFunctionParameter) {

	var (
		voipCarrierID                                   string
		callDirection                                   string
		enterMonthYear                                  string
		yapCustomerID                                   string
		yapCustomerName                                 string
		cdrTag                                          string
		cdrNumberDialled                                string
		cdrDuration                                     string
		cdrDateTime                                     string
		cdrTime                                         string
		cdrMonthYear                                    string
		noMinimumChargeForCallDuration                  string
		minimumChargeForCallDurationBelow60Seconds      string
		rateCardPricePerCall                            string
		yapInvoiceItemSalesTaxRate                      string
		yapInvoiceItemSalesTaxStatus                    string
		totalNoMinimumChargeForCallDuration             string
		totalMinimumChargeForCallDurationBelow60Seconds string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          ╔═════╦════════════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgBlue + textBoldWhite + "[1]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List a CDR from a particular month and year for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦══════════════════════════════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")

	voipCarrierID = voipCarrierIDInput(dbDetail, "option1")
	callDirection = callDirectionInput(dbDetail, "option1")
	enterMonthYear = enterMonthYearInput(dbDetail, "option1")

	option1ASQL, err := dbDetail.connection.Query(`SELECT
							 yap_customer_id,
							 yap_customer_name,
							 cdr_tag,
							 yap_invoice_item_sales_tax_rate,
							 yap_invoice_item_sales_tax_status,
							 cdr_month_year,
							 total_no_minimum_charge_for_call_duration,
							 total_minimum_charge_for_call_duration_below_60_seconds
                                                      FROM
                                                        cdr_importer.view___`+callDirection+`_total_cdr_`+voipCarrierID+`
                                                      WHERE
                                                        cdr_month_year = ?;`, enterMonthYear)

	// Error
	if err != nil {
		// Invalid input message displays to the user
		messageBox("The VoIP carrier does not exist", bgYellow)
		fmt.Print("     Press the enter/return key to continue ", err)
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option1(dbDetail)
		}
	}

	clearScreen()
	fmt.Println("")
	fmt.Println(textBoldBlack)
	fmt.Println("          ╔═════╦════════════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgBlue + textBoldWhite + "[1]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List a CDR from a particular month and year for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩════════════╦═══════════════════════════════════════════════════╝")
	fmt.Println("     ║ Total Call Charge CDR ║")
	fmt.Println("     ╠═══════════════════╦═══╩════════════════════════════╦═══════════════════════╦══════════╦════════════╦═══════════╦═══════════════════════════════╦═══════════════════════════════╗")
	fmt.Println("     ║        YAP        ║               YAP              ║          TAG          ║  Sales   ║   Sales    ║ MonthYear ║    Total no Minimum Charge    ║ Total Minimum Charge For Call ║")
	fmt.Println("     ║    Customer ID    ║          Customer Name         ║ (Same on YAP Invoice) ║ Tax Rate ║ Tax Status ║  MM-YYYY  ║       For Call Duration       ║   Duration Below 60 Seconds   ║")

	for option1ASQL.Next() {

		err = option1ASQL.Scan(
			&yapCustomerID,
			&yapCustomerName,
			&cdrTag,
			&yapInvoiceItemSalesTaxRate,
			&yapInvoiceItemSalesTaxStatus,
			&cdrMonthYear,
			&totalNoMinimumChargeForCallDuration,
			&totalMinimumChargeForCallDurationBelow60Seconds,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Println("     ╠═══════════════════╬════════════════════════════════╬═══════════════════════╬══════════╬════════════╬═══════════╬═══════════════════════════════╬═══════════════════════════════╣")
		// If no customer highlight in red
		if yapCustomerID == "NO CUSTOMER" {
			fmt.Println("     ║ " + bgRed + textBoldWhite + yapCustomerID + resetColour + textBoldBlack + strings.Repeat(" ", 18-len(yapCustomerID)) + "║ " + bgRed + textBoldWhite + yapCustomerName + resetColour + textBoldBlack + strings.Repeat(" ", 31-len(yapCustomerName)) + "║ " + cdrTag + strings.Repeat(" ", 22-len(cdrTag)) + "║ " + yapInvoiceItemSalesTaxRate + strings.Repeat(" ", 9-len(yapInvoiceItemSalesTaxRate)) + "║ " + yapInvoiceItemSalesTaxStatus + strings.Repeat(" ", 11-len(yapInvoiceItemSalesTaxStatus)) + "║ " + cdrMonthYear + strings.Repeat(" ", 10-len(cdrMonthYear)) + "║ " + totalNoMinimumChargeForCallDuration + strings.Repeat(" ", 30-len(totalNoMinimumChargeForCallDuration)) + "║ " + totalMinimumChargeForCallDurationBelow60Seconds + strings.Repeat(" ", 30-len(totalMinimumChargeForCallDurationBelow60Seconds)) + "║")
		} else {
			fmt.Println("     ║ " + yapCustomerID + strings.Repeat(" ", 18-len(yapCustomerID)) + "║ " + yapCustomerName + strings.Repeat(" ", 31-len(yapCustomerName)) + "║ " + cdrTag + strings.Repeat(" ", 22-len(cdrTag)) + "║ " + yapInvoiceItemSalesTaxRate + "%" + strings.Repeat(" ", 8-len(yapInvoiceItemSalesTaxRate)) + "║ " + yapInvoiceItemSalesTaxStatus + strings.Repeat(" ", 11-len(yapInvoiceItemSalesTaxStatus)) + "║ " + cdrMonthYear + strings.Repeat(" ", 10-len(cdrMonthYear)) + "║ " + totalNoMinimumChargeForCallDuration + strings.Repeat(" ", 30-len(totalNoMinimumChargeForCallDuration)) + "║ " + totalMinimumChargeForCallDurationBelow60Seconds + strings.Repeat(" ", 30-len(totalMinimumChargeForCallDurationBelow60Seconds)) + "║")
		}
	}
	fmt.Println("     ╚═══════════════════╩════════════════════════════════╩═══════════════════════╩══════════╩════════════╩═══════════╩═══════════════════════════════╩═══════════════════════════════╝")

	fmt.Println(resetColour)

	option1BSQL, err := dbDetail.connection.Query(`SELECT
                                                        yap_customer_id,
							yap_customer_name,
							cdr_tag,
							cdr_number_dialled,
							cdr_duration,
							cdr_date_time,
							cdr_time,
							no_minimum_charge_for_call_duration,
							minimum_charge_for_call_duration_below_60_seconds,
							rate_card_price_per_call
                                                      FROM
                                                        cdr_importer.view___`+callDirection+`_itemised_cdr_`+voipCarrierID+`
						      WHERE
                                                        cdr_month_year = ?;`, enterMonthYear)

	// Error
	if err != nil {
		// Invalid input message displays to the user
		messageBox("The VoIP carrier does not exist", bgYellow)
		fmt.Print("     Press the enter/return key to continue ", err)
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option1(dbDetail)
		}
	}

	fmt.Println(textBoldBlack)
	fmt.Println("     ╔══════════════════════════╗")
	fmt.Println("     ║ Itemised Call Charge CDR ║")
	fmt.Println("     ╠═══════════════════╦══════╩═════════════════════════╦═══════════════════════╦══════════════════╦══════════╦═══════════════════════╦══════════╦═══════════════════════════╦═══════════════════════════╦══════════╗")
	fmt.Println("     ║        YAP        ║               YAP              ║          TAG          ║    Telephone     ║ Duration ║       Date (Time)     ║   Time   ║     No Minimum Charge     ║  Minimum Charge For Call  ║  Charge  ║")
	fmt.Println("     ║    Customer ID    ║          Customer Name         ║ (Same on YAP Invoice) ║  Number Dialled  ║ HH:MM:SS ║ DD/MM/YYYY (HH:MM:SS) ║ HH:MM:SS ║     For Call Duration     ║ Duration Below 60 Seconds ║ Per Call ║")

	for option1BSQL.Next() {

		err = option1BSQL.Scan(
			&yapCustomerID,
			&yapCustomerName,
			&cdrTag,
			&cdrNumberDialled,
			&cdrDuration,
			&cdrDateTime,
			&cdrTime,
			&noMinimumChargeForCallDuration,
			&minimumChargeForCallDurationBelow60Seconds,
			&rateCardPricePerCall,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Println("     ╠═══════════════════╬════════════════════════════════╬═══════════════════════╬══════════════════╬══════════╬═══════════════════════╬══════════╬═══════════════════════════╬═══════════════════════════╬══════════╣")
		// If no customer highlight in red
		if yapCustomerID == "NO CUSTOMER" {
			fmt.Println("     ║ " + bgRed + textBoldWhite + yapCustomerID + resetColour + textBoldBlack + strings.Repeat(" ", 18-len(yapCustomerID)) + "║ " + bgRed + textBoldWhite + yapCustomerName + resetColour + textBoldBlack + strings.Repeat(" ", 31-len(yapCustomerName)) + "║ " + cdrTag + strings.Repeat(" ", 22-len(cdrTag)) + "║ " + cdrNumberDialled + strings.Repeat(" ", 17-len(cdrNumberDialled)) + "║ " + cdrDuration + strings.Repeat(" ", 9-len(cdrDuration)) + "║ " + cdrDateTime + strings.Repeat(" ", 22-len(cdrDateTime)) + "║ " + cdrTime + strings.Repeat(" ", 9-len(cdrTime)) + "║ " + noMinimumChargeForCallDuration + strings.Repeat(" ", 26-len(noMinimumChargeForCallDuration)) + "║ " + minimumChargeForCallDurationBelow60Seconds + strings.Repeat(" ", 10-len(minimumChargeForCallDurationBelow60Seconds)) + "║ " + rateCardPricePerCall + strings.Repeat(" ", 9-len(rateCardPricePerCall)) + "║")
		} else {
			fmt.Println("     ║ " + yapCustomerID + strings.Repeat(" ", 18-len(yapCustomerID)) + "║ " + yapCustomerName + strings.Repeat(" ", 31-len(yapCustomerName)) + "║ " + cdrTag + strings.Repeat(" ", 22-len(cdrTag)) + "║ " + cdrNumberDialled + strings.Repeat(" ", 17-len(cdrNumberDialled)) + "║ " + cdrDuration + strings.Repeat(" ", 9-len(cdrDuration)) + "║ " + cdrDateTime + strings.Repeat(" ", 22-len(cdrDateTime)) + "║ " + cdrTime + strings.Repeat(" ", 9-len(cdrTime)) + "║ " + noMinimumChargeForCallDuration + strings.Repeat(" ", 26-len(noMinimumChargeForCallDuration)) + "║ " + minimumChargeForCallDurationBelow60Seconds + strings.Repeat(" ", 26-len(minimumChargeForCallDurationBelow60Seconds)) + "║ " + rateCardPricePerCall + strings.Repeat(" ", 9-len(rateCardPricePerCall)) + "║")
		}
	}
	fmt.Println("     ╚═══════════════════╩════════════════════════════════╩═══════════════════════╩══════════════════╩══════════╩═══════════════════════╩══════════╩═══════════════════════════╩═══════════════════════════╩══════════╝")
	fmt.Println(resetColour)
	returnToMainMenu()
}

// Option 2 function
// List all call rates for a VoIP carrier
func option2(dbDetail databaseFunctionParameter) {

	var (
		voipCarrierID string
		callDirection string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          ╔═════╦════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgBlue + textBoldWhite + "[2]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List all call rates for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦══════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")

	voipCarrierID = voipCarrierIDInput(dbDetail, "option2")
	callDirection = callDirectionInput(dbDetail, "option2")

	var (
		description    string
		chargeCode     string
		pricePerMinute string
		pricePerCall   string
	)

	option2SQL, err := dbDetail.connection.Query(`SELECT
                                                        description,
                                                        charge_code,
                                                        price_per_minute,
                                                        price_per_call
                                                      FROM
                                                        cdr_importer.view___` + callDirection + `_rate_card_` + voipCarrierID + `;`)

	// Error
	if err != nil {
		// Invalid input message displays to the user
		messageBox("The VoIP carrier does not exist", bgYellow)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option2(dbDetail)
		}
	}

	clearScreen()
	fmt.Println("")
	fmt.Println(textBoldBlack)
	fmt.Println("          ╔═════╦════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgBlue + textBoldWhite + "[2]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List all call rates for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩════════════════════════════════════════╩═════════════════════════════════════════════════╦═════════════════════════╦════════════════════════╦════════════════════════╗")
	fmt.Println("     ║                                             Description                                             ║       Charge Code       ║    Price Per Minute    ║     Price Per Call     ║")

	for option2SQL.Next() {

		err = option2SQL.Scan(
			&description,
			&chargeCode,
			&pricePerMinute,
			&pricePerCall,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Println("     ╠═════════════════════════════════════════════════════════════════════════════════════════════════════╬═════════════════════════╬════════════════════════╬════════════════════════╣")
		fmt.Println("     ║ " + description + strings.Repeat(" ", 100-len(description)) + "║ " + chargeCode + strings.Repeat(" ", 24-len(chargeCode)) + "║ " + pricePerMinute + strings.Repeat(" ", 23-len(pricePerMinute)) + "║ " + pricePerCall + strings.Repeat(" ", 23-len(pricePerCall)) + "║")

	}
	fmt.Println("     ╚═════════════════════════════════════════════════════════════════════════════════════════════════════╩═════════════════════════╩════════════════════════╩════════════════════════╝")
	fmt.Println(resetColour)
	returnToMainMenu()
}

// Option 3 function
// Log of re-rated CDRs previously inserted into (Yet Another PBX)
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
	fmt.Println("          ╔═════╦══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgBlue + textBoldWhite + "[3]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "Show log of re-rated CDRs previously inserted into YAP (Yet Another PBX)" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩══════╦═════════════════════════════════════════════════════╦═════════════╩═══════╦══════════════════════════════════════╗")
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
		fmt.Println("     ║ " + voipCarrierID + strings.Repeat(" ", 16-len(voipCarrierID)) + "║ " + voipCarrierName + strings.Repeat(" ", 52-len(voipCarrierName)) + "║ " + cdrMonthYear + "             ║ " + yapCDRInsertLogDateTimeAdded + "                  ║")

	}
	fmt.Println("     ╚═════════════════╩═════════════════════════════════════════════════════╩═════════════════════╩══════════════════════════════════════╝")
	fmt.Println(resetColour)
	returnToMainMenu()
}

// Option 4 function
// Add a new VoIP carrier
func option4(dbDetail databaseFunctionParameter) {

	var (
		newName                            string
		callDirection                      string
		rateCardFilePath                   string
		rateCardIgnoreFirstLine            string
		cdrTagColumnNumber                 string
		cdrNumberDialledColumnNumber       string
		cdrDescriptionColumnNumber         string
		cdrChargeCodeColumnNumber          string
		cdrDurationColumnNumber            string
		cdrDateTimeColumnNumber            string
		cdrTimeColumn                      string
		cdrMonthYearColumnNumber           string
		rateCardDescriptionColumnNumber    string
		rateCardChargeCodeColumnNumber     string
		rateCardPricePerMinuteColumnNumber string
		rateCardPricePerCallColumnNumber   string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println(textBoldBlack)
	fmt.Println("          ╔═════╦═══════════════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgGreen + textBoldWhite + "[4]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Add a new VoIP carrier (must have a call rate card in CSV format)" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦═════════════════════════════════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")

	fmt.Println("")
	fmt.Print("     Enter the VoIP carrier name [Valid input - alphanumeric up to 50 characters long]: ")
	fmt.Scanln(&newName)

	// If the user pressed the enter/return key then re-run the main function
	if newName == "" {
		cdrimporter()
	}

	// Return to main menu if menu is input
	mainMenu(newName)

	// Validate new name input
	validateNewName := validateInput(newName, "alphaNum")

	if validateNewName == false {
		// Invalid input message displays to the user
		messageBox("Invalid input, please re-enter a new name that is "+validationMessageAlphaNum, bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option4(dbDetail)
		}
	}

	callDirection = callDirectionInput(dbDetail, "option4")

	// Open database connection
	dbConnection, err := sql.Open("mysql", dbDetail.username+":"+dbDetail.password+"@"+dbDetail.transport+"("+dbDetail.address+":"+dbDetail.port+")/"+dbDetail.database+"?tls="+dbDetail.tls)
	defer dbConnection.Close()

	// Error
	if err != nil {
		panic(err)
	}

	dbDetail.connection = dbConnection

	// Get VoIP carrier name
	dbDetail.column = "name"
	dbDetail.table = "voip_carrier"
	dbDetail.columnWhere = "name"
	dbDetail.columnWhereValue = newName + " (" + callDirection + ")"

	// Check if name already exists in the voip_carrier table
	voipCarrierName := selectWhere(dbDetail)

	if voipCarrierName == newName+" ("+callDirection+")" {
		// VoIP carrier name already exists
		messageBox("VoIP carrier name already exists, please choose another name ", bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option4(dbDetail)
		}
	}

	// Get the rate card path
	rateCardFilePath = csvImportPathInput(dbDetail, "option4")

	// Remove the first line of the rate card option
	rateCardIgnoreFirstLine = csvIgnoreFirstLineOption(dbDetail, "option4")

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the tag, usally the customers phone number but could be a SIP trunk username or ID [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrTagColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrTagColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the phone number dialled [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrNumberDialledColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrNumberDialledColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the call record description [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrDescriptionColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrDescriptionColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the charge code [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrChargeCodeColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrChargeCodeColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the call duration [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrDurationColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrDurationColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the date, if the date and time is combined in one column this option will also take date and time values [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrDateTimeColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrDateTimeColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the time, if the date and time is combined in one column enter n/a (not applicable) for this option [Valid input - Number 1-50, n/a]: ")
	fmt.Scan(&cdrTimeColumn)
	// Return to main menu if menu is input
	mainMenu(cdrTimeColumn)

	fmt.Println("")
	fmt.Print("     Enter the CDR column number for the month year [Valid input - Number 1-50]: ")
	fmt.Scan(&cdrMonthYearColumnNumber)
	// Return to main menu if menu is input
	mainMenu(cdrMonthYearColumnNumber)

	// Validate input for CDR column numbers
	validateCDRTagColumnNumber := validateColumnNumber(cdrTagColumnNumber)
	validateCDRNumberDialledColumnNumber := validateColumnNumber(cdrNumberDialledColumnNumber)
	validateCDRDescriptionColumnNumber := validateColumnNumber(cdrDescriptionColumnNumber)
	validateCDRChargeCodeColumnNumber := validateColumnNumber(cdrChargeCodeColumnNumber)
	validateCDRDurationColumnNumber := validateColumnNumber(cdrDurationColumnNumber)
	validateCDRDateTimeColumnNumber := validateColumnNumber(cdrDateTimeColumnNumber)
	validateCDRTimeColumnNumber := validateColumnNumber(cdrTimeColumn)
	validateCDRMonthYearColumnNumber := validateColumnNumber(cdrMonthYearColumnNumber)

	if validateCDRTagColumnNumber == false || validateCDRNumberDialledColumnNumber == false || validateCDRDescriptionColumnNumber == false || validateCDRDurationColumnNumber == false || validateCDRDateTimeColumnNumber == false || validateCDRChargeCodeColumnNumber == false || validateCDRMonthYearColumnNumber == false {
		// Invalid input message displays to the user
		messageBox("Invalid input for CDR column number, please re-enter a number between 1 and 50 ", bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option4(dbDetail)
		}
	} else if cdrTimeColumn != "n/a" && cdrTimeColumn != "N/A" && validateCDRTimeColumnNumber == false {
		// Invalid input message displays to the user
		messageBox("Invalid input for CDR time column, please re-enter a number between 1 and 50 or n/a", bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option4(dbDetail)
		}
	}

	fmt.Println("")
	fmt.Print("     Enter the rate card column number for the description [Valid input - Number 1-50]: ")
	fmt.Scan(&rateCardDescriptionColumnNumber)
	// Return to main menu if menu is input
	mainMenu(rateCardDescriptionColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the rate card column number for the charge code [Valid input - Number 1-50]: ")
	fmt.Scan(&rateCardChargeCodeColumnNumber)
	// Return to main menu if menu is input
	mainMenu(rateCardChargeCodeColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the rate card column number for the price per minute [Valid input - Number 1-50]: ")
	fmt.Scan(&rateCardPricePerMinuteColumnNumber)
	// Return to main menu if menu is input
	mainMenu(rateCardPricePerMinuteColumnNumber)

	fmt.Println("")
	fmt.Print("     Enter the rate card column number for the price per call [Valid input - Number 1-50]: ")
	fmt.Scan(&rateCardPricePerCallColumnNumber)
	// Return to main menu if menu is input
	mainMenu(rateCardPricePerCallColumnNumber)

	// Validate input for rate card column numbers
	validateRateCardDescriptionColumnNumber := validateColumnNumber(rateCardDescriptionColumnNumber)
	validateRateCardChargeCodeColumnNumber := validateColumnNumber(rateCardChargeCodeColumnNumber)
	validateRateCardPricePerMinuteColumnNumber := validateColumnNumber(rateCardPricePerMinuteColumnNumber)
	validateRateCardPricePerCallColumnNumber := validateColumnNumber(rateCardPricePerCallColumnNumber)

	if validateRateCardDescriptionColumnNumber == false || validateRateCardChargeCodeColumnNumber == false || validateRateCardPricePerMinuteColumnNumber == false || validateRateCardPricePerCallColumnNumber == false {
		// Invalid input message displays to the user
		messageBox("Invalid input for rate card column numbers, please re-enter a number between 1 and 50 ", bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option4(dbDetail)
		}
	}

	// Insert new VoIP carrier
	dbDetail.connection.Exec("INSERT INTO cdr_importer.voip_carrier (name, cdr_month_year_column, rate_card_charge_code_column, rate_card_price_per_minute_column, rate_card_price_per_call_column) VALUES(?, ?, ?, ?, ?);", newName+" ("+callDirection+")", "column_"+cdrMonthYearColumnNumber, "column_"+rateCardChargeCodeColumnNumber, "column_"+rateCardPricePerMinuteColumnNumber, "column_"+rateCardPricePerCallColumnNumber)

	// Get VoIP carrier ID
	dbDetail.column = "id"
	dbDetail.table = "voip_carrier"
	dbDetail.columnWhere = "name"
	dbDetail.columnWhereValue = newName + " (" + callDirection + ")"

	voipCarrierID := selectWhere(dbDetail)

	if voipCarrierID == "" {
		// Inform the user the VoIP carrier was not created
		messageBox("Error creating VoIP carrier", bgRed)
		returnToMainMenu()
	} else {
		var sqlDetail sqlFunctionParameter

		sqlDetail.callDirection = callDirection
		sqlDetail.voipCarrierID = voipCarrierID
		sqlDetail.cdrTimeColumn = cdrTimeColumn

		// Create CDR table
		sqlDetail.tableType = "cdr"
		makeTable(dbDetail, sqlDetail)

		// Create rate card table
		sqlDetail.tableType = "rate_card"
		makeTable(dbDetail, sqlDetail)

		// Create rate card view
		sqlDetail.rateCardDescriptionColumnNumber = rateCardDescriptionColumnNumber
		sqlDetail.rateCardChargeCodeColumnNumber = rateCardChargeCodeColumnNumber
		sqlDetail.rateCardPricePerMinuteColumnNumber = rateCardPricePerMinuteColumnNumber
		sqlDetail.rateCardPricePerCallColumnNumber = rateCardPricePerCallColumnNumber
		makeRateCardView(dbDetail, sqlDetail)

		// Create CDR view
		sqlDetail.cdrTagColumnNumber = cdrTagColumnNumber
		sqlDetail.cdrNumberDialledColumnNumber = cdrNumberDialledColumnNumber
		sqlDetail.cdrDescriptionColumnNumber = cdrDescriptionColumnNumber
		sqlDetail.cdrChargeCodeColumnNumber = cdrChargeCodeColumnNumber
		sqlDetail.cdrDurationColumnNumber = cdrDurationColumnNumber
		sqlDetail.cdrDateTimeColumnNumber = cdrDateTimeColumnNumber
		sqlDetail.cdrMonthYearColumnNumber = cdrMonthYearColumnNumber
		makeCDRView(dbDetail, sqlDetail)

		// Create CDR rate card view
		makeCDRRateCardView(dbDetail, sqlDetail)

		// Create itemised CDR view
		makeItemisedCDRView(dbDetail, sqlDetail)

		// Create total CDR view, used mostly for YAP invoice_item
		makeTotalCDRView(dbDetail, sqlDetail)

		// Import rate card CSV into VoIP carrier rate card table
		dbDetail.table = callDirection + "_rate_card_" + voipCarrierID
		sqlDetail.filePath = rateCardFilePath

		if rateCardIgnoreFirstLine == "yes" || rateCardIgnoreFirstLine == "Yes" || rateCardIgnoreFirstLine == "YES" || rateCardIgnoreFirstLine == "y" || rateCardIgnoreFirstLine == "Y" {
			sqlDetail.ignoreFirstCSVLine = true
		} else {
			sqlDetail.ignoreFirstCSVLine = false
		}

		importCSV(dbDetail, sqlDetail)

		// Inform the user the VoIP carrier was created
		messageBox("VoIP carrier and rate card added ", bgGreen)
		returnToMainMenu()
	}

	fmt.Println(resetColour)
}

// Option 5 function
// Replace or edit an exisitng call rate card for a VoIP carrier
func option5(dbDetail databaseFunctionParameter) {

	var (
		rateCardOption string
		voipCarrierID  string
		callDirection  string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          ╔═════╦═══════════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgPurple + textBoldWhite + "[5]" + resetColour + textBoldBlack + " ║ " + bgPurple + textBoldWhite + "Replace or edit an exisitng call rate card for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦═════════════════════════════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Println("     ╔═════╦════════════════════════════════════╗")
	fmt.Println("     ║ " + bgGreen + textBoldWhite + "[A]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Replace an exisitng call rate card" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬═════════════════════════════════╦══╝")
	fmt.Println("     ║ " + bgPurple + textBoldWhite + "[B]" + resetColour + textBoldBlack + " ║ " + bgPurple + textBoldWhite + "Edit an existing call rate card" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╚═════╩═════════════════════════════════╝")
	fmt.Println("")
	fmt.Print("     " + textBoldBlack + "Select an option [A-B]: ")
	fmt.Scanln(&rateCardOption)
	fmt.Println(resetColour)

	// Values allowed for option
	var rateCardOptionList = []string{"", "A", "B", "a", "b", "menu", "Menu", "MENU"}
	validRateCardOption := slices.Contains(rateCardOptionList, rateCardOption)

	// Conditional statment to determine what happens when an option is input
	if validRateCardOption == false {
		messageBox("Invalid option - enter option [A-B]", bgYellow)
		returnToMainMenu()
	}

	// If the user pressed the enter/return key then re-run the main function
	if rateCardOption == "" {
		cdrimporter()
	}

	clearScreen()
	fmt.Println("")
	voipCarrierID = voipCarrierIDInput(dbDetail, "option5")
	callDirection = callDirectionInput(dbDetail, "option5")

	if rateCardOption == "A" || rateCardOption == "a" {

		var (
			rateCardFilePath        string
			rateCardIgnoreFirstLine string
		)

		clearScreen()
		fmt.Println("")
		fmt.Println("")
		fmt.Println("          ╔═════╦════════════════════════════════════╗")
		fmt.Println("          ║ " + bgGreen + textBoldWhite + "[A]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Replace an exisitng call rate card" + resetColour + textBoldBlack + " ║")
		fmt.Println("     ╔════╩═════╩════════════════════════════════════╩══════════════════════════════════════════════════════════════════╗")
		fmt.Println("     ║ " + textBoldBlack + "The rate card columns must be the same as before, if the columns are different add a new VoIP carrier (Option 4) ║")
		fmt.Println("     ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")

		// Import the new rate card
		rateCardFilePath = csvImportPathInput(dbDetail, "option5")

		// Ignore first line of CSV option
		rateCardIgnoreFirstLine = csvIgnoreFirstLineOption(dbDetail, "option5")

		// Delete the exisitng rate card
		dbDetail.connection.Exec("DELETE FROM cdr_importer." + callDirection + "_rate_card_" + voipCarrierID)

		var sqlDetail sqlFunctionParameter

		// Import rate card CSV into VoIP carrier rate card table
		dbDetail.table = callDirection + "_rate_card_" + voipCarrierID
		sqlDetail.filePath = rateCardFilePath

		if rateCardIgnoreFirstLine == "yes" || rateCardIgnoreFirstLine == "Yes" || rateCardIgnoreFirstLine == "YES" || rateCardIgnoreFirstLine == "y" || rateCardIgnoreFirstLine == "Y" {
			sqlDetail.ignoreFirstCSVLine = true
		} else {
			sqlDetail.ignoreFirstCSVLine = false
		}

		importCSV(dbDetail, sqlDetail)

		// Inform the user the rate card was replaced
		messageBox("Call rate card has been replaced ", bgGreen)
		returnToMainMenu()

	} else if rateCardOption == "B" || rateCardOption == "b" {

		var (
			chargeCode               string
			rateOption               string
			rateCardChargeCodeColumn string
			newRate                  string
		)

		clearScreen()
		fmt.Println("")
		fmt.Println("")
		fmt.Println("          ╔═════╦═════════════════════════════════╗")
		fmt.Println("          ║ " + bgPurple + textBoldWhite + "[B]" + resetColour + textBoldBlack + " ║ " + bgPurple + textBoldWhite + "Edit an existing call rate card" + resetColour + textBoldBlack + " ║")
		fmt.Println("     ╔════╩═════╩═════════════════════════════════╩═════════════════════════════════════════════════════════════════════════╗")
		fmt.Println("     ║ " + textBoldBlack + "The charge code must be known in advance to replace with a new rate, if the charge code is unknown select (Option 2) ║")
		fmt.Println("     ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")

		fmt.Println("")
		fmt.Print("     Enter the charge code [Valid input - alphanumeric up to 50 characters long]: ")
		fmt.Scanln(&chargeCode)

		// Return to main menu if menu is input
		mainMenu(chargeCode)

		// Validate new name input
		validateChargeCode := validateInput(chargeCode, "alphaNum")

		if validateChargeCode == false {
			// Invalid input message displays to the user
			messageBox("Invalid input, please re-enter a charge code that is "+validationMessageAlphaNum, bgYellow)
			fmt.Print(textBoldBlack)
			fmt.Print("     Press the enter/return key to continue ")
			fmt.Print(resetColour)
			var enter string
			fmt.Scanln(&enter)
			if enter == "" || enter != "" {
				option5(dbDetail)
			}
		}

		// Choose wether to update the per minute price or the per call price
		rateOptionList := []string{"minute", "min", "Minute", "Min", "m", "M", "call", "Call", "c", "C"}
		fmt.Println("")
		fmt.Print("     Update the per minute price or update the per call price? [Valid Input - " + strings.Join(rateOptionList, ", ") + " ]: ")
		fmt.Scanln(&rateOption)

		// Check rateOption is contained in the slice
		validateRateOption := slices.Contains(rateOptionList, rateOption)

		if validateRateOption == false {
			// Invalid input message displays to the user
			messageBox("Invalid input, please re-enter with either "+strings.Join(rateOptionList, ", "), bgYellow)
			fmt.Print(textBoldBlack)
			fmt.Print("     Press the enter/return key to continue ")
			fmt.Print(resetColour)
			var enter string
			fmt.Scanln(&enter)
			if enter == "" || enter != "" {
				option5(dbDetail)
			}
		}

		fmt.Println("")
		fmt.Print("     Enter the new rate [Valid input - numeric decimal up to 12 numbers long]: ")
		fmt.Scanln(&newRate)

		// Return to main menu if menu is input
		mainMenu(newRate)

		// Validate the newRate input
		validateNewRate := validateInput(newRate, "rate")

		if validateNewRate == false {
			// Invalid input message displays to the user
			messageBox("Invalid input, please re-enter a new rate that is "+validationMessageNumeric, bgYellow)
			fmt.Print(textBoldBlack)
			fmt.Print("     Press the enter/return key to continue ")
			fmt.Print(resetColour)
			var enter string
			fmt.Scanln(&enter)
			if enter == "" || enter != "" {
				option5(dbDetail)
			}
		}

		// Retrieve rate card chargeCode column number from the voip_carrier table
		dbDetail.column = "rate_card_charge_code_column"
		dbDetail.table = "voip_carrier"
		dbDetail.columnWhere = "id"
		dbDetail.columnWhereValue = voipCarrierID
		rateCardChargeCodeColumn = selectWhere(dbDetail)

		if rateOption == "minute" || rateOption == "Minute" || rateOption == "Min" || rateOption == "min" || rateOption == "m" || rateOption == "M" {

			// Retrieve rate card per minute column number from the voip_carrier table
			dbDetail.column = "rate_card_price_per_minute_column"
			dbDetail.table = "voip_carrier"
			dbDetail.columnWhere = "id"
			dbDetail.columnWhereValue = voipCarrierID
			rateCardPricePerMinuteColumn := selectWhere(dbDetail)

			// Update the charge code with a new rate
			dbDetail.connection.Exec("UPDATE cdr_importer."+callDirection+"_rate_card_"+voipCarrierID+" SET "+rateCardPricePerMinuteColumn+" = ? WHERE "+rateCardChargeCodeColumn+" = ?;", newRate, chargeCode)

			// Inform the user the rate card per minute price has been updated
			messageBox("Call rate card per minute price has been updated ", bgGreen)
			returnToMainMenu()

		} else if rateOption == "call" || rateOption == "Call" || rateOption == "c" || rateOption == "C" {

			// Retrieve rate card per card column number from the voip_carrier table
			dbDetail.column = "rate_card_price_per_call_column"
			dbDetail.table = "voip_carrier"
			dbDetail.columnWhere = "id"
			dbDetail.columnWhereValue = voipCarrierID
			rateCardPricePerCallColumn := selectWhere(dbDetail)

			// Update the charge code with a new rate
			dbDetail.connection.Exec("UPDATE cdr_importer."+callDirection+"_rate_card_"+voipCarrierID+" SET "+rateCardPricePerCallColumn+" = ? WHERE "+rateCardChargeCodeColumn+" = ?;", newRate, chargeCode)

			// Inform the user the rate card per call price has been updated
			messageBox("Call rate card per call price has been updated ", bgGreen)
			returnToMainMenu()
		}
	}
	returnToMainMenu()
}

// Option 6 function
// Import a new CDR into an existing VoIP carrier
func option6(dbDetail databaseFunctionParameter) {

	var (
		voipCarrierID      string
		callDirection      string
		cdrFilePath        string
		cdrIgnoreFirstLine string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          ╔═════╦════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgGreen + textBoldWhite + "[6]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Import a new CDR into an existing VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦══════════════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")

	voipCarrierID = voipCarrierIDInput(dbDetail, "option6")
	callDirection = callDirectionInput(dbDetail, "option6")
	cdrFilePath = csvImportPathInput(dbDetail, "option6")
	cdrIgnoreFirstLine = csvIgnoreFirstLineOption(dbDetail, "option6")

	var sqlDetail sqlFunctionParameter

	// Import rate card CSV into VoIP carrier rate card table
	dbDetail.table = callDirection + "_cdr_" + voipCarrierID
	sqlDetail.filePath = cdrFilePath

	if cdrIgnoreFirstLine == "yes" || cdrIgnoreFirstLine == "Yes" || cdrIgnoreFirstLine == "YES" || cdrIgnoreFirstLine == "y" || cdrIgnoreFirstLine == "Y" {
		sqlDetail.ignoreFirstCSVLine = true
	} else {
		sqlDetail.ignoreFirstCSVLine = false
	}

	importCSV(dbDetail, sqlDetail)

	// Inform the user the VoIP carrier was created
	messageBox("CDR added", bgGreen)
	returnToMainMenu()
}

// Option 7 function
// Delete a previsouly imported CDR from a VoIP carrier
func option7(dbDetail databaseFunctionParameter) {

	var (
		voipCarrierID  string
		callDirection  string
		enterMonthYear string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          ╔═════╦══════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgRed + textBoldWhite + "[7]" + resetColour + textBoldBlack + " ║ " + bgRed + textBoldWhite + "Delete a previsouly imported CDR from a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦════════════════════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")

	voipCarrierID = voipCarrierIDInput(dbDetail, "option7")
	callDirection = callDirectionInput(dbDetail, "option7")
	enterMonthYear = enterMonthYearInput(dbDetail, "option7")

	// Retrieve the column needed
	dbDetail.column = "cdr_month_year_column"
	dbDetail.table = "voip_carrier"
	dbDetail.columnWhere = "id"
	dbDetail.columnWhereValue = voipCarrierID
	monthYearColumn := selectWhere(dbDetail)

	// Delete CDR
	dbDetail.connection.Exec("DELETE FROM cdr_importer."+callDirection+"_cdr_"+voipCarrierID+" WHERE "+monthYearColumn+" = ?;", enterMonthYear)

	// Inform the user the CDR was deleted
	messageBox("If CDR exsisted it has been deleted", bgGreen)
	returnToMainMenu()
}

// Option 8 function
// Delete an existing VoIP carrier, all associated CDRs and inserted into YAP logs
func option8(dbDetail databaseFunctionParameter) {
	returnToMainMenu()
}

// Option 9 function
// Insert a re-rated CDR for a VoIP carrier from a specific month into YAP
func option9(dbDetail databaseFunctionParameter) {

	var (
		voipCarrierID                                   string
		callDirection                                   string
		enterMonthYear                                  string
		perSecondBilling                                string
		yapCustomerID                                   string
		cdrTag                                          string
		yapServiceProductName                           string
		yapInvoiceItemSalesTaxRate                      string
		yapInvoiceItemSalesTaxStatus                    string
		totalNoMinimumChargeForCallDuration             string
		totalMinimumChargeForCallDurationBelow60Seconds string
	)

	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("          ╔═════╦═════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("          ║ " + bgGreen + textBoldWhite + "[9]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Insert a re-rated CDR for a VoIP carrier from a specific month into YAP" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╔════╩═════╩═════════════════════════════╦═══════════════════════════════════════════╝")
	fmt.Println("     ║ Type \"menu\" to return to the main menu ║")
	fmt.Println("     ╚════════════════════════════════════════╝")

	voipCarrierID = voipCarrierIDInput(dbDetail, "option9")
	callDirection = callDirectionInput(dbDetail, "option9")
	enterMonthYear = enterMonthYearInput(dbDetail, "option9")

	yesNoList := yesNoSlice()
	fmt.Println("")
	fmt.Println("     Bill calls with 60/1 billing - if the input is yes the first 60 seconds in a call are always billed as 1 full minute even if the call was less than 60 seconds")
	fmt.Print("     if the input is no on this option all seconds in a call are billed per second even if a call is below 1 minute (60 seconds) [Valid options " + strings.Join(yesNoList, ", ") + "]: ")
	fmt.Scan(&perSecondBilling)
	// Return to main menu if menu is input
	mainMenu(perSecondBilling)

	// Check cdrIgnoreFirstLine is contained in the slice
	validatePerSecondBilling := slices.Contains(yesNoList, perSecondBilling)

	if validatePerSecondBilling == false {
		// Invalid input message displays to the user
		messageBox("Invalid option for 60/1 billing, please re-enter with either"+strings.Join(yesNoList, ", "), bgYellow)
		fmt.Print(textBoldBlack)
		fmt.Print("     Press the enter/return key to continue ")
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option9(dbDetail)
		}
	}

	option9SQL, err := dbDetail.connection.Query(`SELECT
                                                        yap_customer_id,
                                                        cdr_tag,
                                                        yap_service_product_name,
                                                        yap_invoice_item_sales_tax_rate,
                                                        yap_invoice_item_sales_tax_status,
                                                        total_no_minimum_charge_for_call_duration,
                                                        total_minimum_charge_for_call_duration_below_60_seconds
                                                      FROM
                                                        cdr_importer.view___`+callDirection+`_total_cdr_`+voipCarrierID+`
                                                      WHERE
                                                        cdr_month_year = ?;`, enterMonthYear)

	// Error
	if err != nil {
		// Invalid input message displays to the user
		messageBox("The re-rated CDR does not exist", bgYellow)
		fmt.Print("     Press the enter/return key to continue ", err)
		fmt.Print(resetColour)
		var enter string
		fmt.Scanln(&enter)
		if enter == "" || enter != "" {
			option9(dbDetail)
		}
	}

	for option9SQL.Next() {

		err = option9SQL.Scan(
			&yapCustomerID,
			&cdrTag,
			&yapServiceProductName,
			&yapInvoiceItemSalesTaxRate,
			&yapInvoiceItemSalesTaxStatus,
			&totalNoMinimumChargeForCallDuration,
			&totalMinimumChargeForCallDurationBelow60Seconds,
		)

		// Error
		if err != nil {
			panic(err)
		}

		if yapCustomerID == "NO CUSTOMER" {
			//Do Nothing
		} else {

			var serviceProductName string

			if yapServiceProductName == "Inbound Minutes Tag ✎" {
				serviceProductName = "Total Inbound Minutes ✆"
			} else if yapServiceProductName == "Outbound Minutes Tag ✎" {
				serviceProductName = "Total Outbound Minutes ✆"
			}

			var sellPrice string

			if perSecondBilling == "yes" || perSecondBilling == "yes" || perSecondBilling == "Yes" || perSecondBilling == "YES" || perSecondBilling == "y" || perSecondBilling == "Y" {
				sellPrice = totalMinimumChargeForCallDurationBelow60Seconds
			} else if perSecondBilling == "no" || perSecondBilling == "no" || perSecondBilling == "No" || perSecondBilling == "NO" || perSecondBilling == "n" || perSecondBilling == "N" {
				sellPrice = totalNoMinimumChargeForCallDuration
			}

			// Insert re-rated CDRs into YAP
			dbDetail.connection.Exec("INSERT INTO yap.invoice_item (customer_id, pbx_id, tag, service_product_name, sell_price, sales_tax_rate, sales_tax_status, bill_item_once, item_on_hold) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);", yapCustomerID, 1, enterMonthYear+" "+callDirection+" minutes for "+cdrTag, serviceProductName, sellPrice, yapInvoiceItemSalesTaxRate, yapInvoiceItemSalesTaxStatus, "yes", "no")
		}
	}
	// Add log to the yap_cdr_insert_log table
	dbDetail.connection.Exec("INSERT INTO cdr_importer.yap_cdr_insert_log (voip_carrier_id, cdr_month_year) VALUES(?, ?);", voipCarrierID, enterMonthYear)

	// Inform the user the re-rated CDR was inserted into YAP
	messageBox("Re-rated CDR has been inserted into YAP", bgGreen)
	returnToMainMenu()
}

//----------------------------------------------------------------------------------------------------

func cdrimporter() {
	clearScreen()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("     " + bgCyan + textBoldWhite + " ╔═╗ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ╔═╗ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ╚═╝ ┃                                                                                                       ┃ ╚═╝ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┏━━━┛                                                                                                       ┗━━━┓ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃      ██████╗ ██████╗  ██████╗     ██╗ ███╗   ███╗ ██████╗   ██████╗  ██████╗  ████████╗ ███████╗ ██████╗      ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃     ██╔════╝ ██╔══██╗ ██╔══██╗    ██║ ████╗ ████║ ██╔══██╗ ██╔═══██╗ ██╔══██╗ ╚══██╔══╝ ██╔════╝ ██╔══██╗     ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃     ██║      ██║  ██║ ██████╔╝    ██║ ██╔████╔██║ ██████╔╝ ██║   ██║ ██████╔╝    ██║    █████╗   ██████╔╝     ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃     ██║      ██║  ██║ ██╔══██╗    ██║ ██║╚██╔╝██║ ██╔═══╝  ██║   ██║ ██╔══██╗    ██║    ██╔══╝   ██╔══██╗     ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃     ╚██████╗ ██████╔╝ ██║  ██║    ██║ ██║ ╚═╝ ██║ ██║      ╚██████╔╝ ██║  ██║    ██║    ███████╗ ██║  ██║     ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃      ╚═════╝ ╚═════╝  ╚═╝  ╚═╝    ╚═╝ ╚═╝     ╚═╝ ╚═╝       ╚═════╝  ╚═╝  ╚═╝    ╚═╝    ╚══════╝ ╚═╝  ╚═╝     ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃                                                                                                               ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃           Source code for CDR Importer available at https://github.com/yet-another-pbx/cdr-importer           ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃                                                                                                               ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃                            " + resetColour + bgRed + "                                                     " + textBoldWhite + resetColour + bgCyan + textBoldWhite + "                              ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┃                            " + resetColour + bgRed + textBoldWhite + "   Type \"exit\" or \"quit\" to terminate CDR Importer   " + resetColour + bgCyan + textBoldWhite + "                              ┃ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ┗━━━┓                        " + resetColour + bgRed + textBoldWhite + "                                                     " + resetColour + bgCyan + textBoldWhite + "                          ┏━━━┛ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ╔═╗ ┃                                                                                                       ┃ ╔═╗ " + resetColour)
	fmt.Println("     " + bgCyan + textBoldWhite + " ╚═╝ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ╚═╝ " + resetColour)

	fmt.Print(textBoldBlack)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("     ╔═════╦═════════════════════════════════════╗")
	fmt.Println("     ║ " + bgBlue + textBoldWhite + "[0]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List all VoIP carrier IDs and names" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬═════════════════════════════════════╩══════════════════════════╗")
	fmt.Println("     ║ " + bgBlue + textBoldWhite + "[1]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List a CDR from a particular month and year for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬════════════════════════════════════════╦═══════════════════════╝")
	fmt.Println("     ║ " + bgBlue + textBoldWhite + "[2]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "List all call rates for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬════════════════════════════════════════╩═════════════════════════════════╗")
	fmt.Println("     ║ " + bgBlue + textBoldWhite + "[3]" + resetColour + textBoldBlack + " ║ " + bgBlue + textBoldWhite + "Show log of re-rated CDRs previously inserted into YAP (Yet Another PBX)" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬═══════════════════════════════════════════════════════════════════╦══════╝")
	fmt.Println("     ║ " + bgGreen + textBoldWhite + "[4]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Add a new VoIP carrier (must have a call rate card in CSV format)" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬═══════════════════════════════════════════════════════════════╦═══╝")
	fmt.Println("     ║ " + bgPurple + textBoldWhite + "[5]" + resetColour + textBoldBlack + " ║ " + bgPurple + textBoldWhite + "Replace or edit an exisitng call rate card for a VoIP carrier" + resetColour + textBoldBlack + " ║")
	fmt.Println("     ╠═════╬════════════════════════════════════════════════╦══════════════╝")
	fmt.Println("     ║ " + bgGreen + textBoldWhite + "[6]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Import a new CDR into an existing VoIP carrier" + resetColour + " ║")
	fmt.Println("     ╠═════╬════════════════════════════════════════════════╩═════╗")
	fmt.Println("     ║ " + bgRed + textBoldWhite + "[7]" + resetColour + textBoldBlack + " ║ " + bgRed + textBoldWhite + "Delete a previsouly imported CDR from a VoIP carrier" + resetColour + " ║")
	fmt.Println("     ╠═════╬══════════════════════════════════════════════════════╩══════════════════════════╗")
	fmt.Println("     ║ " + bgRed + textBoldWhite + "[8]" + resetColour + textBoldBlack + " ║ " + bgRed + textBoldWhite + "Delete an existing VoIP carrier, all associated CDRs and inserted into YAP logs" + resetColour + " ║")
	fmt.Println("     ╠═════╬═════════════════════════════════════════════════════════════════════════╦═══════╝")
	fmt.Println("     ║ " + bgGreen + textBoldWhite + "[9]" + resetColour + textBoldBlack + " ║ " + bgGreen + textBoldWhite + "Insert a re-rated CDR for a VoIP carrier from a specific month into YAP" + resetColour + " ║")
	fmt.Println("     ╚═════╩═════════════════════════════════════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Print("     " + textBoldBlack + "Select an option [0-9]: ")
	var option string
	fmt.Scanln(&option)
	fmt.Println(resetColour)

	// Values allowed for option
	var optionList = []string{"", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "exit", "Exit", "EXIT", "quit", "Quit", "QUIT"}
	validOption := slices.Contains(optionList, option)

	// Conditional statment to determine what happens when an option is input
	if validOption == false {
		messageBox("Invalid option - enter option [0-9] or exit", bgYellow)
		returnToMainMenu()
	}

	// If user typed exit or quit then stop program
	exit(option)

	// If the user pressed the enter/return key then re-run the main function
	if option == "" {
		cdrimporter()
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

	dbDetail.connection = dbConnection
	dbDetail.username = dbUsername
	dbDetail.password = dbPassword
	dbDetail.database = dbName
	dbDetail.address = dbAddress
	dbDetail.port = dbPort
	dbDetail.transport = dbTransport
	dbDetail.tls = dbTLS

	if option == "0" {
		option0(dbDetail)
	} else if option == "1" {
		option1(dbDetail)
	} else if option == "2" {
		option2(dbDetail)
	} else if option == "3" {
		option3(dbDetail)
	} else if option == "4" {
		option4(dbDetail)
	} else if option == "5" {
		option5(dbDetail)
	} else if option == "6" {
		option6(dbDetail)
	} else if option == "7" {
		option7(dbDetail)
	} else if option == "8" {
		option8(dbDetail)
	} else if option == "9" {
		option9(dbDetail)
	} else {
		messageBox("Invalid option - enter option [0-9] or exit", bgYellow)
		returnToMainMenu()
	}
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("Operating system must be GNU/Linux to work")
	} else {
		cdrimporter()
	}
}

// Contributor(s):
// Elliot Michael Keavney
