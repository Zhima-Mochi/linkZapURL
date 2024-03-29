package models

import "errors"

const (
	base58alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var (
	ErrInvalidID = errors.New("models: invalid id")

	ErrInvalidCode = errors.New("models: invalid code")
)

var (
	base58alphabetMap = map[byte]int64{}

	Encode = func(num int64) (string, error) {
		codes := make([]byte, 0, 7)

		for num > 0 {
			codes = append(codes, base58alphabet[num%58])
			num /= 58
		}

		if len(codes) > 7 {
			return "", ErrInvalidID
		}

		for len(codes) < 7 {
			codes = append(codes, base58alphabet[0])
		}

		return string(codes), nil
	}

	Decode = func(code string) (int64, error) {
		if len(code) != 7 {
			return 0, ErrInvalidCode
		}

		var num int64

		for i := len(code) - 1; i >= 0; i-- {
			if n, ok := base58alphabetMap[code[i]]; !ok {
				return 0, ErrInvalidCode
			} else {
				num = num*58 + n
			}
		}

		return num, nil
	}
)

func init() {
	for i := 0; i < len(base58alphabet); i++ {
		base58alphabetMap[base58alphabet[i]] = int64(i)
	}
}

type URL struct {
	ID       int64  `json:"-" bson:"ID"`
	ShardID  int64  `json:"-" bson:"shardID"`
	Code     string `json:"code" bson:"-"`
	URL      string `json:"url" bson:"url"`
	ExpireAt int64  `json:"expireAt" bson:"expireAt"`
}

func (u *URL) ToBSON() error {
	u.ShardID = getShardID(u)

	return nil
}

func (u *URL) ToJSON() error {
	code, err := getCode(u)
	if err != nil {
		return err
	}

	u.Code = code

	return nil
}

func (u *URL) ToModel() error {
	id, err := getID(u)
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func getCode(u *URL) (string, error) {
	code, err := Encode(u.ID)
	if err != nil {
		return "", err
	}

	return code, nil
}

func getShardID(u *URL) int64 {
	return (u.ID >> 24) % 1000
}

func getID(u *URL) (int64, error) {
	id, err := Decode(u.Code)
	if err != nil {
		return 0, err
	}

	return id, nil
}
