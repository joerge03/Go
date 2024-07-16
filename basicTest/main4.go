package main

import "fmt"

type Product struct {
	Name string
	Price int
}


func productInfo(product ...Product){
	var sum int; 
	lastProduct := product[len(product)-1]

	print("The last item on your lists is: ", lastProduct.Name, "with a price of", lastProduct.Price, "\n")
	print("The total of the item you order is: ", len(product), "\n")
	
	for i:= 0; i < len(product); i++{
		sum += product[i].Price
	}
	print("The total costs is: ", sum, "\n")
}

func main4(){
	products := []Product{
		{
			"safdadssf",23,
		},
		{
			"safdadsfss",23,
		},
		{
			"safdadsf231s",23,
		},
	}

	slicess := make([]Product, 999)

	fmt.Println(slicess)

	products = append(products, products...)
	// products = append(products, Product{"Asdasdf", 30})
	
	// slicess := products[1:3]

	fmt.Println(products)

	// productInfo(products...)

}