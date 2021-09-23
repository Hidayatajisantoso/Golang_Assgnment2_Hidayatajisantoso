package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ordres struct {
	OrderedAt    time.Time `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	Items        orderItems
}

type orderItems struct {
	Item_code   string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/orders_by?parseTime=true")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createData(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	ordered_at := c.PostForm("ordered_at")
	customer_name := c.PostForm("customer_name")
	item_code := c.PostForm("item_code")
	description := c.PostForm("description")
	quantity := c.PostForm("quantity")

	insertOrders, err := db.Exec("insert into orders (customer_name, ordered_at) values (?, ?)", customer_name, ordered_at)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lastId, err := insertOrders.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = db.Exec("insert into items (item_code, description, quantity, order_id) values (?, ?, ?, ?)", item_code, description, quantity, lastId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res := gin.H{
		"result": "Succesfully insert data",
	}
	c.JSON(http.StatusOK, res)
}

func getData(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	var result = ordres{}
	err = db.QueryRow("select ordered_at, customer_name, item_code, description, quantity from orders inner join items on orders.order_id = items.order_id").Scan(&result.OrderedAt, &result.CustomerName, &result.Items.Item_code, &result.Items.Description, &result.Items.Quantity)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
	res := gin.H{
		"result": result,
	}
	c.JSON(http.StatusOK, res)
}

func updateData(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	ordered_at := c.PostForm("ordered_at")
	customer_name := c.PostForm("customer_name")
	item_code := c.PostForm("item_code")
	description := c.PostForm("description")
	quantity := c.PostForm("quantity")

	_, err = db.Exec("update orders set customer_name = ?,  ordered_at = ?", customer_name, ordered_at)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = db.Exec("update items set item_code = ?, description = ?, quantity = ?", item_code, description, quantity)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res := gin.H{
		"result": "Succesfully update data",
	}
	c.JSON(http.StatusOK, res)
}

func deleteData(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("delete from orders")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	res := gin.H{
		"result": "Succesfully delete data",
	}
	c.JSON(http.StatusOK, res)

}

func main() {
	router := gin.Default()
	router.POST("/create", createData)
	router.GET("/orders", getData)
	router.POST("/update", updateData)
	router.POST("/delete", deleteData)

	router.Run(":8080")
}
