package passman

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"os"
	"path/filepath"
	"time"
)

var (
	loadsDone = promauto.NewCounter(prometheus.CounterOpts{
		Name: "passman_file_loads_total",
		Help: "Number of times passman file was loaded from disk",
	})
)

type Passman struct {
	// Absolute path to the passman file
	fpath    string
	contents []*PassmanEntry
}

type PassmanEntry struct {
	Site      string    `json:site`
	Username  string    `json:username`
	Password  string    `json:password`
	CreatedAt time.Time `json:created_at`
}

func NewPassman(dir string) *Passman {
	absdir, _ := filepath.Abs(dir)
	fpath := filepath.Join(absdir, ".passman")
	return &Passman{fpath: fpath, contents: []*PassmanEntry{}}
}

func (p *Passman) Load() error {
	var contents []*PassmanEntry
	raw, err := os.ReadFile(p.fpath)
	if err != nil {
		return err
	}
	loadsDone.Inc()
	err = json.Unmarshal(raw, &contents)
	if err != nil {
		return err
	}
	p.contents = contents
	return nil
}

func (p *Passman) sync() error {
	contentbytes, err := json.Marshal(p.contents)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p.fpath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(contentbytes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Passman) InitOrLoad() error {
	f, err := os.OpenFile(p.fpath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	defer f.Close()

	// File exists, load it
	if err != nil && os.IsExist(err) {
		return p.Load()
	}

	return p.init(f)
}

func (p *Passman) init(f *os.File) error {
	// Create an empty file with no contents
	var contents []PassmanEntry
	contentbytes, err := json.Marshal(contents)
	if err != nil {
		return err
	}
	_, err = f.Write(contentbytes)
	if err != nil {
		return err
	}

	return nil
}

func (p *Passman) GetAllContents() ([]*PassmanEntry, error) {
	// This might be redundant if we have an empty file that's already been loaded
	if len(p.contents) == 0 {
		p.Load()
	}
	return p.contents, nil
}

func (p *Passman) GetForSite(site string) ([]*PassmanEntry, error) {
	if len(p.contents) == 0 {
		p.Load()
	}

	if len(p.contents) == 0 {
		return nil, nil
	}

	var results []*PassmanEntry
	for _, entry := range p.contents {
		// TODO: Normalize site to strip protocol, etc.
		if entry.Site == site {
			results = append(results, entry)
		}
	}

	// TODO: Make "no such site" error
	return results, nil
}

func (p *Passman) GetForSiteAndUser(site, username string) (*PassmanEntry, error) {
	results, err := p.GetForSite(site)
	if err != nil {
		return nil, err
	}
	for _, entry := range results {
		if entry.Username == username {
			return entry, nil
		}
	}
	return nil, nil
}

func (p *Passman) Create(site, username, password string) error {
	if len(p.contents) == 0 {
		p.Load()
	}

	existing, err := p.GetForSiteAndUser(site, username)
	if err != nil {
		return err
	}
	if existing == nil {
		p.contents = append(p.contents, &PassmanEntry{Site: site, Username: username, Password: password, CreatedAt: time.Now()})
	} else {
		// Overwrite if username exists for site
		existing.Password = password
		existing.CreatedAt = time.Now()
	}
	err = p.sync()
	if err != nil {
		return err
	}
	return nil
}
