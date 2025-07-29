package listener

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"news-feed/accessor"
	"news-feed/dto"
	"news-feed/manager"
	"strconv"
	"sync"
	"time"
)

type RestListener struct {
	Engine *gin.Engine
}

var (
	restOnce     sync.Once
	restInstance *RestListener
)

func NewRestListener() *RestListener {
	r := gin.Default()
	return &RestListener{Engine: r}
}

func GetRestListener() *RestListener {
	restOnce.Do(func() {
		restInstance = NewRestListener()
	})
	return restInstance
}

func (rl *RestListener) registerRoutes() {
	rl.Engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// POST /login
	rl.Engine.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		ma, _ := accessor.GetMySQLAccessor()
		id, err := ma.GetUserIDByUsernamePassword(req.Username, req.Password)
		if err != nil || id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})

	// POST /signup
	rl.Engine.POST("/signup", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		ma, _ := accessor.GetMySQLAccessor()
		exists, _ := ma.IsUserExists(req.Username)
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		id, err := ma.InsertUser(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})

	// POST /posts
	rl.Engine.POST("/posts", func(c *gin.Context) {
		var req struct {
			UserID  int64  `json:"user_id"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		ma, _ := accessor.GetMySQLAccessor()
		id, err := ma.InsertPost(req.Content, req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
			return
		}
		km, _ := manager.GetKafkaManager()
		err = km.Produce(context.Background(), dto.KafkaMessage{
			UserID:    req.UserID,
			PostID:    id,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			fmt.Printf("failed to produce message: %v\n", err)
		}
		c.JSON(http.StatusOK, gin.H{"post_id": id})
	})

	// POST /follow
	rl.Engine.POST("/follow", func(c *gin.Context) {
		var req struct {
			FollowerID int64 `json:"follower_id"`
			FolloweeID int64 `json:"followee_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		ma, _ := accessor.GetMySQLAccessor()
		exists, _ := ma.IsFollowExists(req.FollowerID, req.FolloweeID)
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "already following"})
			return
		}
		id, err := ma.InsertFollow(req.FollowerID, req.FolloweeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"follow_id": id})
	})

	// GET /posts?user_id=xxx
	rl.Engine.GET("/posts", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}
		ra := accessor.GetRedisAccessor()
		postIDs, _ := ra.GetUserFeed(context.Background(), userID, 100)
		if len(postIDs) == 0 {
			c.JSON(http.StatusOK, gin.H{"posts": []dto.Post{}})
			return
		}
		ids := make([]int64, 0, len(postIDs))
		for _, pid := range postIDs {
			id, err := strconv.ParseInt(pid, 10, 64)
			if err == nil {
				ids = append(ids, id)
			}
		}
		ma, _ := accessor.GetMySQLAccessor()
		posts, _ := ma.GetPostsByIDs(ids)
		c.JSON(http.StatusOK, gin.H{"posts": posts})
	})
}

func (rl *RestListener) Start(port string) error {
	rl.registerRoutes()
	return rl.Engine.Run(port)
}
