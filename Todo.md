### TODO
- Add tests for ValidateUserScope()
- Decide where you really have a need to pass scope array (in model layer). In most cases, we have a single scope operation
- When the client asks for a resource, should we scope it to a group/user or send all possible info (that user has access to?)
- In the model layer, you have to check if the user & scope combo are valid. 
- LOT OF REFACTORING LEFT ON impl.group

- Delete group & Udpate group info methods
- Scope list is passed automatically from Handler layer, but we have confusion on Group Scope vs user scope
- Scope violation error messages have to be more clear for user to be able to show right error on client
- How does the client know if the scope for an item (like a source) is read only or read-write (scope markers don't go to the client)
- [---]
- Should we verify userID or/and scopeID  (in header) for every handler call? Or should the  JWT enough?
- Updated category tests, remove hardcoded table names
- Add test cases to check missing scope (category, source)
### Multi user approach
- Scope table maps users to various scopes
- Group table optionally decorates scope table
- Txns can be under one scope
- Users getting added/removed from group, can result in a new scope_id for the group (no history mode)
- Sources can be a part of a group or a user
- Txn with source_id not visible to the group will have to be gracefully resolved
- Categories can be a part of a group or a user
- Tags --> A very loose concept, to be decided if to keep or not. 
- GetAllTxns should return txns that a user has access to (across scopes)
- txns can be modified, to move to a different scope (by owner)
- sources can eb modified, to move to another scope (by owner)
- Ephemeral scopes to share a slice of transactions?
### Logging & Monitoring
- Implement structured logging with different log levels (info, error, debug, etc.).
- Integrate with monitoring tools for database metrics.
- Use monitoring tools for database performance metrics and set up alerts for unusual activities.

### Security and safety
- Review database security settings, ensure limited open ports, and consider encryption for sensitive data.
- Ensure that database configurations are optimized for each environment (development, staging, production).
- Plan/scripts for regular backups of the database and test the restoration process.

### API & Middleware
- Consider rate limiting on API endpoints.
- Implement middleware for tasks like logging, CORS handling.

#### DB Performance & Scalability
- Implement caching mechanisms, like Redis, for frequently accessed data.
- Partitioning large tables like `transactions` for more efficient querying.
- Evaluate database connection pooling libraries for better performance.
- Sharding the database by `user_id` for scalability.
- Database migration tool for schema evolution.

#### Operations & Maintenance
- Implement a health check endpoint for periodic service health checks, including the database.
- Archiving old or seldom-used data for a leaner database.
- Balancing use of transactions to ensure data integrity without harming performance.

#### Error Handling & Feedback
- Structured and informative error feedback.
- Implement comprehensive error handling in DB operations.
- Define and monitor SLAs for database performance and uptime.

---
### High priority items

### Must Have:
1. **Validation Middleware:** Validate request payloads to ensure data integrity.
4. **Config Management:** Fetch sensitive information from environment variables or a secure configuration system.
5. **Database Migrations Tool:** Manage database schema changes and migrations.
6. **Health Check Endpoint:** Ensure your service's health can be monitored.
7. **API Versioning:** Protect users from breaking changes by versioning the API.
8. **Backup & Recovery:** Implement a system for regular data backups and a clear recovery plan.
9. **Graceful Shutdown:** Handle service interruptions by finishing tasks and releasing resources properly.
10. **Logging:** Log important events, errors, and other relevant information.

### Should Have:
1. **Rate Limiting & Throttling Middleware:** Protect the API from abuse and manage user requests effectively.
2. **Caching:** Improve performance by caching frequently accessed data.
3. **Metrics & Monitoring Middleware:** Collect metrics for analysis and visualization.
4. **Request ID Middleware:** Aid debugging and tracking with unique IDs for each request.
5. **Content Negotiation:** Support multiple response formats based on client requirements.
6. **Compression Middleware:** Reduce API response payload sizes for faster transmission.
7. **Secure Sensitive Routes:** Enhance protection for routes that modify data or configurations.

### Good to Have:
1. **Documentation:** Use tools like Swagger to auto-generate API documentation.
2. **CORS Middleware:** Handle cross-origin requests if your backend serves multiple frontends.
3. **Performance Profiling Tools:** Integrate tools that can help profile and optimize your API's performance.


