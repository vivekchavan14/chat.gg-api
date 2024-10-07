package controllers

import (
	"log"
	"net/http"
	"server/models"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageController struct {
	DB *gorm.DB
}

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan models.Message)
)

func (mc *MessageController) HandleConnections(ctx *gin.Context, jwtKey []byte) {
	tokenString := ctx.Query("token")
	contactPhone := ctx.Query("contact")

	if tokenString == "" || contactPhone == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing token or contact",
		})
		return
	}

	claims, err := utils.ParseJWT(tokenString, jwtKey)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	username := claims.Username

	var contact models.User
	if err := mc.DB.Where("phone_number = ?", contactPhone).First(&contact).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Contact not found"})
		return
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Websocket Upgrade error : ", err)
		return
	}

	defer ws.Close()

	clients[ws] = username

	for {
		var msg struct {
			Content string `json:"content"`
		}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Websocket Read error : ", err)
			delete(clients, ws)
			break
		}

		newMsg := models.Message{
			Sender:    username,
			Recipient: contact.Username,
			Content:   msg.Content,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		for client, user := range clients {
			if user == contact.Username {
				if err := client.WriteJSON(newMsg); err != nil {
					log.Println("Websocket Write error : ", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func (mc *MessageController) HandleBroadcast() {
	for {
		msg := <-broadcast
		for client, user := range clients {
			if user == msg.Recipient {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Println("Websocket Write error : ", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func (mc *MessageController) GetMessages(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)
	contactPhone := ctx.Query("with")
	if contactPhone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Contact phone number is required",
		})
		return
	}

	var contact models.User
	if err := mc.DB.Where("phone_number = ?", contactPhone).First(&contact).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Contact not found",
		})
		return
	}

	var messages []models.Message
	if err := mc.DB.Where("(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)",
		username, contact.Username, contact.Username, username).Order("created_at asc").Find(&messages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch messages"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
