#!/bin/bash

# Install Script for CDR Importer

#----------------------------------------------------------------------

# Make sure UTF-8 is used for install
export LANG=C.UTF-8;

# Clear Screen
clear_screen="\033[H\033[2J";

# American National Standards Institute (ANSI) reset colour code
reset_colour="\033[0m";

# American National Standards Institute (ANSI) text colour code
text_bold_black="\033[1;30m";
text_bold_white="\033[1;37m";

# American National Standards Institute (ANSI) background colour codes
bg_red="\033[41m";
bg_green="\033[42m";
bg_yellow="\033[43m";
bg_purple="\033[45m";

# Function to search and replace text in files
function string_update {
  sed -i "s,$search_string,$replace_string," $string_update_file;
};

#----------------------------------------------------------------------

# Check user is root otherwise exit script
if [[ "$EUID" -ne 0 ]]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔════════════════════════════════════════════════════╗ \n";
  printf " ║ Please run the CDR Importer install script as root ║ \n";
  printf " ╚════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# Check CDR Importer has been cloned from GitHub
if [[ ! -d "/root/cdrimporter" ]]
then
  printf $clear_screen;
  printf $bg_red;
  printf $text_bold_white;
  printf " ╔═════════════════════════════════════════════════════════════════════════════════════════════╗ \n";
  printf " ║ Directory cdrimporter does not exist in /root.                                              ║ \n";
  printf " ║ Please run commands: \"cd /root && git clone https://github.com/yet-another-pbx/cdrimporter\" ║ \n";
  printf " ║ and run the install script again.                                                           ║ \n";
  printf " ╚═════════════════════════════════════════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;
