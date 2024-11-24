#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN="\033[0;32m"
CYAN="\033[0;36m"
RED="\033[0;31m"
RESET="\033[0m"

handle_error() {
    echo -e "${RED}FAILURE! $1${RESET}"
    exit 1
}

echo -e "${CYAN}Starting setup script...${RESET}"

# Check if the script is being run as root
if [ "$EUID" -ne 0 ]; then
    handle_error "Please run the script as root."
fi

# Prompt for email address
read -p "Enter your email for domain setup (used for SSL certificate alerts): " user_email

# Validate email input
if [[ ! "$user_email" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
    handle_error "Invalid email format. Please run the script again with a valid email."
fi

# Step 1: Install Mole
echo -e "${CYAN}Installing Mole...${RESET}"
curl -O https://raw.githubusercontent.com/zulubit/mole/main/install.sh || handle_error "Failed to download the Mole install script."
chmod +x install.sh || handle_error "Failed to make the install script executable."
./install.sh || handle_error "Mole installation failed."
rm install.sh

echo -e "${GREEN}Mole installation complete.${RESET}"

# Switch to Mole user and setup domains
echo -e "${CYAN}Setting up Mole domains...${RESET}"
su - mole -c "mole domains setup $user_email" || handle_error "Failed to set up Mole domains."

echo -e "${GREEN}Mole domains setup complete.${RESET}"

# Step 2: Install Caddy
echo -e "${CYAN}Installing Caddy...${RESET}"
dnf install 'dnf-command(copr)' -y || handle_error "Failed to install COPR plugin for dnf."
dnf copr enable @caddy/caddy -y || handle_error "Failed to enable Caddy repository."
dnf install caddy -y || handle_error "Failed to install Caddy."

# Set permissions
echo -e "${CYAN}Configuring permissions for Caddy...${RESET}"
usermod -aG mole caddy || handle_error "Failed to add Mole user to the Caddy group."
chmod 750 /home/mole || handle_error "Failed to set permissions on the Mole home directory."

# Enable and start Caddy API service
echo -e "${CYAN}Starting Caddy API service...${RESET}"
systemctl daemon-reload || handle_error "Failed to reload systemd daemon."
systemctl enable --now caddy-api || handle_error "Failed to enable and start the Caddy API service."

echo -e "${GREEN}Caddy installation complete.${RESET}"

# Step 3: Install Docker and Docker Compose
echo -e "${CYAN}Installing Docker...${RESET}"
dnf install -y dnf-plugins-core || handle_error "Failed to install dnf-plugins-core."
dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo || handle_error "Failed to add Docker repository."
dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin -y || handle_error "Failed to install Docker."

# Enable and start Docker
echo -e "${CYAN}Starting Docker...${RESET}"
systemctl start docker || handle_error "Failed to start Docker service."
systemctl enable docker || handle_error "Failed to enable Docker service."

# Add mole user to Docker group
echo -e "${CYAN}Adding mole user to Docker group...${RESET}"
usermod -aG docker mole || handle_error "Failed to add Mole user to Docker group."

# Verify Docker installation
echo -e "${CYAN}Verifying Docker installation...${RESET}"
docker --version || handle_error "Docker installation verification failed."

echo -e "${GREEN}Docker and Docker Compose installation complete.${RESET}"

# Step 4: Configure Firewall
echo -e "${CYAN}Setting up Firewall...${RESET}"
dnf install firewalld -y || handle_error "Failed to install firewalld."
systemctl start firewalld || handle_error "Failed to start firewalld."
systemctl enable firewalld || handle_error "Failed to enable firewalld."

firewall-cmd --permanent --add-service=ssh || handle_error "Failed to allow SSH in the firewall."
firewall-cmd --permanent --add-service=http || handle_error "Failed to allow HTTP in the firewall."
firewall-cmd --permanent --add-service=https || handle_error "Failed to allow HTTPS in the firewall."
firewall-cmd --reload || handle_error "Failed to reload firewall rules."

echo -e "${GREEN}Firewall configuration complete.${RESET}"

# Step 5: Install Git
echo -e "${CYAN}Installing Git...${RESET}"
dnf install git -y || handle_error "Failed to install Git."

echo -e "${GREEN}Git installation complete.${RESET}"

rm -rf install-rocky.sh

echo -e "${CYAN}Setup script completed successfully.${RESET}"

