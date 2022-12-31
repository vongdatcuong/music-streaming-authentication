# music-streaming-authentication

Music Streaming Authentication Service

- Configure 2 containers: service & db with docker-compose
- Write Database migration scripts
- Create makefile with common commands
- Define protobuf files + Generate pb files
- Implement CRUD methods for User/Permission services
- Use grpc-gateway to run Rest server (http endpoints)
- Use jwt, grpc interceptor, http middleware to implement authentication and authorization
- Provide endpoint for other service to do authentication and authorization
- Migrate Authentication operation to API Gateway. One issue was to defined error response for all Grpc endpoints
- Inject jwt token to grpc interceptor context so that grpc methods have access to jwt token
