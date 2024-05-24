package plex

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

type MediaContainer struct {
	XMLName             xml.Name `xml:"MediaContainer"`
	Size                int      `xml:"size,attr"`
	AllowSync           int      `xml:"allowSync,attr"`
	Art                 string   `xml:"art,attr"`
	Identifier          string   `xml:"identifier,attr"`
	LibrarySectionID    int      `xml:"librarySectionID,attr"`
	LibrarySectionTitle string   `xml:"librarySectionTitle,attr"`
	// ... (other MediaContainer attributes)
	Videos []Video `xml:"Video"`
}

type Video struct {
	XMLName       xml.Name `xml:"Video"`
	RatingKey     int      `xml:"ratingKey,attr"`
	Key           string   `xml:"key,attr"`
	Guid          string   `xml:"guid,attr"`
	Slug          string   `xml:"slug,attr"`
	Studio        string   `xml:"studio,attr"`
	Type          string   `xml:"type,attr"`
	Title         string   `xml:"title,attr"`
	ContentRating string   `xml:"contentRating,attr"`
	Summary       string   `xml:"summary,attr"`
}

type VideoShort struct {
	Title         string
	Summary       string
	ContentRating string
}

func GetRecentlyPlayed(c Client, sectionId string) ([]Video, error) {
	log.Println("connecting to Plex...")
	connectionUri := c.Connect(sectionId)
	log.Println("connected")

	log.Println("getting recently watched...")
	resp, err := c.MakeNetworkRequest(connectionUri, http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("received")

	var container MediaContainer
	if err := xml.Unmarshal(bodyBytes, &container); err != nil {
		return nil, err
	}
	log.Printf("total count: %v\n", len(container.Videos))
	log.Println("returning recently watched")

	return container.Videos, nil

}

func GetAllVideos(c Client, sectionId string) ([]VideoShort, error) {
	return nil, nil
}
