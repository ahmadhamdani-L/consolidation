package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"mcash-finance-console-core/configs"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/redis"
	res "mcash-finance-console-core/pkg/util/response"
	"regexp"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	jwtKey := configs.Jwt().SecretKey()

	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, nil).Send(c)
		}

		splitToken := strings.Split(authToken, "Bearer ")
		token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
			}

			return []byte(jwtKey), nil
		})

		if !token.Valid || err != nil {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
		}
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		timeUntilExpiration := expirationTime.Sub(time.Now())

		if timeUntilExpiration <= 3*time.Minute {
			waktu := time.Now()
			newToken := jwt.New(jwt.SigningMethodHS256)
			newClaims := newToken.Claims.(jwt.MapClaims)
			newClaims["id"] = claims["id"]
			newClaims["rid"] = claims["rid"]
			newClaims["exp"] = time.Now().Add(5 * time.Minute).Unix() // Atur waktu kedaluwarsa token yang baru

			newTokenString, err := newToken.SignedString([]byte(jwtKey))
			if err != nil {
				return res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err).Send(c)
			}
			// token, err := data.GenerateToken()
			var id int
			destructID := claims["id"]
			if destructID != nil {
				id = int(destructID.(float64))
			} else {
				id = 0
			}
			msg := kafka.JsonData{
				Data:  newTokenString,
				UserID: id,
				Name: "token",
				Timestamp: &waktu,
			}
			jsonStr, err := json.Marshal(&msg)
			if err != nil {
				return nil
			}
			go kafka.NewService("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))	
		}

		var id int
		destructID := token.Claims.(jwt.MapClaims)["id"]
		if destructID != nil {
			id = int(destructID.(float64))
		} else {
			id = 0
		}

		// var name string
		// destructName := token.Claims.(jwt.MapClaims)["name"]
		// if destructName != nil {
		// 	name = destructName.(string)
		// } else {
		// 	name = ""
		// }

		// var username string
		// destructUsername := token.Claims.(jwt.MapClaims)["username"]
		// if destructUsername != nil {
		// 	username = destructUsername.(string)
		// } else {
		// 	username = ""
		// }

		// var companyID int
		// destructCompanyID := token.Claims.(jwt.MapClaims)["cid"]
		// if destructCompanyID != nil {
		// 	companyID = int(destructCompanyID.(float64))
		// } else {
		// 	companyID = 0
		// }

		var roleID int
		destructRoleID := token.Claims.(jwt.MapClaims)["rid"]
		if destructRoleID != nil {
			roleID = int(destructRoleID.(float64))
		} else {
			roleID = 0
		}

		cc := c.(*abstraction.Context)
		cc.Auth = &abstraction.AuthContext{
			ID: id,
			// Name:      name,
			// Username:  username,
			// CompanyID: companyID,
			RoleID: roleID,
		}
		//checkrole
		var re = regexp.MustCompile(`\/\d+`)
		method := c.Request().Method
		path := re.ReplaceAllString(c.Request().URL.Path, ``)
		check := fmt.Sprintf("role=%d:allow=%s_%s*", roleID, method, path)
		anyData, err := redis.RedisClient.Keys(c.Request().Context(), check).Result()
		if redis.IsNil(err) || len(anyData) == 0 {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, errors.New("Don't have permission to acces this menu")).Send(c)
		} else if err != nil {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, errors.New("Don't have permission to acces this menu")).Send(c)
		}

		return next(cc)
	}
}

func AuthenticationResetPassword(next echo.HandlerFunc) echo.HandlerFunc {
	jwtKey := configs.Jwt().SecretKey()

	return func(c echo.Context) error {
		authToken := c.Param("resetToken")
		if authToken == "" {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, nil).Send(c)
		}

		splitToken := (authToken)
		token, err := jwt.Parse(splitToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
			}

			return []byte(jwtKey), nil
		})

		if !token.Valid || err != nil {
			return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
		}

		var id int
		destructID := token.Claims.(jwt.MapClaims)["id"]
		if destructID != nil {
			id = int(destructID.(float64))
		} else {
			id = 0
		}

		var name string
		destructName := token.Claims.(jwt.MapClaims)["name"]
		if destructName != nil {
			name = destructName.(string)
		} else {
			name = ""
		}

		var email string
		destructEmail := token.Claims.(jwt.MapClaims)["email"]
		if destructEmail != nil {
			email = destructEmail.(string)
		} else {
			email = ""
		}

		var company_id int
		destructCompanyID := token.Claims.(jwt.MapClaims)["company_id"]
		if destructCompanyID != nil {
			company_id = int(destructCompanyID.(float64))
		} else {
			company_id = 0
		}

		cc := c.(*abstraction.Context)
		cc.Auth = &abstraction.AuthContext{
			ID:        id,
			Name:      name,
			Email:     email,
			CompanyID: company_id,
		}

		return next(cc)
	}
}
