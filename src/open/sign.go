package open

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	//"log"
)

//
func Sign(token, timestamp, nonce, appid string) (signature string) {
	// token + timestamp + nonce + appid
	strs := sort.StringSlice{token, timestamp, nonce, appid}
	//strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce)+len(appid))

	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)
	buf = append(buf, strs[3]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}
