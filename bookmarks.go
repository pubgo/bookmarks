package bookmarks

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pubgo/assert"
	"io"
	"strings"
)

type _Dir struct {
	Title                 string               `json:"title"`
	AddDate               string               `json:"add_date"`
	Modified              string               `json:"modified"`
	PERSONALTOOLBARFOLDER string               `json:"personal_toolbar_folder"`
	Dirs                  map[string]*_Dir     `json:"dirs"`
	Links                 map[string]*Bookmark `json:"links"`
}

type Bookmark struct {
	URL      string   `json:"url"`
	Title    string   `json:"title"`
	IconUri  string   `json:"icon_uri"`
	Icon     string   `json:"icon"`
	AddDate  string   `json:"add_date"`
	Modified string   `json:"modified"`
	Category []string `json:"category"`
	Tags     string   `json:"tags"`
}

type Bookmarks struct {
	bks *_Dir
}

func (t *Bookmarks) h3(dir *_Dir) string {
	var _ls []string
	for _, _l := range dir.Links {
		_ls = append(_ls, t.link(_l))
	}

	var _dls []string
	for _, _l := range dir.Dirs {
		_dls = append(_dls, t.h3(_l))
	}

	var _params []string
	if dir.AddDate != "" {
		_params = append(_params, fmt.Sprintf(`ADD_DATE="%s"`, dir.AddDate))
	}
	if dir.Modified != "" {
		_params = append(_params, fmt.Sprintf(`LAST_MODIFIED="%s"`, dir.Modified))
	}
	if dir.PERSONALTOOLBARFOLDER != "" {
		_params = append(_params, fmt.Sprintf(`PERSONAL_TOOLBAR_FOLDER="%s"`, dir.PERSONALTOOLBARFOLDER))
	}

	return fmt.Sprintf(_d, strings.Join(_params, " "), dir.Title, strings.Join(_dls, "\n"), strings.Join(_ls, "\n"))
}

func (t *Bookmarks) link(bk *Bookmark) string {
	var _params []string
	if bk.Modified != "" {
		_params = append(_params, fmt.Sprintf(`LAST_MODIFIED="%s"`, bk.Modified))
	}

	if bk.URL != "" {
		_params = append(_params, fmt.Sprintf(`HREF="%s"`, bk.URL))
	}

	if bk.AddDate != "" {
		_params = append(_params, fmt.Sprintf(`ADD_DATE="%s"`, bk.AddDate))
	}

	if bk.Icon != "" {
		_params = append(_params, fmt.Sprintf(`ICON="%s"`, bk.Icon))
	}

	if bk.Tags != "" {
		_params = append(_params, fmt.Sprintf(`TAGS="%s"`, bk.Tags))
	}

	if bk.IconUri != "" {
		_params = append(_params, fmt.Sprintf(`ICON_URI="%s"`, bk.IconUri))
	}
	return fmt.Sprintf(`<DT><A %s>%s</A>`, strings.Join(_params, " "), bk.Title)
}

func (t *Bookmarks) Import(r io.Reader) {
	assert.Bool(r == nil, "")

	doc, err := goquery.NewDocumentFromReader(r)
	assert.Err(err, "parse html error")

	doc.Find("dt>a").Each(func(_ int, a *goquery.Selection) {

		_bk := &Bookmark{
			Title:    strings.TrimSpace(a.Text()),
			URL:      strings.TrimSpace(a.AttrOr("href", "")),
			Tags:     strings.TrimSpace(a.AttrOr("tags", "")),
			Modified: strings.TrimSpace(a.AttrOr("last_modified", "")),
			AddDate:  strings.TrimSpace(a.AttrOr("add_date", "")),
			IconUri:  strings.TrimSpace(a.AttrOr("icon_uri", "")),
			Icon:     strings.TrimSpace(a.AttrOr("icon", "")),
		}

		var _ct = make(map[string]*_Dir)
		for {
			if a.Is("dl") {
				if a.Prev().Is("h1") {
					break
				}

				_txt := a.Prev().Text()
				_txt = strings.ToLower(_txt)
				_txt = strings.Replace(_txt, " ", "-", -1)
				_bk.Category = append(_bk.Category, _txt)

				_ct[_txt] = &_Dir{
					Title:                 _txt,
					AddDate:               strings.TrimSpace(a.Prev().AttrOr("add_date", "")),
					Modified:              strings.TrimSpace(a.Prev().AttrOr("last_modified", "")),
					PERSONALTOOLBARFOLDER: strings.TrimSpace(a.Prev().AttrOr("personal_toolbar_folder", "")),
					Dirs:                  make(map[string]*_Dir),
					Links:                 make(map[string]*Bookmark),
				}
			}
			a = a.Parent()
		}

		if t.bks.Title == "" {
			t.bks = _ct[_bk.Category[len(_bk.Category)-1]]
		}
		cur := t.bks
		for i := len(_bk.Category) - 2; i >= 0; i-- {

			if cur.Dirs == nil {
				cur.Dirs = make(map[string]*_Dir)
			}

			if cur.Links == nil {
				cur.Links = make(map[string]*Bookmark)
			}

			_name := _bk.Category[i]
			if cur.Dirs[_name] == nil {
				cur.Dirs[_name] = _ct[_name]
			}

			cur = cur.Dirs[_name]
		}
		cur.Links[_bk.Title] = _bk
	})
}

func (t *Bookmarks) Export() string {
	return fmt.Sprintf(header, t.h3(t.bks))
}

func (t *Bookmarks) Json() []byte {
	dd, err := json.Marshal(t.bks)
	assert.MustNotError(err)
	return dd
}
