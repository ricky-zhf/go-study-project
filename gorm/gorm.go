package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id       uint    `gorm:"column:id;AUTO_INCREMENT;primary_key"` // 用户编号
	Username string  `gorm:"column:username;unique;NOT NULL"`      // 用户名
	Password string  `gorm:"column:password;NOT NULL"`             // 密码
	Email    string  `gorm:"column:email;unique;NOT NULL"`         // 邮箱
	Age      uint    `gorm:"column:age;default:18;NOT NULL"`       // 年龄
	Sex      string  `gorm:"column:sex;default:baomi;NOT NULL"`    // 性别
	Tel      string  `gorm:"column:tel;unique;NOT NULL"`           // 电话
	Addr     string  `gorm:"column:addr;default:beijing;NOT NULL"` // 地址
	Card     string  `gorm:"column:card;unique;NOT NULL"`          // 身份证号
	Married  int     `gorm:"column:married;default:0;NOT NULL"`    // 0代表未结婚，1代表已结婚
	Salary   float64 `gorm:"column:salary;default:0;NOT NULL"`     // 薪水
}

func (m *User) TableName() string {
	return "user"
}

type UserV2 struct {
	Id       uint   `gorm:"column:id;AUTO_INCREMENT;primary_key"` // 用户编号
	Username string `gorm:"column:username;unique;NOT NULL"`      // 用户名
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//自动迁移 - 功能就是把结构体和数据表进行匹配
	db.AutoMigrate(&User{})

	//create()
	//update()
	//copyF()
	//findV2()
	findV3()
}

func findV2() {
	var user User
	//db.First(&user)
	//fmt.Println(user)

	var v2 UserV2
	db.Model(&User{}).Where("id=?", 1).Find(&v2)
	fmt.Println(v2)
	fmt.Println(user)
}

func findV3() {
	var v2 UserV2
	sql := `select id, username, age from user where id = ?`
	db.Raw(sql, 1).Scan(&v2)
	fmt.Println(v2)
}

func copyF() {
	var user User
	user.Id = 21
	err := db.First(&user).Error
	fmt.Println(err)
	fmt.Println(user)

	user.Id = 0
	err = db.Create(&user).Error
	fmt.Println(err)
	fmt.Println(user)

}

func create() {
	user := User{Username: "zzz", Age: 18}
	err := db.Create(&user).Error
	fmt.Println(err)
}

func update() {
	user := User{Id: 21, Age: 122}
	if err := db.Model(&user).Update(user).Error; err != nil {
		fmt.Println(err)
	}
}
