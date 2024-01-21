package models

const (
	base58alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var (
	base58alphabetMap = map[byte]int64{}

	encode = func(num int64) string {
		codes := make([]byte, 0, 7)

		for num > 0 {
			codes = append(codes, base58alphabet[num%58])
			num /= 58
		}

		return string(codes)
	}

	decode = func(code string) int64 {
		var num int64

		for i := 0; i < len(code); i++ {
			num = num*58 + base58alphabetMap[code[i]]
		}

		return num
	}
)

func init() {
	for i := 0; i < len(base58alphabet); i++ {
		base58alphabetMap[base58alphabet[i]] = int64(i)
	}
}

type URL struct {
	ID       int64  `json:"-" bson:"_id"`
	Code     string `json:"code" bson:"-"`
	URL      string `json:"url" bson:"url"`
	ExpireAt int64  `json:"expireAt" bson:"expireAt"`
}

func (u *URL) FillCode() string {
	u.Code = encode(u.ID)
	return u.Code
}

func (u *URL) FillID() int64 {
	u.ID = decode(u.Code)
	return u.ID
}
