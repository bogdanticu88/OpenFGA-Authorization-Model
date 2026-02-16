# Getting Started with OpenFGA

Complete tutorial for implementing OpenFGA authorization in your application.

## Prerequisites

- Docker and Docker Compose
- Git
- curl or Postman for API testing
- Basic understanding of authorization concepts

## Step 1: Start OpenFGA

```bash
cd deployment/docker
docker-compose up -d
```

Services running:
- OpenFGA: http://localhost:8080
- Playground: http://localhost:3000
- Postgres: localhost:5432

## Step 2: Create Authorization Store

Create a store to hold your authorization data:

```bash
curl -X POST http://localhost:8080/stores \
  -H 'Content-Type: application/json' \
  -d '{"name":"my-store"}'
```

Response:
```json
{
  "id": "store-id-abc123",
  "name": "my-store"
}
```

Save the store ID.

## Step 3: Write Authorization Model

Choose a model from models/ directory based on your needs.

For RBAC example:

```bash
curl -X POST http://localhost:8080/stores/{store-id}/authorization-models \
  -H 'Content-Type: application/json' \
  -d @models/rbac/model.fga
```

## Step 4: Create Relationships

Add relationships defining who has what role:

```bash
curl -X POST http://localhost:8080/stores/{store-id}/write \
  -H 'Content-Type: application/json' \
  -d '{
    "writes": [
      {
        "key": {
          "user": "user:alice",
          "relation": "admin",
          "object": "organization:acme"
        }
      },
      {
        "key": {
          "user": "user:bob",
          "relation": "member",
          "object": "organization:acme"
        }
      }
    ]
  }'
```

## Step 5: Check Authorization

Check if user has permission:

```bash
curl -X POST http://localhost:8080/stores/{store-id}/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "admin",
      "object": "organization:acme"
    }
  }'
```

Response:
```json
{
  "allowed": true
}
```

## Step 6: List Permissions

List all resources user can access:

```bash
curl -X POST http://localhost:8080/stores/{store-id}/list-objects \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user:alice",
    "relation": "admin",
    "type": "organization"
  }'
```

Response:
```json
{
  "objects": [
    "organization:acme"
  ]
}
```

## Using the Playground

OpenFGA Playground provides visual interface:

1. Open http://localhost:3000
2. Create authorization model
3. Add relationships
4. Test authorization checks
5. Visualize permission structure

## API Endpoints

### Core Endpoints

- POST /stores - Create store
- GET /stores/{id} - Get store details
- POST /stores/{id}/authorization-models - Write model
- GET /stores/{id}/authorization-models/{model-id} - Get model
- POST /stores/{id}/write - Write relationships
- POST /stores/{id}/check - Check authorization
- POST /stores/{id}/expand - Expand permissions
- POST /stores/{id}/list-objects - List accessible objects
- POST /stores/{id}/list-users - List users with access

## Code Examples

### Go Client

```go
import "github.com/openfga/go-sdk/client"

// Create client
configuration := configuration.NewConfiguration(...)
client := client.NewSdkClient(configuration)

// Check access
response := client.Check(ctx).
  StoreId(storeID).
  Body(client.CheckRequest{...}).
  Execute()
```

See examples/go/ for complete example.

### Python Client

```python
from openfga_sdk import OpenFgaClient

client = OpenFgaClient(
  api_url="http://localhost:8080",
  store_id="store-id"
)

response = client.check({
  "user": "user:alice",
  "relation": "admin",
  "object": "organization:acme"
})
```

See examples/python/ for complete example.

### Node.js Client

```javascript
const { OpenFgaClient } = require('@openfga/sdk');

const client = new OpenFgaClient({
  apiUrl: 'http://localhost:8080',
  storeId: 'store-id'
});

const response = await client.check({
  user: 'user:alice',
  relation: 'admin',
  object: 'organization:acme'
});
```

See examples/nodejs/ for complete example.

## Common Patterns

### User Assignment to Role

```bash
# Add user to admin role
curl -X POST http://localhost:8080/stores/{id}/write \
  -d '{
    "writes": [{
      "key": {
        "user": "user:alice",
        "relation": "admin",
        "object": "organization:acme"
      }
    }]
  }'
```

### Hierarchical Permissions

Define parent-child relationships:

```fga
type project
  relations
    define organization: [organization]
    define editor: [user] or editor from parent_team
```

### Conditional Access

Implement time-based or context-based access:
- Business hours access
- Location-based access
- Device-based access
- Require additional conditions in application

## Troubleshooting

### Authorization always returns false

1. Check store ID is correct
2. Verify authorization model is written
3. Confirm relationship exists in database
4. Check user and object identifiers match exactly

### Performance issues

1. Monitor database query times
2. Add indexes on frequently queried relations
3. Implement caching for repeated checks
4. Use batch operations

### Model validation errors

1. Verify DSL syntax is correct
2. Check type definitions are valid
3. Ensure relation names are properly defined

## Next Steps

1. Choose appropriate model for your use case
2. Implement in your application
3. Integrate with your backend
4. Add audit logging
5. Test thoroughly before production deployment

See ARCHITECTURE.md for advanced topics.
