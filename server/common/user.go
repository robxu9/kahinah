package common

// user.go contains structures that link Kahinah's base users with
// api keys, persona session keys, etc.

type UserToken struct {
	Id int64

	UserId  int64
	Token   []byte
	Comment string
	Revoked bool
}

func GenerateUserToken(k *kahinah.Kahinah, user int64, comment string) string {

}
