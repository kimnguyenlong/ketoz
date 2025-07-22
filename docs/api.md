# Ketoz REST API Documentation

- [Endpoints](#endpoints)
  - [Identity](#identity)
    - [Create Identity](#create-identity)
    - [Get Identity](#get-identity)
    - [List Identities](#list-identities)
    - [Add Child Identity](#add-child-identity)
    - [List Child Identities](#list-child-identities)
    - [List Permissions for Identity](#list-permissions-for-identity)
  - [Resource](#resource)
    - [Create Resource](#create-resource)
    - [Get Resource](#get-resource)
    - [List Resources](#list-resources)
    - [Add Child Resource](#add-child-resource)
    - [List Child Resources](#list-child-resources)
  - [Permission](#permission)
    - [Check Permission](#check-permission)
    - [Grant Permission](#grant-permission)
    - [Revoke Permission](#revoke-permission)
    - [Deny Permission](#deny-permission)
    - [Delete Denied Permission](#delete-denied-permission)
- [Schemas](#schemas)

## Endpoints

### Identity
Manage identities and their hierarchical relationships.

#### Create Identity
Create a new identity.

- **Method**: POST
- **URL**: `/identities`
- **Request Body**:
  | Field | Datatype | Description |
  |-------|----------|-------------|
  | `id`  | String   | The unique identifier for the new identity. |
- **Description**: Creates a new identity with the specified ID.
- **cURL Example**:
  ```bash
  curl -X POST http://ketoz/api/identities \
    -H "Content-Type: application/json" \
    -d '{"id": "swe"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if the ID is invalid or already exists.

#### Get Identity
Retrieve details of a specific identity.

- **Method**: GET
- **URL**: `/identities/:id`
- **Parameters**:
  - `id`: The identity ID (e.g., `swe`) (String)
- **Description**: Fetches details of the identity specified by `id`.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/identities/swe
  ```
- **Response**:
  - `200 OK` on success
  - `404 Not Found` if the identity does not exist.

#### List Identities
Retrieve a list of all identities.

- **Method**: GET
- **URL**: `/identities`
- **Description**: Returns a list of all identities.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/identities
  ```
- **Response**:
  - `200 OK` on success.

#### Add Child Identity
Add a child identity to a parent identity.

- **Method**: POST
- **URL**: `/identities/:id/children`
- **Parameters**:
  - `id`: The parent identity ID (e.g., `swe`) (String)
- **Request Body**:
  | Field      | Datatype | Description |
  |------------|----------|-------------|
  | `child_id` | String   | The ID of the child identity to associate with the parent. |
- **Description**: Associates a child identity with the specified parent identity.
- **cURL Example**:
  ```bash
  curl -X POST http://ketoz/api/identities/swe/children \
    -H "Content-Type: application/json" \
    -d '{"child_id": "backend"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if the child ID is invalid.
  - `404 Not Found` if the parent identity does not exist.

#### List Child Identities
Retrieve a list of child identities for a specific identity.

- **Method**: GET
- **URL**: `/identities/:id/children`
- **Parameters**:
  - `id`: The parent identity ID (e.g., `swe`) (String)
- **Description**: Returns a list of child identities for the specified identity.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/identities/swe/children
  ```
- **Response**:
  - `200 OK` on success, 
  - `404 Not Found` if the parent identity does not exist.

#### List Permissions for Identity
Retrieve permissions associated with a specific identity.

- **Method**: GET
- **URL**: `/identities/:id/permissions`
- **Parameters**:
  - `id`: The identity ID (e.g., `swe`) (String)
- **Description**: Returns a list of permissions assigned to the specified identity.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/identities/swe/permissions
  ```
- **Response**:
  - `200 OK` on success
  - `404 Not Found` if the identity does not exist.

### Resource
Manage resources and their hierarchical relationships.

#### Create Resource
Create a new resource.

- **Method**: POST
- **URL**: `/resources`
- **Request Body**:
  | Field | Datatype | Description |
  |-------|----------|-------------|
  | `id`  | String   | The unique identifier for the new resource. |
- **Description**: Creates a new resource with the specified ID.
- **cURL Example**:
  ```bash
  curl -X POST http://ketoz/api/resources \
    -H "Content-Type: application/json" \
    -d '{"id": "projects"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if the ID is invalid or already exists.

#### Get Resource
Retrieve details of a specific resource.

- **Method**: GET
- **URL**: `/resources/:id`
- **Parameters**:
  - `id`: The resource ID (e.g., `projects`) (String)
- **Description**: Fetches details of the resource specified by `id`.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/resources/projects
  ```
- **Response**:
  - `200 OK` on success.
  - `404 Not Found` if the resource does not exist.

#### List Resources
Retrieve a list of all resources.

- **Method**: GET
- **URL**: `/resources`
- **Description**: Returns a list of all resources.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/resources
  ```
- **Response**:
  - `200 OK` on success.

#### Add Child Resource
Add a child resource to a parent resource.

- **Method**: POST
- **URL**: `/resources/:id/children`
- **Parameters**:
  - `id`: The parent resource ID (e.g., `projects`) (String)
- **Request Body**:
  | Field      | Datatype | Description |
  |------------|----------|-------------|
  | `child_id` | String   | The ID of the child resource to associate with the parent. |
- **Description**: Associates a child resource with the specified parent resource.
- **cURL Example**:
  ```bash
  curl -X POST http://ketoz/api/resources/projects/children \
    -H "Content-Type: application/json" \
    -d '{"child_id": "projects_iam"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if the child ID is invalid.
  - `404 Not Found` if the parent resource does not exist.

#### List Child Resources
Retrieve a list of child resources for a specific resource.

- **Method**: GET
- **URL**: `/resources/:id/children`
- **Parameters**:
  - `id`: The parent resource ID (e.g., `projects`) (String)
- **Description**: Returns a list of child resources for the specified resource.
- **cURL Example**:
  ```bash
  curl -X GET http://ketoz/api/resources/projects/children
  ```
- **Response**:
  - `200 OK` on success.
  - `404 Not Found` if the parent resource does not exist.

### Permission
Manage permissions for identities and resources.

#### Grant Permission
Grant a permission to an identity for a resource.

- **Method**: POST
- **URL**: `/permissions/granted`
- **Request Body**:
  | Field         | Datatype | Description |
  |---------------|----------|-------------|
  | `identity_id` | String   | The ID of the identity to grant the permission to. |
  | `resource_id` | String   | The ID of the resource the permission applies to. |
  | `permission`  | String   | The permission to grant (e.g., `owners`). Either `owners`, `editors`, `child_creators` or `viewers`|
- **Description**: Grants the specified permission to the identity for the resource.
- **cURL Example**:
  ```bash
  curl -X POST http://ketoz/api/permissions/granted \
    -H "Content-Type: application/json" \
    -d '{"identity_id": "swe", "resource_id": "projects", "permission": "owners"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if parameters are invalid.
  - `404 Not Found` if identity or resource does not exist.

#### Revoke Permission
Revoke a permission from an identity for a resource.

- **Method**: DELETE
- **URL**: `/permissions/granted`
- **Request Body**:
  | Field         | Datatype | Description |
  |---------------|----------|-------------|
  | `identity_id` | String   | The ID of the identity to revoke the permission from. |
  | `resource_id` | String   | The ID of the resource the permission applies to. |
  | `permission`  | String   | The permission to revoke (e.g., `owners`). Either `owners`, `editors`, `child_creators` or `viewers`|
- **Description**: Removes the specified permission from the identity for the resource.
- **cURL Example**:
  ```bash
  curl -X DELETE http://ketoz/api/permissions/granted \
    -H "Content-Type: application/json" \
    -d '{"identity_id": "swe", "resource_id": "projects", "permission": "owners"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if parameters are invalid.
  - `404 Not Found` if the permission does not exist.

#### Deny Permission
Deny a permission to an identity for a resource.

- **Method**: POST
- **URL**: `/permissions/denied`
- **Request Body**:
  | Field         | Datatype | Description |
  |---------------|----------|-------------|
  | `identity_id` | String   | The ID of the identity to deny the permission to. |
  | `resource_id` | String   | The ID of the resource the permission applies to. |
  | `permission`  | String   | The permission to deny (e.g., `viewers`). Either `owners`, `editors`, `child_creators` or `viewers`|
- **Description**: Explicitly denies the specified permission to the identity for the resource.
- **cURL Example**:
  ```bash
  curl -X POST http://ketoz/api/permissions/denied \
    -H "Content-Type: application/json" \
    -d '{"identity_id": "swe", "resource_id": "projects", "permission": "viewers"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if parameters are invalid.
  - `404 Not Found` if identity or resource does not exist.

#### Delete Denied Permission
Remove a denied permission from an identity for a resource.

- **Method**: DELETE
- **URL**: `/permissions/denied`
- **Request Body**:
  | Field         | Datatype | Description |
  |---------------|----------|-------------|
  | `identity_id` | String   | The ID of the identity to remove the denied permission from. |
  | `resource_id` | String   | The ID of the resource the permission applies to. |
  | `permission`  | String   | The permission to remove from denial (e.g., `viewers`). Either `owners`, `editors`, `child_creators` or `viewers`|
- **Description**: Removes the denial of the specified permission for the identity and resource.
- **cURL Example**:
  ```bash
  curl -X DELETE http://ketoz/api/permissions/denied \
    -H "Content-Type: application/json" \
    -d '{"identity_id": "swe", "resource_id": "projects", "permission": "viewers"}'
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if parameters are invalid.
  - `404 Not Found` if the denied permission does not exist.

#### Check Permission
Check if an identity has a specific permission on a resource.

- **Method**: GET
- **URL**: `/permissions/check`
- **Query Parameters**:
  - `identity_id`: The identity ID (e.g., `swe`) (String)
  - `action`: The action to check (e.g., `view`) (String, either `view`, `edit`, `create_child` or `delete`)
  - `resource_id`: The resource ID (e.g., `projects`) (String)
- **Description**: Verifies if the specified identity has the given permission for the resource.
- **cURL Example**:
  ```bash
  curl -X GET "http://ketoz/api/permissions/check?identity_id=swe&action=view&resource_id=projects"
  ```
- **Response**:
  - `200 OK` on success.
  - `400 Bad Request` if parameters are invalid.

## Schemas

### Response

| Field        | Type     | Description                                         |
|--------------|----------|-----------------------------------------------------|
| `code`       | int      | Status code of the response.                        |
| `message`    | String   | Description of the response.                        |
| `data`       | Object   | Response's detail (Nullable).                       |

### List Response

| Field        | Type         | Description                                         |
|--------------|--------------|-----------------------------------------------------|
| `code`       | int          | Status code of the response.                        |
| `message`    | String       | Description of the response.                        |
| `records`    | List\<Object\> | List of items.                                    |

### Identity

| Field        | Type     | Description                                         |
|--------------|----------|-----------------------------------------------------|
| `id`         | String   | Unique identifier of the identity.                  |

### Resource

| Field        | Type     | Description                                         |
|--------------|----------|-----------------------------------------------------|
| `id`         | String   | Unique identifier of the resource.                  |

### Permission

| Field        | Type     | Description                                               |
|--------------|----------|-----------------------------------------------------------|
| `identity_id`| String   | Unique identifier of the identity.                        |
| `resource_id`| String   | Unique identifier of the resource.                        |
| `permission` | String   | Either `owners`, `editors`, `child_creators` or `viewers` |