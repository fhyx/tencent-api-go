package wxwork

import (
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
)

var validFilenameRegex = regexp.MustCompile(`filename\*=utf-8''(.*?)("|\;)`)

type MediaFunc func(filename string, body io.Reader, cl int64) error

func (a *API) GetMedia(mediaID string, mf MediaFunc) error {
	err := a.c.Do("GET", UriPrefix+"/media/get?media_id="+mediaID, nil, func(hdr http.Header, r io.Reader, cl int64) error {
		logger().Infow(" header info ", "header ", hdr, "media id ", mediaID)

		filename := mediaID

		_, params, err := mime.ParseMediaType(hdr.Get("content-disposition"))
		if err != nil {
			logger().Infow("get content disposition fail", "err", err)
		}
		if f, ok := params["filename"]; ok {
			filename = f
		} else {
			if filename == mediaID {
				// 通过正则匹配替换文件名
				if name := GetFilenameRegex(hdr.Get("content-disposition")); name != "" {
					filename = name
				}
			}
		}

		if err := mf(filename, r, cl); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger().Infow("get media fail", "err", err)
		return err
	}
	return nil
}

func GetFilenameRegex(s string) (filename string) {
	matchs := validFilenameRegex.FindStringSubmatch(s)
	if matchs == nil {
		logger().Infow("failed to match", "content-disposition", s)
	}
	for k, m := range matchs {
		if k == 1 && m != "" {
			if str, err := url.QueryUnescape(m); err == nil {
				filename = str
			}
		}
	}
	logger().Infow("regexp matched", "filename", filename)
	return filename
}
