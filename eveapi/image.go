package eveapi

import (
	"fmt"
	"path"
)

const imageServerHost = "https://imageserver.eveonline.com/"

type imageKind string

const (
	ImageAllianceLogo      imageKind = "Alliance"
	ImageCorpLogo                    = "Corporation"
	ImageCharacterPortrait           = "Character"
	ImageItemTypeIcon                = "Type"
	ImageShipRender                  = "Render"
)

func ImageURL(kind imageKind, id int, width int) string {
	ext := "png"
	if kind == ImageCharacterPortrait {
		ext = "jpg"
	}
	return imageServerHost + path.Join(string(kind), fmt.Sprintf("%d_%d.%s", id, width, ext))
}
