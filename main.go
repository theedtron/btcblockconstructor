package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Transaction struct {
    TxID        string
    Fee         int
    Weight      int
    ParentTxIDs []string
}

func main() {
    // Read the mempool.csv file
    file, err := os.Open("mempool.csv")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // Parse the CSV data
    reader := csv.NewReader(file)
    lines, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        return
    }

    // Parse transactions
    var transactions []Transaction
    for _, line := range lines {
		fmt.Println(line[1])
		os.Exit(0)
        fee, _ := strconv.Atoi(line[1])
        weight, _ := strconv.Atoi(line[2])
        parentTxIDs := strings.Split(line[3], ";")

        tx := Transaction{
            TxID:        line[0],
            Fee:         fee,
            Weight:      weight,
            ParentTxIDs: parentTxIDs,
        }
        transactions = append(transactions, tx)
    }

    // Select transactions for the block
    block := selectTransactions(transactions)

    // Output the block
    for _, tx := range block {
        fmt.Println(tx.TxID)
    }
}

func selectTransactions(transactions []Transaction) []Transaction {
    // Sort transactions by fee (descending order)
    sortByFee(transactions)

    // Initialize a map to track transaction inclusion
    included := make(map[string]bool)

    // Initialize block weight and total fee
    blockWeight := 0
    totalFee := 0

    // Initialize the block
    var block []Transaction

    // Iterate through transactions
    for _, tx := range transactions {
        // Check if transaction can be included in the block
        if canIncludeTransaction(tx, included) && blockWeight+tx.Weight <= 4000000 {
            // Include transaction in the block
            block = append(block, tx)
            included[tx.TxID] = true
            blockWeight += tx.Weight
            totalFee += tx.Fee
        }
    }

    fmt.Println("Total fee collected:", totalFee)

    return block
}

func sortByFee(transactions []Transaction) {
    // Sort transactions by fee (descending order)
    for i := range transactions {
        maxIndex := i
        for j := i + 1; j < len(transactions); j++ {
            if transactions[j].Fee > transactions[maxIndex].Fee {
                maxIndex = j
            }
        }
        transactions[i], transactions[maxIndex] = transactions[maxIndex], transactions[i]
    }
}

func canIncludeTransaction(tx Transaction, included map[string]bool) bool {
    // Check if all parent transactions are included
    for _, parentTxID := range tx.ParentTxIDs {
        if !included[parentTxID] {
            return false
        }
    }
    return true
}
