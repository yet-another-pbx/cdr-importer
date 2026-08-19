#!/bin/bash

# Uninstall Script for CDR Importer

#----------------------------------------------------------------------

# Clear Screen
clear_screen="\033[H\033[2J";

# American National Standards Institute (ANSI) reset colour code
reset_colour="\033[0m";

# American National Standards Institute (ANSI) text colour code
text_bold_black="\033[1;30m";
text_bold_white="\033[1;37m";

# American National Standards Institute (ANSI) background colour codes
bg_green="\033[42m";
bg_yellow="\033[43m";
bg_purple="\033[45m";

#----------------------------------------------------------------------

# Check the user is root otherwise exit the script

if [ "$EUID" -ne 0 ]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔════════════════════╗ \n";
  printf " ║ Please run as root ║ \n";
  printf " ╚════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# Remove the CDR Importer binary

printf $clear_screen;
printf $bg_purple;
printf $text_bold_white;
printf " ╔═══════════════════════════════════════════════╗ \n";
printf " ║ Remove the CDR Importer binary from /usr/bin? ║ \n";
printf " ╚═══════════════════════════════════════════════╝ \n";
printf $reset_colour;
printf "\n";
printf $text_bold_black;
read -p "   Would you like to continue? [Yes/No]: " response;
if [ $response == "Yes" ] || [ $response == "yes" ] || [ $response == "YES" ] || [ $response == "Y" ] || [ $response == "y" ]
then
  printf $reset_colour;
  rm /usr/bin/cdrimporter;
  printf $clear_screen;
  printf $bg_green;
  printf $text_bold_white;
  printf " ╔════════════════════════════════════════════════════════╗ \n";
  printf " ║ The CDR Importer binary has been removed from /usr/bin ║ \n";
  printf " ╚════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
else
  printf $reset_colour;
  printf $clear_screen;
  printf $bg_green;
  printf $text_bold_white;
  printf " ╔════════════════════════════════════════════════════════════╗ \n";
  printf " ║ Uninstall script exited, CDR Importer has not been removed ║ \n";
  printf " ╚════════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;
