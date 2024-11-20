# Install Mole 

## This is a detailed install script for Mole. This should ideally be performed on a fresh server.

This guide aims to produce an environment that is minimalistic and reasonably secure by default.

---

### Prerequisites

#### This quickstart guide assumes a freshly provisioned VPS with the following:

- Running a distribution with **systemd** (minimum version 240).
- SSH root access.
- Ability to install Caddy.
- Ability to reload the Caddy systemd service without root access.
- Ability to install Git.
- Ability to install Docker and Docker Compose.

**Note**: **Rocky Linux** is recommended for its stability and strong support for server environments as well as support for **Docker**. If you'd like to use a different distribution, you may need to make some adjustments to the steps below. Notes on choosing a different distro can be found at the bottom of this document.

---

### 1. SSH into your newly created VPS as a root user

---

### 2. Install Mole

#### Step 1: Download and Run the Install Script

Download the installation script and execute it:

```bash
curl -O https://raw.githubusercontent.com/zulubit/mole/main/install.sh
chmod +x install.sh
./install.sh
```

#### Step 2: Verify the Installation

Once logged in as the `mole` user, check that the `mole` CLI is properly installed and accessible:

```bash
mole version
```

You should see the version number of Mole displayed.

#### Step 3: Set Up Domains

This step is required for Caddy not to fail on start. Make sure you don't set up your domains with the root user to avoid permission errors.

```bash
su mole
```

```bash
mole domains setup your@email.com
```

```bash
exit
```

---

### 3. Install Caddy

#### Step 1: Install Caddy - Reverse Proxy to Manage Your Domains

Run the following commands to install Caddy and its dependencies:

```bash
dnf install 'dnf-command(copr)' -y
dnf copr enable @caddy/caddy -y
dnf install caddy -y
```

#### Step 2: Set Permissions

```bash
usermod -aG mole caddy
chmod 750 /home/mole
```

#### Step 3: Enable and Start Caddy API Service

Reload the systemd daemon and enable the Caddy API service:

```bash
systemctl daemon-reload
systemctl enable --now caddy-api
```

#### Additional Resources

- [Caddy RHEL Installation Documentation](https://caddyserver.com/docs/install#fedora-redhat-centos)

---

### 4. Install Docker and Docker Compose

#### Step 1: Install Docker

Run the following commands to install Docker:

```bash
dnf install -y dnf-plugins-core
dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin -y
```

#### Step 2: Enable and Start Docker

Start the Docker service and enable it to run on boot:

```bash
systemctl start docker
systemctl enable docker
```

#### Step 3: Add the `mole` User to the `docker` Group

To allow the `mole` user to manage Docker without `sudo`:

```bash
usermod -aG docker mole
```

Log out and log back in as the `mole` user for the group changes to take effect. Verify the `mole` user can run Docker commands without `sudo`:

```bash
su mole
docker ps
```

#### Step 4: Install Docker Compose

Run the following command to install Docker Compose:

```bash
dnf install docker-compose -y
```

Verify the installation:

```bash
docker-compose --version
```

---

### 5. Setting Up a Firewall (firewalld)

To ensure your server is secure, you can set up a firewall using **firewalld**. This will restrict incoming traffic to SSH (port 22), HTTP (port 80), and HTTPS (port 443).

#### Step 1: Install firewalld

```bash
dnf install firewalld
```

#### Step 2: Start and Enable the Firewall

```bash
systemctl start firewalld
systemctl enable firewalld
```

#### Step 3: Configure Allowed Ports

1. Allow SSH (port 22):

   ```bash
   firewall-cmd --permanent --add-service=ssh
   ```

2. Allow HTTP (port 80):

   ```bash
   firewall-cmd --permanent --add-service=http
   ```

3. Allow HTTPS (port 443):

   ```bash
   firewall-cmd --permanent --add-service=https
   ```

#### Step 4: Reload Firewall Rules

```bash
firewall-cmd --reload
```

#### Step 5: Verify Active Firewall Rules

To check if the correct rules are in place:

```bash
firewall-cmd --list-all
```

This should show that SSH, HTTP, and HTTPS (next to services:) are allowed for incoming traffic.

---

### 6. Install Git

```bash
dnf install git
```

---

## Choosing a Different Distribution

This guide is optimized for **Rocky Linux** and other **RHEL-based** distros (like **CentOS**, **AlmaLinux**, and **Fedora**). If you'd like to use a different distribution, here are the adjustments you’ll need to make:

### Adjustments for Ubuntu/Debian-based Distros:

- **Package Manager**: Replace `dnf` with `apt` in the installation commands:
  - For example, use `apt install docker` instead of `dnf install docker`.

- **Firewall**: If using **UFW** instead of `firewalld`, replace `firewall-cmd` commands with UFW commands:
  - `ufw allow ssh`
  - `ufw allow http`
  - `ufw allow https`

- **Docker Installation**:
  - Use the `apt` commands to install Docker:
    ```bash
    apt update
    apt install -y docker.io
    ```
  - Install Docker Compose:
    ```bash
    apt install -y docker-compose
    ```

