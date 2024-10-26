# Graceful Deployment with Systemd and Service Grouping

This guide describes how to achieve graceful, zero-downtime deployments using `systemd` for applications that require replication and sequential restarts. By structuring `systemd` units with dependencies, we can ensure a controlled, ordered deployment process that minimizes downtime and keeps services up if part of the deployment fails.

## Overview

The goal is to set up two instances of a service to deploy in a controlled sequence. If one instance fails to restart, the subsequent instance will remain running, maintaining availability for users.

## Key Concepts

1. **Service Instances**: Each instance is defined as a separate `systemd` service unit (e.g., `my-service-1`, `my-service-2`).
2. **Dependencies**: We use `After` and `Requires` directives to enforce sequential starting and stopping, ensuring each instance only starts if the previous one is running.
3. **Service Grouping**: A "master" unit (`my-service-group.service`) controls the group of instances, making it easy to initiate a rolling restart across all instances.

## Systemd Unit Configuration

Each instance has its own `.service` file with specific dependencies on the previous instance.

### Example Service Unit Files

#### `my-service-1.service`
```ini
[Unit]
Description=My Go Application Instance 1
After=network.target
Requires=network.target

[Service]
ExecStart=/path/to/my-binary --port=8081
Restart=on-failure
TimeoutStopSec=10
RestartSec=2  # Delay before restart if failed
```

#### `my-service-2.service`
```ini
[Unit]
Description=My Go Application Instance 2
After=my-service-1.service
Requires=my-service-1.service

[Service]
ExecStart=/path/to/my-binary --port=8082
Restart=on-failure
TimeoutStopSec=10
RestartSec=2
```

### Grouping Unit (Optional Master Unit)

To control both instances as a single group, we define a master unit that links to all instances using `PartOf`. Restarting or reloading this master unit will apply the action to both instances in sequence.

#### `my-service-group.service`
```ini
[Unit]
Description=My Go Application Group

[Service]
PartOf=my-service-1.service my-service-2.service
```

## Deployment Behavior

1. **Sequential Start and Stop**: Each instance only starts after the previous one reaches an "active" state. If `my-service-1` fails, `my-service-2` will not start, preventing a cascading failure.

2. **Failure Isolation**: If a restart attempt fails for `my-service-1`, `my-service-2` remains unaffected and will continue to run, preserving service availability.

3. **Graceful Rollback**: If a deployment fails, services that were already running remain unaffected. This provides resilience, allowing you to attempt another deployment without disrupting users.
