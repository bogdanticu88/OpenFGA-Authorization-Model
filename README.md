# OpenFGA Authorization Model

Production-ready OpenFGA authorization models and implementation guide for RBAC, ABAC, API, and SaaS applications. Designed for teams adopting fine-grained access control with security best practices.

## Overview

OpenFGA is a high-performance authorization service built on Zanzibar architecture. This repository provides comprehensive authorization models, deployment configurations, and implementation examples for different use cases.

## Features

- Production-grade authorization models (RBAC, ABAC, API, SaaS)
- Fine-grained access control patterns
- Multi-tenancy support
- API authorization examples
- Security best practices
- Docker and Kubernetes deployment
- Example implementations in Go, Python, Node.js
- Testing and validation patterns
- Migration strategies

## Directory Structure

```
.
├── models/
│   ├── rbac/                      # Role-Based Access Control
│   │   ├── model.fga              # DSL authorization model
│   │   ├── relations.txt          # Relationship definitions
│   │   └── README.md
│   │
│   ├── abac/                      # Attribute-Based Access Control
│   │   ├── model.fga
│   │   ├── relations.txt
│   │   └── README.md
│   │
│   ├── api/                       # API Authorization
│   │   ├── model.fga
│   │   ├── relations.txt
│   │   └── README.md
│   │
│   └── saas/                      # SaaS Multi-tenant Authorization
│       ├── model.fga
│       ├── relations.txt
│       └── README.md
│
├── examples/
│   ├── go/                        # Go client examples
│   │   ├── main.go
│   │   ├── rbac_client.go
│   │   ├── abac_client.go
│   │   ├── api_client.go
│   │   └── go.mod
│   │
│   ├── python/                    # Python client examples
│   │   ├── requirements.txt
│   │   ├── rbac_client.py
│   │   ├── abac_client.py
│   │   ├── api_client.py
│   │   └── main.py
│   │
│   └── nodejs/                    # Node.js client examples
│       ├── package.json
│       ├── rbac_client.js
│       ├── abac_client.js
│       ├── api_client.js
│       └── main.js
│
├── deployment/
│   ├── docker/
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   └── .env.example
│   │
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       ├── configmap.yaml
│       └── README.md
│
├── tests/
│   ├── test_rbac.yaml
│   ├── test_abac.yaml
│   ├── test_api.yaml
│   ├── test_saas.yaml
│   └── README.md
│
├── migrations/
│   ├── v1_initial_rbac.txt
│   ├── v2_add_api_authorization.txt
│   └── README.md
│
├── docs/
│   ├── ARCHITECTURE.md
│   ├── BEST_PRACTICES.md
│   ├── SECURITY.md
│   ├── GETTING_STARTED.md
│   └── MIGRATION_GUIDE.md
│
└── README.md
```

## Quick Start

### Prerequisites

- OpenFGA server running (v1.0+)
- Docker (for local development)
- Client SDK for your language (Go, Python, Node.js)

### Run Locally

1. Start OpenFGA with Docker Compose

```bash
cd deployment/docker
docker-compose up -d
```

2. Create authorization model

```bash
cd models/rbac
openfga model write --api-url http://localhost:8080 model.fga
```

3. Create relationships

```bash
openfga relationship write --api-url http://localhost:8080 relations.txt
```

4. Test authorization

```bash
openfga check --api-url http://localhost:8080 \
  --user "user:alice" \
  --relation "admin" \
  --object "organization:acme"
```

## Authorization Models

### RBAC - Role-Based Access Control

Traditional role-based access control. Users are assigned roles, roles have permissions.

- Simple and well-understood
- Good for small to medium organizations
- Limited to predefined roles

See: models/rbac/

### ABAC - Attribute-Based Access Control

Attribute-based access control using user and resource attributes.

- Flexible and scalable
- Policy-based decisions
- Complex policies possible
- Higher maintenance overhead

See: models/abac/

### API Authorization

Fine-grained API endpoint authorization.

- Resource-based access control
- Scope-based permissions
- Token validation patterns
- API key management

See: models/api/

### SaaS Multi-tenant Authorization

Multi-tenant SaaS authorization patterns.

- Tenant isolation
- Per-tenant resource ownership
- Organization hierarchies
- Cross-tenant access rules

See: models/saas/

## Implementation Examples

### Go

Full-featured Go client with authorization checks:

```bash
cd examples/go
go run main.go
```

Demonstrates:
- Connecting to OpenFGA
- Creating relationships
- Checking access
- Listing permissions

### Python

Python client implementation:

```bash
cd examples/python
pip install -r requirements.txt
python main.py
```

Demonstrates:
- Authentication
- Authorization checks
- Relationship management
- Error handling

### Node.js

Node.js client implementation:

```bash
cd examples/nodejs
npm install
node main.js
```

Demonstrates:
- Async/await patterns
- API integration
- Permission checks
- Relationship writes

## Deployment

### Docker

Development and testing:

```bash
cd deployment/docker
docker-compose up -d
```

Includes:
- OpenFGA server
- Postgres database
- Pre-configured models

### Kubernetes

Production deployment:

```bash
cd deployment/kubernetes
kubectl apply -f .
```

Includes:
- OpenFGA deployment
- Service configuration
- ConfigMap for models
- High availability setup

See: deployment/kubernetes/README.md

## Security Best Practices

### Model Design

- Keep models simple and maintainable
- Use explicit deny patterns where needed
- Document all relationships
- Regularly audit permission models

### Access Control

- Implement principle of least privilege
- Regular access reviews
- Revoke unused permissions
- Audit permission changes

### API Security

- Validate all inputs
- Use TLS for communication
- Implement rate limiting
- Add request authentication

### Database

- Enable encryption at rest
- Use connection pooling
- Regular backups
- Audit access logs

### Network

- Network isolation
- VPC segmentation
- Firewall rules
- DDoS protection

See: docs/SECURITY.md for detailed security guidelines

## Testing

Comprehensive test suites for each model:

```bash
cd tests
openfga test test_rbac.yaml
openfga test test_abac.yaml
openfga test test_api.yaml
openfga test test_saas.yaml
```

Each test file includes:
- Positive test cases
- Negative test cases
- Edge cases
- Performance checks

## Migration

Strategies for migrating from existing authorization systems:

- From LDAP/Active Directory
- From custom role systems
- From OAuth2 scopes
- From legacy ACLs

See: migrations/ and docs/MIGRATION_GUIDE.md

## Documentation

- ARCHITECTURE.md - System architecture and design decisions
- BEST_PRACTICES.md - Authorization design patterns and recommendations
- SECURITY.md - Security guidelines and hardening steps
- GETTING_STARTED.md - Step-by-step tutorial
- MIGRATION_GUIDE.md - Migrating from existing systems

## Performance Considerations

- Response time: <5ms for authorization checks
- Throughput: 10,000+ checks per second
- Latency: Network-bound, typically <10ms
- Storage: Efficient tuple storage with indexing

Optimization tips:
- Use caching for repeated checks
- Batch operations where possible
- Monitor database performance
- Regular index optimization

## Troubleshooting

### Authorization Denied Issues

1. Verify relationship exists in database
2. Check model for correct relation definition
3. Validate user and object identifiers
4. Review namespace formatting

### Performance Problems

1. Check database connection pool
2. Monitor query execution time
3. Verify index usage
4. Analyze hot relationships

### Integration Issues

1. Verify API credentials
2. Check network connectivity
3. Validate model format
4. Review error logs

## Contributing

Guidelines for contributions:
- Models should be production-ready
- Include comprehensive tests
- Update documentation
- Follow security best practices

## References

- OpenFGA Documentation: https://openfga.dev/docs
- Zanzibar Architecture: https://research.google/pubs/zanzibar/
- Fine-grained Authorization: https://openfga.dev/blog
- Authorization Best Practices: https://en.wikipedia.org/wiki/Access_control

## Support

For issues and questions:
1. Check documentation in docs/
2. Review model examples in models/
3. Test with included test cases
4. Consult OpenFGA official documentation

## License

This repository is provided as-is for educational and enterprise use.
