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

		if len(codes) != 7 {
			return "", ErrInvalidID
		}

		return string(codes), nil
	}

	Decode = func(code string) (int64, error) {
		if len(code) != 7 {
			return 0, ErrInvalidCode
		}

		var num int64

		for i := 0; i < len(code); i++ {
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
	Seq      uint8  `json:"-" bson:"seq"`
	Code     string `json:"code" bson:"-"`
	URL      string `json:"url" bson:"url"`
	ExpireAt int64  `json:"expireAt" bson:"expireAt"`
}

func (u *URL) ToBSON() error {
	u.Seq = getSeq(u)

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

func (u *URL) ToMOodel() error {
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

func getSeq(u *URL) uint8 {
	return uint8(u.ID & 0xFF)
}

func getID(u *URL) (int64, error) {
	id, err := Decode(u.Code)
	if err != nil {
		return 0, err
	}

	return id, nil
}
