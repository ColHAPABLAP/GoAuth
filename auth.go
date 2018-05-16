package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"encoding/json"
)

const clientId = ""
const clientSecret = ""
const callbackEndpoint = ""
const authEndpoint = ""
const tokenEndpoint = ""
const tokenCallback = ""

func main() {
	router := gin.Default()
	oid := router.Group("/oid")
	{
		oid.GET("/login", loginHandler)
		oid.GET("/callback", callbackHandler)
	}

	router.Run()
}

func loginHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, authEndpoint)
}

func callbackHandler(c *gin.Context) {
	var code = c.Query("code")
	fmt.Println("Code: " + code)

	var tokens = getTokens(code)
	c.Redirect(http.StatusMovedPermanently, tokenCallback + "?id_token=" + tokens.IdToken)
}

type tokens struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int `json:"expires_in"`
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func getTokens(code string) (tokens) {
	client := &http.Client{}
	//postParams := url.Values{"code": {code}, "client_id": {clientId}, "redirect_uri": {callbackEndpoint}, "grant_type": {"authorization_code"}}
	postParams := url.Values{}
	postParams.Set("code", code)
	postParams.Add("client_id", clientId)
	postParams.Add("redirect_uri", callbackEndpoint)
	postParams.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(postParams.Encode()))
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
	//fmt.Println(resp.Body)

	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	tokens := tokens{}
	json.NewDecoder(resp.Body).Decode(&tokens)
	fmt.Println(tokens)
	fmt.Println(tokens.TokenType)
	fmt.Println(tokens.ExpiresIn)
	fmt.Println(tokens.AccessToken)
	fmt.Println(tokens.RefreshToken)
	fmt.Println(tokens.IdToken)

	return tokens
}
