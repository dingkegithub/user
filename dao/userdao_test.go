package dao

import "testing"

// ToDo: mysql 初始化公共

func TestUserDaoImpl_SelectByEmail(t *testing.T) {

	err := MysqlInit("127.0.0.1", 6379, "root", "root", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userDao := &UserDaoImpl{}
	user, err := userDao.SelectByEmail("dk_mcu@163.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("username is: %s\n", user.UserName)
}

func TestUserDaoImpl_Create(t *testing.T) {
	err := MysqlInit("127.0.0.1", 6379, "root", "root", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userDao := &UserDaoImpl{}
	user := &UserEntity{
		UserName: "dk_mcu",
		Email:    "dk_mcu@163.com",
		Password: "123456",
	}
	err = userDao.Save(user)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is: ", user.ID)
}
