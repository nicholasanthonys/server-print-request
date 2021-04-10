package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	//* init logger with timestamp
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.Level = logrus.DebugLevel
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware

	e.Use(middleware.Recover())

	// Routes
	e.POST("/", printRequest)
	e.PUT("/", printRequest)

	e.POST("/error/bad-request", returnBadRequest)
	e.PUT("/error/bad-request", returnBadRequest)

	e.POST("/dummy-response", returnDummyResponse)

	// Routes
	// e.POST("/:anything", printRequest)
	// e.PUT("/:anything", printRequest)

	// Start server
	e.Logger.Fatal(e.Start(":3001"))
}

type dummyResponse struct {
	DummyObject map[string]interface{}
}

func returnDummyResponse(c echo.Context) error {
	logrus.Info("triggered")
	dummyResponse := &dummyResponse{
		DummyObject: (map[string]interface{}{
			"loop": 3,
		}),
	}
	return c.JSON(http.StatusOK, dummyResponse)
}

func returnBadRequest(c echo.Context) error {
	resMap := map[string]string{
		"message": "Bad Requestzz",
	}

	return c.JSON(http.StatusBadRequest, resMap)
}

func printRequest(c echo.Context) error {
	contentType := c.Request().Header["Content-Type"][0]

	// logrus.Info("queryparmas")
	// logrus.Info(c.QueryParams())

	// logrus.Info("header")
	// logrus.Info(c.Request().Header)
	logrus.Info("query is")
	for key, val := range c.QueryParams() {
		logrus.Info("key : ", key, " value ", val)
	}

	switch contentType {
	case "application/json":
		logrus.Info("content type application json")
		reqByte, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		logrus.Info(string(reqByte))
		return c.JSONBlob(http.StatusOK, reqByte)
	case "application/xml":
		logrus.Info("content type application xml")
		reqByte, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		logrus.Info(string(reqByte))
		return c.XMLBlob(http.StatusOK, reqByte)
	case "application/x-www-form-urlencoded":
		myMap := make(map[string]interface{})
		logrus.Info("Content type is x-www-form-urlencoded")
		c.Request().ParseForm()
		for key, value := range c.Request().Form { // range over map
			logrus.Info("key is ", key, " value is ", value, "length is ", len(value))

			if len(value) > 1 {
				logrus.Warn("KEY ", key, " LENGTH IS ", len(value))
				myMap[key] = value
			} else {
				myMap[key] = c.FormValue(key)
			}
		}

		logrus.Info("body is")
		logrus.Info(myMap)
		return c.JSON(http.StatusOK, myMap)

	}

	return nil
}
