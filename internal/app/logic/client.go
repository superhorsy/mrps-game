package logic

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"

	"mrps-game/internal/app/logic/message"
	"mrps-game/internal/app/logic/model"
	"mrps-game/internal/app/logic/model/game"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

const channelBufSize = 100

// Client struct holds client-specific variables.
type Client struct {
	Id    uint32
	Name  string
	Funds model.Funds

	// ActiveChallenge is a challenge that is currently being played.
	ActiveChallenge *model.Challenge
	// PendingChallenges is a map of challenges that client has received from other clients.
	PendingChallenges map[uint32]*model.Challenge

	ws     *websocket.Conn
	ch     chan []byte
	doneCh chan bool
	server *GameServer
}

// NewClient initializes a new Client struct with given websocket.
func NewClient(id uint32, username string, ws *websocket.Conn, server *GameServer) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan []byte, channelBufSize)
	doneCh := make(chan bool)

	return &Client{
		Id:                id,
		Name:              username,
		PendingChallenges: make(map[uint32]*model.Challenge),

		ws:     ws,
		ch:     ch,
		doneCh: doneCh,
		server: server,
	}
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	defer func() {
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening write to client")
	for {
		select {
		case txt := <-c.ch:
			err := c.ws.WriteMessage(websocket.TextMessage, txt)
			if err != nil {
				log.Println(err)
			}
		case <-c.doneCh:
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	defer func() {
		c.leave()
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening read from client")
	for {
		select {
		case <-c.doneCh:
			c.doneCh <- true
			return
		default:
			c.readAndProcessFromWebSocket()
		}
	}
}

func (c *Client) leave() {
	c.doneCh <- true
	c.server.Clients.Remove(c.Id)
}

func (c *Client) readAndProcessFromWebSocket() {
	messageType, data, err := c.ws.ReadMessage()

	if errors.Is(err, websocket.ErrCloseSent) {
		log.Println("Client closed connection")
		c.leave()
		return
	} else if err != nil {
		log.Println("Error reading from websocket:", err)
		c.leave()
		return
	}

	if messageType != websocket.TextMessage {
		log.Println("Non text message received, ignoring")
		return
	}

	rawJSON, err := c.readMessage(data)
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}
	c.processMessage(rawJSON)
}

func (c *Client) readMessage(data []byte) (map[string]interface{}, error) {
	rawJSON := map[string]interface{}{}
	err := json.Unmarshal(data, &rawJSON)
	if err != nil {
		return rawJSON, err
	}
	return rawJSON, nil
}
func (c *Client) processMessage(rawJSON map[string]interface{}) {
	eventType := rawJSON["type"].(string)
	switch eventType {
	case "me":
		c.sendAsJSON(message.NewClientInfo(c.Name, c.Funds.Amount, c.Id))
	case "opponents":
		c.sendAsJSON(message.NewOpponents(c.server.Clients.GetOpponents(c.Id)))
	case "funds.add":
		c.addFunds(rawJSON)
	case "funds.withdraw":
		c.withdrawFunds(rawJSON)
	case "challenge.send":
		c.sendChallenge(rawJSON)
	case "challenge.cancel":
		c.cancelChallenge()
	case "challenge.list":
		c.listPendingChallenges()
	case "challenge.decline":
		c.declineChallenge(cast.ToUint32(rawJSON["opponent"]))
	case "challenge.accept":
		c.acceptChallenge(rawJSON)
	case "transaction.list":
		c.listTransactions()
	default:
		log.Println("Unknown message type: ", eventType)
	}
}
func (c *Client) declineChallenge(opponentId uint32) {
	// check if client has challenge from opponent
	if _, ok := c.PendingChallenges[opponentId]; !ok {
		c.sendError(errors.New("no challenge from opponent"))
		return
	}
	if opponent, ok := c.server.Clients.Get(opponentId); ok {
		opponent.removeActiveChallengeAndUnblockFunds()
		opponent.sendDeclineMessage(c.Id)
	}
}

func (c *Client) removeActiveChallengeAndUnblockFunds() {
	if c.ActiveChallenge != nil {
		_ = c.Funds.Unblock(c.ActiveChallenge.Amount)
		c.ActiveChallenge = nil
	}
	if opponent, ok := c.server.Clients.Get(c.ActiveChallenge.OpponentId); ok {
		opponent.PendingChallenges[c.Id] = nil
	}
}

func (c *Client) sendDeclineMessage(opponentId uint32) {
	c.sendMessage([]byte(`{"type":"challenge.declined", "opponent":` + cast.ToString(opponentId) + `}`))
}

func (c *Client) listPendingChallenges() {
	challenges := make([]message.Challenge, 0, len(c.PendingChallenges))
	for _, challenge := range c.PendingChallenges {
		if challenge != nil {
			challenges = append(challenges, message.Challenge{
				ChallengerId: challenge.ChallengerId,
				OpponentId:   challenge.OpponentId,
				Amount:       challenge.Amount,
				Game:         challenge.Game.Name(),
			})
		}
	}
	c.sendAsJSON(challenges)
}

func (c *Client) cancelChallenge() {
	if c.ActiveChallenge != nil {
		c.ActiveChallenge = nil
		return
	}
	if opponent, ok := c.server.Clients.Get(c.ActiveChallenge.OpponentId); ok {
		opponent.PendingChallenges[c.Id] = nil
	}
	err := c.Funds.Unblock(c.ActiveChallenge.Amount)
	if err != nil {
		c.sendError(err)
	}
}

func (c *Client) sendChallenge(rawJSON map[string]interface{}) {
	challengeSend := &message.ChallengeSend{}
	challengeSendBytes, _ := json.Marshal(rawJSON)
	err := json.Unmarshal(challengeSendBytes, challengeSend)
	if err != nil {
		log.Println("Error parsing challenge.send message:", err)
		return
	}

	// check if opponent is not the same as client
	if c.Id == cast.ToUint32(challengeSend.Opponent) {
		c.sendError(errors.New("opponent is the same as client"))
		return
	}

	// check if client has active challenge
	if c.ActiveChallenge != nil {
		c.sendError(errors.New("already has active challenge"))
		return
	}

	// check if opponent exists
	opponent, ok := c.server.Clients.Get(cast.ToUint32(challengeSend.Opponent))
	if !ok {
		c.sendError(errors.New("opponent not found"))
		return
	}
	// check if client has enough funds
	if !c.Funds.HasAvailableAmount(cast.ToInt64(challengeSend.Amount)) {
		c.sendError(errors.New("not enough funds"))
		return
	}

	// block sum of money and don't allow to withdraw it
	err = c.Funds.Block(cast.ToInt64(challengeSend.Amount))
	if err != nil {
		c.sendError(err)
		return
	}

	var g model.Game
	switch challengeSend.Game {
	case "rps":
		g = game.NewRPSGame()
		err := g.SetMoveForPlayer1(challengeSend.Move)
		if err != nil {
			c.sendError(err)
			return
		}
	case "dice":
		g = game.NewDiceGame()
		// dice is just random game, so we set both moves randomly before accepting challenge
		_ = g.SetMoveForPlayer1(rand.Intn(6) + 1)
		_ = g.SetMoveForPlayer2(rand.Intn(6) + 1)
	default:
		c.sendError(errors.New("unknown game"))
	}
	ch := &model.Challenge{
		ChallengerId: c.Id,
		OpponentId:   cast.ToUint32(challengeSend.Opponent),
		Amount:       cast.ToInt64(challengeSend.Amount),
		Move:         challengeSend.Move,
		Game:         g,
	}
	opponent.PendingChallenges[c.Id], c.ActiveChallenge = ch, ch
}

func (c *Client) acceptChallenge(rawJSON map[string]interface{}) {
	opponentId := cast.ToUint32(rawJSON["opponent"])

	challenge := c.PendingChallenges[opponentId]
	if challenge == nil {
		c.sendError(errors.New("client doesn't have challenge from this client"))
		return
	}
	if !c.Funds.HasAvailableAmount(challenge.Amount) {
		c.sendError(errors.New("not enough funds"))
		return
	}

	if rawJSON["move"] != nil {
		move := rawJSON["move"]
		err := challenge.Game.SetMoveForPlayer2(move)
		if err != nil {
			c.sendError(err)
			return
		}
	}

	c.playGame(opponentId, challenge)
}

func (c *Client) playGame(opponentId uint32, challenge *model.Challenge) {
	opponent, _ := c.server.Clients.Get(opponentId)
	opponent.removeActiveChallengeAndUnblockFunds()
	// ideally better implement transaction-like system, but for now it's enough
	switch challenge.Game.Result() {
	case 1:
		_, _ = c.makeTransaction(float64(challenge.Amount), "game win vs "+opponent.Name)
		_, _ = opponent.makeTransaction(float64(-challenge.Amount), "game lose vs "+c.Name)
		c.sendAsJSON(message.NewGameResult("win", challenge.Amount))
		opponent.sendAsJSON(message.NewGameResult("lose", challenge.Amount))
	case -1:
		_, _ = c.makeTransaction(float64(-challenge.Amount), "game lose vs "+opponent.Name)
		_, _ = opponent.makeTransaction(float64(challenge.Amount), "game win vs "+c.Name)
		c.sendAsJSON(message.NewGameResult("lose", challenge.Amount))
		opponent.sendAsJSON(message.NewGameResult("win", challenge.Amount))
	case 0:
		c.sendAsJSON(message.NewGameResult("draw", challenge.Amount))
		opponent.sendAsJSON(message.NewGameResult("draw", challenge.Amount))
	}
}

func (c *Client) makeTransaction(amount float64, reason string) (int64, error) {
	err := c.server.transService.LogTransaction(c.Id, amount, reason)
	if err != nil {
		return 0, err
	}

	var newTotal int64
	if amount < 0 {
		newTotal, err = c.Funds.Subtract(int64(amount))
		if err != nil {
			return 0, err
		}
	} else {
		newTotal = c.Funds.Add(int64(amount))
	}

	return newTotal, nil
}

func (c *Client) addFunds(rawJSON map[string]interface{}) {
	amount := rawJSON["amount"].(float64)
	newTotal, err := c.makeTransaction(amount, "add")
	if err != nil {
		log.Println("Error adding funds:", err)
		c.sendError(err)
		return
	}
	c.sendAsJSON(message.NewFundsAdded(int64(amount), newTotal))
}
func (c *Client) withdrawFunds(rawJSON map[string]interface{}) {
	amount := rawJSON["amount"].(float64)
	newTotal, err := c.makeTransaction(amount, "withdraw")
	if err != nil {
		log.Println("Error withdrawing funds:", err)
		c.sendError(err)
		return
	}
	c.sendAsJSON(message.NewFundsWithdrawn(int64(amount), newTotal))
}

func (c *Client) sendError(err error) {
	if err != nil {
		errMsg, _ := json.Marshal(message.NewError(err.Error()))
		c.sendMessage(errMsg)
	}
}

func (c *Client) sendAsJSON(msg interface{}) {
	messageJSON, _ := json.Marshal(msg)
	c.sendMessage(messageJSON)
}

// sendMessage sends game state to the client.
func (c *Client) sendMessage(msg []byte) {
	select {
	case c.ch <- msg:
	}
}

func (c *Client) listTransactions() {
	transactions, err := c.server.transService.GetLastTransactions(c.Id)
	if err != nil {
		log.Println("Error getting transactions:", err)
		c.sendError(err)
		return
	}

	// Convert transactions to a format suitable for sending to the client
	transactionMessages := make([]message.Transaction, len(transactions))
	for i, transaction := range transactions {
		transactionMessages[i] = message.Transaction{
			Type:      transaction.Type,
			Amount:    transaction.Amount,
			Reason:    transaction.Reason,
			CreatedAt: transaction.CreatedAt,
		}
	}

	// Send transactions to the client
	c.sendAsJSON(transactionMessages)
}
