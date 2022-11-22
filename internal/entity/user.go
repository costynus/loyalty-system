package entity


type (
    UserAuth struct {
        login string `json:"login"`
        password string `json:"password"`
    }

    User struct {
        ID int `json:"id" db:"id"`
        login string `json:"login" db:"login"`
        password_hash string `db:"password_hash"`
    }
)
