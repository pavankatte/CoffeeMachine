package main

import (
	"fmt"
	"flag"
	"sync"
	"encoding/json"
)

// CoffeeMachine struct which contains Machine Filed which will be unmarshalled from input JSON
type CoffeeMachine struct {
	Machine Machine `json:"machine"`
}

// Machine struct which contains fields for Number of Outlets, Total Items quantity and Beverages
type Machine struct {
	Olet Outlets `json:"outlets"`
	TotalItemsQuantity map[string]int `json:"total_items_quantity"`
	Beverages map[string]interface{} `json:"beverages"` 
}

//Beverage struct will have ingredients which is a map of Item string to quantity needed for that beverage
type Beverage struct {
	Ingredients map[string]int 

}
// Outlet struct which contains count of outlets n
type Outlets struct {
	Count int `json:"count_n"`
}

// DoWork function will be called by a go routine to perform the jaob of making a drink/beverage provided in the input
func DoWork(drink string, ingredientIntf interface{}, inventory *Inventory) string {

	ingredients, ok := ingredientIntf.(map[string]interface{}) //Unmarshalling the ingredients needed for given dring
	if !ok {
	  return "Error: Processing ingredients"
	}

	var mutex = &sync.Mutex{}
	mutex.Lock()

	for item,v:=range ingredients {
		value, err := v.(float64)
		if !err {
			return "Error: Processing value of ingredient "+ item 
		}
		currentValue, ok := inventory.Items[item]
		if !ok { //checking if the item needed for the drink is present in the inventory 
			return drink+" cannot be prepared because " + item + " is not available"  //if not ok then we return item not available
		} else if currentValue < int(value) { //checking if quantity of item in inventory is less than it is needed for the drink 
			return drink+" cannot be prepared because item " + item + " is not sufficient" //if yes we return item quantity insufficient
		}
	}
	//We reach this step when we have enough quantity of all items needed for drink in the inventory
	//And we deduct them form the inventory and return drink prepared status
	//error conditions are ignored below as they are checked in the above loop aleady
	for item,v:=range ingredients {
		value, _ := v.(float64)
		currentValue, _ := inventory.Items[item]
		inventory.Items[item] = currentValue - int(value)
	}
	mutex.Unlock()
	return drink+" is prepared"
}

type Inventory struct {
	sync.RWMutex
	Items map[string]int
}


func ExecuteCoffeeMachine(byteValue []byte) {

    var coffeeMachine CoffeeMachine

    err := json.Unmarshal(byteValue, &coffeeMachine)
    if err != nil {
    	fmt.Println("Error: UnMarshelling JSON")
    	return
    }

    var inventory = Inventory{Items: coffeeMachine.Machine.TotalItemsQuantity}

	maxNbConcurrentGoroutines := flag.Int("maxNbConcurrentGoroutines", coffeeMachine.Machine.Olet.Count, "the number of goroutines that are allowed to run concurrently")
	nbJobs := len(coffeeMachine.Machine.Beverages)
	flag.Parse()

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)

	var wg sync.WaitGroup

	answer := make([]string,0)

	//Creating n go routines form the outlet count
	var mutex = &sync.Mutex{}

	for i := 1; i <= nbJobs; i++ {
		wg.Add(1)
		go func(i int) {

			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			var drinkName string
			var ingredients interface{}

			mutex.Lock()
			defer mutex.Unlock()
			//Picking one beverage from the map and pass that to one of the go routine to preapare and 
			for key, value := range  coffeeMachine.Machine.Beverages {
				drinkName = key
				ingredients = value
				break
			}
			//deleting that from map so that no other routine picks up same beverage again
			delete(coffeeMachine.Machine.Beverages, drinkName)
			//Calling Do Work to perform the drink Preparation
			op := DoWork(drinkName, ingredients, &inventory)
			//adding  status of the drink to the output array
			answer = append(answer,op)

			<-concurrentGoroutines

		}(i)

	}
	wg.Wait()
	
	//Printing All the Results
	for _,op := range answer{
		fmt.Println(op)
	}
}