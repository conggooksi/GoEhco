package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/users", GetUsers)

	e.GET("/myid/:id", getId)

	e.GET("/show", show)

	e.POST("/save", save)

	e.POST("/save_file", fileSave)

	e.POST("/connect_mysql", connectMySQL)

	e.POST("/automigration", connectMySQLbyORM)

	e.POST("/create_user", createUser)

	e.GET("/select_user", selectUser)

	e.PUT("/update_user", updateUser)

	e.DELETE("delete_user", deleteUser)

	e.POST("/todo", createTodo)

	e.GET("/todo", selectTodo)

	e.PUT("/todo", updateTodo)

	e.DELETE("/todo", deleteTodo)

	e.Logger.Fatal(e.Start(":1323"))
}

func GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "user 정보")
}

func getId(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

type User struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}

func save(c echo.Context) error {
	// Get name and email

	u := new(User)

	if err := c.Bind(u); err != nil {
		return err
	}

	return c.String(http.StatusOK, "name:"+u.Name+", email:"+u.Email)
}

func fileSave(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
}

func connectMySQL(c echo.Context) error {

	db, err := sql.Open("mysql", "root:1234@tcp(192.168.35.152:3306)/smart_media?charset=utf8mb4")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("select * from User")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var Id int
	var Name string
	var Email string

	for rows.Next() {
		err := rows.Scan(&Id, &Name, &Email)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Id: %d User: %s, Email: %s \n", Id, Name, Email)
	}

	return c.HTML(http.StatusOK, "<b>Thank you!</b>")
}

func connectMySQLbyORM(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	return c.HTML(http.StatusOK, "<b>Created!</b>")
}

func createUser(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	user := User{
		Name:  "kuki",
		Email: "kuki@test.com",
	}

	result := db.Create(&user)

	fmt.Println(result.RowsAffected)

	return c.HTML(http.StatusOK, "<b>Created!</b>")

}

func selectUser(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	user := map[string]interface{}{}

	result := db.Model(&User{}).First(&user)

	fmt.Println(result.RowsAffected)

	return c.HTML(http.StatusOK, "<b>user:"+fmt.Sprintf("%v", user)+"</b>")
}

func updateUser(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	user := new(User)

	db.Where("name=?", "kuki").First(&user)

	user.Email = "update@test.com"

	db.Where("name=?", "kuki").Save(&user)

	return c.HTML(http.StatusOK, "<b>updated</b>")

}

func deleteUser(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	db.Where("name=?", "kuki").Delete(&User{})

	return c.HTML(http.StatusOK, "<b>deleted</b>")

}

// --------------------------------------------------------
type Todos struct {
	gorm.Model
	Userid     string
	Start_date time.Time
	End_date   time.Time
	Title      string
	Status     string
}

func createTodo(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todos{})

	s := "2022-07-15 14:37:00"
	startTime, _ := time.Parse("2006-01-02 15:04:05", s)
	endTime, _ := time.Parse("2006-01-02 15:04:05", s)

	user := Todos{
		Userid:     "kuki",
		Start_date: startTime,
		End_date:   endTime,
		Title:      "Hello",
		Status:     "eating",
	}

	result := db.Create(&user)

	fmt.Println(result.RowsAffected)

	return c.HTML(http.StatusOK, "<b>Created!</b>")

}

func selectTodo(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todos{})

	todo := new(Todos)

	result := db.First(&todo)

	fmt.Println(result.RowsAffected)

	return c.HTML(http.StatusOK, "<b>todo:"+fmt.Sprintf("%v", todo)+"</b>")
}

func updateTodo(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todos{})

	todo := new(Todos)

	db.Where("name=?", "kuki").First(&todo)

	todo.Status = "Sleeping"

	db.Where("name=?", "kuki").Save(&todo)

	return c.HTML(http.StatusOK, "<b>updated</b>")

}

func deleteTodo(c echo.Context) error {
	dsn := "root:1234@tcp(localhost:3306)/smart_media?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todos{})

	db.Where("Userid=?", "kuki").Delete(&Todos{})

	return c.HTML(http.StatusOK, "<b>deleted</b>")

}
