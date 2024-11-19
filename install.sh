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

# 4. Grant 'mole' access to read Caddy logs
echo -e "\n\033[0;32m### Step 3: Add mole to the systemd-journal group for reading logs ###\033[0m"
usermod -aG systemd-journal mole

# 5. Install the Mole binary
echo -e "\n\033[0;32m### Step 4: Install the Mole binary ###\033[0m"

MOLE_VERSION="0.0.1" # Update with the correct version
MOLE_BINARY_URL="https://github.com/zulubit/mole/releases/download/${MOLE_VERSION}/mole"
MOLE_INSTALL_PATH="/usr/local/bin/mole"

echo -e "Downloading Mole binary from ${MOLE_BINARY_URL}..."
curl -L -o "$MOLE_INSTALL_PATH" "$MOLE_BINARY_URL"

echo "Setting executable permissions on the binary..."
chmod +x "$MOLE_INSTALL_PATH"

echo "Changing ownership of the binary to 'mole' user..."
chown mole:mole "$MOLE_INSTALL_PATH"

echo -e "\n\033[0;32m### Installation complete. Mole is installed at ${MOLE_INSTALL_PATH}. ###\033[0m"
