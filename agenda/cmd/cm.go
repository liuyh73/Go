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
	"strings"
	"time"

	"github.com/liuyh73/Go/agenda/models"
	"github.com/spf13/cobra"
)

var meeting models.Meeting
var timeLayout = "2006-01-02 15:04:05"

// cmCmd represents the cm command
var cmCmd = &cobra.Command{
	Use:   "cm",
	Short: "This command can create one meeting",
	Long:  `You can use agenda login to create one meeting and manage it`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime, _ := cmd.Flags().GetString("startTime")
		endTime, _ := cmd.Flags().GetString("endTime")

		// 查询title是否重复
		usermeeting := new(models.UserMeeting)
		usermeeting.Title = meeting.Title
		has, _ := models.Engine.Get(usermeeting)
		if has {
			checkErr("The meeting" + meeting.Title + "has been initiated!")
		}
		// 获取发起者参与的所有会议，进行时间上的判断
		// TO-DO

		// 获取参与者参与的所有会议，进行时间上的判断
		// TO-DO
		for _, participator := range strings.Split(meeting.Participants, ",") {
			// 获取参与者参与的会议title
			participate_meetings := &models.UserMeeting{Participator: participator}
			models.Engine.Get(participate_meetings)
			// 根据title查询会议
			meeting := &models.Meeting{Title: participate_meeting.Title}
			models.Engine.Get(meeting)

		}
		var err1, err2 error
		loc, _ := time.LoadLocation("Local") //获取时区
		meeting.StartTime, err1 = time.ParseInLocation(timeLayout, startTime, loc)
		meeting.EndTime, err2 = time.ParseInLocation(timeLayout, endTime, loc)
		if err1 != nil || err2 != nil {
			checkErr("The StartTime or EndTime must conform to the following format: " + timeLayout)
		}
		_, err := models.Engine.Insert(&meeting)
		if err != nil {
			checkErr("Fail to insert the meeting " + meeting.Title + " to database")
		}
	},
}

func init() {
	rootCmd.AddCommand(cmCmd)
	cmCmd.Flags().StringVarP(&meeting.Title, "title", "t", "", "The Meeting's Title")
	cmCmd.Flags().StringVarP(&meeting.Moderator, "Moderator", "m", "", "The Meeting's Moderator")
	cmCmd.Flags().StringP("startTime", "s", "", "The Meeting's startTime")
	cmCmd.Flags().StringP("endTime", "e", "", "The Meeting's endTime")
	cmCmd.Flags().StringVarP(&meeting.Participants, "participants", "p", "", "The Meeting's Participants")
	models.Logger.SetPrefix("[agenda cm]")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
