![Mole Logo](mole.svg)
   
# Micro PaaS for Efficient Git-based Deployments with Systemd

## Navigation

- [Server Installation Guide](/docs/install.md)
- [Graceful Deployment Guide](/docs/graceful_deployment.md)
- [Mole CLI Documentation](/docs/cli/mole.md)

## TLDR

Mole is a micro PaaS that squeezes every bit of power from low-cost VPSes while keeping things simple and scalable. It handles Git-based deployments with systemd, SSH, Podman, and Caddy—minimal fuss required.

Check out the navigation above to learn how to configure your server or deploy your project with Mole!

## Quick install - Rocky linux

**Requirements:**

- Freshly provisioned Rocky linux VPS.
- SSH access to the root account.

```bash
curl -O https://raw.githubusercontent.com/zulubit/mole/main/install-rocky.sh
chmod +x install-rocky.sh
./install-rocky.sh
rm install-rocky.sh
```

For detailed install guide navigate to the [Server Installation Guide](/docs/install.md)

## Project Status

Mole is in early development and may require some effort to set up. Most of the CLI is still subject to change.

Currently, I’m focused on improving and testing the CLI. Next, I’ll tackle a more streamlined setup process and a web dashboard for users who prefer not to use SSH.

## TODO

- [ ] Gather Feedback
- [ ] Detailed Testing
- [ ] CLI Improvements
- [ ] Streamline Installation Process
- [ ] Premade Projects Registry (Quick Deploys)
- [ ] Web Dashboard
