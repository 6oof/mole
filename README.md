![Mole Logo](mole.svg)

# Micro PaaS Inspired by Laravel Forge for Git-based Deployments

## Navigation

- [Server Installation Guide](/docs/install.md)
- [Prepare Projects for Deployment](/docs/deployments.md)
    - [Project secrets](/docs/secrets.md)
    - [Docker compose (mole-compose.yaml)](/docs/compose.md)
- [Mole CLI Documentation](/docs/cli/mole.md)

## TL;DR

Mole is a micro PaaS inspired by **Laravel Forge**. Designed to get the most out of low-cost VPSes, Mole leverages **Docker Compose**, **SSH**, and **Caddy** for seamless deployment and management, all with minimal configuration.

Mole stands out with its flexible template system, allowing you to define and reuse dynamic configurations for Docker compose. This ensures you can deploy the same service on the same server seamlessly multiple times, avoiding conflicts and reducing setup time.

Check out the navigation above to learn how to configure your server or deploy your project with Mole!

## Quick Install - Rocky Linux

**Requirements:**

- Freshly provisioned Rocky Linux VPS.
- SSH access to the root account.

```bash
curl -O https://raw.githubusercontent.com/zulubit/mole/main/install-rocky.sh
chmod +x install-rocky.sh
./install-rocky.sh
rm install-rocky.sh
```

For a detailed installation guide, navigate to the [Server Installation Guide](/docs/install.md).

## Project Status

Mole is in early development and may require some effort to set up. Most of the CLI is still subject to change.

Currently, I’m focused on improving and testing the CLI. The next steps include creating a more streamlined setup process and building a web dashboard for users who prefer not to use SSH.

## TODO

- [ ] Gather Feedback
- [ ] Detailed Testing
- [ ] CLI Improvements
- [ ] Premade Projects Registry (Quick Deploys)
- [ ] API
- [ ] Web Dashboard
