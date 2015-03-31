package common

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/robxu9/kahinah/kahinah"
)

// user.go contains structures that link Kahinah's base users with
// api keys, persona session keys, etc.

var (
	ErrNoSuchToken  = errors.New("kahinah: no such token")
	ErrBadToken     = errors.New("kahinah: bad token")
	ErrExpiredToken = errors.New("kahinah: token has expired")

	CacheUserTokenExp = 12 * time.Hour
)

type UserToken struct {
	Id int64

	// Used for validation
	UserId     int64  // user
	Comment    string // comment
	IssueTime  int64  // issuetime, utc, unix
	ExpiryTime int64  // expirytime, utc, unix

	Token   string
	Revoked bool
}

type UserHandler func(http.ResponseWriter, *http.Request, *UserToken)

func (c *Common) UserWrapHandler(u UserHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// check two places for the authentication token - query parameter "token=<token>" & header "Authorization: Bearer <token>"
		token := r.URL.Query().Get("token")
		if token == "" {
			// try via header
			if header := r.Header.Get("Authorization"); header != "" {
				split := strings.Split(header, " ")
				if len(split) == 2 && strings.ToLower(split[0]) == "Bearer" { // based on auth0/go-jwt-middleware
					token = split[1]
				}
			}
		}

		if token == "" { // if token is __still__ nothing, then we have no user
			u(rw, r, nil)
			return
		}

		// validate token. if okay, keep going. if bad, bail out.
		parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(c.C.SecretKey), nil
		})

		if err == nil && parsed.Valid {
			// valid token, look it up
			id, err := c.FindUserToken(token)
			if err != nil { // valid token, but not in database
				log.Printf("[warn] got a token from %s that's valid but not in db!", r.RemoteAddr)
				http.Error(rw, ErrBadToken.Error(), http.StatusUnauthorized)
				return
			}

			userToken, err := c.RetrieveUserToken(id)
			if err != nil {
				panic(err)
			}

			// FIXME could the user be deleted?

			if userToken.Revoked {
				// revoked. bad token.
				http.Error(rw, ErrBadToken.Error(), http.StatusUnauthorized)
				return
			}

			if userToken.ExpiryTime != 0 && time.Unix(userToken.ExpiryTime, 0).Before(time.Now()) {
				// already expired, do _not_ accept
				http.Error(rw, ErrExpiredToken.Error(), http.StatusUnauthorized)
				return
			}

			u(rw, r, userToken)
			return
		}

		// bad token
		http.Error(rw, ErrBadToken.Error(), http.StatusUnauthorized)
	}
}

func usertokenCacheId(id int64) string {
	return "usertoken/" + strconv.FormatInt(id, 10)
}

func (c *Common) RetrieveUserToken(id int64) (*UserToken, error) {
	// get user token cache id
	if cached, found := c.cache.Get(usertokenCacheId(id)); found {
		record := cached.(*UserToken)
		return record, nil
	}

	record := &UserToken{}

	if c.K.DB().First(record, id).RecordNotFound() {
		return nil, ErrNoSuchToken
	}

	// store in the cache
	c.cache.Set(usertokenCacheId(record.Id), record, CacheUserTokenExp)

	return record, nil
}

func (c *Common) FindUserToken(token string) (int64, error) {
	var id []int64

	// FIXME WHY DOES GORM NOT KNOW THE TABLE NAME?
	if err := c.K.DB().Model(&UserToken{}).Where(&UserToken{Token: token}).Limit(1).Pluck("id", &id).Error; err != nil {
		return 0, err
	}

	if len(id) == 0 {
		return 0, ErrNoSuchToken
	}

	return id[0], nil
}

func (c *Common) ListUserTokens(user int64) ([]int64, error) {
	// verify the user exists
	if _, err := c.K.RetrieveUser(user); err != nil {
		return nil, kahinah.ErrNoSuchUser
	}

	// lookup in database
	var records []int64

	if err := c.K.DB().Model(&UserToken{}).Where(&UserToken{UserId: user}).Order("created_at desc").Pluck("id", &records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (c *Common) GenerateUserToken(user int64, comment string, expires bool) (int64, string, error) {
	// verify the user exists
	if _, err := c.K.RetrieveUser(user); err != nil {
		return 0, "", kahinah.ErrNoSuchUser
	}

	jwtToken := jwt.New(jwt.SigningMethodHS512)
	jwtToken.Claims["user"] = user
	jwtToken.Claims["issue"] = time.Now().UTC().Unix()
	jwtToken.Claims["expiry"] = 0

	if expires {
		jwtToken.Claims["expiry"] = time.Now().Add(336 * time.Hour).UTC().Unix() // two weeks
	}

	generatedToken, err := jwtToken.SignedString(c.C.SecretKey)
	if err != nil {
		return 0, "", err
	}

	userToken := &UserToken{
		UserId:     user,
		Comment:    comment,
		IssueTime:  jwtToken.Claims["issue"].(int64),
		ExpiryTime: jwtToken.Claims["expiry"].(int64),
		Token:      generatedToken,
	}

	// insert into db
	if err := c.K.DB().Save(userToken).Error; err != nil {
		return 0, "", err
	}

	return userToken.Id, generatedToken, nil
}
