#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN="\033[0;32m"
CYAN="\033[0;36m"
RESET="\033[0m"

echo -e "${CYAN}Starting setup script...${RESET}"

# Check if the script is being run as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${CYAN}Please run as root.${RESET}"
    exit 1
fi

# Prompt for email address
read -p "Enter your email for domain setup (used for SSL certificate alerts): " user_email

# Validate email input
if [[ ! "$user_email" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
    echo -e "${CYAN}Invalid email format. Please run the script again with a valid email.${RESET}"
    exit 1
fi

# Step 1: Install Mole
echo -e "${CYAN}Installing Mole...${RESET}"
curl -O https://raw.githubusercontent.com/zulubit/mole/main/install.sh
chmod +x install.sh
./install.sh
rm install.sh

echo -e "${GREEN}Mole installation complete.${RESET}"

# Switch to Mole user and setup domains
echo -e "${CYAN}Setting up Mole domains...${RESET}"
su - mole -c "mole domains setup $user_email"

echo -e "${GREEN}Mole domains setup complete.${RESET}"

# Step 2: Install Caddy
echo -e "${CYAN}Installing Caddy...${RESET}"
dnf install 'dnf-command(copr)' -y
dnf copr enable @caddy/caddy -y
dnf install caddy -y

# Set permissions
echo -e "${CYAN}Configuring permissions for Caddy...${RESET}"
usermod -aG mole caddy
chmod 750 /home/mole

# Enable and start Caddy API service
echo -e "${CYAN}Starting Caddy API service...${RESET}"
systemctl daemon-reload
systemctl enable --now caddy-api

echo -e "${GREEN}Caddy installation complete.${RESET}"

# Step 3: Install Podman
echo -e "${CYAN}Installing Podman...${RESET}"
dnf copr enable rhcontainerbot/podman-next -y
dnf install podman -y

# Verify Podman installation
echo -e "${CYAN}Verifying Podman installation...${RESET}"
podman --version

echo -e "${GREEN}Podman installation complete.${RESET}"

# Step 4: Configure Firewall
echo -e "${CYAN}Setting up Firewall...${RESET}"
dnf install firewalld -y
systemctl start firewalld
systemctl enable firewalld

firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload

echo -e "${GREEN}Firewall configuration complete.${RESET}"

# Step 5: Install Git
echo -e "${CYAN}Installing Git...${RESET}"
dnf install git -y

echo -e "${GREEN}Git installation complete.${RESET}"

echo -e "${CYAN}Setup script completed successfully.${RESET}"
