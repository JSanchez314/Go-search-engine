package routes

import (
	"fmt"
	"log"
	"os"
	"pro/jet/db"
	"pro/jet/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

type login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := login{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Error: Unauthorized</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong logging in</h2>")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

func LogOutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	fmt.Println("Cookie 'admin' received:", cookie)

	if cookie == "" {
		fmt.Println("The 'admin' cookie is empty or not found.")
		return c.Redirect("/login", 302)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not configured")
	}

	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println("Error parsing JWT token:", err)
		return c.Redirect("/login", 302)
	}

	claims, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		fmt.Println("Valid token. Claims:", claims)
		return c.Next()
	}

	fmt.Println("Invalid token or an issue occurred.")
	return c.Redirect("/login", 302)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot get settings</h2>")
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return c.Render("index", fiber.Map{
		"amount":   amount,
		"searchOn": settings.SearchOn,
		"addNew":   settings.AddNew,
	})
}

type settingsForm struct {
	Amount   int    `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

func DashboarPostdHandler(c *fiber.Ctx) error {
	input := settingsForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot parse settings</h2>")
	}

	addNew := false
	if input.AddNew == "on" {
		addNew = true
	}
	searchOn := false
	if input.SearchOn == "on" {
		searchOn = true
	}

	settings := &db.SearchSetting{}
	settings.Amount = uint(input.Amount)
	settings.SearchOn = searchOn
	settings.AddNew = addNew

	err := settings.Update()
	if err != nil {
		fmt.Println("Error updating settings:", err)
		return c.SendString("<h2>Error: Cannot update settings</h2>")
	}

	c.Append("HX-Refresh", "true")
	return c.SendStatus(200)
}
