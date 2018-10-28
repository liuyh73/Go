// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/liuyh73/Go/agenda/models"
	"github.com/spf13/cobra"
)

var user models.User

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "This command can register user",
	Long:  `You can use agenda register to sign up one user`,
	Run: func(cmd *cobra.Command, args []string) {
		// args 存储没有被flag解析的参数
		// fmt.Println(args)
		if user.Username == "" {
			checkErr("The username is required")
		}

		if user.Username == "" {
			checkErr("The password is required")
		}
		rows, err := models.DBSql.Query("SELECT * FROM user WHERE username=?", user.Username)
		if err != nil {
			checkErr("Fail to query the databse!")
		}
		if rows.Next() {
			checkErr("The username " + user.Username + " had been register")
		}

		stmt, err := models.DBSql.Prepare("INSERT into user(username, password, email, telephone, login) value(?,?,?,?,?)")
		if err != nil {
			checkErr("Fail to insert into table")
		}
		_, err = stmt.Exec(user.Username, user.Password, user.Email, user.Telephone, false)
		if err != nil {
			checkErr("Fail to insert into table")
		}
		fmt.Println("Register", user.Username, "successfully!")
		models.Logger.Println("Register", user.Username, "successfully!")
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	registerCmd.Flags().StringVarP(&user.Username, "username", "u", "", "The User's Username")
	registerCmd.Flags().StringVarP(&user.Password, "password", "p", "", "The User's Password")
	registerCmd.Flags().StringVarP(&user.Email, "email", "e", "", "The User's Email")
	registerCmd.Flags().StringVarP(&user.Telephone, "telephone", "P", "", "The User's telephone")
	models.Logger.SetPrefix("[agenda register]")
}
