package opengraph_test

import (
	"strings"
	"testing"
	"time"

	"github.com/isacikgoz/go-opengraph/opengraph"
)

const html = `
  <!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" dir="ltr" lang="en-US">
<head profile="http://gmpg.org/xfn/11">
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>WordPress &#8250;   WordPress 4.3 &#8220;Billie&#8221;</title>

<!-- Jetpack Open Graph Tags -->
<meta property="og:type" content="article" />
<meta property="og:title" content="WordPress 4.3 &quot;Billie&quot;" />
<meta property="og:url" content="https://wordpress.org/news/2015/08/billie/" />
<meta property="og:description" content="Version 4.3 of WordPress, named &quot;Billie&quot; in honor of jazz singer Billie Holiday, is available for download or update in your WordPress dashboard. New features in 4.3 make it even easier to format y..." />
<meta property="og:article:published_time" content="2015-08-18T19:12:38+00:00" />
<meta property="og:article:modified_time" content="2015-08-19T13:10:24+00:00" />
<meta property="og:site_name" content="WordPress News" />
<meta property="og:image" content="https://www.gravatar.com/avatar/2370ea5912750f4cb0f3c51ae1cbca55?d=mm&amp;s=180&amp;r=G" />
<meta property="og:locale" content="en_US" />
<meta name="twitter:site" content="@WordPress" />
<meta name="twitter:card" content="summary" />
<meta name="twitter:creator" content="@WordPress" />
  `

func BenchmarkOpenGraph_ProcessHTML(b *testing.B) {
	og := opengraph.NewOpenGraph()
	b.ReportAllocs()
	b.SetBytes(int64(len(html)))
	for i := 0; i < b.N; i++ {
		if err := og.ProcessHTML(strings.NewReader(html)); err != nil {
			b.Fatal(err)
		}
	}
}

func TestOpenGraphForwardDetails(t *testing.T) {
	const sample = `<!DOCTYPE html>
	<html lang="en" xmlns="http://www.w3.org/1999/xhtml">
	<meta property="og:title" content="Humpback"/>
	<meta property="og:type" content="video.other" />
	<meta property="og:image"             content="https://i.imgur.com/h5FZ72Yh.jpg" />
	<meta property="og:video:width"       content="720" />
	<meta property="og:video:height"      content="720" />
	<meta property="og:video"             content="https://i.imgur.com/h5FZ72Y.mp4" />
	<meta property="og:video:secure_url"  content="https://i.imgur.com/h5FZ72Y.mp4" />
	<meta property="og:video:type"        content="video/mp4" />
	<meta property="og:description" content="Imgur: The magic of the Internet"/>

	<meta property="og:video"             content="https://i.imgur.com/h5FZ72Y.ogg" />
	<meta property="og:video:type"        content="video/ogg" />
	<meta property="og:video:width"       content="360" />

	<meta property="og:audio:type"       content="audio/mp3" />
	<meta property="og:audio"            content="https://i.imgur.com/h5FZ72Y.mp3" />
	`

	og := opengraph.NewOpenGraph()
	err := og.ProcessHTML(strings.NewReader(sample))

	if err != nil {
		t.Fatal(err)
	}

	if og.Type != "video.other" {
		t.Error("type parsed incorrectly")
	}

	if len(og.Videos) != 2 {
		t.Error("videos parsed incorrectly")
	} else {
		if og.Videos[0].URL != "https://i.imgur.com/h5FZ72Y.mp4" {
			t.Error("video url parsed incorrectly")
		}
		if og.Videos[0].Width != 720 {
			t.Error("video width parsed incorrectly")
		}
		if og.Videos[1].URL != "https://i.imgur.com/h5FZ72Y.ogg" {
			t.Error("video url parsed incorrectly")
		}
		if og.Videos[1].Width != 360 {
			t.Error("video width parsed incorrectly")
		}
	}

	if len(og.Audios) != 1 {
		t.Error("audio parsed incorrectly")
	} else {
		if og.Audios[0].URL != "https://i.imgur.com/h5FZ72Y.mp3" {
			t.Error("audio url parsed incorrectly")
		}
		if og.Audios[0].Type != "audio/mp3" {
			t.Error("audio type parsed incorrectly")
		}
	}
}

func TestOpenGraphProcessHTML(t *testing.T) {
	og := opengraph.NewOpenGraph()
	err := og.ProcessHTML(strings.NewReader(html))

	if err != nil {
		t.Fatal(err)
	}

	if og.Type != "article" {
		t.Error("type parsed incorrectly")
	}

	if len(og.Title) == 0 {
		t.Error("title parsed incorrectly")
	}

	if len(og.URL) == 0 {
		t.Error("url parsed incorrectly")
	}

	if len(og.Description) == 0 {
		t.Error("description parsed incorrectly")
	}

	if len(og.Images) == 0 {
		t.Error("images parsed incorrectly")
	} else {
		if len(og.Images[0].URL) == 0 {
			t.Error("image url parsed incorrectly")
		}
	}

	if len(og.Locale) == 0 {
		t.Error("locale parsed incorrectly")
	}

	if len(og.SiteName) == 0 {
		t.Error("site name parsed incorrectly")
	}

	if og.Article == nil {
		t.Error("articles parsed incorrectly")
	} else {
		ev, _ := time.Parse(time.RFC3339, "2015-08-18T19:12:38+00:00")
		if !og.Article.PublishedTime.Equal(ev) {
			t.Error("article published time parsed incorrectly")
		}
	}
}

func TestOpenGraphProcessHTML_YouTube(t *testing.T) {
	const youTubeHTML = `
	<!DOCTYPE html>
	<html style="font-size: 10px;font-family: Roboto, Arial, sans-serif;" lang="en-GB">
	<head>
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	</head>
	<body dir="ltr" no-y-overflow>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta property="og:site_name" content="YouTube">
	<meta property="og:url" content="https://www.youtube.com/watch?v=test123">
	<meta property="og:title" content="Test 123">
	<meta property="og:image" content="https://i.ytimg.com/vi/test123/hqdefault.jpg">
	<meta property="og:image:width" content="480">
	<meta property="og:image:height" content="360">
	<meta property="og:description" content="This is a test.">
	<meta property="og:type" content="video.other">
	<meta property="og:video:url" content="https://www.youtube.com/embed/test123">
	<meta property="og:video:secure_url" content="https://www.youtube.com/embed/test123">
	<meta property="og:video:type" content="text/html">
	<meta property="og:video:width" content="480">
	<meta property="og:video:height" content="360">
	<meta property="og:video:tag" content="Test">
	<meta property="og:video:tag" content="Testing">
	</body>
	</html>
	`

	og := opengraph.NewOpenGraph()
	err := og.ProcessHTML(strings.NewReader(youTubeHTML))

	if err != nil {
		t.Fatal(err)
	}

	if og.Type != "video.other" {
		t.Error("type parsed incorrectly")
	}

	if len(og.Title) == 0 {
		t.Error("title parsed incorrectly")
	}

	if len(og.URL) == 0 {
		t.Error("url parsed incorrectly")
	}

	if len(og.Description) == 0 {
		t.Error("description parsed incorrectly")
	}

	if len(og.Images) == 0 {
		t.Error("images parsed incorrectly")
	} else {
		if len(og.Images[0].URL) == 0 {
			t.Error("image url parsed incorrectly")
		}
	}

	if len(og.SiteName) == 0 {
		t.Error("site name parsed incorrectly")
	}

	if len(og.Videos) != 1 {
		t.Error("videos parsed incorrectly")
	} else {
		if len(og.Videos[0].URL) == 0 {
			t.Error("video url parsed incorrectly")
		}
	}
}

func TestOpenGraphProcessMeta(t *testing.T) {
	og := opengraph.NewOpenGraph()

	og.ProcessMeta(map[string]string{"property": "og:type", "content": "book"})

	if og.Type != "book" {
		t.Error("wrong og:type processing")
	}

	og.ProcessMeta(map[string]string{"property": "og:book:isbn", "content": "123456"})

	if og.Book == nil {
		t.Error("wrong book type processing")
	} else {
		if og.Book.ISBN != "123456" {
			t.Error("wrong book isbn processing")
		}
	}

	og.ProcessMeta(map[string]string{"property": "og:article:section", "content": "testsection"})

	if og.Article != nil {
		t.Error("article processed when it should not be")
	}

	og.ProcessMeta(map[string]string{"property": "og:book:author", "content": "https://site.com/author/VitaliDeatlov"})
	og.ProcessMeta(map[string]string{"property": "og:book:author", "content": "https://site.com/author/JohnDoe"})

	if og.Book != nil {
		if len(og.Book.Authors) != 2 {
			t.Error("Incorrect amount of book authors")
		}
	} else {
		t.Error("Book data wasn't processed")
	}

	og.ProcessMeta(map[string]string{"property": "og:music:creator", "content": "https://site.com/author/JohnDoe"})

	if og.Music == nil {
		t.Error("Incorrectly processed music creator")
	}

	og.ProcessMeta(map[string]string{"property": "og:music:song", "content": "https://site.com/song/1"})
	og.ProcessMeta(map[string]string{"property": "og:music:musician", "content": "https://site.com/author/VitaliDeatlov"})
	og.ProcessMeta(map[string]string{"property": "og:music:song", "content": "https://site.com/song/2"})

	if len(og.Music.Songs) != 2 {
		t.Error("Incorrectly parsed music song urls")
	}
	if len(og.Music.Musicians) != 1 {
		t.Error("Incorrectly parsed song musicians")
	}
}
