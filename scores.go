package emg

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

const (
	obfKey       = "beepboopBADGOOSE"
	obfNonce     = `don'teatcats`
	obfScoreData = `KWFEMuyAM4jev3c7/7J3Zq3phtXkGIZZrmCbsBtUhaSvH4tWK6KhIeShXzv5/mI/TH/9WCqb1mhtomQob8s6OR/wCppuy8Ksb8N2FTW/4bMdD1ft/qqUjj1t8Pv84eYtz1+i+tuq8o2tcgLEkC4`
)

func deobfScores() ([]HighScore, error) {
	c, err := aes.NewCipher([]byte(obfKey))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	b, err := base64.RawStdEncoding.DecodeString(obfScoreData)
	if err != nil {
		return nil, err
	}

	b, err = gcm.Open(b[:0], []byte(obfNonce), b, nil)
	if err != nil {
		return nil, err
	}
	hs := []HighScore{}
	if err := json.Unmarshal(b, &hs); err != nil {
		return nil, err
	}
	return hs, nil
}
