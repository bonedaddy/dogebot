package imgur

import (
	"fmt"
	"math/rand"
	"time"
)

func (c *Client) Random(opt SearchOptions) (*Image, error) {
	gallery, _, err := c.Search(opt)
	if err != nil {
		return nil, err
	}
	nbImg := len(gallery.Data)
	if nbImg == 0 {
		return nil, fmt.Errorf("No data found for %+v", opt)
	}

	rand.Seed(time.Now().UnixNano())
	var image Image
	// Try 10 times to grab a sfw picture
	for i := 0; i < 10; i++ {
		image = gallery.Data[rand.Intn(nbImg)]
		if !image.Nsfw {
			return &image, nil
		}
	}
	return &image, nil
}
