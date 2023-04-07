package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianolmedo/go-restapi/api"

	"github.com/gofiber/fiber/v2"
)

// signUpUser handler POST: /users
func signUpUser(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := UserSignUpForm{}
		err := c.BodyParser(&form)
		if err != nil {
			resp := api.RespJSON(api.MsgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.SignUp(&User{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if err != nil {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := api.RespJSON(api.MsgOK, "user created", UserProfileDTO{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// findUser handler GET: /users/:id
func findUser(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		user, err := s.Find(int64(id))
		if errors.Is(err, ErrUserNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp)
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := api.RespJSON(api.MsgOK, "", UserProfileDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// updateUser handler PUT: /users/:id
func updateUser(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form := UserUpdateForm{}
		err = c.BodyParser(&form)
		if err != nil {
			resp := api.RespJSON(api.MsgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		form.ID = int64(id)

		err = s.Update(User{
			ID:        form.ID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
			Password:  form.Password,
		})

		if errors.Is(err, ErrUserNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := api.RespJSON(api.MsgOK, "user updated", User{
			ID:        form.ID,
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// listUsers handler GET: /users
func listUsers(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := s.List()
		if err != nil {
			resp := api.RespJSON(api.MsgError, "", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if len(users) == 0 {
			resp := api.RespJSON(api.MsgOK, "there are not users", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		assemble := func(u *User) UserProfileDTO {
			return UserProfileDTO{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
			}
		}

		list := make(UserList, 0, len(users))
		for _, v := range users {
			list = append(list, assemble(v))
		}

		resp := api.RespJSON(api.MsgOK, "", list)
		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// deleteUser handler DELETE: /users/:id
func deleteUser(s Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := api.RespJSON(api.MsgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.Remove(int64(id))
		if errors.Is(err, ErrUserNotFound) {
			resp := api.RespJSON(api.MsgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := api.RespJSON(api.MsgError, fmt.Sprintf("could not delete user: %s", err), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		// TO-DO: Averiguar como poner el logger de Fiber.

		resp := api.RespJSON(api.MsgOK, "user deleted", nil)
		return c.Status(http.StatusOK).JSON(resp)

	}
}
