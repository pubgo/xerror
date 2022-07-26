package cryptoutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	ErrSignExpired   = errors.New("sign expired")
	ErrSignInvalid   = errors.New("sign invalid")
	ErrExpireInvalid = errors.New("expire invalid")
	ErrExpireMissing = errors.New("expire missing")
)

type HMACSign struct {
	SecretKey []byte
}

func (s HMACSign) Sign(data string, expire int64) string {
	h := hmac.New(sha256.New, s.SecretKey)
	expireTimeStamp := strconv.FormatInt(expire, 10)
	_, err := io.WriteString(h, data+":"+expireTimeStamp)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)) + ":" + expireTimeStamp
}

func (s HMACSign) Verify(data, sign string) error {
	signSlice := strings.Split(sign, ":")
	// check whether contains expire time
	if signSlice[len(signSlice)-1] == "" {
		return ErrExpireMissing
	}
	// check whether expire time is expired
	expires, err := strconv.ParseInt(signSlice[len(signSlice)-1], 10, 64)
	if err != nil {
		return ErrExpireInvalid
	}
	// if expire time is expired, return error
	if expires < time.Now().Unix() && expires != 0 {
		return ErrSignExpired
	}
	// verify sign
	if s.Sign(data, expires) != sign {
		return ErrSignInvalid
	}
	return nil
}
