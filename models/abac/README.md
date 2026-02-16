# ABAC - Attribute-Based Access Control

Attribute-based authorization model providing fine-grained, policy-driven access control.

## Model Overview

ABAC decisions are based on attributes of:
- Subject: User attributes (department, role, clearance level)
- Resource: Resource attributes (type, classification, owner)
- Action: Requested action (read, write, delete)
- Environment: Context attributes (time, location, network)

## Core Concepts

### Subjects
Users with attributes:
- department
- role
- clearance_level
- location

### Resources
Protected resources with attributes:
- resource_type: document, database, api, etc.
- classification: public, internal, confidential, restricted
- owner
- created_date
- department

### Actions
Operations on resources:
- read
- write
- delete
- execute
- admin

### Policies
Rules matching subject, resource, and action attributes:

```
IF subject.department == resource.department AND
   subject.clearance >= resource.classification AND
   action IN ["read", "write"] AND
   time.hour BETWEEN 9 AND 17
THEN allow
```

## Relationship Examples

### Define user attributes

```
user:bob#department@organization:engineering
user:bob#role@string:senior-engineer
user:bob#clearance@string:confidential
```

### Define resource attributes

```
resource:api#resource_type@string:api
resource:api#classification@string:confidential
resource:api#owner@user:alice
```

### Create permission

```
permission:p1#subject@user:bob
permission:p1#resource@resource:api
permission:p1#action@string:read,write
```

### Define policy

```
policy:engineers_api#subject_type@string:engineer
policy:engineers_api#resource_type@string:api
policy:engineers_api#action@string:read,write
policy:engineers_api#classification@string:confidential
policy:engineers_api#effect@string:allow
```

## Access Decision Logic

1. Collect subject attributes
2. Collect resource attributes
3. Collect action and environment attributes
4. Match against all applicable policies
5. Combine results (allow wins if any match)

## Use Cases

### Time-based access

Allow access during business hours only:

```
permission#conditions@string:business_hours_9_to_17
```

### Department-based access

Allow access within department:

```
policy#department@department:engineering
```

### Classification-based access

Allow access to resources of certain classification:

```
policy#classification@string:internal,public
```

### Conditional access

Multiple conditions required:

```
permission#conditions@string:vpn_required,mfa_enabled
```

## Advantages

- Highly flexible and dynamic
- Easy to add new attributes
- Policies defined separately from code
- Supports complex conditions
- Audit trail friendly

## Disadvantages

- More complex to understand
- Higher computational overhead
- Requires careful policy design
- Policy conflicts possible

## Best Practices

1. Start with simple attributes
2. Define clear attribute meanings
3. Document all policies
4. Regular policy reviews
5. Test policy changes
6. Monitor policy effectiveness

## Performance Considerations

- Attribute lookup: O(1)
- Policy matching: O(n) where n = policies
- Condition evaluation: Depends on complexity
- Cache frequently accessed policies

## Common Mistakes

1. Overly complex policies
2. Conflicting policies
3. Missing attribute validation
4. Inefficient condition evaluation
5. Poor audit logging

## Testing

Test each policy with:
- Positive cases (should allow)
- Negative cases (should deny)
- Edge cases (boundary conditions)
- Performance benchmarks
