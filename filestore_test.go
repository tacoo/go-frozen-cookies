package cookiejar

import (
	"net/http"
	urlutil "net/url"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	filepath := "./test_save.log"
	jarBefore, err := New(&Options{FilePath: filepath})
	if err != nil {
		t.Fatal("new jar error", err)
	}
	urlobj, _ := urlutil.Parse("https://example.com/a/b/c")

	addCookies(urlobj, jarBefore, "ASESSION=aaaaa; Path=/; HttpOnly; Secure")
	addCookies(urlobj, jarBefore, "BSESSION=bbbbb; Path=/; Domain=example.com; HttpOnly; Secure; SameSite=lax")
	addCookies(urlobj, jarBefore, "CCC=ccccc; Path=/; Domain=example.com; Max-Age=10800; Secure")
	addCookies(urlobj, jarBefore, "DDD=ddddd; Path=/; Domain=example.com; Secure")
	addCookies(urlobj, jarBefore, "EEE=eeeee-=; Path=/; Expires=Sun, 10 Jan 2000 00:00:00 GMT; Max-Age=461005548; HttpOnly; Secure")
	jarBefore.Save()

	jarAfter, err := New(&Options{FilePath: filepath})
	if err != nil {
		t.Fatal("load jar error", err)
	}

	compareCookies(t, jarBefore.Cookies(urlobj), jarAfter.Cookies(urlobj))
}

func TestSaveExpired(t *testing.T) {
	filepath := "./test_save_expired.log"
	jarBefore, err := New(&Options{FilePath: filepath})
	if err != nil {
		t.Fatal("new jar error", err)
	}
	urlobj, _ := urlutil.Parse("https://example.com/a/b/c")

	addCookies(urlobj, jarBefore, "FFF=fffff; Path=/; Domain=example.com; Max-Age=2; Secure")
	cookieBefore := jarBefore.Cookies(urlobj)
	jarBefore.Save()

	<-time.After(3 * time.Second)

	jarAfter, err := New(&Options{FilePath: filepath})
	if err != nil {
		t.Fatal("load jar error", err)
	}

	cookieBefore8SecAfter := jarBefore.Cookies(urlobj)
	cookieAfter := jarAfter.Cookies(urlobj)

	if len(cookieBefore) != 1 {
		t.Fatalf("[cookieBefore] cookie not found")
	}
	if len(cookieBefore8SecAfter) != 0 {
		t.Fatalf("[cookieBefore8SecAfter] cookie found: %#v", cookieBefore8SecAfter)
	}
	if len(cookieAfter) != 0 {
		t.Fatalf("[cookieAfter] cookie found: %#v", cookieAfter)
	}
}

func compareCookies(t *testing.T, expectedCookies, actualCookies []*http.Cookie) {
	for i := range expectedCookies {
		expected := expectedCookies[i]
		var actual *http.Cookie
		for j := range actualCookies {
			if actualCookies[j].Name == expected.Name {
				actual = actualCookies[j]
				break
			}
		}
		if actual == nil {
			t.Fatalf("cookie not found name=%s", expected.Name)
		}
		if expected.Name != actual.Name {
			t.Fatalf("cookie name=%s Name is invalid actual=%v expected=%v", expected.Name, actual.Name, expected.Name)
		}
		if expected.Value != actual.Value {
			t.Fatalf("cookie name=%s Value is invalid actual=%v expected=%v", expected.Name, actual.Value, expected.Value)
		}
		if expected.Path != actual.Path {
			t.Fatalf("cookie name=%s Path is invalid actual=%v expected=%v", expected.Name, actual.Path, expected.Path)
		}
		if expected.Domain != actual.Domain {
			t.Fatalf("cookie name=%s Domain is invalid actual=%v expected=%v", expected.Name, actual.Domain, expected.Domain)
		}
		if expected.Expires != actual.Expires {
			t.Fatalf("cookie name=%s Expires is invalid actual=%v expected=%v", expected.Name, actual.Expires, expected.Expires)
		}
		if expected.MaxAge != actual.MaxAge {
			t.Fatalf("cookie name=%s MaxAge is invalid actual=%v expected=%v", expected.Name, actual.MaxAge, expected.MaxAge)
		}
		if expected.Secure != actual.Secure {
			t.Fatalf("cookie name=%s Secure is invalid actual=%v expected=%v", expected.Name, actual.Secure, expected.Secure)
		}
		if expected.HttpOnly != actual.HttpOnly {
			t.Fatalf("cookie name=%s HttpOnly is invalid actual=%v expected=%v", expected.Name, actual.HttpOnly, expected.HttpOnly)
		}
		if expected.SameSite != actual.SameSite {
			t.Fatalf("cookie name=%s SameSite is invalid actual=%v expected=%v", expected.Name, actual.SameSite, expected.SameSite)
		}
	}
	for i := range actualCookies {
		actual := actualCookies[i]
		var expected *http.Cookie
		for j := range expectedCookies {
			if expectedCookies[j].Name == actual.Name {
				expected = expectedCookies[j]
				break
			}
		}
		if expected == nil {
			t.Fatalf("unknown cookie found name=%s", actual.Name)
		}
	}
}

func addCookies(url *urlutil.URL, jar *Jar, c string) {
	header := http.Header{}
	header.Add("Set-Cookie", c)
	r := http.Response{Header: header}
	jar.SetCookies(url, r.Cookies())
}
