package entity


type (
    UserAuth struct {
        Login string `json:"login"`
        Password string `json:"password"`
    }

    User struct {
        ID int `json:"id" db:"id"`
        Login string `json:"login" db:"login"`
        Password_hash string `db:"password_hash"`
    }
)
