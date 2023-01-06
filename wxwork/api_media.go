package wxwork

import (
	"io"
	"mime"
	"net/http"
)

type MediaFunc func(filename string, body io.Reader, cl int64) error

func (a *API) GetMedia(mediaID string, mf MediaFunc) error {
	err := a.c.Do("GET", UriPrefix+"/media/get?media_id="+mediaID, nil, func(hdr http.Header, r io.Reader, cl int64) error {
		contentDisposition := hdr.Get("content-disposition")
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			logger().Infow("get content disposition fail", "err", err)
			return err
		}

		filename := mediaID
		if f, ok := params["filename"]; ok {
			filename = f
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
