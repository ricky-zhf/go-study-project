package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID   int64  `json:"id" form:"id" gorm:"primarykey"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

//TableName
func (p *User) TableName() string {
	return "user"
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("mysql", "root:123@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//自动迁移 - 功能就是把结构体和数据表进行匹配
	db.AutoMigrate(&User{})

	//create()
	//update()
	copyF()
}

func copyF() {
	var user User
	user.ID = 21
	err := db.First(&user).Error
	fmt.Println(err)
	fmt.Println(user)

	user.ID = 0
	err = db.Create(&user).Error
	fmt.Println(err)
	fmt.Println(user)

}

func create() {
	user := User{Name: "zzz", Age: 18}
	err := db.Create(&user).Error
	fmt.Println(err)
}

func update() {
	user := User{ID: 21, Age: 122}
	if err := db.Model(&user).Update(user).Error; err != nil {
		fmt.Println(err)
	}
}
