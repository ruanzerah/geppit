package server

import (
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/google/uuid"
	"github.com/ruanzerah/geppit/cmd/web"
	"github.com/ruanzerah/geppit/internal/repository"
	"github.com/ruanzerah/geppit/internal/utils"

	"github.com/a-h/templ"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))
	staticFiles, _ := fs.Sub(web.Files, "assets")
	r.StaticFS("/assets", http.FS(staticFiles))

	r.GET("/", s.Home)
	r.GET("/web", func(c *gin.Context) {
		templ.Handler(web.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/health", s.healthHandler)
	r.GET("/login", s.Login)
	r.DELETE("/logout", s.Logout)
	r.POST("/signin", s.Signin)

	api := r.Group("/api")
	{
		user := api.Group("/user", s.AuthMiddleware)
		{
			user.POST("/", s.CreateUser)
			user.GET("/:id", s.ListUser)
			user.DELETE("/:id", s.DeleteUser)
			user.PUT("/rename/:id", s.RenameUser)
			user.PUT("/credential/:id", s.ChangeHash)
		}

		transaction := api.Group("/transaction", s.AuthMiddleware)
		{
			transaction.POST("/", s.CreateTransaction)
			transaction.GET("/:id", s.ListTransaction)
		}
	}
	return r
}

func (s *Server) Home(c *gin.Context) {
}

func (s *Server) Login(c *gin.Context) {
}

func (s *Server) Logout(c *gin.Context) {
}

func (s *Server) Signin(c *gin.Context) {
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) CreateUser(c *gin.Context) {
	var user repository.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Username == "" || user.Email == "" {
		http.Error(c.Writer, "Username or email can't be empty", http.StatusBadRequest)
		return
	}
	hash, err := utils.Hash(user.Hash)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	qUser, err := s.queries.CreateUser(c, repository.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Hash:     hash,
	})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, qUser)
}

func (s *Server) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.queries.DeleteUser(c, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete Successful"})
}

func (s *Server) ListUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := s.queries.GetUserByID(c, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) RenameUser(c *gin.Context) {
	var user repository.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(c.Writer, "Username is required", http.StatusBadRequest)
		return
	}

	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	renamedUser, err := s.queries.RenameUser(c, repository.RenameUserParams{
		ID:       userID,
		Username: user.Username,
	})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, renamedUser)
}

func (s *Server) ChangeHash(c *gin.Context) {
	// TODO: Implement this later
}

func (s *Server) CreateTransaction(c *gin.Context) {
}

func (s *Server) ListTransaction(c *gin.Context) {
}

func (s *Server) AuthMiddleware(c *gin.Context) {
}
