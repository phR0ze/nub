package tmpl

import (
	"testing"

	"github.com/phR0ze/n/pkg/sys"
	"github.com/stretchr/testify/assert"
)

var tmpDir = "../../test/temp"
var tmpFile = "../../test/temp/.tmp"

func TestLoad(t *testing.T) {
	clearTmpDir()

	data := `labels:
  chart: {{ name }}:{{ version }}
  release: {{ Release.Name }}
  heritage: {{ Release.Service }}`
	sys.WriteBytes(tmpFile, []byte(data))

	expected := `labels:
  chart: foo:1.0.2
  release: babble
  heritage: fish`

	result, err := Load(tmpFile, "{{ ", " }}", map[string]string{
		"name":            "foo",
		"version":         "1.0.2",
		"Release.Name":    "babble",
		"Release.Service": "fish",
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestNoTemplatVariablesFound(t *testing.T) {
	// When no template varibles are found the file should simple be returned as is
	clearTmpDir()

	data := `labels:
  chart: foo
  release: bar
  heritage: foobar`
	sys.WriteBytes(tmpFile, []byte(data))

	result, err := Load(tmpFile, "{{ ", " }}", map[string]string{
		"name": "foo",
	})
	assert.Nil(t, err)
	assert.Equal(t, data, result)

}

func TestTagSpaces(t *testing.T) {
	{
		// README.md example
		data := `labels:
	  chart: {{ name }}:{{ version }}
	  release: {{ Release.Name }}
	  heritage: {{ Release.Service }}`
		expected := `labels:
	  chart: foo:1.0.2
	  release: babble
	  heritage: fish`
		{
			// Tags with spaces
			tpl, err := New(data, "{{ ", " }}")
			assert.Nil(t, err)
			assert.Equal(t, []string{"name", "version", "Release.Name", "Release.Service"}, tpl.tags)
			result, err := tpl.Process(map[string]string{
				"name":            "foo",
				"version":         "1.0.2",
				"Release.Name":    "babble",
				"Release.Service": "fish",
			})
			assert.Nil(t, err)
			assert.Equal(t, expected, result)
		}
		{
			// Tags with no spaces
			tpl, err := New(data, "{{", "}}")
			assert.Equal(t, []string{"name", "version", "Release.Name", "Release.Service"}, tpl.tags)
			assert.NotNil(t, tpl)
			assert.Nil(t, err)
			result, err := tpl.Process(map[string]string{
				"name":            "foo",
				"version":         "1.0.2",
				"Release.Name":    "babble",
				"Release.Service": "fish",
			})
			assert.Nil(t, err)
			assert.Equal(t, expected, result)
		}
	}
}

func TestNoStartTag(t *testing.T) {
	tpl, err := New("foobar", "", "}}")
	assert.Nil(t, tpl)
	assert.NotNil(t, err)
}

func TestNoEndTag(t *testing.T) {
	tpl, err := New("foobar", "{{", "")
	assert.Nil(t, tpl)
	assert.NotNil(t, err)
}

func TestDataHasNoTags(t *testing.T) {
	tpl, err := New("foobar", "{{", "}}")
	assert.Nil(t, tpl)
	assert.NotNil(t, err)
}

func TestEmptyTag(t *testing.T) {
	tpl, err := New("foo{{}}bar", "{{", "}}")
	assert.Nil(t, err)
	result, err := tpl.Process(map[string]string{"": "111", "aaa": "bbb"})
	assert.Nil(t, err)
	assert.Equal(t, "foo111bar", result)
}

func TestSpaceTag(t *testing.T) {
	tpl, err := New("foo{{ }}bar", "{{", "}}")
	assert.Nil(t, err)
	result, err := tpl.Process(map[string]string{"": "111", "aaa": "bbb"})
	assert.Nil(t, err)
	assert.Equal(t, "foo111bar", result)
}

func TestOnlyTag(t *testing.T) {
	tpl, err := New("[foo]", "[", "]")
	assert.Nil(t, err)
	result, err := tpl.Process(map[string]string{"foo": "111", "aaa": "bbb"})
	assert.Nil(t, err)
	assert.Equal(t, "111", result)
}

func TestStartWithTag(t *testing.T) {
	tpl, err := New("[foo]barbaz", "[", "]")
	assert.Nil(t, err)
	result, err := tpl.Process(map[string]string{"foo": "111", "aaa": "bbb"})
	assert.Nil(t, err)
	assert.Equal(t, "111barbaz", result)
}

func TestEndWithTag(t *testing.T) {
	tpl, err := New("foobar[foo]", "[", "]")
	assert.Nil(t, err)
	result, err := tpl.Process(map[string]string{"foo": "111", "aaa": "bbb"})
	assert.Nil(t, err)
	assert.Equal(t, "foobar111", result)
}

func TestProcess(t *testing.T) {
	testProcess(t, "", "", true)
	testProcess(t, "a", "a", true)
	testProcess(t, "abc", "abc", true)
	testProcess(t, "{foo}", "xxxx", false)
	testProcess(t, "a{foo}", "axxxx", false)
	testProcess(t, "{foo}a", "xxxxa", false)
	testProcess(t, "a{foo}bc", "axxxxbc", false)
	testProcess(t, "{foo}{foo}", "xxxxxxxx", false)
	testProcess(t, "{foo}bar{foo}", "xxxxbarxxxx", false)

	// unclosed tag
	testProcess(t, "{unclosed", "{unclosed", true)
	testProcess(t, "{{unclosed", "{{unclosed", true)
	testProcess(t, "{un{closed", "{un{closed", true)

	// test unknown tag (they get removed from the final output)
	testProcess(t, "{unknown}", "", false)
	testProcess(t, "{foo}q{unexpected}{missing}bar{foo}", "xxxxqbarxxxx", false)
}

func testProcess(t *testing.T, template, expected string, shouldErr bool) {
	tpl, err := New(template, "{", "}")
	if shouldErr {
		assert.Nil(t, tpl)
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
		result, err := tpl.Process(map[string]string{"foo": "xxxx"})
		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	}
}

func clearTmpDir() {
	if sys.Exists(tmpDir) {
		sys.RemoveAll(tmpDir)
	}
	sys.MkdirP(tmpDir)
}
