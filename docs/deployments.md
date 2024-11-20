# Prepare Projects for Deployment by Mole

This section explains the entire process of deploying a project with Mole, as well as the steps involved in the deployment cycle.

---

## Requirements

Mole is primarily designed to deploy two types of projects:

- **Static Websites**
- **Docker Compose Projects**

### Minimum Requirements for Deployment

To deploy a project with Mole, the following conditions must be met:

1. The project is hosted in an accessible Git repository.
2. The root of the project contains a `.gitignore` file.
3. The `.gitignore` file includes the following entries:
   - `.env`
   - `mole-compose-ready.yaml`
   - `mole-deploy-ready.sh`
4. The root of the project contains:
   - `mole-compose.yaml`
   - `mole-deploy.sh`

### Directory Structure Example
```plaintext
project-root/
├── .gitignore           # Contains necessary exclusions
├── mole-compose.yaml    # Docker Compose template
├── mole-deploy.sh       # Deployment script template
├── src/                 # Project source files
└── README.md            # Documentation for the project
```

A bootstrapping script to initialize your project with these requirements is available [here](#).

---

## Adding Projects to Mole

You can add a project to Mole by running the following command on the server as the `mole` user:

```bash
mole projects add <project-name> -b <branch-name> -r <remote-repository-url>
```

### Automatic `.env` Generation

When a project is successfully added, Mole automatically generates a `.env` file for the project. The file will look like this:

```env
# Auto-generated environment configuration for <project-name>.
# DO NOT DELETE OR MODIFY THIS SECTION.
# This configuration is necessary for the project to work properly.

# Project name
MOLE_PROJECT_NAME=<project-name>

# Project root path
MOLE_ROOT_PATH=/home/mole/projects/<project-name>

# Reserved ports for this deployment
MOLE_PORT_APP=9000
MOLE_PORT_TWO=9001
MOLE_PORT_THREE=9002

# Random string used as a key when necessary
MOLE_APP_KEY=2fqhh3GeY3s2O3HPj2Qj9H7P5XKZ6DLW

# Database credentials
MOLE_DB_NAME=wpdb0BfdDV9r
MOLE_DB_USER=wpuserJyj7SP
MOLE_DB_PASS=cYEDMVF4vYFaHTKrsht64lnw

# User-defined environment variables can be added below:
# Add your custom variables here:
```

#### Merging Custom Variables
If you want to merge additional environment variables into the generated `.env` file, include them in a file named `.env.mole` in the root of the repository. Mole will automatically merge them during deployment.

Learn more about the `.env` file mole generates [here](#).

---

## Configurations as Templates

Both `mole-compose.yaml` and `mole-deploy.sh` are treated as Go `text/template` files. When a deployment is triggered, they are transformed into:

- `mole-compose-ready.yaml`
- `mole-deploy-ready.sh`

### Template Transformation

During transformation, Mole reads the `.env` file and injects the values into the templates using Go's `{{.}}` syntax. For example:

**Template: `mole-deploy.sh`**
```bash
#!/bin/bash

echo "{{.MOLE_PROJECT_NAME}}"
```

**Transformed File: `mole-deploy-ready.sh`**
```bash
#!/bin/bash

echo "choosen-project-name"
```

This transformation assumes the `.env` file contains:
```env
MOLE_PROJECT_NAME=choosen-project-name
```

#### No Replacements Necessary
If no template placeholders (`{{.}}`) are present, the files are simply copied over to their "ready" state.

---

## Deployment Cycle

You can trigger deployments using the following command:

```bash
mole deploy <project-name-or-id>
```

### Steps in the Deployment Cycle

1. **Git Pull**: Mole pulls the latest changes from the specified Git branch.
2. **Template Transformation**: Mole processes all templates into their "ready" state (`mole-compose-ready.yaml` and `mole-deploy-ready.sh`).
3. **Deployment Script Execution**: Mole runs the `mole-deploy-ready.sh` script to execute the deployment process.

---

This guide ensures that your project is prepared and deployed seamlessly with Mole while adhering to its requirements and workflows.
