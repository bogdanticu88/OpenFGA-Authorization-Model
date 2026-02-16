# SaaS Multi-Tenant Authorization

Multi-tenant SaaS authorization model with complete isolation between customers.

## Model Overview

Complete multi-tenant authorization hierarchy:

- Tenant: Customer organization (strict isolation)
- Workspace: Organization unit within tenant
- Project: Project within workspace
- Resource: Protected resource within project
- Role: Tenant-specific role definitions
- Audit Log: Complete access audit trail

## Tenant Isolation

Complete data and access isolation between tenants:

```
Tenant A (Company A) - Completely isolated
  ├── Workspace
  ├── Projects
  └── Resources

Tenant B (Company B) - Completely isolated
  ├── Workspace
  ├── Projects
  └── Resources
```

User in Tenant A cannot access Tenant B resources.

## Role Hierarchy

### Tenant Level

- owner: Full control over entire tenant
- admin: Administrative functions
- member: Standard member access

### Workspace Level

- owner: Full workspace control
- admin: Workspace administration
- member: Member access

### Project Level

- owner: Project ownership
- admin: Project administration
- member: Project participation

### Resource Level

- owner: Resource ownership
- editor: Can modify resource
- viewer: Read-only access

## Relationship Examples

### Create tenant with owner

```
tenant:acme#owner@user:alice
tenant:acme#admin@user:bob
tenant:acme#member@user:charlie
```

### Create workspace within tenant

```
workspace:main#tenant@tenant:acme
workspace:main#owner@user:alice
workspace:main#member@user:charlie
```

### Share resource

```
resource:doc1#project@project:api
resource:doc1#owner@user:alice
resource:doc1#editor@user:bob
resource:doc1#viewer@user:charlie
```

## Invitation System

Invite users to tenant with role:

```
invitation:inv_1#tenant@tenant:acme
invitation:inv_1#inviter@user:alice
invitation:inv_1#invitee@user:new_user
invitation:inv_1#role@string:workspace_member
invitation:inv_1#expires_at@string:2025-01-15T00:00:00Z
```

Invitation workflow:
1. Existing member creates invitation
2. Invite sent to new user
3. User accepts invitation
4. User added to tenant with role
5. Invitation expires after acceptance or date

## Custom Roles

Define custom roles per tenant:

```
role:custom_analyst#tenant@tenant:acme
role:custom_analyst#name@string:Data Analyst
role:custom_analyst#permissions@string:view_reports,export_data,create_queries
```

Each tenant can define their own roles.

## Audit Logging

Track all access and changes:

```
audit_log:log_1#tenant@tenant:acme
audit_log:log_1#user@user:alice
audit_log:log_1#action@string:created_resource
audit_log:log_1#resource@string:resource:api_keys
audit_log:log_1#timestamp@string:2024-12-15T10:30:00Z
```

Audit trail includes:
- Who performed action
- What action was performed
- What resource was affected
- When action occurred

## Cross-Tenant Operations

Prevent cross-tenant access with resource relationships:

```
resource:api_keys#tenant@tenant:acme
```

User from tenant:techcorp cannot access this resource.

## Workspace Management

Multiple workspaces within tenant:

```
workspace:production#tenant@tenant:acme
workspace:staging#tenant@tenant:acme
workspace:development#tenant@tenant:acme
```

Users can be members of multiple workspaces within tenant.

## Project Scoping

Projects scoped to both workspace and tenant:

```
project:backend_api#workspace@workspace:acme_main
project:backend_api#tenant@tenant:acme
```

Double-scoping ensures:
- Tenant isolation
- Workspace boundaries
- Project ownership clarity

## Member Lifecycle

### Adding Member

1. Create invitation with role
2. User receives invitation
3. User accepts
4. User added to tenant/workspace
5. Audit log entry created

### Modifying Permissions

```
workspace:main#admin@user:charlie
```

Change user from member to admin.

### Removing Member

1. Remove all relationships
2. Revoke tokens/sessions
3. Audit log removal
4. Notify user (optional)

## Billing and Usage

Track resource usage per tenant:

```
audit_log#action@string:created_resource
audit_log#resource@string:workspace
```

Use audit logs to:
- Calculate tenant usage
- Implement usage limits
- Generate billing reports

## Performance Patterns

### Fast tenant check

Check if user belongs to tenant:

```
openfga check --user user:alice --relation member --object tenant:acme
```

### List user's tenants

```
openfga list --relation member --from user:alice
```

### Check resource access

```
openfga check --user user:alice --relation editor --object resource:doc1
```

## Security Considerations

1. Enforce tenant isolation at application layer
2. Validate all relationships against tenant
3. Implement request-level tenant context
4. Audit all cross-tenant operations
5. Implement hard delete for tenant offboarding
6. Regular audit log reviews

## Common Operations

### User joins workspace

```
workspace:main#member@user:alice
```

### User becomes admin

```
workspace:main#admin@user:alice
```

### Grant resource access

```
resource:doc#viewer@user:bob
```

### Create audit entry

```
audit_log:new#user@user:alice
audit_log:new#action@string:read_resource
audit_log:new#resource@string:resource:doc
```

## Scaling Considerations

- Tenant isolation reduces query scope
- Workspace scoping further improves performance
- Use indexes on tenant_id and user_id
- Implement caching for role definitions
- Archive old audit logs regularly

## Compliance Features

- Complete audit trail
- User access history
- Change tracking
- Compliance reports
- Data retention policies

## Integration Patterns

### Single Sign-On (SSO)

Tenant users authenticate via SSO.
Add authenticated user to tenant:

```
tenant:acme#member@user:{sso_user_id}
```

### External Identity Providers

Map external IDs to internal users:

```
tenant:acme#member@user:okta-user-123
```

### Webhook Notifications

Notify on permission changes:
- User added to tenant
- Role changed
- Access revoked
