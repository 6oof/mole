#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# 1. User creation and permission setup for Caddy
echo -e "\n\033[0;32m### Step 1: Create user 'mole' for managing Caddy ###\033[0m"
useradd -m -s /bin/bash mole

# 2. Copy SSH keys to the 'mole' user
echo -e "\n\033[0;32m### Step 2: Copy SSH keys to the 'mole' user ###\033[0m"

# Create the .ssh directory for 'mole' and set correct permissions
mkdir -p /home/mole/.ssh
chmod 700 /home/mole/.ssh
chown mole:mole /home/mole/.ssh

# Copy the SSH keys (from root or another user, depending on where they're stored)
cp /root/.ssh/authorized_keys /home/mole/.ssh/authorized_keys

# Set correct ownership and permissions for the authorized_keys file
chown mole:mole /home/mole/.ssh/authorized_keys
chmod 600 /home/mole/.ssh/authorized_keys

# 3. Set password for the 'mole' user
echo -e "\n\033[0;32m### Step 3: Set a password for user 'mole' ###\033[0m"
passwd mole

# 4. Grant 'mole' permission to manage the Caddy service without password
echo -e "\n\033[0;32m### Step 4: Grant mole permission to manage Caddy service ###\033[0m"
echo "mole ALL=(root) NOPASSWD: /bin/systemctl restart caddy, /bin/systemctl reload caddy, /bin/systemctl status caddy" | sudo tee /etc/sudoers.d/caddy-management

# 5. Grant 'mole' access to read Caddy logs
echo -e "\n\033[0;32m### Step 5: Add mole to the systemd-journal group for reading logs ###\033[0m"
usermod -aG systemd-journal mole

echo -e "\n\033[0;32m### We're done! The user 'mole' can now manage the Caddy service, read logs, and use SSH. ###\033[0m"
