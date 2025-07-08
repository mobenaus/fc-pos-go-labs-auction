package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"io"
	"net/http"
	"os"
	"time"
)

type Test struct {
	baseUrl string
}

func main() {

	baseUrl := os.Getenv("BASE_URL")

	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}

	logger.Info(fmt.Sprintf("Testing %s", baseUrl))

	time.Sleep(time.Second * 10)

	test := Test{
		baseUrl: baseUrl,
	}

	userId, err := test.addUser()
	if err != nil {
		logger.Error("Error createing user", err)
		return
	}

	err = test.createAuction()
	if err != nil {
		logger.Error("Error createing auction", err)
		return
	}

	auctionId, err := test.getAuctionId()
	if err != nil {
		logger.Error("Error createing auction", err)
		return
	}

	test.createBids(auctionId, userId, 30)

	test.getAuction(auctionId)

	test.getWinnerAuction(auctionId)

}

func (t *Test) addUser() (string, error) {
	data := user_usecase.UserInputDTO{
		Name: "Teste Usuário",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error create user json", err)
		return "", err
	}
	res, err := http.Post(fmt.Sprintf("%s/user", t.baseUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error create user request", err)
		return "", err
	}
	logger.Info(fmt.Sprintf("Create user Result %s", res.Status))
	if res.StatusCode != 201 {
		err = fmt.Errorf("result: %s", res.Status)
		logger.Error("Error create user request", err)
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Error("Error read user request", err)
		return "", err
	}
	var user user_usecase.UserOutputDTO
	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Error("Error read user result", err)
		return "", err
	}

	return user.Id, nil
}

func (t *Test) createAuction() error {
	data := auction_usecase.AuctionInputDTO{
		ProductName: "Janela",
		Category:    "casa",
		Description: "Janela para casa da cidade",
		Condition:   1,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error create auction json", err)
		return err
	}
	res, err := http.Post(fmt.Sprintf("%s/auction", t.baseUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error create auction request", err)
		return err
	}
	logger.Info(fmt.Sprintf("Create Auction Result %s", res.Status))
	if res.StatusCode != 201 {
		err = fmt.Errorf("result: %s", res.Status)
		logger.Error("Error create auction request", err)
		return err
	}

	return nil
}

func (t *Test) getAuctionId() (string, error) {

	auctions, err := t.getAuctions()
	if err != nil {
		logger.Error("Error get auctions", err)
		return "", err
	}

	return auctions[len(auctions)-1].Id, nil

}

func (t *Test) getAuctions() ([]auction_usecase.AuctionOutputDTO, error) {
	logger.Info("Wait for query auction results...")

	var res *http.Response
	var err error

	for {
		time.Sleep(time.Second * 1)

		res, err = http.Get(fmt.Sprintf("%s/auction?status=1", t.baseUrl))
		if err != nil {
			logger.Error("Error get auction request", err)
			return nil, err
		}

		if res.ContentLength > 4 {
			break
		}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Error("Error read auction request", err)
		return nil, err
	}
	var auctions []auction_usecase.AuctionOutputDTO
	err = json.Unmarshal(body, &auctions)
	if err != nil {
		logger.Error("Error read auction result", err)
		return nil, err
	}
	return auctions, nil
}

func (t *Test) createBids(auctionId, userId string, bids int) error {

	logger.Info(fmt.Sprintf("auctionId: %s", auctionId))
	logger.Info(fmt.Sprintf("userId: %s", userId))
	for i := range bids {
		err := t.createBid(auctionId, userId, i)
		if err != nil {
			logger.Error("Error create bids", err)
			return err
		}
		time.Sleep(time.Second)
	}
	return nil

}

func (t *Test) createBid(auctionId, userId string, i int) error {
	data := bid_usecase.BidInputDTO{
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    float64(100 + i),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error create bid json", err)
		return err
	}
	res, err := http.Post(fmt.Sprintf("%s/bid", t.baseUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error create bid request", err)
		return err
	}
	logger.Info(fmt.Sprintf("Create Bid %d Result %s", i, res.Status))
	if res.StatusCode != 201 {
		err = fmt.Errorf("result: %s", res.Status)
		logger.Error("Error create bid request", err)
		return err
	}

	return nil
}

func (t *Test) getAuction(auctionId string) error {

	res, err := http.Get(fmt.Sprintf("%s/auction/%s", t.baseUrl, auctionId))
	if err != nil {
		logger.Error("Error get auction request", err)
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error read auction request", err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println("Auction:")
	fmt.Println(prettyJSON.String())

	return nil
}

func (t *Test) getWinnerAuction(auctionId string) error {

	res, err := http.Get(fmt.Sprintf("%s/auction/winner/%s", t.baseUrl, auctionId))
	if err != nil {
		logger.Error("Error get auction request", err)
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error read auction request", err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println("Winner:")
	fmt.Println(prettyJSON.String())

	return nil
}
