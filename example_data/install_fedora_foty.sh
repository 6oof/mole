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
systemctl daemon-reload
systemctl enable --now caddy

echo -e "\n\033[0;32m### Step 4: Check Caddy Status ###\033[0m"
if systemctl is-active --quiet caddy; then
    echo "Caddy is running."
else
    echo "Caddy is not running or failed to start."
fi

# 3. Install Podman
echo -e "\n\033[0;32m### Step 5: Install Podman ###\033[0m"
dnf -y install podman

echo -e "\n\033[0;32m### Step 6: Verify the Installation ###\033[0m"
podman --version

# 4. Setting Up a Firewall (firewalld)
echo -e "\n\033[0;32m### Step 7: Install firewalld ###\033[0m"
dnf install -y firewalld

echo -e "\n\033[0;32m### Step 8: Start and Enable the Firewall ###\033[0m"
systemctl start firewalld
systemctl enable firewalld

echo -e "\n\033[0;32m### Step 9: Configure Allowed Ports ###\033[0m"
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https

echo -e "\n\033[0;32m### Step 10: Reload Firewall Rules ###\033[0m"
firewall-cmd --reload

echo -e "\n\033[0;32m### Step 11: Verify Active Firewall Rules ###\033[0m"
firewall-cmd --list-all

# 5. User creation and permission setup for Caddy

echo -e "\n\033[0;32m### Step 12: Create user 'mole' for managing Caddy ###\033[0m"
useradd -m -s /bin/bash mole

# 6. Copy SSH keys to the 'mole' user
echo -e "\n\033[0;32m### Step 13: Copy SSH keys to the 'mole' user ###\033[0m"

# Create the .ssh directory for 'mole' and set correct permissions
mkdir -p /home/mole/.ssh
chmod 700 /home/mole/.ssh

# Copy the SSH keys (from root or another user, depending on where they're stored)
cp /root/.ssh/authorized_keys /home/mole/.ssh/authorized_keys

# Set correct ownership and permissions for the authorized_keys file
chown mole:mole /home/mole/.ssh/authorized_keys
chmod 600 /home/mole/.ssh/authorized_keys

# 7. Grant 'mole' permission to manage the Caddy service without password
echo -e "\n\033[0;32m### Step 14: Grant mole permission to manage Caddy service ###\033[0m"
echo "mole ALL=(root) NOPASSWD: /bin/systemctl start caddy, /bin/systemctl stop caddy, /bin/systemctl restart caddy, /bin/systemctl reload caddy, /bin/systemctl status caddy" | sudo tee /etc/sudoers.d/caddy-management

# 8. Grant 'mole' access to read Caddy logs
echo -e "\n\033[0;32m### Step 15: Add mole to the systemd-journal group for reading logs ###\033[0m"
usermod -aG systemd-journal mole

echo -e "\n\033[0;32m### We're done! The user 'mole' can now manage the Caddy service, read logs, and use SSH. ###\033[0m"
