# Prepare Projects for Deployment by Mole

This section explains the entire process of deploying a project with Mole, as well as the steps involved in the deployment cycle.

---

## Requirements

Mole is primarily designed to deploy two types of projects:

- **Static Websites**
- **Docker Compose Projects**

### Minimum Requirements for Deployment

To deploy a project with Mole, the following conditions must be met:

1. The project is a Git repository.
2. The root of the project contains:
   - `mole.sh`

### Directory Structure Example
```plaintext
project-root/
├── mole.sh             # Deployment script template
├── mole-compose.yaml   # Optional: Docker Compose file
├── env.example         # Optional: Env file example to be copied to .env
├── .gitignore          # Optional: If using .env it should be ignored
...
```

A bootstrapping script to initialize your project with these requirements is available [here](#).

---

## Adding Projects to Mole

You can add a project to Mole by running the following command on the server as the `mole` user:

```bash
mole projects add <project-name> -b <branch-name> -r <remote-repository-url>
```

### Automatic `secrets` Generation

When a project is successfully added, Mole automatically generates `project secrets` for the project. The secrets available are:

```txt
	EnvFilePath   - Absolute path to the .env file
	RootDirectory - Absolute Path to the project root
	LogDirectory  - Absolute Path to the log directory
	ProjectName   - Project Name given at mole add
	AppKey        - Generated App Key if needed for the app
	PortApp       - Primary port to be used
	PortTwo       - Alternate port if needed
	PortThree     - Alternate port if needed
	DatabaseName  - Database Name if needed
	DatabaseUser  - Database user if needed
	DatabasePass  - Database password if needed
```

You can read more about them [here](/docs/secrets.md).

### Creating a Base .env File

If a `.env.example` file is found in the root of the repository, Mole will copy it directly to `.env` when the project is added. This provides a simple way to include predefined environment variables in your project.

---

## Configurations as Templates

Both `mole-compose.yaml` and `mole.sh` are treated as Go `text/template` files. When a deployment is triggered, `mole.sh` is turned into:

- `mole-ready.sh`

### Template Transformation

During transformation, Mole reads the `project secrets` file and injects the values into the templates using Go's `{{.}}` syntax. For example:

**Template: `mole.sh`**
```bash
#!/bin/bash

echo "{{.ProjectName}}"
```

**Transformed File: `mole-ready.sh`**
```bash
#!/bin/bash

echo "choosen-project-name"
```

#### No Replacements Necessary
If no template placeholders (`{{.}}`) are present, the files are simply copied over to their "ready" state.

---

## Deployment Cycle

your `mole.sh` should probably look something like this:

```bash
#!/bin/bash

git pull

docker compose -f mole-compose-ready.yaml up -d --build
```

You can trigger deployments using the following command:

```bash
mole deploy <project-name-or-id>
```

### Steps in the Deployment Cycle

1. **Deployment script is tranfromed**: Mole transforms the `mole.sh` to `mole-ready.sh`.
2. **Deployment script execution**: Mole runs the `mole-ready.sh` script to execute the deployment process.

---

This guide ensures that your project is prepared and deployed seamlessly with Mole while adhering to its requirements and workflows.
