#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# 1. SSH into your newly created VPS as a root user

# 2. Setting Up Caddy on Your VPS

echo -e "\n\033[0;32m### Step 1: Install Caddy - reverse proxy to manage your domains ###\033[0m"
dnf install -y 'dnf-command(copr)'
dnf copr enable -y @caddy/caddy
dnf install -y caddy

echo -e "\n\033[0;32m### Step 2: Configure Caddy ###\033[0m"

echo -e "\n\033[0;32m### Step 3: Enable and Start Caddy Service ###\033[0m"
sudo systemctl daemon-reload
sudo systemctl enable --now caddy

echo -e "\n\033[0;32m### Step 4: Check Caddy Status ###\033[0m"
if systemctl is-active --quiet caddy; then
    echo "Caddy is running."
else
    echo "Caddy is not running or failed to start."
fi

# 3. Install Podman
echo -e "\n\033[0;32m### Step 5: Install Podman ###\033[0m"
sudo dnf -y install podman

echo -e "\n\033[0;32m### Step 6: Verify the Installation ###\033[0m"
podman --version

# 4. Setting Up a Firewall (firewalld)
echo -e "\n\033[0;32m### Step 7: Install firewalld ###\033[0m"
sudo dnf install -y firewalld

echo -e "\n\033[0;32m### Step 8: Start and Enable the Firewall ###\033[0m"
sudo systemctl start firewalld
sudo systemctl enable firewalld

echo -e "\n\033[0;32m### Step 9: Configure Allowed Ports ###\033[0m"
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https

echo -e "\n\033[0;32m### Step 10: Reload Firewall Rules ###\033[0m"
sudo firewall-cmd --reload

echo -e "\n\033[0;32m### Step 11: Verify Active Firewall Rules ###\033[0m"
sudo firewall-cmd --list-all

echo -e "\n\033[0;32m### We're done! ###\033[0m"
