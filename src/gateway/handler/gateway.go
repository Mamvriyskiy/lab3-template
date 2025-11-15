package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	modelGateway "github.com/Mamvriyskiy/lab3-template/src/gateway/model"
	"github.com/Mamvriyskiy/lab3-template/src/gateway/rollback"
	"github.com/gin-gonic/gin"
)

func forwardRequest(c *gin.Context, method, targetURL string, headers map[string]string, body []byte) (int, []byte, http.Header, error) {
	if len(c.Request.URL.RawQuery) > 0 {
		targetURL = fmt.Sprintf("%s?%s", targetURL, c.Request.URL.RawQuery)
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewReader(body))
	if err != nil {
		return 0, nil, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Копируем оригинальный Content-Type, если есть
	if c.Request.Header.Get("Content-Type") != "" {
		req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, err
	}

	return resp.StatusCode, respBody, resp.Header, nil
}

func (h *Handler) GetInfoAboutFlight(c *gin.Context) {
	status, body, headers, err := forwardRequest(c, "GET", "http://flight:8060/flight", nil, nil)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Data(status, headers.Get("Content-Type"), body)
}

func (h *Handler) GetInfoAboutUserTicket(c *gin.Context) {
	ticketUid := c.Param("ticketUid")
	if ticketUid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticketUid is required"})
		return
	}

	// 1️⃣ Запрашиваем билет
	ticketURL := "http://ticket:8070/ticket/" + ticketUid
	status, body, _, err := forwardRequest(c, "GET", ticketURL, nil, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if status != http.StatusOK {
		c.Data(status, "application/json", body)
		return
	}

	var ticket modelGateway.Ticket
	if err := json.Unmarshal(body, &ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse ticket response"})
		return
	}

	// 2️⃣ Запрашиваем данные о рейсе
	flightURL := "http://flight:8060/flight/" + ticket.FlightNumber
	flightStatus, flightBody, _, err := forwardRequest(c, "GET", flightURL, nil, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if flightStatus != http.StatusOK {
		c.Data(flightStatus, "application/json", flightBody)
		return
	}

	var flight modelGateway.Flight
	if err := json.Unmarshal(flightBody, &flight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse flight response"})
		return
	}

	// 3️⃣ Объединяем оба ответа
	response := modelGateway.TicketInfo{
		TicketUID:    ticket.TicketUID,
		FlightNumber: flight.FlightNumber,
		FromAirport:  flight.FromAirport,
		ToAirport:    flight.ToAirport,
		Date:         flight.Datetime,
		Price:        flight.Price,
		Status:       ticket.Status,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetInfoAboutAllUserTickets(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	headers := map[string]string{"X-User-Name": username}
	status, body, respHeaders, err := forwardRequest(c, "GET", "http://ticket:8070/tickets", headers, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if status != http.StatusOK {
		c.Data(status, respHeaders.Get("Content-Type"), body)
		return
	}

	var tickets []modelGateway.Ticket
	if err := json.Unmarshal(body, &tickets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse tickets"})
		return
	}

	var ticketInfos []modelGateway.TicketInfo
	for _, ticket := range tickets {
		if ticket.FlightNumber == "" {
			continue
		}
		flightURL := "http://flight:8060/flight/" + ticket.FlightNumber
		flightStatus, flightBody, _, err := forwardRequest(c, "GET", flightURL, nil, nil)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		if flightStatus != http.StatusOK {
			c.Data(flightStatus, "application/json", flightBody)
			return
		}

		var flight modelGateway.Flight
		if err := json.Unmarshal(flightBody, &flight); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse flight response"})
			return
		}

		ticketInfos = append(ticketInfos, modelGateway.TicketInfo{
			TicketUID:    ticket.TicketUID,
			FlightNumber: flight.FlightNumber,
			FromAirport:  flight.FromAirport,
			ToAirport:    flight.ToAirport,
			Date:         flight.Datetime,
			Price:        flight.Price,
			Status:       ticket.Status,
		})
	}

	// 4️⃣ Отправляем массив в ответе
	c.JSON(http.StatusOK, ticketInfos)
}

func (h *Handler) GetInfoAboutUserPrivilege(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	headers := map[string]string{"X-User-Name": username}
	status, body, respHeaders, err := forwardRequest(c, "GET", "http://bonus:8050/privilege", headers, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
        	"message": "Bonus Service unavailable",
    	})
		return
	}

	c.Data(status, respHeaders.Get("Content-Type"), body)
}

type CombinedResponse struct {
	Tickets   []modelGateway.TicketInfo `json:"tickets"`
	Privilege struct {
		Balance int    `json:"balance"`
		Status  string `json:"status"`
	} `json:"privilege"`
}

func (h *Handler) GetInfoAboutUser(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	headers := map[string]string{"X-User-Name": username}
	status, body, respHeaders, err := forwardRequest(c, "GET", "http://ticket:8070/tickets", headers, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if status != http.StatusOK {
		c.Data(status, respHeaders.Get("Content-Type"), body)
		return
	}

	var tickets []modelGateway.Ticket
	if err := json.Unmarshal(body, &tickets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse tickets"})
		return
	}

	var ticketInfos []modelGateway.TicketInfo
	for _, ticket := range tickets {
		if ticket.FlightNumber == "" {
			continue
		}
		flightURL := "http://flight:8060/flight/" + ticket.FlightNumber
		flightStatus, flightBody, _, err := forwardRequest(c, "GET", flightURL, nil, nil)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		if flightStatus != http.StatusOK {
			c.Data(flightStatus, "application/json", flightBody)
			return
		}

		var flight modelGateway.Flight
		if err := json.Unmarshal(flightBody, &flight); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse flight response"})
			return
		}

		ticketInfos = append(ticketInfos, modelGateway.TicketInfo{
			TicketUID:    ticket.TicketUID,
			FlightNumber: flight.FlightNumber,
			FromAirport:  flight.FromAirport,
			ToAirport:    flight.ToAirport,
			Date:         flight.Datetime,
			Price:        flight.Price,
			Status:       ticket.Status,
		})
	}

	status, BonusBody, respHeaders, err := forwardRequest(c, "GET", "http://bonus:8050/privilege", headers, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"tickets": tickets,
			"privilege": gin.H{},
		})
		return
	}

	var bonus modelGateway.PrivilegeResponse
	if err := json.Unmarshal(BonusBody, &bonus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp CombinedResponse
	resp.Tickets = ticketInfos
	resp.Privilege.Balance = bonus.Balance
	resp.Privilege.Status = bonus.Status

	c.JSON(http.StatusOK, resp)
}

type BuyTicket struct {
	FlightNumber    string `json:"flightNumber"`
	Price           int    `json:"price"`
	PaidFromBalance bool   `json:"paidFromBalance"`
}

type TicketResponse struct {
	TicketUID     string `json:"ticketUid"`
	FlightNumber  string `json:"flightNumber"`
	FromAirport   string `json:"fromAirport"`
	ToAirport     string `json:"toAirport"`
	Date          string `json:"date"` // Можно оставить string, если нужен формат "YYYY-MM-DD HH:MM"
	Price         int    `json:"price"`
	PaidByMoney   int    `json:"paidByMoney"`
	PaidByBonuses int    `json:"paidByBonuses"`
	Status        string `json:"status"`
	Privilege     struct {
		Balance int    `json:"balance"`
		Status  string `json:"status"`
	} `json:"privilege"`
}

type PrivilegeInfo struct {
	Status      string `db:"status"`
	Balance     int    `db:"balance"`
	BalanceDiff int    `db:"balance_diff"`
}

func (h *Handler) BuyTicketUser(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}
	headers := map[string]string{"X-User-Name": username}

	// Читаем тело запроса
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Получаем информацию о рейсе
	var reqData struct {
		FlightNumber    string `json:"flightNumber"`
		Price           int    `json:"price"`
		PaidFromBalance bool   `json:"paidFromBalance"`
	}

	if err := json.Unmarshal(bodyBytes, &reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Покупаем билет
	status, body, _, err := forwardRequest(c, "POST", "http://ticket:8070/ticket", headers, bodyBytes)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	uid := strings.TrimSpace(string(body))

	var privilege = PrivilegeInfo{
		Status:      "GOLD",
		Balance:     1650,
		BalanceDiff: 0,
	}

	if reqData.PaidFromBalance {
		curlBouns := "http://bonus:8050/bonus/" + uid + "/" + strconv.Itoa(reqData.Price)
		// Получаем данные с бонусного счета
		statusBonus, bodyBonus, _, err := forwardRequest(c, "PATCH", curlBouns, headers, nil)
		if err != nil || status >= 400 {
			// rollback
			// _, _, _, rollbackErr := forwardRequest(c, "DELETE", "http://ticket:8070/ticket/", headers, bodyBytes)
			// if rollbackErr != nil {
			// 	log.Printf("Rollback failed: %v", rollbackErr)
			// }

			// ставим в Redis очередь на повтор
			if err := rollback.EnqueueRetry(rollback.RetryRequest{
				Method:  "POST",
				URL:     "http://localhost:8080/ticket",
				Headers: headers,
				Body:    bodyBytes,
			}); err != nil {
				log.Printf("Failed to enqueue request: %v", err)
			}

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Bonus Service unavailable",
			})
			return
		}

		if statusBonus != http.StatusOK {
			c.Status(statusBonus)
			return
		}

		if err := json.Unmarshal(bodyBonus, &privilege); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	curlUpdateBouns := "http://bonus:8050/bonusUpdate/" + uid + "/" + strconv.Itoa(reqData.Price)
	_, _, _, err = forwardRequest(c, "PATCH", curlUpdateBouns, headers, nil)
	if err != nil || status >= 400 {
		// rollback
		// _, _, _, rollbackErr := forwardRequest(c, "DELETE", "http://ticket:8070/ticket/", headers, bodyBytes)
		// if rollbackErr != nil {
		// 	log.Printf("Rollback failed: %v", rollbackErr)
		// }

		// ставим в Redis очередь на повтор
		if err := rollback.EnqueueRetry(rollback.RetryRequest{
			Method:  "POST",
			URL:     "http://localhost:8080/ticket",
			Headers: headers,
			Body:    bodyBytes,
		}); err != nil {
			log.Printf("Failed to enqueue request: %v", err)
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "queued for retry"})
		return
	}

	flightNumber := reqData.FlightNumber

	flightURL := "http://flight:8060/flight/" + flightNumber
	flightStatus, flightBody, _, err := forwardRequest(c, "GET", flightURL, nil, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if flightStatus != http.StatusOK {
		c.Data(flightStatus, "application/json", flightBody)
		return
	}

	var flight modelGateway.Flight
	if err := json.Unmarshal(flightBody, &flight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse flight response"})
		return
	}

	var paidByBonuses int
	if reqData.PaidFromBalance {
		paidByBonuses = privilege.BalanceDiff
	}

	resultUID := strings.Trim(uid, `"\n\r `)

	result := TicketResponse{
		TicketUID:     resultUID,
		FlightNumber:  flight.FlightNumber,
		FromAirport:   flight.FromAirport,
		ToAirport:     flight.ToAirport,
		Date:          flight.Datetime,
		Price:         reqData.Price,
		PaidByMoney:   flight.Price - paidByBonuses,
		PaidByBonuses: paidByBonuses,
		Status:        "PAID",
		Privilege: struct {
			Balance int    `json:"balance"`
			Status  string `json:"status"`
		}{
			Balance: privilege.Balance,
			Status:  privilege.Status,
		},
	}

	c.JSON(status, result)
}

func (h *Handler) DeleteTicketUSer(c *gin.Context) {
	ticketUid := c.Param("ticketUid")
	if ticketUid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticketUid is required"})
		return
	}

	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	ticketURL := "http://ticket:8070/ticket/" + ticketUid
	status, body, _, err := forwardRequest(c, "PATCH", ticketURL, nil, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if status != http.StatusOK {
		c.Data(status, "application/json", body)
		return
	}

	c.Status(http.StatusNoContent)
}
