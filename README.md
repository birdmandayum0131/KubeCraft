# KubeCraft

KubeCraft is a project that sets up a Minecraft server on Kubernetes and provides an external API interface via the KubeCraft Gateway. For a simple front-end implementation, you can check out the [CreeperKeeper](https://github.com/birdmandayum0131/CreeperKeeper) project.

If you're unsure whether you need KubeCraft, here's some guidance.   

If you have a manageable host and you want to run a Minecraft server on it, and all you need is to manage the server's status manually through commands or a dashboard, you may not need the KubeCraft project. You can achieve this using Minecraft server files or a Minecraft server Docker image.

However, if you have a Kubernetes-managed host and you want to run a Minecraft server on it while providing internal management APIs, then KubeCraft can meet your technical requirements.

## About The Project

The origin of this project stems from my experiences with friends in gaming sessions. I'm usually the one responsible for running the server, and when we're all playing together, there's no need to handle additional server access issues.

However, after the initial excitement fades, sometimes only a few people still want to play, or our playing times become more scattered. Out of consideration to avoid bothering me, they usually won't ask me to start or stop the server just for a bit of interest.

In such cases, sharing the server files isn't an ideal solution, and hosting the server on the cloud isn't cost-effective. With a second machine available to serve as a self-hosted server, I started this project.

**KubeCraft** can be deployed on Kubernetes using Helm and provides APIs for managing internal server deployment, rcon commands, config editing, and more.

## Features

KubeCraft consists of two primary subprojects:

- **Minecraft Bridge**
  - A Python FastAPI server.
  - Communicates with the Minecraft server using mcstatus and rcon (in progress).

- **KubeCraft Gateway**
  - A Golang Gin server.
  - Exposes APIs for starting/stopping the server, fetching server status, etc.
  - Internally communicates with Kubernetes via client-go.
  - Internally communicates with the Minecraft server via Minecraft Bridge.

## Getting Started

At this time, the fully customized Helm chart templates are not complete, but you can modify the `deploy/helm` chart to suit your needs.

### Prerequisites

- A Kubernetes environment with Helm installed (I use K3s).
- If you have an Ingress Controller, configure it to expose the Minecraft server's TCP port (25565).
- Configure your router to ensure port 25565 is publicly accessible.
- If you have existing server files, place them in the persistent volume path and **ensure** the files have the appropriate read-write permissions.

### Deploy the Project

**Modify the `deploy/helm` chart to fit your usage scenario.**  
Then install the chart into your Kubernetes cluster using the following command:

```bash
helm upgrade --install kubecraft deploy/helm -n minecraft \ 
      --set image.gateway.repository="$GATEWAY_IMAGE" \ 
      --set image.gateway.tag="$GATEWAY_VERSION" \ 
      --set image.bridge.repository="$BRIDGE_IMAGE" \ 
      --set image.bridge.tag="$BRIDGE_VERSION"
```

## Timeline

- Project start: `2024/09/01`
- Hardware setup: `2024/09/15`
- KubeCraft Gateway subproject: `2024/10/01`
- Minecraft Bridge subproject: `2024/10/23`
- GitLab CI/CD pipeline set up: `2024/10/24`

## Roadmap

- [x] FastAPI with mcstatus
- [x] Golang backend server implementation
- [x] GitLab CI/CD pipeline
- [x] Kubernetes deployment with Helm chart
- [ ] Add more server control actions
- [ ] Rcon to Minecraft server
- [ ] Manage multiple servers
- [ ] Complete Helm's chart template
