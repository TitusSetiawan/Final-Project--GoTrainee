package svr

import (
	"fmt"
	"newfinalproject/db"
	"newfinalproject/entity"
	"time"
)

func UserUpdate(uname string, email string, id int) bool {
	sqlSt := `update users set u_email = $1, u_username = $2, u_updated_date = $3 where u_id = $4`
	_, err := db.Db.Exec(sqlSt,
		&email,
		&uname,
		time.Now(),
		&id,
	)
	if err != nil {
		fmt.Errorf("Error Update User: " + err.Error())
		return false
	}
	return true
}

func UserDelete(uname string) error {
	sqlSt := `delete from users where username = $1`
	_, err := db.Db.Exec(sqlSt, uname)
	if err != nil {
		return err
	}
	return err

}

func UserGetById(id int) {
	var NewUser entity.Users
	sqlSt := `Select id, username, email, age, updated_at from users where id = $1;`
	row, err := db.Db.Query(sqlSt, id)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer row.Close()
	for row.Next() {
		if err = row.Scan(
			&NewUser.Id,
			&NewUser.Username,
			&NewUser.Email,
			&NewUser.Age,
			&NewUser.Update_at,
		); err != nil {
			fmt.Println("No Data", err)
		}
	}
	fmt.Println(NewUser)
}
