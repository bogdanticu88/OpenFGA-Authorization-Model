# RBAC - Role-Based Access Control

Traditional role-based authorization model suitable for organizational hierarchies.

## Model Overview

This model defines a hierarchy of objects with role-based permissions:

- Organization: Top-level container
- Team: Team within organization
- Project: Project within team
- Resource: Resource within project
- Document: Document within resource

Roles: admin, member, owner, editor, viewer

## Roles and Permissions

### Organization Level

- admin: Full control over organization
- member: Standard member access
- viewer: Read-only access

### Team Level

- owner: Can manage team and resources
- member: Can access team resources
- viewer: Read-only team access

### Project/Resource/Document Level

- owner: Full control
- editor: Can modify content
- viewer: Read-only access

## Relationships

Relationships are inherited from parent objects:

```
organization -> team -> project -> resource -> document
   admin         owner     owner      owner       owner
   member        member    editor     editor      editor
   viewer        viewer    viewer     viewer      viewer
```

## Usage Examples

### Check if user is admin of organization

```
openfga check --user user:alice --relation admin --object organization:acme
```

Expected response: True (alice is admin)

### Check if user can edit document

```
openfga check --user user:bob --relation editor --object document:user-service-spec
```

Expected response: True (bob has editor role through project)

### List all users who can view resource

```
openfga list --relation viewer --object resource:user-service
```

Expected response: alice, bob, charlie (inherited from project)

## Best Practices

1. Use specific roles matching your organization structure
2. Define clear role responsibilities
3. Regularly audit permissions
4. Use inheritance to reduce redundancy
5. Document all relationships

## Common Patterns

### Add user to team

```
team:platform#member@user:alice
```

### Promote user to owner

```
project:backend-api#owner@user:bob
```

### Add entire team as viewers

```
document:spec#viewer@team:platform#member
```

## Performance Characteristics

- Authorization check: O(1) with indexing
- List users: O(n) where n = number of users
- Audit trail: All changes logged

## Limitations

- Fixed role hierarchy
- No dynamic attributes
- Limited conditional logic
- Requires predefined roles

For more flexibility, see ABAC model.
