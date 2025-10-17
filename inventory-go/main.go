package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

func main() {
	app := tview.NewApplication()
	LoadInv()
	inventoryList := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	inventoryList.SetBorder(true).SetTitle("Inventory Items")

	// refresh Inv
	refreshInv := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "No Items Found.")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	nameInput := tview.NewInputField().SetLabel("Item Name: ")
	stockInput := tview.NewInputField().SetLabel("Stock: ")
	idInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	form := tview.NewForm().
		AddFormItem(nameInput).
		AddFormItem(stockInput).
		AddFormItem(idInput).
		AddButton("Add Item", func() {
			name := nameInput.GetText()
			stock := stockInput.GetText()
			if name != "" && stock != "" {
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Println(inventoryList, "Invalid Stock value")
					return
				}
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				SaveInv()
				refreshInv()
				nameInput.SetText("")
				stockInput.SetText("")
			}
		}).
		AddButton("Delete Item", func() {
			idStr := idInput.GetText()
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Enter ID to delete.")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprint(inventoryList, "Invalid item ID")
				return
			}
			DelItem(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted. \n", id)
			refreshInv()
			idInput.SetText("")
		}).
		AddButton("Edit", func() {
			idStr := idInput.GetText()
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Enter ID to edit.")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
				return
			}

			name := nameInput.GetText()
			stock := stockInput.GetText()

			if name == "" && stock == "" {
				fmt.Fprintln(inventoryList, "Enter at least one field to update.")
				return
			}
			EditItem(id-1, name, stock)
			fmt.Fprintf(inventoryList, "Item [%d] updated.\n", id)
			refreshInv()
			nameInput.SetText("")
			stockInput.SetText("")
			idInput.SetText("")
		}).
		AddButton("Exit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false).
		AddItem(form, 0, 1, true)

	refreshInv()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func LoadInv() {
	if _, err := os.Stat(inventoryFile); err == nil {
		{
			data, err := os.ReadFile(inventoryFile)
			if err != nil {
				log.Fatal("Error Reading File: - ", err)
			}
			json.Unmarshal(data, &inventory)
		}
	}
}

func SaveInv() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("error saving : - ", err)
	}
	os.WriteFile(inventoryFile, data, 0o644)
}

func DelItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid Index!")
		return
	}
	inventory = append(inventory[:index], inventory[index+1:]...)
	SaveInv()
}

func EditItem(index int, name, stock string) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid Index!")
		return
	}

	if name != "" {
		inventory[index].Name = name
	}

	if stock != "" {
		qty, err := strconv.Atoi(stock)
		if err != nil {
			fmt.Print("Invalid Stock value!")
			return
		}
		inventory[index].Stock = qty
	}

	SaveInv()
}
