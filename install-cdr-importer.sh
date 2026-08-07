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
if [[ ! -d "/root/cdr-importer" ]]
then
  printf $clear_screen;
  printf $bg_red;
  printf $text_bold_white;
  printf " ╔══════════════════════════════════════════════════════════════════════════════════════════════╗ \n";
  printf " ║ Directory cdrimporter does not exist in /root.                                               ║ \n";
  printf " ║ Please run commands: \"cd /root && git clone https://github.com/yet-another-pbx/cdr-importer\" ║ \n";
  printf " ║ and run the install script again.                                                            ║ \n";
  printf " ╚══════════════════════════════════════════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  exit;
fi;

#----------------------------------------------------------------------

# CDR Importer Install Title
printf $clear_screen;
printf "\n";
printf "          ██████╗ ██████╗  ██████╗     ██╗ ███╗   ███╗ ██████╗   ██████╗  ██████╗  ████████╗ ███████╗ ██████╗\n";
printf "         ██╔════╝ ██╔══██╗ ██╔══██╗    ██║ ████╗ ████║ ██╔══██╗ ██╔═══██╗ ██╔══██╗ ╚══██╔══╝ ██╔════╝ ██╔══██╗\n";
printf "         ██║      ██║  ██║ ██████╔╝    ██║ ██╔████╔██║ ██████╔╝ ██║   ██║ ██████╔╝    ██║    █████╗   ██████╔╝\n";
printf "         ██║      ██║  ██║ ██╔══██╗    ██║ ██║╚██╔╝██║ ██╔═══╝  ██║   ██║ ██╔══██╗    ██║    ██╔══╝   ██╔══██╗\n";
printf "         ╚██████╗ ██████╔╝ ██║  ██║    ██║ ██║ ╚═╝ ██║ ██║      ╚██████╔╝ ██║  ██║    ██║    ███████╗ ██║  ██║\n";
printf "          ╚═════╝ ╚═════╝  ╚═╝  ╚═╝    ╚═╝ ╚═╝     ╚═╝ ╚═╝       ╚═════╝  ╚═╝  ╚═╝    ╚═╝    ╚══════╝ ╚═╝  ╚═╝\n";
printf "";
printf "   ██╗ ███╗   ██╗ ███████╗ ████████╗  █████╗  ██╗      ██╗         ███████╗  ██████╗ ██████╗  ██╗ ██████╗  ████████╗\n";
printf "   ██║ ████╗  ██║ ██╔════╝ ╚══██╔══╝ ██╔══██╗ ██║      ██║         ██╔════╝ ██╔════╝ ██╔══██╗ ██║ ██╔══██╗ ╚══██╔══╝\n";
printf "   ██║ ██╔██╗ ██║ ███████╗    ██║    ███████║ ██║      ██║         ███████╗ ██║      ██████╔╝ ██║ ██████╔╝    ██║   \n";
printf "   ██║ ██║╚██╗██║ ╚════██║    ██║    ██╔══██║ ██║      ██║         ╚════██║ ██║      ██╔══██╗ ██║ ██╔═══╝     ██║   \n";
printf "   ██║ ██║ ╚████║ ███████║    ██║    ██║  ██║ ███████╗ ███████╗    ███████║ ╚██████╗ ██║  ██║ ██║ ██║         ██║   \n";
printf "   ╚═╝ ╚═╝  ╚═══╝ ╚══════╝    ╚═╝    ╚═╝  ╚═╝ ╚══════╝ ╚══════╝    ╚══════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═╝ ╚═╝         ╚═╝   \n";
printf "\n";
printf "  $bg_purple $text_bold_white╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════╗ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ To install and use the CDR Importer the following prerequisites must be completed in advance of the install: ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                              ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ - YAP (Yet Another PBX) must be installed                                                                    ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║   (To install YAP - https://github.com/yet-another-pbx/yap)                                                  ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                              ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ - The Go programming language (minimum version 1.26.0) must be installed                                     ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║   (To install the Go programming language - https://go.dev/doc/install)                                      ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║   (The YAP install script also automatically installs the Go programming langauge)                           ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                              ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                             - TYPE EXIT TO STOP THE CDR IMPORTER INSTALL SCRIPT -                            ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝ $reset_colour\n";
printf "\n";
