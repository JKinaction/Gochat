package session

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveSession(c *gin.Context, key interface{}, val interface{}) error {
	session := sessions.Default(c)
	session.Set(key, val)
	err := session.Save()
	return err
}

func GetSession(c *gin.Context, key interface{}) interface{} {
	session := sessions.Default(c)
	get := session.Get(key)
	return get
}
func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}
