package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {

	e := echo.New()

	viper.SetConfigFile(`config.json`)

	err := viper.ReadInConfig()

	if err != nil {
		e.Logger.Panic(err.Error())
	}

	e.Logger.Print(viper.Get("supersecret"))

	dbtype := viper.GetString("mysql.type")
	dbhost := viper.GetString("mysql.host")
	dbport := viper.GetString("mysql.port")
	dbusername := viper.GetString("mysql.username")
	dbpassword := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.dbname")

	sqlConnection, err := sql.Open(dbtype, dbusername+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname)

	if err != nil {
		e.Logger.Panic(err.Error())
	}

	defer func() {
		sqlConnection.Close()
		e.Close()
	}()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{""},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	e.GET("/meow", func(c echo.Context) error {
		return c.String(200, "Meowwwwwwwwwwwww!!!!!!!!!!!!!!!!!!!!")
	})

	e.GET("/test", func(c echo.Context) error {

		row, err := sqlConnection.QueryContext(context.TODO(), "SELECT * FROM `meow__lmw`")

		if err != nil {
			e.Logger.Panic(err.Error())
		}

		var getMeow []map[string]interface{}

		allColumn, _ := row.Columns()

		log.Print(allColumn)

		getType, _ := row.ColumnTypes()

		for row.Next() {
			newMeow := make([]interface{}, len(allColumn))

			newMeowPointer := make([]interface{}, len(allColumn))

			for a := range newMeow {
				newMeowPointer[a] = &newMeow[a]
			}

			err = row.Scan(newMeowPointer...)

			if err != nil {
				e.Logger.Panic(err.Error())
			}

			resultMeow := make(map[string]interface{}, 0)

			for a := range getType {
				columnType := getType[a].DatabaseTypeName()
				e.Logger.Print(columnType)
				if columnType == "VARCHAR" || columnType == "NVARCHAR" || columnType == "TEXT" {
					resultMeow[allColumn[a]] = string(newMeow[a].([]byte))
				} else if columnType == "BOOL" || columnType == "TINYINT" {
					resultMeow[allColumn[a]], _ = strconv.ParseBool(string(newMeow[a].([]byte)))
				} else {
					resultMeow[allColumn[a]], _ = strconv.Atoi(string(newMeow[a].([]byte)))
				}
			}

			getMeow = append(getMeow, resultMeow)
		}

		// resultByte, _ := json.Marshal(getMeow)
		// var result []interface{}
		// json.Unmarshal(resultByte, &result)
		e.Logger.Print(getMeow[0]["name"])
		return c.JSON(200, getMeow)
	})

	e.GET("/secret", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"secret": "meowret",
		})
	})

	e.GET("/sql/connection", func(c echo.Context) error {
		err = sqlConnection.Ping()

		if err == nil {
			return c.JSON(200, map[string]interface{}{
				"status": "online",
			})
		}
		return c.JSON(200, map[string]interface{}{
			"status": "offline",
		})
	})

	e.GET("/health-check", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"Status": "Online",
		})
	})

	e.Logger.Print(map[string]interface{}{
		"secret": "meowret",
	})

	e.Logger.Fatal(e.Start(":10800"))
}
