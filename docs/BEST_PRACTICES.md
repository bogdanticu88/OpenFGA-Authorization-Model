# Security Best Practices

Guidelines for secure OpenFGA deployment and usage.

## Core Principles

1. Principle of Least Privilege
   - Grant minimum necessary permissions
   - Regular permission audits
   - Revoke unused access immediately

2. Defense in Depth
   - Multiple layers of authorization
   - Application-level checks
   - Database-level constraints

3. Audit Everything
   - Log all authorization changes
   - Monitor access patterns
   - Track permission grants and revocations

## API Security

### Authentication

Authenticate all API requests:

- Use mutual TLS (mTLS) for service-to-service communication
- Implement API key management with rotation
- Use OAuth2 for user authentication
- Never expose tokens in logs

### Authorization

Protect OpenFGA API endpoints:

- Restrict write operations to authorized admins
- Implement read-only mode for applications
- Use separate credentials for different environments
- Implement rate limiting

### Network Security

- Run OpenFGA in private network
- Use VPC security groups
- Implement firewall rules
- Disable public internet access
- Use TLS for all communication

## Database Security

### Connection Security

```bash
export OPENFGA_DATASTORE_URI=postgres://user:pass@localhost/db?sslmode=require
```

- Require TLS for database connections
- Use strong credentials
- Implement connection pooling
- Use separate DB user per environment

### Data Protection

- Enable encryption at rest (EBS encryption)
- Enable encryption in transit (TLS)
- Implement row-level security
- Regular backups to secure storage

### Access Control

- Restrict database access to OpenFGA service only
- No direct application access
- Implement audit logging
- Monitor query execution

## Model Security

### Model Design

```fga
# Good: Explicit relationships
type resource
  relations
    define owner: [user]
    define editor: [user]
    define viewer: [user]

# Avoid: Over-permissive wildcards
type resource
  relations
    define member: [user, team, organization]
```

Guidelines:
- Keep models simple and understandable
- Use explicit types instead of wildcards
- Document all relationships
- Regular model reviews
- Version control for models

### Permission Model

- Use explicit allow, no implicit deny
- Implement hierarchical roles
- Regular permission audits
- Test edge cases

## Application Integration

### Correct Integration Pattern

```go
// Authenticate user
user := authenticateUser(request)

// Create authorization context
ctx := context.WithUser(request.Context(), user)

// Check authorization
canAccess, err := checkAuthorization(ctx, user, "edit", resource)
if !canAccess {
  return unauthorized()
}

// Perform operation
performOperation(resource)

// Log the action
auditLog(user, "edit", resource, "allowed")
```

### Common Mistakes to Avoid

1. Caching authorization results too long
2. Checking authorization after operation
3. Trusting client-provided user ID
4. No error handling for OpenFGA failures
5. Skipping authorization for admin users

## Tenant Isolation

### Validation

Always validate tenant membership:

```go
// BAD: Trust resource tenant
resource := getResource(resourceID)
user.checkPermission(resource.tenant, "edit")

// GOOD: Validate request tenant
requestTenant := getCurrentTenant(request)
resource := getResource(resourceID)
if resource.tenant != requestTenant {
  return forbidden()
}
user.checkPermission(requestTenant, "edit")
```

### Prevent Leakage

- Never expose relationships from other tenants
- Filter queries by tenant
- Implement hard tenant boundary
- Regular cross-tenant penetration testing

## Audit and Logging

### Log Everything

```json
{
  "timestamp": "2024-12-15T10:30:00Z",
  "user": "user:alice",
  "action": "check_authorization",
  "resource": "organization:acme",
  "relation": "admin",
  "allowed": true,
  "duration_ms": 2,
  "request_id": "req-123"
}
```

### Log Retention

- Retain logs for compliance period (typically 1-3 years)
- Archive old logs to cold storage
- Implement log integrity checks
- Monitor for suspicious patterns

### Audit Trail

- Who created/modified relationships
- When changes occurred
- What changed
- Why (if available)

## Compliance

### Security Reviews

- Annual security audits
- Penetration testing
- Code reviews of integration points
- Model design reviews

### Compliance Standards

- SOC2 compliance
- GDPR data handling
- HIPAA audit requirements
- Industry-specific standards

## Operational Security

### Secrets Management

- Store credentials in vault (HashiCorp, Azure Key Vault, AWS Secrets Manager)
- Rotate credentials regularly
- Never commit secrets to source control
- Use separate credentials per environment

### Deployment Security

- Use signed Docker images
- Scan for vulnerabilities
- Update base images regularly
- Implement admission controllers (K8s)

### Incident Response

- Document authorization-related incidents
- Implement incident detection
- Have playbooks for:
  - Unauthorized access detection
  - Permission escalation attempts
  - Data leakage prevention

## Testing Security

### Authorization Testing

```yaml
test: "unauthorized_user_cannot_edit"
entity: organization:acme
relations:
  - user:alice -> admin
operations:
  - user: user:bob
    relation: editor
    expected: false
```

### Security Test Cases

- Unauthorized access attempts
- Cross-tenant access attempts
- Permission inheritance correctness
- Boundary conditions
- Performance under load

### Regular Reviews

- Quarterly permission audits
- Annual security assessment
- Vulnerability scanning
- Dependency updates

## Monitoring

### Key Metrics

- Authorization check latency (should be <5ms)
- Cache hit rate
- Database query performance
- Error rates
- Suspicious patterns

### Alerts

Set up alerts for:
- Authorization failures
- Unusual permission grants
- Bulk operations
- Error rate spikes
- Performance degradation

## Deployment Security

### Production Deployment

- Use HTTPS for all communication
- Enable authentication/authorization on API
- Implement rate limiting
- Use network policies
- Deploy in private networks

### Kubernetes Security

- Non-root containers
- Read-only filesystems
- Resource limits
- Network policies
- Pod security policies
- RBAC for service accounts

## References

- OpenFGA Security: https://openfga.dev/docs/security
- OWASP Authorization: https://owasp.org/www-community/Authorization
- Zanzibar Paper: https://research.google/pubs/zanzibar/
