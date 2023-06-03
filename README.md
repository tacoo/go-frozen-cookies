# go-frozen-cookies

This implementation enhances the functionality of net/http/cookiejar by providing the ability to store cookies, including session cookies, in a JSON file.

```
import (
    cookiejar "github.com/tacoo/go-frozen-cookies"
)

jar, _ := cookiejar.New(&cookiejar.Options{FilePath: "/path/to/cookie.json"})

jar.Save()
```