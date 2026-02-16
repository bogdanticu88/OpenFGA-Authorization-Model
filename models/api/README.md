# API Authorization

Fine-grained API endpoint authorization using keys, tokens, and scopes.

## Model Overview

This model implements OAuth2-like authorization for APIs:

- Application: Top-level API application
- API Key: Long-lived credentials for service accounts
- API Token: Short-lived tokens with scopes
- Scope: Grouping of permissions
- Endpoint: Protected API endpoints
- Service Account: Application-to-application authentication

## Authentication vs Authorization

### Authentication
Verify that the caller is who they claim to be:
- API key validation
- Token signature verification
- MTLS certificate validation

### Authorization
Verify that the caller has permission for the requested operation:
- Scope validation
- Rate limit checking
- Endpoint access control

## API Key Model

Long-lived credentials for service-to-service communication:

```
api_key:prod_key#application@application:api_v1
api_key:prod_key#owner@user:alice
api_key:prod_key#scopes@string:users:read,users:write
api_key:prod_key#status@string:active
```

Properties:
- Owned by user or application
- Multiple scopes
- Active/inactive status
- Creation date
- Last used date

## Token Model

Short-lived tokens with expiration:

```
api_token:token_abc#key@api_key:prod_key
api_token:token_abc#scopes@api_scope:users_read
api_token:token_abc#expires_at@string:2025-01-01T00:00:00Z
```

Properties:
- Derived from API key
- Subset of key scopes
- Expiration time
- Revocation capability

## Scope Model

Logical grouping of permissions:

```
api_scope:users_read#name@string:users:read
api_scope:users_read#permissions@string:list_users,get_user

api_scope:users_write#name@string:users:write
api_scope:users_write#permissions@string:create_user,update_user
```

Properties:
- Human-readable name
- List of permissions
- Application context

## Endpoint Model

Protected API endpoints with access rules:

```
api_endpoint:get_users#method@string:GET
api_endpoint:get_users#authenticated@user:alice
api_endpoint:get_users#authorized@user:alice
api_endpoint:get_users#ratelimit@string:100_per_minute
```

Properties:
- HTTP method
- Public/authenticated/authorized
- Rate limiting
- Owner management

## Authorization Flow

1. Client presents API key or token
2. Server validates credentials (authentication)
3. Server extracts scopes from key/token
4. Server checks if scopes include required permission
5. Server checks rate limits
6. Server checks endpoint authorization
7. Allow or deny request

## Usage Examples

### Check if API key has scope

```
openfga check \
  --user api_key:prod_key \
  --relation "scopes" \
  --object "api_scope:users_read"
```

### Check if user can manage endpoint

```
openfga check \
  --user user:alice \
  --relation "owner" \
  --object "api_endpoint:get_users"
```

### List scopes for API key

```
openfga list \
  --relation "scopes" \
  --from api_key:prod_key
```

## Rate Limiting Strategy

Define rate limits on endpoints:

```
api_endpoint:create_user#ratelimit@string:10_per_minute
api_endpoint:list_users#ratelimit@string:100_per_minute
```

Enforce in application:
- Track requests per key/token
- Reject when limit exceeded
- Return 429 Too Many Requests

## Service Accounts

Application-to-application authentication:

```
service_account:ci_system#keys@api_key:ci_key
service_account:ci_system#roles@string:read-only
service_account:ci_system#environment@string:production
```

Use cases:
- CI/CD systems
- Batch jobs
- System integrations

## Security Practices

1. Generate strong API keys (32+ bytes)
2. Rotate keys regularly (annually minimum)
3. Use short-lived tokens (15-60 min expiration)
4. Scope tokens to minimum required permissions
5. Audit all API key usage
6. Revoke unused keys immediately
7. Use HTTPS for all API calls
8. Implement rate limiting per key
9. Add request signing for sensitive operations

## Token Generation Workflow

```
Client request + API Key
    |
    v
Validate API Key (authentication)
    |
    v
Check key is active
    |
    v
Extract scopes from key
    |
    v
Generate short-lived token with scopes
    |
    v
Return token to client
```

## Debugging Issues

### Token rejected
1. Verify token not expired
2. Check token signature
3. Confirm scopes include required permission
4. Check endpoint authorization

### Rate limit exceeded
1. Increase key rate limit
2. Distribute load across keys
3. Implement client-side throttling

### Scope insufficient
1. Add required scope to key
2. Generate new token with scope
3. Use alternative key/token
