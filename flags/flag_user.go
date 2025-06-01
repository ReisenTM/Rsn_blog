package flags

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/pwd"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type FlagUser struct {
}

// Create 创建用户
func (FlagUser) Create() {
	var role enum.RoleType
	fmt.Println("选择角色     1 普通用户  2  超级管理员    3 访客")
	_, err := fmt.Scan(&role)
	if err != nil {
		logrus.Errorf("角色输入错误 %s", err)
		return
	}
	if !(role == 1 || role == 2 || role == 3) {
		logrus.Errorf("输入角色错误")
		return
	}

	var username string
	fmt.Println("请输入用户名:")
	_, err = fmt.Scan(&username)
	if err != nil {
		logrus.Errorf("用户名输入错误 %s", err)
		return
	}

	// 查用户名是否存在
	var m model.UserModel
	err = global.DB.Take(&m, "username = ?", username).Error
	if err == nil {
		logrus.Errorf("此用户名已存在")
		return
	}
	fmt.Println("请输入密码:")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("读取密码时出错:", err)
		return
	}
	fmt.Println("请再次输入密码:")
	rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("读取密码时出错:", err)
		return
	}
	if string(password) != string(rePassword) {
		fmt.Println("两次密码不一致")
		return
	}

	hashPwd, _ := pwd.GenerateFromPassword(string(password))
	// 创建用户
	err = global.DB.Create(&model.UserModel{
		Username:  username,
		Nickname:  "用户001",
		RegSource: enum.RegisterTerminalSourceType,
		Password:  hashPwd,
		Role:      role,
	}).Error
	if err != nil {
		fmt.Println("创建用户失败", err)
		return
	}
	logrus.Infof("创建用户成功")
}
