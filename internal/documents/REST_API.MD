
Uniform Resource Identifier (URI)

### restfulapi
| HTTP Method | Behavior | Idempotent | Safe | Cacheable |
|-------------|------------------------------------------------|------------|------------|-----------|
| **GET** | Retrieves a resource. | Yes | Yes | Yes |
| **POST** | Sends data to create/update a resource. | No | No | No |
| **PUT** | Replaces or updates a resource. | Yes | No | No |
| **DELETE** | Deletes a resource. | Yes | No | No |
| **HEAD** | Retrieves headers only, without body. | Yes | Yes | Yes |
| **OPTIONS** | Describes communication options for a resource. | Yes | Yes | No |
| **TRACE** | Echoes received request for diagnostic purposes. | Yes | No | No |
| **CONNECT** | Establishes a tunnel, often used for HTTPS. | Yes | No | No |


https://restfulapi.net/rest-api-best-practices/

The reserved characters such as /, ?, and # must be used only in accordance with their defined roles.
For example, the (/) character should separate hierarchical components, 
while query parameters (?) are intended for non-hierarchical data filtering and sorting purposes.
````
// Correct use of '/' to denote hierarchical relationships:
GET /users/123/orders/456  
POST /books/789/reviews

// Correct use of '?' to specify query parameters:
GET /products?category=electronics&sort=price_desc

// Correct use of '#' for fragment identifiers:
GET /documentation#section2  
````

'/' to denote hierarchical relationships:
'?' to specify query parameters:
'#' for fragment identifiers

| Category | Range | Description | Examples |
|---|---|---|---|
| **Informational** | `1xx` | Indicates that the request has been received and is being processed. | `100 Continue`, `101 Switching Protocols` |
| **Successful** | `2xx` | Indicates that the request was successfully received, understood, and accepted. | `200 OK`, `201 Created`, `204 No Content` |
| **Redirection** | `3xx` | Indicates that further action is needed to fulfill the request. | `301 Moved Permanently`, `302 Found`, `304 Not Modified` |
| **Client Error** | `4xx` | Indicates that there was an error with the request. | `400 Bad Request`, `401 Unauthorized`, `404 Not Found` |
| **Server Error** | `5xx` | Indicates that the server failed to fulfill a valid request. | `500 Internal Server Error`, `502 Bad Gateway`, `503 Service Unavailable` |




