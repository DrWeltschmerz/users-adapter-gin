
# users-adapter-gin

Gin HTTP adapter for [users-core](https://github.com/DrWeltschmerz/users-core).

## Install

```sh
go get github.com/DrWeltschmerz/users-adapter-gin@v1.2.0
```

## Features
- RESTful endpoints for user and role management
- JWT authentication middleware
- Swagger/OpenAPI documentation (see `/swagger/index.html`)
- Easily extensible for custom business logic

## Usage

This module is intended to be used with [users-core](https://github.com/DrWeltschmerz/users-core) and a JWT implementation (e.g., [jwt-auth](https://github.com/DrWeltschmerz/jwt-auth)).

See the users-core README for wiring instructions. This repo provides only the HTTP adapter and API documentation.

Quickstart wiring:

```go
import (
	users "github.com/DrWeltschmerz/users-core"
	gormadapter "github.com/DrWeltschmerz/users-adapter-gorm/gorm"
	ginadapter "github.com/DrWeltschmerz/users-adapter-gin/ginadapter"
	"github.com/DrWeltschmerz/jwt-auth/pkg/authjwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
db.AutoMigrate(&gormadapter.GormUser{}, &gormadapter.GormRole{})
svc := users.NewService(
	gormadapter.NewGormUserRepository(db),
	gormadapter.NewGormRoleRepository(db),
	authjwt.NewBcryptHasher(),
	authjwt.NewJWTTokenizer(),
)
r := gin.Default()
ginadapter.RegisterRoutes(r, svc, authjwt.NewJWTTokenizer())
r.Run(":8080")
```

## Endpoints

- `POST   /register` — Register a new user
- `POST   /login` — Login and receive JWT
- `GET    /user/profile` — Get current user profile (JWT required)
- `PUT    /user/profile` — Update current user profile (JWT required)
- `POST   /user/change-password` — Change password (JWT required)
- `GET    /users` — List all users
- `DELETE /users/:id` — Delete user by ID
- `GET    /roles` — List all roles
- `POST   /users/:id/assign-role` — Assign role to user
- `POST   /users/:id/reset-password` — Reset user password

## Extending

- Add or override handlers as needed for your project.
- See `ginadapter/handler.go` for all endpoints and extension points.

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0). See the LICENSE file for details.

See also:
- [users-adapter-gorm](https://github.com/DrWeltschmerz/users-adapter-gorm)
- [users-core](https://github.com/DrWeltschmerz/users-core)
- [jwt-auth](https://github.com/DrWeltschmerz/jwt-auth)
- Tests: [users-tests](https://github.com/DrWeltschmerz/users-tests)


