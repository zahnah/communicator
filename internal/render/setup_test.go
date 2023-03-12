package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct {
}

func (m myWriter) Header() http.Header {
	return http.Header{}
}

func (m myWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (m myWriter) WriteHeader(statusCode int) {

}
