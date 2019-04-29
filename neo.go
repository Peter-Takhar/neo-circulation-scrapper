package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//to get the NEO client node, run `docker run -it -d -p 10332:10332 petertakhar/ubuntu-neo-cli`
//ubuntu-neo-cli can be found at petertakhar/ubuntu-neo-cli

var (
	//address of your local node
	clientNode     = "localhost:10332"
	holdingAddress = "AQVh2pG732YvtNaxEGkQUei3YA4cvo7d2i"
)

//MAXSUPPLY is the total cap of NEO coins to be released
const MAXSUPPLY = 100000000

//set a timer on the length of http connections
var netClient = &http.Client{
	Timeout: time.Second * 15,
}

//Account is a representation of the account's assets and their values
type Account struct {
	Result struct {
		Balances []struct {
			Asset string `json:"asset"`
			Value string `json:"value"`
		} `json:"balances"`
	} `json:"result"`
}

//calculates the circulating supply of NEO
func main() {

	//list of static rpc nodes found at https://github.com/"CityOfZion/neo-api-js/wiki
	//used to verify if local blockchain is fully synchronized with the main net
	var staticNodes [4]string
	staticNodes[0] = "seed2.neo.org:10332"
	staticNodes[1] = "seed3.neo.org:10332"
	staticNodes[2] = "seed4.neo.org:10332"
	staticNodes[3] = "seed5.neo.org:10332"

	var node string

	//sets the node to pull data depending on synchronization status of local node
	previousAmount := 0.0
	for {
		if isBlockchainSynchronized(staticNodes[0]) {
			node = clientNode
		} else {
			fmt.Println("Local NEO blockchain is not synchronized yet." +
				" Pulling data from another fully synchronized node.\n")
			node = staticNodes[0]
		}

		url := fmt.Sprintf("http://%s?jsonrpc=2.0&method=getaccountstate&params=[\"%s\"]&id=1",
			node, holdingAddress)

		resp, err := netClient.Get(url)
		if err != nil {
			log.Fatalf("Failed to connect to the local node: %v", err)
		}

		accountData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to retrive json from node: %v", err)
		}
		resp.Body.Close()

		var account Account
		err = json.Unmarshal(accountData, &account)
		if err != nil {
			log.Fatalf("Failed to unmarshal Account data: %v.", err)
		}

		/*the index is calculated using account.Result.Balances-1 because if there are
		  multiple assets, then the balance of NEO coins are always the last index
		*/
		accountBalanceNEO, err := strconv.ParseFloat(
			account.Result.Balances[len(account.Result.Balances)-1].Value, 64)
		if err != nil {
			log.Fatalf("Failed to convert account balancee to float: %v.", err)
		}

		circSupplyAmount := MAXSUPPLY - accountBalanceNEO
		fmt.Printf("Circrulating Supply Amount: %f.\n\n", circSupplyAmount)

		if previousAmount < circSupplyAmount {
			fmt.Println("Updated amount.")
			previousAmount = circSupplyAmount
		} else {
			fmt.Println("The previous amount is the same as the scraped data.")
		}
		time.Sleep(10 * time.Second)
	}

}

/*determines if local node has up-to-date blockchain by comparing block heights
with a fully synchonized node
*/
func isBlockchainSynchronized(otherNode string) bool {

	otherNodeHeight := getClientBlockHeight(otherNode)
	myHeight := getClientBlockHeight(clientNode)
	fmt.Printf("Current node block height: %f. Other node block height: %f.\n",
		myHeight, otherNodeHeight)
	return myHeight >= otherNodeHeight
}

//pulls the block height from a NEO node
func getClientBlockHeight(ip string) float64 {
	query := "?jsonrpc=2.0&method=getblockcount&params=[]&id=1"

	resp, err := netClient.Get(fmt.Sprintf("http://%s%s", ip, query))
	if err != nil {
		log.Fatalf("Failed to connect to node: %v.", err)
	}

	blockHeightContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to retrive json from node: %v.", err)
	}
	resp.Body.Close()

	blockHeightData := make(map[string]interface{})
	err = json.Unmarshal(blockHeightContent, &blockHeightData)
	if err != nil {
		log.Fatalf("Failed to unmarshal block height data: %v.", err)
	}

	height, ok := blockHeightData["result"].(float64)
	if !ok {
		log.Fatalf("JSON value must be float64.")
	}

	return height
}
