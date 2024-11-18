# Install Mole

## This is a detailed install script for Mole. This should Ideally be performed on a fresh server.

This guide should produce an enviroment that is minimalistic and reasonably secure by default.

### Prerequisites

<!--TODO: make sure you add git install isntructions-->
#### This quicstart guide assumes a freshly provisioned VPS with the following:

- Running a distribution wiht **systemd** (minimum version 240).
- SSH root access.
- Ability to install Caddy
- Ability to reload caddy systemd service without root
- Ability to isntall git
- Ability to isntall Podman

**Note**: **Rocky linux** is recommended for its stability and strong support for server environments as well as support for recent versions of **Podman**. If you'd like to use a different distribution, you may need to some adjustments to the steps below. notes on choosing a different distro can be found on the bottom of this document.

----

### 1. SSH into your newly created VPS as a root user

### 2. Install Mole

#### Step 1: Download and Run the Install Script

Download the installation script and execute it:

```bash
curl -O https://raw.githubusercontent.com/zulubit/mole/main/install.sh
chmod +x install.sh
sudo ./install.sh
```

Replace `<your-server-ip>` with your server's IP address.

#### Step 2: Verify the Installation

Once logged in as the `mole` user, check that the `mole` CLI is properly installed and accessible:

```bash
mole version
```

You should see the version number of Mole displayed.

### 3. Install Caddy

#### Step 1: Install Caddy - Reverse Proxy to Manage Your Domains

Run the following commands to install Caddy and its dependencies:

```bash
dnf install 'dnf-command(copr)'
dnf copr enable @caddy/caddy
dnf install caddy
```

#### Step 2: Setup domains

this step is required for Caddy not to fail on start. Make sure you don't setup your domains with root user to avoid permission errors

```bash
su mole
```

```bash
mole domains setup your@email.com
```

```bash
exit
```

this domain is used for ssl certificate alerts

#### Step 3: Setup permissions

```bash
chown -R caddy:caddy /home/mole/domains
chown -R caddy:caddy /home/mole/caddy
```

#### Step 4: Configure Caddy

1. **Edit the Caddyfile**: Open the main configuration file located at `/etc/caddy/Caddyfile` with your preferred text editor.

   ```bash
   sudo vi /etc/caddy/Caddyfile
   ```

2. **Set Up the Caddyfile**: Replace the contents with the following:

   ```caddyfile
   import /etc/caddy/setup/main.caddy
   ```

#### Step 5: Enable and Start Caddy Service

Reload the systemd daemon and enable the Caddy service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now caddy
```

#### Step 6: Check Caddy Status

To verify that Caddy is running:

```bash
sudo systemctl status caddy --no-pager
```

**You should see output indicating that the Caddy service is active and running.**

#### Step 5: Allow User to Reload Caddy Without a Password Prompt

If you need to allow a non-root user to reload the Caddy service without a password prompt, follow these steps:

1. **Create the `caddygroup`**:
   Create a new group for Caddy:

   ```bash
   sudo groupadd caddygroup
   ```

2. **Create or Edit the `polkit` Rule**:
   Create a new rule in the `/etc/polkit-1/rules.d/` directory to allow users in the `caddygroup` to reload the Caddy service without authentication.

   ```bash
   sudo vi /etc/polkit-1/rules.d/99-caddygroup.rules
   ```

3. **Add the Rule for `systemd` Permissions**:
   Add the following JavaScript code to this file:

   ```js
   polkit.addRule(function(action, subject) {
       if (action.id == "org.freedesktop.systemd1.manage-units" &&
           action.lookup("unit") == "caddy.service" &&
           subject.isInGroup("caddygroup")) {
               return polkit.Result.YES;
       }
   });
   ```

4. **Ensure the User is in the Group**:
   Make sure your user is part of the `caddygroup`. You can check this with:

   ```bash
   sudo usermod -aG caddygroup mole
   sudo usermod -aG caddygroup caddy
   ```

   Then log out and log back in to apply the group membership.

5. **Restart the `polkit` Service**:
   After creating the rule, restart the `polkit` service for the changes to take effect:

   ```bash
   sudo systemctl restart polkit
   ```

#### Additional Resources

- [Caddy RHEL Installation Documentation](https://caddyserver.com/docs/install#fedora-redhat-centos)

### 4. Install Podman

#### Step 1: Install Podman

Run the following command to install Podman:

```bash
sudo dnf copr enable rhcontainerbot/podman-next -y
sudo dnf install podman
```

**NOTE: this uses the copr (testing) version of podman to get all features. If this is too bleeding edge for you, only run the second command from the two provided above.**

#### Step 2: Verify the Installation

Check that Podman is installed correctly:

```bash
podman --version
```

### 5. Setting Up a Firewall (firewalld)

To ensure your server is secure, you can set up a firewall using **firewalld**. This will restrict incoming traffic to SSH (port 22), HTTP (port 80), and HTTPS (port 443).

#### Step 1: Install firewalld

```bash
sudo dnf install firewalld
```

#### Step 2: Start and Enable the Firewall

```bash
sudo systemctl start firewalld
sudo systemctl enable firewalld
```

#### Step 3: Configure Allowed Ports

1. Allow SSH (port 22):

   ```bash
   sudo firewall-cmd --permanent --add-service=ssh
   ```

2. Allow HTTP (port 80):

   ```bash
   sudo firewall-cmd --permanent --add-service=http
   ```

3. Allow HTTPS (port 443):

   ```bash
   sudo firewall-cmd --permanent --add-service=https
   ```

#### Step 4: Reload Firewall Rules

```bash
sudo firewall-cmd --reload
```

#### Step 5: Verify Active Firewall Rules

To check if the correct rules are in place:

```bash
sudo firewall-cmd --list-all
```

This should show that SSH, HTTP, and HTTPS (next to services:) are allowed for incoming traffic.

## Choosing a different distribution

This guide is optimized for **Rocky Linux** and other **RHEL-based** distros (like **CentOS**, **AlmaLinux**, and **Fedora**). If you'd like to use a different distribution, here are the adjustments you’ll need to make:

### Adjustments for Ubuntu/Debian-based Distros:

- **Package Manager**: Replace `dnf` with `apt` in the installation commands:
  - For example, use `sudo apt install caddy` instead of `dnf install caddy`.
  
- **Firewall**: If using **UFW** instead of `firewalld`, replace `firewall-cmd` commands with UFW commands:
  - `sudo ufw allow ssh`
  - `sudo ufw allow http`
  - `sudo ufw allow https`

- **Podman**: Use the `apt` equivalent for installing Podman:
  - `sudo apt install podman` (or follow [Podman installation for Ubuntu](https://podman.io/getting-started/installation)).

- **Polkit**: On **Ubuntu/Debian-based** systems, **polkit** is typically used differently. You may need to manage permissions via **sudoers** instead of using `polkit` for user access. For example:
  - Edit the `sudoers` file to allow the user to reload Caddy without a password:
    ```bash
    sudo visudo
    ```
    Add this line to the file:
    ```bash
    username ALL=NOPASSWD: /usr/bin/systemctl reload caddy
    ```
