package id

import (
    "math/rand"
    "strings"

    "git.fuwafuwa.moe/x3/ngfshare/config"
)

var idChars = []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789")

func genRandStr(length int) string {
    var b strings.Builder
    b.Grow(length)
    for i := 0; i < length; i++ {
        rndIdx := rand.Intn(len(idChars))
        b.WriteRune(idChars[rndIdx])
    }
    return b.String()
}

func GenFileId() string {
    return genRandStr(config.Conf.IdLen)
}

func GenAuthKey() string {
    return genRandStr(config.Conf.AuthKeyLen)
}
