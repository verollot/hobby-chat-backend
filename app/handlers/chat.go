package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Message struct {
	ID        string    `json:"id"`
	Msg       string    `json:"message"`
	Recipient uuid.UUID `json:"recipient"`
}

type connexionGetter interface {
	AddConn(id uuid.UUID, conn *websocket.Conn) error
	RemoveConn(id uuid.UUID) error
	GetConn(id uuid.UUID) (*websocket.Conn, error)
}

type HandlerChat struct {
	cg connexionGetter
}

func NewHandlerChat(cg connexionGetter) HandlerChat {
	return HandlerChat{cg: cg}
}

func (hc *HandlerChat) RegisterClientSocket(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}

	log.Infof("upgrading connection [%s]", userIDStr)
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Errorf("fail to upgrade %s [%s],", err, userIDStr)
	}
	defer func() {
		conn.Close(websocket.StatusInternalError, "socket on error")
		hc.cg.RemoveConn(userID)
	}()

	err = hc.cg.AddConn(userID, conn)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	var msg Message
	for {
		err := wsjson.Read(c.Request.Context(), conn, &msg)
		if err != nil {
			log.Errorf("fail to read closing conn... %+v [%s]", err.Error(), userID.String())
			conn.Close(websocket.StatusNormalClosure, "everything's fine")
			break
		}
		log.Infof("message from %s => %+v", userIDStr, msg)

		//todo : store in db

		// find recipient
		rConn, err := hc.cg.GetConn(msg.Recipient)
		if err != nil {
			log.Errorf("recipient not found %s [%s]", msg.Recipient, userID.String())
			continue
		}

		msg.ID = fmt.Sprintf("%s-%d", userIDStr, time.Now().Nanosecond())
		outMsg, _ := json.Marshal(msg)
		rConn.Write(c.Request.Context(), websocket.MessageText, outMsg)
	}
}
