package hero

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var replacer *regexp.Regexp

func init() {
	replacer, _ = regexp.Compile("\\s")
}

func TestWriteToFile(t *testing.T) {
	path := "/tmp/hero.test"
	content := "hello, hero"

	buffer := bytes.NewBufferString(content)
	writeToFile(path, buffer)

	defer os.Remove(path)

	if _, err := os.Stat(path); err != nil {
		t.Fail()
	}

	if c, err := ioutil.ReadFile(path); err != nil || string(c) != content {
		t.Fail()
	}
}

func TestGenAbsPath(t *testing.T) {
	dir, _ := filepath.Abs("./")

	parts := strings.Split(dir, "/")
	parent := strings.Join(parts[:len(parts)-1], "/")

	cases := []struct {
		in  string
		out string
	}{
		{in: "/", out: "/"},
		{in: ".", out: dir},
		{in: "../", out: parent},
	}

	for _, c := range cases {
		if genAbsPath(c.in) != c.out {
			t.Fail()
		}
	}
}

func TestGenerate(t *testing.T) {
	Generate(rootDir, rootDir, "template")

	cases := []struct {
		file string
		code string
	}{
		{file: "index.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/index.html
// DO NOT EDIT!
package template
`},
		{file: "item.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/item.html
// DO NOT EDIT!
package template
`},
		{file: "list.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/list.html
// DO NOT EDIT!
package template

import (
	"bytes"

	"github.com/shiyanhui/hero"
)

func Add(a, b int) int {
	return a + b
}
func UserList(userList []string, buffer *bytes.Buffer) {
	buffer.WriteString(` + "`" + `
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    ` + "`" + `)
	for _, user := range userList {
		buffer.WriteString(` + "`" + `
<div>
    <a href="/user/` + "`" + `)
		hero.EscapeHTML(user, buffer)
		buffer.WriteString(` + "`" + `">
        ` + "`" + `)
		buffer.WriteString(user)
		buffer.WriteString(` + "`" + `
    </a>
</div>
` + "`" + `)

	}

	buffer.WriteString(` + "`" + `
  </body>
</html>
` + "`" + `)

}
		`},
		{file: "listwriter.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/listwriter.html
// DO NOT EDIT!
package template

import (
	"io"

	"github.com/shiyanhui/hero"
)

func UserListToWriter(userList []string, w io.Writer) {
	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	buffer.WriteString(` + "`" + `
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    ` + "`" + `)
	for _, user := range userList {
		buffer.WriteString(` + "`" + `
<div>
    <a href="/user/` + "`" + `)
		hero.EscapeHTML(user, buffer)
		buffer.WriteString(` + "`" + `">
        ` + "`" + `)
		buffer.WriteString(user)
		buffer.WriteString(` + "`" + `
    </a>
</div>
` + "`" + `)

	}

	buffer.WriteString(` + "`" + `
  </body>
</html>
` + "`" + `)
	w.Write(buffer.Bytes())
}
		`},
		{file: "listwriterresult.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/listwriterresult.html
// DO NOT EDIT!
package template

import (
	"io"

	"github.com/shiyanhui/hero"
)

func UserListToWriterWithResult(userList []string, w io.Writer) (n int, err error) {
	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	buffer.WriteString(` + "`" + `
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    ` + "`" + `)
	for _, user := range userList {
		buffer.WriteString(` + "`" + `
<div>
    <a href="/user/` + "`" + `)
		hero.EscapeHTML(user, buffer)
		buffer.WriteString(` + "`" + `">
        ` + "`" + `)
		buffer.WriteString(user)
		buffer.WriteString(` + "`" + `
    </a>
</div>
` + "`" + `)

	}

	buffer.WriteString(` + "`" + `
  </body>
</html>
` + "`" + `)
	return w.Write(buffer.Bytes())
}
		`},
	}

	for _, c := range cases {
		content, err := ioutil.ReadFile(filepath.Join(rootDir, c.file))
		if err != nil || !reflect.DeepEqual(
			replacer.ReplaceAll(content, nil),
			[]byte(replacer.ReplaceAllString(c.code, "")),
		) {
			t.Fail()
		}
	}
}

func TestGen(t *testing.T) {
	root := parseFile(rootDir, "list.html")
	buffer := new(bytes.Buffer)

	gen(root, buffer)

	if buffer.String() == replacer.ReplaceAllString(
		`for _, user := range userList {}`, "") {
		t.Fail()
	}
}
