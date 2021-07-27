package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

func main() {
	e := echo.New()
	repo := NewInMemoryUserRepository()

	e.GET("/api/users", getAllUsers(repo))
	e.GET("/api/users/:id", getUser(repo))
	e.POST("/api/users", createNewUsers(repo))
	e.DELETE("/api/users/:id", deleteUser(repo))

	log.Fatal(e.Start(":4000"))
}

func getAllUsers(repo UserRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		users, err := repo.GetAllUsers()
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, users)
	}
}

func getUser(repo UserRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		user, err := repo.GetUserById(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusOK, user)
	}
}

func createNewUsers(repo UserRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request CreateUserRequest
		err := ctx.Bind(&request)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		user, err := repo.CreateUser(request.Email, request.FirstName, request.LastName, request.Phone)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusOK, user)
	}
}

func deleteUser(repo UserRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = repo.DeleteUser(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return nil
	}
}
