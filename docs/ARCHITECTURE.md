# Architecture and Design

OpenFGA architecture patterns and design decisions.

## System Architecture

```
+-------------------+
| Applications      |
+-------------------+
        |
        | gRPC/HTTP
        |
+-------------------+
| OpenFGA API       |
+-------------------+
        |
        | Query
        |
+-------------------+
| Authorization     |
| Engine            |
+-------------------+
        |
        | Read/Write
        |
+-------------------+
| Datastore         |
| (Postgres)        |
+-------------------+
```

## Authorization Decision Flow

```
User Request
    |
    v
[Application Layer]
    - Authenticate user
    - Extract context
    |
    v
[Authorization Layer]
    - Build OpenFGA request
    - Execute check
    |
    +-> Cache hit?
    |   └-> Return cached result
    |
    +-> Database query
    |   └-> Execute tuple search
    |
    v
[Result]
    - Allow / Deny
    |
    v
[Application Layer]
    - Execute operation or reject
    - Log audit entry
```

## Model Patterns

### Direct Assignment

User directly assigned to role:

```
user:alice -> admin -> organization:acme
```

Check: Can alice admin acme? YES

### Transitive Relations

Permission through hierarchy:

```
user:bob -> member -> team:eng
team:eng -> member -> project:api
project:api -> admin -> ...

Check: Can bob admin project:api? NO (must be explicit)
Check: Can bob view project:api? YES (inherited from team)
```

### Group Membership

User inherits permissions through group:

```
user:alice -> member -> group:admins
group:admins -> admin -> organization:acme

Check: Can alice admin acme? YES (through group)
```

### Role Hierarchy

Permission granted through higher role:

```
user:charlie -> editor -> document:spec
document:spec -> editor -> viewer -> ...

Check: Can charlie view document:spec? YES (editor includes view)
```

## Scalability Considerations

### Performance Optimization

- Authorization checks: O(1) with proper indexing
- Relationship lookups: O(log n) binary search
- Cache frequently checked permissions
- Batch operations when possible

### Database Indexes

Create indexes on:
- user_id (for listing user permissions)
- object_id (for listing object users)
- relation_id (for relation queries)
- user_id + relation_id + object_id (compound index)

### Caching Strategy

```go
// Cache authorization results
cache := NewPermissionCache()
cache.Set(authKey, allowed, ttl=5min)

// Check cache first
if cached, ok := cache.Get(authKey); ok {
  return cached
}

// Query OpenFGA
result := checkOpenFGA(user, relation, object)
cache.Set(authKey, result, ttl=5min)
return result
```

TTL considerations:
- Hot permissions: 5-15 minutes
- Cold permissions: 1-5 minutes
- Sensitive operations: No cache

### Batch Operations

```bash
# Write multiple relationships
curl -X POST http://localhost:8080/stores/{id}/write \
  -d '{
    "writes": [
      {"key": {"user": "user:alice", "relation": "admin", "object": "org:acme"}},
      {"key": {"user": "user:bob", "relation": "member", "object": "org:acme"}}
    ]
  }'
```

## High Availability

### Deployment Architecture

```
Load Balancer
    |
    +-- OpenFGA Pod 1
    +-- OpenFGA Pod 2
    +-- OpenFGA Pod 3
    |
    v
Database Cluster
    +-- Primary
    +-- Replica 1
    +-- Replica 2
```

### Failover Handling

- Connection pooling with retries
- Read replicas for scale-out
- Primary for writes
- Automatic failover

### State Management

- Stateless OpenFGA instances
- Shared database state
- No local caching of relationships

## Multi-Tenancy Patterns

### Separate Stores

```
Database Cluster
    |
    +-- Store: tenant_a
    +-- Store: tenant_b
    +-- Store: tenant_c
```

Advantages:
- Complete isolation
- Easy to manage per-tenant
- Billing per store

Disadvantages:
- More infrastructure
- Cross-tenant queries harder

### Shared Store with Tenant ID

```
Store
    |
    +-- Relationships with tenant_id
    +-- tenant_a:org:acme -> admin -> user:alice
    +-- tenant_b:org:techcorp -> admin -> user:bob
```

Advantages:
- Single infrastructure
- Cost efficient
- Cross-tenant queries possible

Disadvantages:
- Requires validation layer
- Higher complexity

## Integration Patterns

### Middleware Pattern

```go
func AuthorizationMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    user := extractUser(r)
    action := extractAction(r)
    resource := extractResource(r)
    
    allowed, err := checkAuth(user, action, resource)
    if !allowed {
      http.Error(w, "Unauthorized", http.StatusForbidden)
      return
    }
    
    next.ServeHTTP(w, r)
  })
}
```

### Policy Engine Pattern

```go
func applyPolicy(user, action, resource string) (bool, error) {
  // Check multiple authorization layers
  
  // Layer 1: Role-based access
  if hasRole(user, "admin") {
    return true, nil
  }
  
  // Layer 2: Resource-based access
  if isOwner(user, resource) {
    return true, nil
  }
  
  // Layer 3: OpenFGA check
  return openFGACheck(user, action, resource)
}
```

### API Gateway Pattern

```
Client Request
    |
    v
API Gateway
    - Extract context
    - Enrich request
    |
    +-> Authorization Engine
    |   └-> OpenFGA Service
    |
    v
Backend Service
```

## Debugging and Troubleshooting

### Expand Relation

Visualize permission hierarchy:

```bash
curl -X POST http://localhost:8080/stores/{id}/expand \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "admin",
      "object": "organization:acme"
    }
  }'
```

Returns tree of relationships.

### Audit Trail

Query relationship history:
- Who created relationship
- When it was created
- Who modified it
- When it was deleted

### Performance Profiling

- Monitor check latency
- Profile database queries
- Analyze cache hit rates
- Identify hot permissions

## Evolution and Maintenance

### Model Versioning

```
v1: Initial RBAC model
v2: Add ABAC attributes
v3: Add team-based permissions
v4: Add time-based access
```

Track changes and deprecate gradually.

### Migration Strategy

1. Deploy new model version
2. Run both models in parallel
3. Migrate relationships gradually
4. Monitor for discrepancies
5. Deprecate old model

### Testing Strategy

- Unit tests for model logic
- Integration tests with real data
- Performance benchmarks
- Security tests
- Chaos engineering tests
