package bookmarks

import (
	"bytes"
	"github.com/pubgo/assert"
	"github.com/pubgo/gotry"
	"io/ioutil"
	"testing"
)

func TestName1(t *testing.T) {
	dt, err := ioutil.ReadFile("bookmarks_2019_4_18.html")
	assert.MustNotError(err)

	bk := &Bookmarks{bks: &_Dir{}}
	bk.Import(bytes.NewReader(dt))
	gotry.Try(func() {
		bk.Import(bytes.NewReader(dt))
	}).P()

	{
		dt := bk.Export()
		assert.MustNotError(err)
		assert.MustNotError(ioutil.WriteFile("test.html", []byte(dt), 0755))
		dt1 := bk.Json()
		assert.MustNotError(ioutil.WriteFile("test.json", dt1, 0755))
	}

}
