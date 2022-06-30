package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"newfinalproject/db"
	"newfinalproject/entity"
	"newfinalproject/svr"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func PostRegister(w http.ResponseWriter, r *http.Request) {
	var NewUser entity.Users
	var NewResponseUser entity.UsersResult
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&NewUser)
	if err != nil {

		w.Write([]byte(err.Error()))
	}

	/////////// Validation Email, Uname, Password ///////////////
	_, checkMail := svr.ValidationUsersEmail(NewUser.Email)
	checkPass := svr.ValidationUsersPass(NewUser.Password)
	checkUname := svr.ValidationUsersUname(NewUser.Username)
	checkAge := svr.ValidationUsersAge(NewUser.Age)

	log.Println(NewUser) // for debug
	hashPass, errHash := bcrypt.GenerateFromPassword([]byte(NewUser.Password), 12)
	if errHash != nil {
		log.Fatal(errHash)
	}
	NewUser.Password = string(hashPass)

	// INSERT to DB
	if checkMail == checkPass == checkUname == checkAge {
		sqlSt := `
		insert into users (
			username,
			email,
			password,
			age,
			created_at)
			values($1,$2,$3,$4,now())
		returning age, email, id, username;
		`
		err = db.Db.QueryRow(sqlSt, NewUser.Username, NewUser.Email, NewUser.Password, NewUser.Age).
			Scan(&NewResponseUser.Age, &NewResponseUser.Email, &NewResponseUser.Id, &NewResponseUser.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
		}
	} else {
		fmt.Println("Check your Input!")
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(NewResponseUser)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var NewUserPostLogin entity.UsersPostLogin
	var result entity.UsersPostLogin1
	var NewUserToken entity.UserToken
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&NewUserPostLogin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	_, checkMail := svr.ValidationUsersEmail(NewUserPostLogin.Email)
	checkPass := svr.ValidationUsersPass(NewUserPostLogin.Password)
	if checkMail == checkPass {
		sqlState := "select email, password, username from users where email=$1"
		rows, err := db.Db.Query(sqlState, NewUserPostLogin.Email)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(
				&result.Email,
				&result.Password,
				&result.Username,
			); err != nil {
				fmt.Println("No Data", err)
			}
		}

		///// check email password if exist /////
		if NewUserPostLogin.Email == result.Email {
			err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(NewUserPostLogin.Password))
			if err == nil {
				strResult, flag := svr.SignIn(result.Username)
				if flag {
					w.WriteHeader(http.StatusOK)
					// NewUserToken.Username = result.Username
					NewUserToken.Token = strResult
					json.NewEncoder(w).Encode(NewUserToken)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Password wrong!"))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Email doesnt match"))
		}

	}
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	var UserUpdate entity.UserUpdateIn
	var NewUser entity.Users
	var NewUserRes entity.UserUpdateRes
	if r.Method == "PUT" {

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&UserUpdate)
		if err != nil {

			w.Write([]byte(err.Error()))
		}
		log.Println(UserUpdate)
		id := r.URL.Query().Get("userid")
		idInt, _ := strconv.Atoi(id)
		fmt.Println(idInt)
		if idInt != 0 {
			_, checkMail := svr.ValidationUsersEmail(UserUpdate.Email)
			checkUname := svr.ValidationUsersPass(UserUpdate.Username)
			if checkMail == checkUname {
				var JwtKey = []byte(entity.JwtKey)
				authHeader := r.Header.Get("Authorization")
				if !strings.Contains(authHeader, "Bearer") {
					http.Error(w, "invalid token", http.StatusBadRequest)
					return
				}
				tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
				// fmt.Println("ini token string:", tokenString)
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Signing method invalid")
					}
					return JwtKey, nil
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok || !token.Valid {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				////////// Select /////////
				fmt.Println(claims["exp"])
				UserUpdate.Username = claims["username"].(string)
				current_time := time.Now()
				sqlSt := `update users set email = $1, username = $2, updated_at = $3 where id = $4`
				_, err = db.Db.Exec(sqlSt,
					&UserUpdate.Email,
					&UserUpdate.Username,
					current_time.Format("2006-01-02 15:04:05"),
					&id,
				)
				if err != nil {
					fmt.Errorf("Error Update User: " + err.Error())
				}

				sqlSt = `Select id, username, email, age, updated_at from users where id = $1;`
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
				if svr.UserUpdate(UserUpdate.Username, UserUpdate.Email, idInt) {
					NewUserRes.Id = idInt
					NewUserRes.Username = UserUpdate.Username
					NewUserRes.Email = UserUpdate.Email
					NewUserRes.UserUpdateAt = NewUser.Update_at
					NewUserRes.Age = NewUser.Age
					jsonData, _ := json.Marshal(NewUserRes)
					w.Header().Add("Content-Type", "application/json")
					w.WriteHeader(200)
					w.Write(jsonData)
				}
			}
		}
	}
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	var NewDelete entity.DeleteData
	if r.Method == "DELETE" {
		var JwtKey = []byte(entity.JwtKey)
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
		fmt.Println("ini token string:", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			}
			return JwtKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		usernameClaim := claims["username"].(string)
		if svr.UserDelete(usernameClaim) != nil {
			fmt.Println("Can't deleted")
		}
		NewDelete.Message = "Your account has been successfully deleted"
		prettyJSON, err := json.MarshalIndent(NewDelete, "", "  ")
		if err != nil {
			log.Fatal("Failed to generate json", err)
		}
		w.Write([]byte(prettyJSON))

	}
}
