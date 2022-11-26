package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/costynus/loyalty-system/internal/entity"
	"github.com/costynus/loyalty-system/internal/usecase"
	"github.com/costynus/loyalty-system/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/theplant/luhn"
)

func NewRouter(handler *chi.Mux, uc usecase.Gophermart, l logger.Interface) {
    tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)


    handler.Use(middleware.Logger)

    // checker
    handler.Get("/healthz", healthzHandler())
    handler.Get("/ping", pingHandler(uc, l))

    // auth
    handler.Group(func(r chi.Router) {
        r.Post("/api/user/register", registrationUser(uc, l, tokenAuth))
        r.Post("/api/user/login", loginUser(uc, l, tokenAuth))
    })

    // Protected routes
    handler.Group(func(r chi.Router) {
        r.Use(jwtauth.Verifier(tokenAuth))
        r.Use(jwtauth.Authenticator)

        r.Post("/api/user/orders", uploadOrder(uc, l, tokenAuth))
        r.Get("/api/user/orders", getOrderInfoList(uc, l, tokenAuth))
        
        r.Get("/api/user/balance", getCurrentBalance(uc, l, tokenAuth))
        r.Post("/api/user/balance/withdraw", withdraw(uc, l, tokenAuth))

        r.Get("/api/user/withdrawals", getWithdrawInfoList(uc, l, tokenAuth))
    })
}

func getWithdrawInfoList(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, claims, _ := jwtauth.FromContext(r.Context())
        withdrawList, err := uc.GetWithdrawList(r.Context(), int(claims["user_id"].(float64)))
        if err != nil {
            errorHandler(w, err)
            return
        }

        jsonResp, err := json.Marshal(withdrawList)
        if err != nil {
            l.Error(err)
            errorHandler(w, err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResp)
    }
}

func withdraw(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var withdrawal entity.Withdrawal

        if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        _, claims, _ := jwtauth.FromContext(r.Context())
        err := uc.Withdraw(r.Context(), int(claims["user_id"].(float64)), withdrawal)
        if err != nil {
            errorHandler(w, err)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}

func getCurrentBalance(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, claims, _ := jwtauth.FromContext(r.Context())
        balance, err := uc.GetCurrentBalance(r.Context(), int(claims["user_id"].(float64)))
        if err != nil {
            errorHandler(w, err)
            return
        }

        jsonResp, err := json.Marshal(balance)
        if err != nil {
            l.Error(err)
            errorHandler(w, err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResp)
    }
}

func getOrderInfoList(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, claims, _ := jwtauth.FromContext(r.Context())
        orderList, err := uc.GetOrderList(r.Context(), int(claims["user_id"].(float64)))
        if err != nil {
            errorHandler(w, err)
            return
        }

        if orderList == nil {
            http.Error(w, "empty order list", http.StatusNoContent)
            return
        }

        jsonResp, err := json.Marshal(orderList)
        if err != nil {
            l.Error(err)
            errorHandler(w, err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResp)
    }
}

func uploadOrder(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        orderNum, _ := strconv.Atoi(string(body))  
        if !luhn.Valid(orderNum){
            http.Error(w, "bad request", http.StatusUnprocessableEntity)
            return
        }

        _, claims, _ := jwtauth.FromContext(r.Context())
        isDouble, err := uc.UploadOrder(r.Context(), int(claims["user_id"].(float64)), strconv.Itoa(orderNum))
        if err != nil {
            errorHandler(w, err)
            return
        }

        if isDouble {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusAccepted)
        }
    }
}


func registrationUser(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var userAuth entity.UserAuth

        if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
            l.Error(err)
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        user, err := uc.CreateNewUser(r.Context(), userAuth)
        if err != nil {
            l.Error(err)
            errorHandler(w, err)
            return
        }

        _, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
        w.Header().Set("Authorization", "Bearer " + tokenString)
        w.WriteHeader(http.StatusOK)
    }
}

func loginUser(uc usecase.Gophermart, l logger.Interface, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        var userAuth entity.UserAuth

        if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        user, err := uc.CheckUser(r.Context(), userAuth)
        if err != nil {
            errorHandler(w, err)
            return
        }

        _, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
        w.Header().Set("Authorization", "Bearer " + tokenString)
        w.WriteHeader(http.StatusOK)
    }
}


func healthzHandler() http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK) }
}

func pingHandler(uc usecase.Gophermart, l logger.Interface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := uc.PingRepo(r.Context()); err != nil {
            http.Error(w, "repo error", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}
