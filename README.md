• ▌ ▄ ·.       ▄▄▌  ▄▄▄ .
·██ ▐███▪▪     ██•  ▀▄.▀·
▐█ ▌▐▌▐█· ▄█▀▄ ██▪  ▐▀▀▪▄
██ ██▌▐█▌▐█▌.▐▌▐█▌▐▌▐█▄▄▌
▀▀  █▪▀▀▀ ▀█▄▀▪.▀▀▀  ▀▀▀

# Git Pull Deployment Server Manager ("micro PaaS") Focused on systemd

## Quickstart



## Detailed setup

### Prerequisites

#### This quicstart guide assumes a freshly provisioned VPS with the following:

- Running **Fedora 40**.
- SSH root access.

**Note**: **Fedora** was chosen due to better support for recent versions of **Podman**. If you'd like to use a different distribution, you may need to adjust the `install` and `enable` commands below. Note that currently Ubuntu installs a very out of date **Podman** package. As support improves this guide will probably migrate towards **Rocky** or even ubuntu.

## 1. SSH into your newly created VPS as a root user

## 2. Setting Up Caddy on Your VPS

### Step 1: Install Caddy - reverse proxy to manage your domains

Run the following commands to install Caddy and its dependencies:

```bash
dnf install 'dnf-command(copr)'
dnf copr enable @caddy/caddy
dnf install caddy
```

### Step 2: Configure Caddy

1. **Edit the Caddyfile**: Open the main configuration file located at `/etc/caddy/Caddyfile` with your preferred text editor.

   ```bash
   sudo vi /etc/caddy/Caddyfile
   ```

2. **Set Up the Caddyfile**: Replace the contents with the following (Make sure to replace the example email):

   ```caddyfile
   {
       email your_email@example.com  # Replace with your email for TLS notifications
       servers {
           protocol {
               experimental_http3
           }
       }
   }

   # Enable automatic HTTPS and request compression globally
   tls {
       on_demand  # Use on-demand TLS certificate provisioning
   }

   header {
       # Enable gzip and brotli compression
       Accept-Encoding gzip, br
       Content-Type * gzip
       Content-Type * brotli
   }

   # Include domain configurations from partials
   import ./partials/*.caddy
   ```

   - Replace `your_email@example.com` with your actual email address.

### Step 3: Enable and Start Caddy Service

Reload the systemd daemon and enable the Caddy service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now caddy
```

### Step 4: Check Caddy Status

To verify that Caddy is running:

```bash
sudo systemctl status caddy --no-pager
```

You should see output indicating that the Caddy service is active and running.

### Additional Resources

- [Caddy Fedora Linux Installation Documentation](https://caddyserver.com/docs/install#fedora-redhat-centos)

---

## 3. Install Podman

### Step 1: Install Podman

Run the following command to install Podman:

```bash
sudo dnf -y install podman
```

### Step 2: Verify the Installation

Check that Podman is installed correctly:

```bash
podman --version
```

---

## 4. Setting Up a Firewall (firewalld)

To ensure your server is secure, you can set up a firewall using **firewalld**. This will restrict incoming traffic to SSH (port 22), HTTP (port 80), and HTTPS (port 443).

### Step 1: Install firewalld

```bash
sudo dnf install firewalld
```

### Step 2: Start and Enable the Firewall

```bash
sudo systemctl start firewalld
sudo systemctl enable firewalld
```

### Step 3: Configure Allowed Ports

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

### Step 4: Reload Firewall Rules

```bash
sudo firewall-cmd --reload
```

### Step 5: Verify Active Firewall Rules

To check if the correct rules are in place:

```bash
sudo firewall-cmd --list-all
```

This should show that SSH, HTTP, and HTTPS are allowed for incoming traffic.

--- 

## Important Files Created

- `mole.json`: Placed in `/etc/mole/mole.json`.
- Project repositories will be pulled to `/var/mole/{projectId}/`.
- Caddy reverse proxy configuration in `/etc/caddy/Caddyfile`.

**Note**: Ideally, these files should not be manually edited.
