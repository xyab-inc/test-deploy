# Docker Stack with Embedded Deployment Action

This project demonstrates how to use an embedded GitHub Action to deploy a Docker Compose stack to remote servers. It includes:

1. A sample web stack using nginx, Node.js, and Redis
2. A custom GitHub Action for remote deployment
3. Workflows for staging and production deployments

## Project Structure

```
.
├── .github/
│   ├── actions/
│   │   └── docker-deploy/     # Embedded GitHub Action
│   │       ├── action.yml
│   │       ├── Dockerfile
│   │       └── ... (Go files)
│   └── workflows/
│       └── deploy.yml         # Deployment workflow
└── docker-compose.yml         # Sample stack definition
```

## Stack Components

The sample stack includes:
- **Web Server**: Nginx reverse proxy
- **API**: Simple Node.js HTTP server
- **Cache**: Redis instance

Each service uses the commit SHA as its image tag via the `DOCKER_TAG` environment variable.

## Deployment Action

The embedded action (`.github/actions/docker-deploy`) handles:
1. SSH connection to remote servers
2. Transfer of docker-compose.yml and .env files
3. Pulling updated images
4. Starting/updating containers

### Action Inputs

- `ssh_user`: SSH username
- `ssh_key`: SSH private key
- `ssh_host`: Remote host
- `ssh_port`: SSH port (default: "22")
- `compose_file`: Path to docker-compose.yml
- `docker_tag`: Docker image tag (usually 7-char commit SHA)

## Usage

1. Set up GitHub Environments:
   - Create `staging` environment
   - Create `production` environment with required approvals

2. Add secrets to each environment:
   ```
   SSH_USER     # Remote server username
   SSH_KEY      # SSH private key
   SSH_HOST     # Remote server hostname
   ```

3. The workflow will:
   - Deploy to staging on every push/PR
   - Deploy to production when merging to main (requires approval)

## Local Development

To run the stack locally:
```bash
# Set a test tag
export DOCKER_TAG=test

# Start the stack
docker compose up -d

# Test the API
curl http://localhost:8080

# Stop the stack
docker compose down
```

## Security Notes

- Use GitHub Environments to manage deployment secrets
- Enable required reviewers for production deployments
- The action uses `ssh.InsecureIgnoreHostKey()` for SSH connections
- Ensure your remote Docker daemon is properly secured
- Review image versions in docker-compose.yml regularly
