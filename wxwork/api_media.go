package wxwork

import (
	"io"
	"mime"
	"net/http"
)

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
			if ext, err := mime.ExtensionsByType(hdr.Get("content-type")); err == nil && len(ext) > 0 {
				filename = filename + ext[0]
			} else {
				logger().Infow("extensions by type fail", "err", err)
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
