package plex

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

const allMovies = true

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
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	ContentRating string `json:"content_rating"`
}

func (v VideoShort) String() string {
	return "Title: " + v.Title + "\nSummary: " + v.Summary + "\nContent Rating: " + v.ContentRating
}

func fullToShort(vids []Video, limit int) []VideoShort {
	shorts := make([]VideoShort, 0, limit)
	for i, vid := range vids {
		if i >= limit {
			break
		}
		shorts = append(shorts, VideoShort{Title: vid.Title, Summary: vid.Summary, ContentRating: vid.ContentRating})
	}

	return shorts
}

func GetRecentlyPlayed(c Client, sectionId string, limit int) ([]VideoShort, error) {
	log.Println("connecting to Plex...")
	connectionUri := c.Connect(sectionId, !allMovies)
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

	shorts := fullToShort(container.Videos, limit)
	log.Printf("returning %v recently watched\n", limit)

	return shorts, nil
}

func GetAllVideos(c Client, sectionId string) ([]VideoShort, error) {
	log.Println("connecting to Plex...")
	connectionUri := c.Connect(sectionId, allMovies)
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
	shorts := fullToShort(container.Videos, len(container.Videos))
	log.Println("returning all movies")

	return shorts, nil
}
