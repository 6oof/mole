# Usign Mole with docker compose

This section provides some general guidelines on using mole with docker compose.

---

## mole-compose.yaml

To differantiate local compose files from the ones meant to be used with `mole`, you can include a file named `mole-compose.yaml`.

Much in the same way as the `mole.sh` file discussed in the [preparing for deployment](/docs/deployments.md) section, `mole-compose.yaml` is treated as a Go `text/template` file.

That means you can use `{{.}}` placeholders to inject `project secrets`. Read more about project secrets [here](/docs/secrets.md).

`mole-compose.yaml` *DOES NOT* get transformed automatically, you have to always explicitly call `mole templates compose [project name]` to transform the tamplets.

The most appropriate place to do so it proabably in `mole.sh`. For example:

```bash
#!/bin/bash

mole templates compose {{.ProjectName}}
```

`mole-compose.yaml` is transformed into `mole-compose-ready.yaml` when the command successfuly runs.

`mole-compose-ready.yaml` can then be used with `docker compose`:

```bash
#!/bin/bash

mole templates compose {{.ProjectName}}

docker compose -f mole-compose-ready.yaml up -d --build
```

## Important considerations when creating `mole-compose.yaml` files

When creating `mole-compose.yaml` files, it’s essential to follow security best practices to ensure your application is safe in production.

### Avoid Directly Exposing Ports
You should **never directly expose ports** to all network interfaces (e.g., using `PORT:INTERNAL_PORT`). Always bind services to `127.0.0.1` to restrict access to the local machine. This prevents Docker from bypassing firewall rules and inadvertently exposing your application to the internet.

#### WRONG Example (Directly Exposed Ports):
```yaml
services:
  app:
    image: my-app:latest
    ports:
      - "8080:8080" # Open to the world, dangerous in production
    environment:
      - NODE_ENV=production
```

#### RIGHT Example (Bind to Localhost):
```yaml
services:
  app:
    image: my-app:latest
    ports:
      - "127.0.0.1:8080:8080" # Safely bound to localhost
    environment:
      - NODE_ENV=production
```

By binding to `127.0.0.1`, your application can only be accessed from the server itself. Any external access should be managed via `mole domains`.

### Logging Best Practices
Properly configuring logging for your `mole-compose.yaml` files is crucial for monitoring application behavior and diagnosing issues in production. Docker Compose provides flexible logging options that you can use to manage log retention and file sizes efficiently.

#### Configuring JSON File Logging
One recommended approach is to use Docker's `json-file` logging driver with size and file limits. This prevents logs from consuming excessive disk space while ensuring that critical logs are retained.

#### Example Logging Configuration:
```yaml
services:
  app:
    image: my-app:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"  # Maximum size of each log file
        max-file: "2"    # Maximum number of log files to retain
```

In this configuration:
- **`max-size`**: Limits the size of each log file to 10 MB.
- **`max-file`**: Retains up to 2 rotated log files, effectively storing 20 MB of logs before older logs are discarded.

### Accessing Logs with JSON-File Logging Driver

If you're using the `json-file` logging driver, accessing logs is straightforward and can be done directly from the project's directory relative to your SSH access point.

#### Steps to Access Logs
1. **SSH into the Host Machine**:
   Log in to the server using SSH:
   ```bash
   ssh mole@<host>
   ```

2. **Navigate to the Project Directory**:
   Change into the directory where your Docker Compose project is located:
   ```bash
   cd /projects/[projectname]
   ```
   Replace `[projectname]` with the name of your project.

3. **View Logs for All Services**:
   Use the following command to see the logs for all services defined in your `docker-compose.yml`:
   ```bash
   docker compose logs
   ```

By following these steps, you can efficiently access your container logs from your SSH entry point and monitor application behavior.

### Saving App Logs to Mole's Log Directory
If your app is writing logs to a directory, it's recommended to integrate Mole's templating capabilities, you can dynamically specify the log storage directory using the `{{.LogDirectory}}` placeholder. This ensures that logs are copied tp a centralized and organized location managed by Mole.

#### Example with Log Directory:
```yaml
services:
  app:
    image: my-app:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "2"
    volumes:
      - {{.LogDirectory}}:/var/log/app # Mount Mole's log directory
```

By using the `{{.LogDirectory}}` placeholder:
- Logs are stored in a directory managed by Mole, ensuring consistency across deployments.
- You maintain easy access to logs for debugging or monitoring.

### Restart Policies

To ensure reliable service availability, it is recommended to use the `always` restart policy for most production scenarios. This policy ensures that services restart automatically under all circumstances, providing maximum reliability for critical applications.

#### Example Configuration
```yaml
services:
  app:
    image: my-app:latest
    restart: always # Ensures the app restarts automatically in all scenarios

  db:
    image: postgres:14
    restart: always # Ensures the database restarts automatically in all scenarios
```

This policy is simple and effective for production environments where services need to stay running at all times. It ensures that services recover from failures or host restarts without requiring manual intervention.

For more information on restart policies, refer to the [Docker Compose documentation](https://docs.docker.com/engine/containers/start-containers-automatically/).

## Example

Below is a complete example of a `mole-compose.yaml` file. This example demonstrates how to utilize Mole's project secrets to dynamically inject values into your deployment configuration, ensuring consistency, security, and maintainability.

```yaml
version: '3.9'
services:
  app:
    image: my-app:latest
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "127.0.0.1:{{.PortApp}}:8080" # Bind to a dynamically allocated port
    environment:
      - NODE_ENV=production
      - APP_KEY={{.AppKey}} # Inject the generated app key
      - DATABASE_URL=postgres://{{.DatabaseUser}}:{{.DatabasePass}}@db/{{.DatabaseName}}
      - LOG_DIRECTORY={{.LogDirectory}} # Path to the log directory
    volumes:
      - {{.LogDirectory}}:/var/log/app # Log storage
      - {{.EnvFilePath}}:/app/.env # Inject the .env file
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "2"
    restart: always
    depends_on:
      - db

  db:
    image: postgres:14
    environment:
      POSTGRES_USER: {{.DatabaseUser}}
      POSTGRES_PASSWORD: {{.DatabasePass}}
      POSTGRES_DB: {{.DatabaseName}}
    volumes:
      - db-data:/var/lib/postgresql/data
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "2"
    restart: always

volumes:
  db-data: {}
```

### Key Features in the Example:
1. **Port Configuration**:
   - Uses `{{.PortApp}}` to dynamically allocate a primary port for the application.
   - Ensures secure binding with `127.0.0.1`.

2. **Secrets**:
   - Injects secrets directly into the environment.

3. **Logging**:
   - Configures `json-file` logging for both the `app` and `db` services.
   - Stores logs in the directory defined by `{{.LogDirectory}}`.

4. **Volumes**:
   - Mounts the `.env` file from `{{.EnvFilePath}}` for app configuration.
   - Stores database data persistently in `db-data`.

5. **Database Configuration**:
   - Dynamically injects the `DatabaseUser`, `DatabasePass`, and `DatabaseName` into the `db` service.

6. **Restart Policies**:
   - Ensures both `app` and `db` services automatically restart in case of failure with `restart: always`.

## Additional Resources

1. **[Docker Compose](https://docs.docker.com/compose/)**  
   Explore the official Docker documentation for all configuration options available in `docker-compose.yaml`.
2. **[Are you completely new to docker? Watch this Fireship video](https://www.youtube.com/watch?v=gAkwW2tuIqE)**
