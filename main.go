package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
)

func main() {
	router := gin.Default()
	router.GET("/rate", getRrateHandler)
	router.POST("/subscribe", subscribeHandler)
	router.POST("/sendEmails", sendEmailHandler)
	router.Run("0.0.0.0:9090")
}

type email struct {
	Email string `json:"email"`
}

type rate struct {
	Rate float64 `json:"rate"`
}

func getRate() (float64, error) {
	var response rate

	url := "https://rest.coinapi.io/v1/exchangerate/BTC/UAH"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-CoinAPI-Key", "A43BCDB6-1495-4612-9208-870E47D0ECCB") // will be expired soon

	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return response.Rate, nil
}

func sendRateToEmail(email string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	to := []string{email}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	rate, err := getRate()
	if err != nil {
		return errors.New("Failed to fetch rate")
	}

	subject := "Subject: Current rate for BTC/UAH\n"
	body := "Current rate for BTC/UAH is " + strconv.FormatFloat(rate, 'f', -1, 64)
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)
	return smtp.SendMail(address, auth, from, to, message)
}
func getRrateHandler(c *gin.Context) {
	rate, err := getRate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rate",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"rate": rate})
}

func isEmailAlreadySubscribed(email string) (bool, error) {
	file, err := os.Open("subscribers.txt")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subscribedEmail := scanner.Text()
		if subscribedEmail == email {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
func subscribeHandler(c *gin.Context) {
	var payload email

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON payload",
		})
		return
	}

	email := payload.Email
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
		return
	}

	alreadySubscribed, err := isEmailAlreadySubscribed(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check email",
		})
		return
	}

	if alreadySubscribed {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email is already subscribed",
		})
		return
	}

	file, err := os.OpenFile("subscribers.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save email",
		})
		return
	}
	defer file.Close()

	if _, err := file.WriteString(email + "\n"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email subscribed successfully",
	})
}

func sendEmailHandler(c *gin.Context) {
	file, err := os.Open("subscribers.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read subscribers file",
		})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		sendRateToEmail(email)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Emails sent successfully",
	})
}
