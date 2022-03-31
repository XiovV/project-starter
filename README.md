# Project starter template for Go backend projects

## Project structure
### `jwt` package:
Package with functions for creating a new token with custom claims, validating a token, and getting claims from a token.

### `repository` package:
- `errors.go` contains predefined errors, such as the NotFoundErr and AlreadyExistsErr. They are structs which implement the Error interface.
- `postgres.go` contains code for establishing a new database connection pool.
- `user.go` and `post.go` contain basic CRUD operations for users and posts, respectively.

### `server` package:
- `helpers.go`: basic helper methods, such as getUserFromContext and validatePageAndLimit.
- `middleware.go`: basic authentication, CORS and error handling middleware (explained below).
- `responses.go`: frequently used responses, such as successResponse (200 OK + message), badRequestResponse (400 Bad Request + message) etc...
- `server.go`: code for instantiating a new server. All endpoints are defined there, along with the middleware which the server will be using.
- `user.go`: user related handlers, such as handlers for registering an account and logging in.

### `validator` package
Very basic package for validating user input. At the moment, it only has a method that checks if a string has been provided, but can be easily extended to do more complex validation.

## Endpoint grouping convention
There are two types of groups, "public" and "auth". In case a collection has both public and protected routes (such as `users`, the register and login endpoints are public,
while the endpoint for fetching personal posts is protected), we create two types of groups. In this example, we have created `usersPublic` and `usersAuth` groups:

```Go
v1 := router.Group("/v1")

usersPublic := v1.Group("/users")
{
	usersPublic.POST("/register", s.registerUserHandler)
	usersPublic.POST("/login", s.loginUserHandler)
}

usersAuth := v1.Group("/users")
usersAuth.Use(s.userAuth)
{
	usersAuth.GET("/posts", s.getPostsHandler)
}
```

## Error handling
Errors which occur frequently are handled in the errorHandler middleware. These errors are in most cases database errors. 
This is done in order to reduce boilerplate error handling in handlers. Example:\
Let' say we are calling the Create() method in the registerUserHandler:


```Go
// server/users.go

createdUser, err := s.UserRepository.Create(user)
if err != nil {
    c.Error(fmt.Errorf("registerUserHandler: %w", err))
    return
}
```

We can call c.Error() and pass the error (in this case we wrapped it for tracing purposes), the errorHandler middleware will then proceed to handle the db error:

```Go
func (s *Server) errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()

		if err != nil {
			var alreadyExistsErr *repository.AlreadyExistsErr
			var notFoundErr *repository.NotFoundErr

			switch {
			case errors.As(err, &alreadyExistsErr):
				c.JSON(http.StatusConflict, gin.H{"error": alreadyExistsErr.Error()})
			case errors.As(err, &notFoundErr):
				c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
			default:
				s.internalServerErrorResponse(c, err)
			}
		}
	}
}
```