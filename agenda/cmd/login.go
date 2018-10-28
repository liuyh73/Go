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

var loginUser models.User

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "This command can login user",
	Long:  `You can use agenda login to sign in one user`,
	Run: func(cmd *cobra.Command, args []string) {
		rows, err := models.DBSql.Query("SELECT * FROM user WHERE username=?", user.Username)
		if err != nil {
			checkErr("Fail to query the databse!")
		}
		if rows.Next() {
			var password string
			rows.Scan(nil, &password)
			if password == loginUser.Password {
				stmt, err1 := models.DBSql.Prepare("UPDATE user SET login=? where username=?")
				_, err2 := stmt.Exec(true, loginUser.Username)
				if err1 != nil || err2 != nil {
					checkErr("Fail to update the login info")
				} else {
					fmt.Println("Login", loginUser.Username, "successfully!")
					models.Logger.Println("Login", loginUser.Username, "successfully!")
				}
			} else {
				checkErr("Please check the password!")
			}
		} else {
			checkErr("Please check the username!")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&user.Username, "username", "u", "", "The User's Username")
	loginCmd.Flags().StringVarP(&user.Password, "password", "p", "", "The User's Password")
	models.Logger.SetPrefix("[agenda login]")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
