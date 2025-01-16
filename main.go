package main

import (
	"bytes"
	"flag"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog/log"
)

func forwardHandler(baseUrl string) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		method := req.Method
		url := baseUrl + req.RequestURI

		log.Printf("Method: %s\n", method)
		log.Printf("Forwarding to %s\n", url)

		var body io.Reader
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				log.Printf("Error reading request body: %v\n", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read body")
			}
			body = bytes.NewReader(bodyBytes)
			req.Body = io.NopCloser(body)
		}

		forwardReq, err := http.NewRequest(method, url, body)
		if err != nil {
			log.Printf("Error creating forward request: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create forward request")
		}

		for key, values := range req.Header {
			for _, value := range values {
				forwardReq.Header.Add(key, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(forwardReq)
		if err != nil {
			log.Printf("Error forwarding request: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to forward request")
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read response body")
		}

		return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
	}
}

func main() {
	port := flag.String("port", "1323", "Port to listen on")
	baseUrl := flag.String("baseurl", "http://localhost:8000", "Base URL for the application forwarding")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Any("/*", forwardHandler(*baseUrl))
	e.Logger.Fatal(e.Start(":" + *port))
}
