package cookiejar

import (
	"encoding/json"
	"os"
	"time"
)

type fileEntry struct {
	Name       string    `json:"name"`
	Value      string    `json:"value"`
	Domain     string    `json:"domain"`
	Path       string    `json:"path"`
	SameSite   string    `json:"sameSite"`
	Secure     bool      `json:"secure"`
	HttpOnly   bool      `json:"httpOnly"`
	Persistent bool      `json:"persistent"`
	HostOnly   bool      `json:"hostOnly"`
	Expires    time.Time `json:"expires"`
	Creation   time.Time `json:"creation"`
	LastAccess time.Time `json:"lastAccess"`
	SeqNum     uint64    `json:"seqNum"`
}

func loadCookiesFromFile(path string) (nextSeqNum uint64, cookies map[string]map[string]entry, err error) {
	now := time.Now()
	cookies = make(map[string]map[string]entry)
	_, fileNotExistsErr := os.Stat(path)
	if fileNotExistsErr != nil {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	filecookies := make(map[string]map[string]fileEntry)
	err = json.NewDecoder(f).Decode(&filecookies)
	if err != nil {
		return
	}
	for k0 := range filecookies {
		for k1 := range filecookies[k0] {
			c := filecookies[k0][k1]
			if !c.Expires.After(now) {
				continue
			}
			if cookies[k0] == nil {
				cookies[k0] = make(map[string]entry)
			}
			cookies[k0][k1] = entry{
				Name:       c.Name,
				Value:      c.Value,
				Domain:     c.Domain,
				Path:       c.Path,
				SameSite:   c.SameSite,
				Secure:     c.Secure,
				HttpOnly:   c.HttpOnly,
				Persistent: c.Persistent,
				HostOnly:   c.HostOnly,
				Expires:    c.Expires,
				Creation:   c.Creation,
				LastAccess: c.LastAccess,
				seqNum:     c.SeqNum,
			}
			if c.SeqNum >= nextSeqNum {
				nextSeqNum = c.SeqNum + 1
			}
		}
	}
	return
}

func (j *Jar) Save() (err error) {
	if j.filePath == "" {
		return
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	err = saveCookiesToFile(j.filePath, j.entries)
	return
}

func saveCookiesToFile(path string, cookies map[string]map[string]entry) (err error) {
	now := time.Now()
	filecookies := make(map[string]map[string]fileEntry)
	for k0 := range cookies {
		for k1 := range cookies[k0] {
			c := cookies[k0][k1]
			if !c.Expires.After(now) {
				continue
			}
			if filecookies[k0] == nil {
				filecookies[k0] = make(map[string]fileEntry)
			}
			filecookies[k0][k1] = fileEntry{
				Name:       c.Name,
				Value:      c.Value,
				Domain:     c.Domain,
				Path:       c.Path,
				SameSite:   c.SameSite,
				Secure:     c.Secure,
				HttpOnly:   c.HttpOnly,
				Persistent: c.Persistent,
				HostOnly:   c.HostOnly,
				Expires:    c.Expires,
				Creation:   c.Creation,
				LastAccess: c.LastAccess,
				SeqNum:     c.seqNum,
			}
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.Encode(filecookies)
	return
}
