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
printf "         ██████╗ ██████╗  ██████╗     ██╗ ███╗   ███╗ ██████╗   ██████╗  ██████╗  ████████╗ ███████╗ ██████╗\n";
printf "        ██╔════╝ ██╔══██╗ ██╔══██╗    ██║ ████╗ ████║ ██╔══██╗ ██╔═══██╗ ██╔══██╗ ╚══██╔══╝ ██╔════╝ ██╔══██╗\n";
printf "        ██║      ██║  ██║ ██████╔╝    ██║ ██╔████╔██║ ██████╔╝ ██║   ██║ ██████╔╝    ██║    █████╗   ██████╔╝\n";
printf "        ██║      ██║  ██║ ██╔══██╗    ██║ ██║╚██╔╝██║ ██╔═══╝  ██║   ██║ ██╔══██╗    ██║    ██╔══╝   ██╔══██╗\n";
printf "        ╚██████╗ ██████╔╝ ██║  ██║    ██║ ██║ ╚═╝ ██║ ██║      ╚██████╔╝ ██║  ██║    ██║    ███████╗ ██║  ██║\n";
printf "         ╚═════╝ ╚═════╝  ╚═╝  ╚═╝    ╚═╝ ╚═╝     ╚═╝ ╚═╝       ╚═════╝  ╚═╝  ╚═╝    ╚═╝    ╚══════╝ ╚═╝  ╚═╝\n";
printf "\n";
printf "  ██╗ ███╗   ██╗ ███████╗ ████████╗  █████╗  ██╗      ██╗         ███████╗  ██████╗ ██████╗  ██╗ ██████╗  ████████╗\n";
printf "  ██║ ████╗  ██║ ██╔════╝ ╚══██╔══╝ ██╔══██╗ ██║      ██║         ██╔════╝ ██╔════╝ ██╔══██╗ ██║ ██╔══██╗ ╚══██╔══╝\n";
printf "  ██║ ██╔██╗ ██║ ███████╗    ██║    ███████║ ██║      ██║         ███████╗ ██║      ██████╔╝ ██║ ██████╔╝    ██║   \n";
printf "  ██║ ██║╚██╗██║ ╚════██║    ██║    ██╔══██║ ██║      ██║         ╚════██║ ██║      ██╔══██╗ ██║ ██╔═══╝     ██║   \n";
printf "  ██║ ██║ ╚████║ ███████║    ██║    ██║  ██║ ███████╗ ███████╗    ███████║ ╚██████╗ ██║  ██║ ██║ ██║         ██║   \n";
printf "  ╚═╝ ╚═╝  ╚═══╝ ╚══════╝    ╚═╝    ╚═╝  ╚═╝ ╚══════╝ ╚══════╝    ╚══════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═╝ ╚═╝         ╚═╝   \n";
printf "\n";
printf "  $bg_purple $text_bold_white╔═══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ NOTE: This install script presumes you will be installing CDR Importer on a server with YAP (Yet Another PBX) already ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║       installed. If you are planning on using CDR Importer on a seprate computer please exit this install script and  ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║       follow the manual install guide at - https://github.com/yet-another-pbx/cdr-importer                            ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                                       ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ To install and use the CDR Importer the following prerequisites must be completed in advance of the install:          ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                                       ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ - YAP (Yet Another PBX) must be installed                                                                             ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║   (To install YAP - https://github.com/yet-another-pbx/yap)                                                           ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                                       ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║ - The Go programming language (minimum version 1.26.0) must be installed                                              ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║   (The YAP install script also automatically installs the Go programming langauge)                                    ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                                                                                                       ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white║                                 - TYPE EXIT TO STOP THE CDR IMPORTER INSTALL SCRIPT -                                 ║ $reset_colour\n";
printf "  $bg_purple $text_bold_white╚═══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝ $reset_colour\n";
printf "\n";
read -p "  Has all the prerequisites been completed? (yes/no): " prerequisites;
printf "\n";
if [[ $prerequisites = "exit" ]] || [[ $prerequisites = "Exit" ]] || [[ $prerequisites = "EXIT" ]]
then
  exit;
elif [[ $prerequisites = "" ]]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔══════════════════════════════════════╗ \n";
  printf " ║ Prerequisites answer cannot be empty ║ \n";
  printf " ║ Press return to continue             ║ \n";
  printf " ╚══════════════════════════════════════╝ \n";
  printf $reset_colour;
  read -p "";
  source ./install-cdr-importer.sh;
elif [[ $prerequisites != "yes" ]] && [[ $prerequisites != "Yes" ]] && [[ $prerequisites != "YES" ]] && [[ $prerequisites != "y" ]] && [[ $prerequisites != "Y" ]]
then
  printf $clear_screen;
  printf $bg_yellow;
  printf $text_bold_white;
  printf " ╔═════════════════════════════════════════════════════════════════════════╗ \n";
  printf " ║ Complete prerequisites first and re-run the CDR Importer install script ║ \n";
  printf " ║ Press return to continue                                                ║ \n";
  printf " ╚═════════════════════════════════════════════════════════════════════════╝ \n";
  printf $reset_colour;
  read -p "";
  exit;
fi;

#----------------------------------------------------------------------

# Generate strong passwords using the OpenSSL cryptographic libary
mariadb_cdr_importer_password=(`openssl rand -base64 40 | tr "/" a | tr "=" a | tr "+" a`);

#----------------------------------------------------------------------

# Create a system user named cdr-importer with no shell, no home directory and lock the account
useradd -r -s /bin/false cdr-importer;
usermod -L cdr-importer;

# Create Go directories in root home directory for compiling the source code
mkdir -p /root/go/{bin,pkg,src/cdrimporter};

# Copy CDR Importer source code
cp /root/cdr-importer/go/cdrimporter.go /root/go/src/cdrimporter/cdrimporter.go;

# Export Go
export GOPATH=/root/go;
export PATH=$PATH:/usr/local/go/bin;

# Remove old go.mod and create a Go mod file for CDR Importer
rm /root/go/src/cdrimporter/go.mod;
cd /root/go/src/cdrimporter;
go mod init root/go/src/cdrimporter;
go mod tidy;

# Compile cdrimporter.go
cd /root/go/src/cdrimporter;
go build cdrimporter.go;

# Create directory used for configuration file
mkdir /etc/cdr-importer;

# Copy CDR Importer configuration file
cp /root/cdr-importer/env/cdr-importer.env /etc/cdr-importer/cdr-importer.env;

# Add the MariaDB password to the CDR Importer configuration file
string_update_file="/etc/cdr-importer/cdr-importer.env";
search_string="<REPLACE_DB_CDR_IMPORTER_PASSWORD>";
replace_string="$mariadb_cdr_importer_password";
string_update;

# Change executables owner, group, file permissions and move executable
chown root:cdr-importer /root/go/src/cdrimporter/cdrimporter;
chmod 550 /root/go/src/cdrimporter/cdrimporter;
mv /root/go/src/cdrimporter/cdrimporter /usr/bin/cdrimporter;

# Change CDR Importer owner, group and file permissions
chown -R root:cdr-importer /etc/cdr-importer;
chmod 550 /etc/cdr-importer;
chmod 660 /etc/cdr-importer/cdr-importer.env;

# Change directroy to /root
cd /root;

#----------------------------------------------------------------------

# MariaDB database and tables setup

# Create database
mysql -u root -e "CREATE DATABASE cdr_importer;";
mysql -u root -e "FLUSH PRIVILEGES;";

# Drop any previous cdr-importer user and create a cdr-importer MariaDB user
mysql -u root -e "DROP USER IF EXISTS 'cdr-importer'@'localhost';";
mysql -u root -e "CREATE USER 'cdr-importer'@'localhost' IDENTIFIED BY '$mariadb_cdr_importer_password';";
mysql -u root -e "FLUSH PRIVILEGES;";

# Create cdr_importer tables
mysql -u root -D yap -e "SOURCE /root/cdrimporter/mariadb/cdr-importer.sql;";

# Grant privileges for cdr-importer MariaDB user
mysql -u root -e "GRANT SELECT, INSERT, UPDATE, CREATE, ALTER, REFERENCES, INDEX ON cdr_importer.* TO 'cdr-importer'@'localhost';";
mysql -u root -e "FLUSH PRIVILEGES;";

#----------------------------------------------------------------------

printf $bg_green;
printf $text_bold_white;
printf " ╔═════════════════════════════════════════════════════════════╗ \n";
printf " ║ CDR Importer has been installed in /usr/bin                 ║ \n";
printf " ║ To use CDR Importer type \"sudo -u cdr-importer cdrimporter\" ║ \n"; 
printf " ╚═════════════════════════════════════════════════════════════╝ \n";
printf $reset_colour;
exit;
