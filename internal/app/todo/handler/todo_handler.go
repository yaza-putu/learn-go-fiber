package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yaza-putu/crud-fiber/internal/app/todo/service"
	"github.com/yaza-putu/crud-fiber/internal/database"
	"github.com/yaza-putu/crud-fiber/internal/http/request"
	"github.com/yaza-putu/crud-fiber/internal/http/response"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type TodoRequest struct {
	ID   int    `json:"id"`
	Todo string `json:"todo" form:"todo" validate:"required,unique=todos:name:ID"`
	Done bool   `json:"done" form:"done" validate:"boolean"`
	//File multipart.File `validate:"required,filetype=image/png image/jpeg image/jpg"`
}

func Create(ctx *fiber.Ctx) error {
	todo := TodoRequest{}

	//f, err := ctx.FormFile("file")
	//if err == nil {
	//	todo.File, err = f.Open()
	//	if err != nil {
	//		return err
	//	}
	//}

	err := ctx.BodyParser(&todo)
	if err != nil {
		return ctx.Status(400).JSON(response.Api(
			response.SetCode(400),
			response.SetMessage("bad request"),
			response.SetError(err.Error()),
		))
	}

	//f, err := ctx.FormFile("file")
	//if err != nil {
	//	return err
	//}
	//
	//allowedTypes := []string{"image/jpeg", "image/png"}
	//file1, _ := f.Open()
	//if !isValidMIMEType(file1, allowedTypes) {
	//	return ctx.Status(422).JSON(response.Api(response.SetCode(422)))
	//}
	val, err := request.Validation(&todo)

	if err != nil {
		return ctx.Status(val.Code).JSON(val)
	}

	//form, err := ctx.FormFile("file")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//path, err := file.ToPublic(form, "", true)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(path)
	//}
	_, err = database.Conn().Exec("INSERT INTO todos (name,done) VALUES (?,?)", todo.Todo, todo.Done)
	defer database.Conn().Close()

	if err != nil {
		return ctx.Status(200).JSON(response.Api(
			response.SetCode(500),
			response.SetMessage("internal server error"),
			response.SetError(err),
		))
	}
	return ctx.JSON(response.Api(response.SetCode(200), response.SetMessage("Success")))
}

func Update(ctx *fiber.Ctx) error {
	param := ctx.Params("id")
	todo := TodoRequest{}

	idInt, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}
	todo.ID = idInt

	err = ctx.BodyParser(&todo)
	if err != nil {
		return ctx.Status(400).JSON(response.Api(response.SetCode(400), response.SetMessage("Bad request")))
	}

	query, err := database.Conn().Exec("UPDATE todos SET name = ? , done = ? WHERE id = ? ", todo.Todo, todo.Done, param)
	defer database.Conn().Close()

	if err != nil {
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}

	id, err := query.LastInsertId()
	return ctx.JSON(response.Api(response.SetMessage(fmt.Sprintf("success update with last id %d", id))))
}

func FindById(ctx *fiber.Ctx) error {
	param := ctx.Params("id")

	todo := TodoRequest{}
	if err := database.Conn().QueryRow("SELECT * FROM todos where id = ?", param).Scan(&todo.ID, &todo.Todo, &todo.Done); err != nil {
		return ctx.Status(404).JSON(response.Api(
			response.SetCode(404),
			response.SetError(err.Error()),
		))
	}
	defer database.Conn().Close()

	return ctx.JSON(response.Api(
		response.SetCode(200),
		response.SetMessage("ok"),
		response.SetData(todo),
	))

}

func All(ctx *fiber.Ctx) error {
	q := ctx.Query("q")
	take := ctx.Query("take")
	page := ctx.Query("page")

	q = "%" + q + "%"
	t, err := strconv.Atoi(take)
	if err != nil {
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}

	offset := (p - 1) * t

	var totalRows int64
	// total rows
	err = database.Conn().QueryRow("SELECT count(*) from todos WHERE name LIKE ? OR done LIKE ?", q, q).Scan(&totalRows)
	defer database.Conn().Close()
	if err != nil {
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}

	rows, err := database.Conn().Query("SELECT * FROM todos WHERE name LIKE ? OR done LIKE ? LIMIT ? OFFSET ?", q, q, t, offset)
	defer database.Conn().Close()
	if err != nil {
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}

	todos := []TodoRequest{}
	for rows.Next() {
		todo := TodoRequest{}
		err = rows.Scan(&todo.ID, &todo.Todo, &todo.Done)
		if err != nil {
			return ctx.Status(500).JSON(response.Api(
				response.SetCode(500),
				response.SetError(err.Error()),
			))
		}

		todos = append(todos, todo)
	}

	return ctx.Status(200).JSON(response.Api(
		response.SetCode(200),
		response.SetData(struct {
			Page      int           `json:"page"`
			Take      int           `json:"take"`
			TotalRows int           `json:"total_rows"`
			TotalPage int           `json:"total_page"`
			Rows      []TodoRequest `json:"rows"`
		}{
			Page:      p,
			Take:      t,
			TotalPage: int(math.Ceil(float64(totalRows) / float64(t))),
			TotalRows: int(totalRows),
			Rows:      todos,
		}),
	))
}

func Delete(ctx *fiber.Ctx) error {
	param := ctx.Params("id")

	_, err := database.Conn().Exec("DELETE from todos WHERE id = ?", param)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(404).JSON(response.Api(
				response.SetCode(404),
				response.SetError(err.Error()),
			))
		}
		return ctx.Status(500).JSON(response.Api(
			response.SetCode(500),
			response.SetError(err.Error()),
		))
	}

	return ctx.Status(200).JSON(response.Api(
		response.SetCode(200),
		response.SetMessage("deleted todo successfully"),
	))
}

func isValidMIMEType(file multipart.File, validMIMETypes []string) bool {
	// Mendapatkan tipe MIME dari berkas
	buffer := make([]byte, 512) // 512 bytes cukup untuk menentukan tipe MIME
	_, err := file.Read(buffer)
	if err != nil {
		return false
	}
	fileType := http.DetectContentType(buffer)
	// Memeriksa apakah tipe MIME berkas valid
	for _, validType := range validMIMETypes {
		if strings.HasPrefix(fileType, validType) {
			return true
		}
	}
	return false
}

func GenReport(ctx *fiber.Ctx) error {
	err := service.NewExcel()

	if err != nil {
		return ctx.JSON(err)
	}

	err = ctx.SendFile("public/test.xlsx")
	if err == nil {
		err = os.Remove("public/test.xlsx")
	}

	return err
}
