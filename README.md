# Org Chart Builder & Real-Time Chat Application

## 🛠 Technology Stack

The application is built on a high-speed, concurrent foundation using **Go**, and leverages a **dual-database approach** for optimized performance in both relational (PostgreSQL) and real-time (ScyllaDB) workloads.

| Category | Technology | Purpose |
| :--- | :--- | :--- |
| **Backend** | GIN (Go) | High-performance HTTP web framework. |
| **Relational DB** | PostgreSQL | Stores Employee, Role, and User data (Org Chart). |
| **Relational ORM** | GORM | Object-Relational Mapper for PostgreSQL. |
| **NoSQL DB** | ScyllaDB | Real-time, scalable storage for Chat Messages and Connection status (Cassandra-compatible). |
| **Auth** | JWT | JSON Web Token authentication for securing APIs. |
| **Communication** | WebSockets & http requests |Websockets for chat application and http requests for org chart builder|

---

## Key Features

### Organizational Chart & Employee Management (PostgreSQL/GORM)

This module handles the core hierarchical structure and employee data, relying on PostgreSQL for data integrity and complex relationship modeling.

* **Role-Based Hierarchy:** Defines the formal structure (e.g., CEO ->Manager -> Team Lead).
* **Full CRUD Support:** Comprehensive **C**reate, **R**ead, **U**pdate, and **D**elete operations for both **Roles** and **Employees**.
* **Hierarchy Mapping:** The `Employee` schema uses `manager_id` (a self-referencing foreign key) to build and maintain the Manager-Subordinate tree structure.
* **Nested Data Fetching:** Utilizes `DB.Preload("Role")` in GORM to fetch complex, nested JSON objects (e.g., an Employee with their associated Role details and Manager information) in a single, efficient query.

### Real-Time Chat System (GIN/ScyllaDB/WebSockets)

This module is designed for scalability and high availability, supporting instant messaging across multiple devices with low latency.

* **Secure Connection:** **JWT validation** is required to establish a WebSocket connection.
* **Scalable Message Store:** Messages are stored in **ScyllaDB** (a high-performance, Cassandra-compatible NoSQL database) for high-volume, low-latency performance and high availability.
* **Multi-Device Support:** The system tracks and manages multiple active connections per user via the `Connections` schema, allowing users to chat seamlessly from different devices.

---

## API Endpoints

### PUBLIC ENDPOINTS (For Authentication)

| Resource | Method | Endpoint | Description |
| :--- | :--- | :--- | :--- |
| User | `POST` | `/api/user/register` | Create a new user account. |
| User | `POST` | `/api/user/login` | Log in and receive JWT **Access** and **Refresh** tokens. |
| Token | `POST` | `/api/token/refresh` | Obtain a new Access Token using a valid Refresh Token. |
| Token | `POST` | `/api/token/isvalid` | Validate the status and expiry of a Token. |

### PROTECTED ENDPOINTS (Requires `Authorization: Bearer <JWT>` Header [or] Token as query parameter)

| Resource | Method | Endpoint | Description |
| :--- | :--- | :--- | :--- |
| Users | `GET` | `/api/users` | List all registered user accounts. |
| Roles | `POST` | `/api/roles` | Create a new organizational role. |
| Roles | `GET` | `/api/roles` | Retrieve all organizational roles. |
| Roles | `GET` | `/api/roles/:id` | Get a specific role by ID. |
| Roles | `PATCH` | `/api/roles/:id` | Update a role's details by ID. |
| Roles | `DELETE` | `/api/roles/:id` | Delete a role by ID. |
| Employees | `POST` | `/api/employees` | Create a new employee record. |
| Employees | `GET` | `/api/employees` | Retrieve all employee records. |
| Employees | `GET` | `/api/employees/summary` | Retrives a formatted employee record with id and name (firstname + lastname) |
| Employees | `PATCH` | `/api/employees/:id` | Update an employee's details (e.g., change manager, role). |
| Employees | `DELETE` | `/api/employees/:id` | Delete an employee record. |
| **WebSocket** | `GET` | `/ws` | Establishes the real-time WebSocket connection. |
| Messages | `GET` | `/api/messages` | Fetch chat conversation history using query parameters (e.g., conversation key). |

---

## Database Schemas

### PostgreSQL Schemas (GORM)

These schemas manage the core user accounts and the organizational hierarchy.

#### 1. User Schema (Authentication)
```go
User {
    id        uint PRIMARY KEY UNIQUE,
    username  string,
    email     string,
    password  string  // Stored as a hashed value
}

2. Role Schema (Organization Roles)GoRole {
    id    uint PRIMARY KEY UNIQUE,
    role  string  // e.g., "CEO", "Manager", "Developer"
}

3. Employee Schema (Org Chart Hierarchy)GoEmployee {
    id          uint PRIMARY KEY UNIQUE,
    first_name  string,
    last_name   string,
    role_id     uint,   // Foreign Key reference to Role(id)
    manager_id  uint    // Self-referencing Foreign Key to Employee(id)
}
```
#### ScyllaDB Schemas
These schemas are optimized for high-volume, low-latency, real-time messaging and connection tracking.
```go
1. Connections Schema (Multi-Device Support)SQLCREATE TABLE connections (
    id              UUID PRIMARY KEY,       -- Unique connection/device ID
    client_id       INT,                    -- The User ID (sender_id)
    is_active       BOOL,                   -- Connection status flag
    connected_at    TIMESTAMP,
    disconnected_at TIMESTAMP
);

2. Messages Schema (Chat History)SQLCREATE TABLE messages (
    id               UUID PRIMARY KEY,
    conversation_key TEXT NOT NULL,          -- Key for finding chat history (e.g., sorted user IDs)
    sender_id        INT NOT NULL,           -- ID of the user who sent the message
    receiver_id      INT NOT NULL,
    content          TEXT,
    sent_at          TIMESTAMP,
    read_at          TIMESTAMP,
    PRIMARY KEY (conversation_key, sent_at) -- Primary key for efficient time-series queries
) WITH CLUSTERING ORDER BY (sent_at DESC);
```
#### Go WebSocket Hub Structure
The central Hub manages all real-time traffic and client connections, ensuring concurrent and efficient handling of chat events using Go's concurrency features.Client StructRepresents a single active connection (device) for a user
```go 
type Client struct {
    ID      uint                // User ID (matches Employee/User ID)
    ConnID  gocql.UUID          // Unique Connection ID from ScyllaDB
    Conn    *websocket.Conn     // The actual WebSocket connection
    Receive chan []byte         // Channel to receive messages from the client
    Hub     *Hub
}
```
Hub Struct\
Manager for all active users
```go
type Hub struct {
    // Stores active clients: Map[UserID] -> Map[ConnID] -> Client instance
    Clients map[uint]map[gocql.UUID]*Client

    // Channel for incoming messages from any client
    Forward chan []byte

    // Channels for client registration/deregistration
    Join    chan *Client
    Leave   chan *Client
}
```
### Example Data
#### Role Table

| r_id | Role      |
| ---- | --------- |
| 1    | CEO       |
| 2    | Manager   |
| 3    | Team Lead |
| 4    | Developer |

#### Employee Table
| e_id | First Name | Last Name | Role ID       | Manager ID      |
| ---- | ---------- | --------- | ------------- | --------------- |
| 1    | John       | Doe       | 1 (CEO)       | NULL            |
| 2    | Priya      | Nair      | 2 (Manager)   | 1 (John Doe)    |
| 3    | Rahul      | Mehta     | 3 (Team Lead) | 2 (Priya Nair)  |
| 4    | Meena      | Iyer      | 4 (Developer) | 3 (Rahul Mehta) |
| 5    | David      | Paul      | 4 (Developer) | 3 (Rahul Mehta) |
