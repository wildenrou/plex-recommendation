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

func GetRecentlyPlayed(c Client, sectionId string) ([]Video, error) {
	connectionUri := c.Connect(sectionId)
	log.Println("connectionUri: " + connectionUri)

	resp, err := c.MakeNetworkRequest(connectionUri, http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var container MediaContainer
	if err := xml.Unmarshal(bodyBytes, &container); err != nil {
		return nil, err
	}

	return container.Videos, nil

}
