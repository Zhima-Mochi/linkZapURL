package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Zhima-Mochi/linkZapURL/config"
	"github.com/Zhima-Mochi/linkZapURL/handlers/redirection"
	"github.com/Zhima-Mochi/linkZapURL/handlers/shortening"
	"github.com/Zhima-Mochi/linkZapURL/pkg/cache/redis"
	"github.com/Zhima-Mochi/linkZapURL/pkg/database/mongodb"
	"github.com/gin-gonic/gin"

	doc "github.com/Zhima-Mochi/linkZapURL/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	machineID uint8

	endpoint string
)

type Handler struct {
	shortening  shortening.Shortening
	redirection redirection.Redirection
}

func setEnv() {
	machineIDStr := os.Getenv("MACHINE_ID")
	if machineIDStr == "" {
		log.Fatal("MACHINE_ID is not set")
		panic("MACHINE_ID is not set")
	}

	machineIDInt64, err := strconv.ParseUint(machineIDStr, 10, 8)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if machineIDInt64 > 0xFF {
		log.Fatal("MACHINE_ID is too large")
		panic("MACHINE_ID is too large")
	}

	machineID = uint8(machineIDInt64)
	log.Printf("Machine ID: %d", machineID)

	endpoint = os.Getenv("ENDPOINT")
	if endpoint == "" {
		log.Fatal("ENDPOINT is not set")
		panic("ENDPOINT is not set")
	}
	log.Printf("Endpoint: %s", endpoint)

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	setEnv()

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	database, err := mongodb.NewMongodb(config.Mongodb)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	cache, err := redis.NewRedis(config.Redis)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	shortening := shortening.NewShortening(machineID, database, cache)

	redirection := redirection.NewRedirection(cache, database)

	handler := &Handler{
		shortening:  shortening,
		redirection: redirection,
	}

	router := gin.Default()
	router.GET("/:code", handler.Redirect)

	apiV1 := router.Group("/api/v1")
	apiV1.POST("/urls", handler.Shorten)

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// swagger doc
	doc.SwaggerInfo.Title = "linkZapURL"
	doc.SwaggerInfo.Description = "A URL shortener service."
	doc.SwaggerInfo.Version = "1.0"

	router.Run(":8080")
}

type ShortenRequest struct {
	URL           string `json:"url" binding:"required" example:"https://example.com"`
	ExpireAt      string `json:"expireAt" binding:"required" example:"2021-02-08T09:20:41Z"`
	ExpireAtInt64 int64  `json:"-" binding:"-"`
}

func (r *ShortenRequest) Bind(g *gin.Context) error {
	if err := g.ShouldBindJSON(r); err != nil {
		return err
	}

	expireAtTime, err := time.Parse(time.RFC3339, r.ExpireAt)
	if err != nil {
		return err
	}

	r.ExpireAtInt64 = expireAtTime.Unix()

	return nil
}

type ShortenResponse struct {
	ID       string `json:"id" example:"5abcABC"`
	ShortURL string `json:"shortURL" example:"https://localhost/5abcABC"`
}

// @Summary Shorten a URL
// @Description Shortens a given URL and provides an expiration time.
// @Tags URL Shortening
// @Accept json
// @Produce json
// @Param body body ShortenRequest true "Shorten Request"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/v1/urls [post]
func (h *Handler) Shorten(g *gin.Context) {
	ctx := g.Request.Context()

	req := ShortenRequest{}
	if err := req.Bind(g); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.shortening.Shorten(ctx, req.URL, req.ExpireAtInt64)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ShortenResponse{
		ID:       url.Code,
		ShortURL: endpoint + "/" + url.Code,
	}

	g.JSON(http.StatusCreated, resp)
}

// @Summary Redirect to a URL
// @Description Redirects to the original URL.
// @Tags URL Redirection
// @Param code path string true "Shortened URL Code"
// @Success 301 string string "Moved Permanently"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /{code} [get]
func (h *Handler) Redirect(g *gin.Context) {
	ctx := g.Request.Context()

	code := g.Param("code")

	url, err := h.redirection.Redirect(ctx, code)
	if err == redirection.ErrNotFound || err == redirection.ErrExpired || err == redirection.ErrInvalidCode {
		g.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.Redirect(http.StatusMovedPermanently, url.URL)
}
