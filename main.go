package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Transaction struct {
	TxID    string
	Fee     int
	Weight  int
	Parents []string
}

func main() {
	transactions, err := readCSV("mempool.csv")
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Sort transactions by parents
	sortedByParents := sortTrxByParents(transactions)

	// Sort transactions by fee
	sortedByFee := getFeeArray(sortedByParents)

	// Select transactions by weight
	selectedByWeight := selectTransactionsByWeight(sortedByFee)

	totalFee := 0
	totalWeight := 0
	for _, blockTrx := range selectedByWeight {
		fmt.Println(blockTrx.TxID)
		totalFee += blockTrx.Fee
		totalWeight += blockTrx.Weight
	}
	fmt.Println("Total Fee:", totalFee)
	fmt.Println("Total Weight:", totalWeight)
}

func readCSV(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	for _, row := range rows {
		line := strings.Split(row[0],",")
		fee, _ := strconv.Atoi(line[1])
		weight, _ := strconv.Atoi(line[2])
		var parents []string
		if line[3] != "" {
			parents = strings.Split(line[3], ",")
		}
		transaction := Transaction{
			TxID:    line[0],
			Fee:     fee,
			Weight:  weight,
			Parents: parents,
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func sortTrxByParents(transactions []Transaction) []Transaction {
	// Create a map to store transaction by TxID
	txMap := make(map[string]Transaction)
	for _, trx := range transactions {
		txMap[trx.TxID] = trx
	}

	// Remove transactions which parents are not included
	for _, trx := range transactions {
		for _, parent := range trx.Parents {
			parentTx, ok := txMap[parent]
			if !ok {
				delete(txMap, trx.TxID)
				break
			}
			if parentTx.TxID != trx.TxID {
				delete(txMap, trx.TxID)
				break
			}
		}
	}

	// Create a slice to store transactions ordered by parents
	var sortedTransactions []Transaction
	for _, trx := range txMap {
		sortedTransactions = append(sortedTransactions, trx)
	}

	return sortedTransactions
}

func getFeeArray(transactions []Transaction) map[int]Transaction {
	feeArray := make(map[int]Transaction)
	for _, transaction := range transactions {
		feeArray[transaction.Fee] = transaction
	}
	return feeArray
}

func selectTransactionsByWeight(transactions map[int]Transaction) []Transaction {
	// Create a slice to store transactions ordered by fee
	var sortedByFee []Transaction
	for _, trx := range transactions {
		sortedByFee = append(sortedByFee, trx)
	}

	// Sort transactions by fee in descending order
	sort.Slice(sortedByFee, func(i, j int) bool {
		return sortedByFee[i].Fee > sortedByFee[j].Fee
	})

	// Select transactions by weight until total weight reaches 4000000
	totalWeight := 0
	var selectedTransactions []Transaction
	for _, trx := range sortedByFee {
		if totalWeight+trx.Weight <= 4000000 {
			selectedTransactions = append(selectedTransactions, trx)
			totalWeight += trx.Weight
		}
	}

	return selectedTransactions
}
