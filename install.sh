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

# Copy the SSH authorized_keys file to the 'mole' user
cp /root/.ssh/authorized_keys /home/mole/.ssh/authorized_keys

# Append a comment below the keys in the 'mole' user's authorized_keys file
echo -e "\n# Keys above were copied from the root user at install" >> /home/mole/.ssh/authorized_keys

# Set correct ownership and permissions for the authorized_keys file
chown mole:mole /home/mole/.ssh/authorized_keys
chmod 600 /home/mole/.ssh/authorized_keys

# 3. Install the Mole binary
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

# 4. Configure .bashrc to execute 'mole' on login
echo -e "\n\033[0;32m### Step 5: Configure .bashrc for the 'mole' user ###\033[0m"

echo -e "\n# Automatically launch the Mole CLI on login" >> /home/mole/.bashrc
echo "$MOLE_INSTALL_PATH" >> /home/mole/.bashrc

# Ensure proper ownership of .bashrc
chown mole:mole /home/mole/.bashrc

echo -e "\n\033[0;32m### Installation complete. Mole is installed at ${MOLE_INSTALL_PATH}. ###\033[0m"

