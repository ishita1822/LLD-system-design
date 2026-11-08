# System Design Fundamentals

System design involves designing systems that are scalable, reliable, available, efficient, and capable of handling failures.

These are the fundamental concepts used when designing distributed systems, backend services, and large-scale applications.

## 1. Latency

### Definition
Latency is the time taken by a system to process a request and return a response.

**Example:**

A user requests their profile from an API.

- Request sent: 10:00:00.000
- Response received: 10:00:00.100
- Latency: 100 milliseconds

Lower latency means faster responses.

### How to reduce latency
- Use caching (Redis).
- Optimize database queries.
- Add database indexes.
- Use CDNs.
- Reduce unnecessary network calls.
- Optimize application code.
- Use asynchronous processing where appropriate.

### Real-world example
A payment API should respond quickly so users don't experience unnecessary delays during transactions.

---

## 2. Availability

### Definition
Availability is the percentage of time a system remains operational and accessible to users when needed.

High availability means minimizing downtime and ensuring users can access the application even when individual components fail.

### Availability percentages

| Availability | Approximate downtime/year |
|---|---|
| 99% | 3.65 days |
| 99.9% | 8.76 hours |
| 99.99% | 52.6 minutes |
| 99.999% | 5.26 minutes |

These are commonly called the number of nines.

### How to achieve high availability
- Multiple application servers.
- Load balancing.
- Database replication.
- Automatic failover.
- Health checks.
- Multiple availability zones.
- Monitoring and alerting.

### Real-world example
If one application server crashes, a load balancer redirects incoming requests to another healthy server.

The application remains accessible.

---

## 3. Reliability

### Definition
Reliability is the ability of a system to perform its intended functions correctly and consistently over a period of time.

A system can be available but unreliable.

**Example:**

Imagine a banking application that is always accessible but occasionally transfers ₹10,000 instead of ₹1,000.

The application is available, but it is not reliable.

### How to improve reliability
- Proper error handling.
- Input validation.
- Automated testing.
- Monitoring.
- Data integrity checks.
- Fault isolation.
- Safe deployment practices.

### Availability vs Reliability

| Availability | Reliability |
|---|---|
| Is the system accessible? | Does the system work correctly? |
| Focuses on uptime. | Focuses on correct operation over time. |
| Measured using uptime percentage. | Can be measured using failure rates and successful operation over time. |

---

## 4. Consistency

### Definition
Consistency means that data operations respect the consistency guarantees defined by a system.

In distributed systems, multiple database replicas may contain the same data.

However, updates may not reach all replicas immediately.

### Example

Suppose a user has ₹10,000 in their account.

They transfer ₹2,000.

The primary database is updated:

- Primary DB: ₹8,000
- Replica DB: ₹10,000

If a request reads from the replica before replication completes, it may return stale data.

This is a consistency problem.

### Types of consistency

#### Strong consistency

After a successful write, subsequent reads are guaranteed to return the latest committed value according to the consistency model.

Example:

A user transfers ₹2,000 and immediately checks their balance.

They see ₹8,000.

Useful for:
- Banking balances.
- Financial transactions.
- Inventory reservations.

#### Eventual consistency

Replicas may temporarily contain different values, but they eventually converge if no new updates occur.

Example:

A user updates their profile picture.

One server displays the new picture while another temporarily displays the old picture.

After replication, both display the updated picture.

Useful for:
- Social media likes.
- View counts.
- Product reviews.
- Distributed caches.

### Strong consistency vs Eventual consistency

| Strong consistency | Eventual consistency |
|---|---|
| Latest committed data is visible to subsequent reads. | Temporary stale reads are possible. |
| Requires coordination in distributed systems. | Can reduce coordination overhead. |
| Useful for financial operations. | Useful for many social features. |

---

## 5. Scalability

### Definition
Scalability is the ability of a system to handle increasing workloads by adding resources.

Example:

An application initially supports 1,000 users.

After six months, it needs to support 1 million users.

The system should be able to handle this growth.

### Types of scalability

#### Vertical scaling (Scale Up)

Increasing the capacity of an existing server.

Examples:
- Increasing CPU.
- Increasing RAM.
- Increasing storage.

Advantages:
- Relatively simple.
- Less distributed-system complexity.

Disadvantages:
- Hardware limitations.
- Can become expensive.

#### Horizontal scaling (Scale Out)

Adding more servers to distribute the workload.

Example:

Instead of one server handling all requests, multiple servers handle incoming traffic through a load balancer.

Advantages:
- Supports increased traffic.
- Easier to add capacity incrementally.
- Can improve availability through redundancy.

Disadvantages:
- Requires load balancing.
- Introduces distributed-system complexity.
- Requires handling shared state and replication.

---

## 6. Fault Tolerance

### Definition
Fault tolerance is the ability of a system to continue operating even when some components fail.

### Example

An application has two backend servers.

Server 1 crashes.

Server 2 continues processing incoming requests.

The application remains functional.

This demonstrates fault tolerance.

### Techniques to achieve fault tolerance
- Redundancy.
- Replication.
- Automatic failover.
- Health checks.
- Retry mechanisms.
- Circuit breakers.
- Graceful degradation.
- Fault isolation.

### Availability vs Fault Tolerance

| Availability | Fault Tolerance |
|---|---|
| Measures how accessible the system is. | Describes the ability to withstand component failures. |
| Focuses on uptime. | Focuses on handling failures. |
| Measured using uptime percentage. | Achieved through fault-handling mechanisms. |

**Important:** Fault tolerance is one of the ways to achieve high availability, but the two concepts are not identical.

A system can have high measured availability while still having poor fault tolerance.

---

## 7. Durability

### Definition
Durability means that once data is successfully committed, it should survive failures.

### Example

A user places an order.

The application confirms that the order was successfully placed.

Then the application server crashes.

When the user opens the application again, the order should still exist.

That is durability.

### How databases achieve durability
- Write-Ahead Logging (WAL).
- Persistent storage.
- Database replication.
- Backups.
- Crash recovery mechanisms.

PostgreSQL uses WAL to support crash recovery.

### Consistency vs Durability

| Consistency | Durability |
|---|---|
| Data respects defined correctness guarantees. | Committed data survives failures. |
| Focuses on correctness. | Focuses on persistence. |

---

## 8. Throughput

### Definition
Throughput is the number of requests or operations a system can process within a given amount of time.

It is commonly measured in:
- Requests Per Second (RPS).
- Transactions Per Second (TPS).

Example:

| Server | Throughput |
|---|---|
| Server A | 100 requests/second |
| Server B | 1,000 requests/second |
| Server C | 10,000 requests/second |

Higher throughput means the system can process more work per unit of time.

### Latency vs Throughput

| Latency | Throughput |
|---|---|
| Time taken for a request. | Number of requests processed per unit time. |
| Measured in milliseconds or seconds. | Measured in RPS or TPS. |
| Focuses on response time. | Focuses on processing capacity. |

Example:

A restaurant:
- Latency: How long one customer waits for food.
- Throughput: How many customers are served per hour.

High throughput does not automatically mean low latency.

---

## 9. CAP Theorem

### Definition

The CAP theorem states that during a network partition, a distributed system cannot guarantee both strong consistency and availability for every request.

CAP stands for:

- C – Consistency
- A – Availability
- P – Partition Tolerance

### Understanding network partition

Imagine two database nodes:

- Node A in Mumbai.
- Node B in London.

Both store customer account information.

Suddenly, the network connection between them breaks.

Both nodes are operational, but they cannot communicate.

This is a network partition.

### CP: Consistency and Partition Tolerance

The system prioritizes consistency.

If a node cannot verify whether its data is up to date, it may reject or delay certain requests.

- Strong consistency is prioritized.
- Some requests may fail.
- Availability is sacrificed during the partition.

### AP: Availability and Partition Tolerance

The system prioritizes availability.

Both nodes continue accepting requests despite being unable to communicate.

- Requests continue to be served.
- Data may temporarily become inconsistent.
- Conflicts may need to be resolved later.

### CAP summary

| Property | Meaning |
|---|---|
| Consistency | Reads respect the latest successful write under the chosen consistency model. |
| Availability | Every request to a non-failing node receives a response. |
| Partition Tolerance | The system continues operating despite network communication failures. |

**Important interview point:**

CAP does not mean that a system permanently chooses any two properties.

The trade-off between consistency and availability becomes relevant when a network partition occurs.

---

## 10. Other Important System Design Properties

| Concept | Definition |
|---|---|
| Maintainability | How easily a system can be modified and maintained. |
| Resilience | Ability to recover and continue operating after disruptions. |
| Elasticity | Ability to automatically adjust resources based on workload. |
| Security | Protecting systems and data from unauthorized access. |
| Isolation | Preventing concurrent operations from interfering incorrectly. |
| Redundancy | Having additional components to avoid single points of failure. |

---

# 11. Real-world Example: Messaging Application

Imagine designing a WhatsApp-like messaging application.

| Requirement | Why is it needed? |
|---|---|
| Low Latency | Messages should be delivered quickly. |
| High Availability | Users should be able to access the service anytime. |
| Scalability | Application should support millions of users. |
| Consistency | Message ordering and conversation state should follow defined guarantees. |
| Durability | Messages should not disappear after server failures. |
| Fault Tolerance | One server failure should not bring down the messaging service. |
| High Throughput | System should process millions of messages. |

# 12. Key Takeaways

- Latency = How quickly does the system respond?
- Availability = Is the system accessible?
- Reliability = Does the system work correctly?
- Consistency = Does data respect the defined consistency guarantees?
- Scalability = Can the system handle more workload?
- Fault Tolerance = Can it continue operating when components fail?
- Durability = Will committed data survive failures?
- Throughput = How much work can the system process?
- CAP = What consistency and availability trade-offs arise during network partitions?

**System design is about understanding these properties and making appropriate trade-offs based on business requirements.**