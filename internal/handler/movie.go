package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"movie-web/internal/db"
	"movie-web/internal/model"
)

var movies = []model.Movie{
	{ID: 1, Title: "Inception", Year: 2010, Rating: 8.8, Genres: []string{"Sci-Fi", "Action"}},
	{ID: 2, Title: "Parasite", Year: 2019, Rating: 8.6, Genres: []string{"Thriller", "Drama"}},
	{ID: 3, Title: "Spirited Away", Year: 2001, Rating: 8.6, Genres: []string{"Animation"}},
}

const movieAPIURL = "https://vsmov.com/api/danh-sach/phim-moi-cap-nhat"

var movieAPIClient = &http.Client{Timeout: 10 * time.Second}

type movieAPIResponse struct {
	Status bool `json:"status"`
	Items  []struct {
		ID         int    `json:"_id"`
		Name       string `json:"name"`
		OriginName string `json:"origin_name"`
		Slug       string `json:"slug"`
		PosterURL  apiURL `json:"poster_url"`
		ThumbURL   apiURL `json:"thumb_url"`
		Year       int    `json:"year"`
		TMDB       struct {
			VoteAverage string `json:"vote_average"`
		} `json:"tmdb"`
	} `json:"items"`
}

type apiURL string

func (url *apiURL) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		*url = apiURL(value)
		return nil
	}

	var object map[string]string
	if err := json.Unmarshal(data, &object); err != nil {
		return nil
	}
	for _, key := range []string{"url", "src", "path"} {
		if value, ok := object[key]; ok {
			*url = apiURL(value)
			break
		}
	}
	return nil
}

func fetchMovies(page int) ([]model.Movie, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?page=%d", movieAPIURL, page), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	response, err := movieAPIClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("movie API returned status %s", response.Status)
	}

	var payload movieAPIResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if !payload.Status {
		return nil, fmt.Errorf("movie API returned an unsuccessful response")
	}

	result := make([]model.Movie, 0, len(payload.Items))
	for _, item := range payload.Items {
		rating, err := strconv.ParseFloat(item.TMDB.VoteAverage, 64)
		if err != nil {
			rating = 0
		}
		result = append(result, model.Movie{
			ID: item.ID, Title: item.Name, OriginName: item.OriginName,
			Year: item.Year, Rating: rating, PosterURL: string(item.PosterURL),
			ThumbURL: string(item.ThumbURL), Slug: item.Slug, Genres: []string{},
		})
	}

	return result, nil
}

func GetMovies(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page không hợp lệ"})
		return
	}

	apiMovies, _ := fetchMovies(page)

	combined := make([]model.Movie, 0)
	seenSlugs := make(map[string]struct{})
	seenTitlesAndPosters := make(map[string]struct{})
	addUnique := func(movie model.Movie) {
		slug := strings.ToLower(strings.TrimSpace(movie.Slug))
		title := strings.ToLower(strings.TrimSpace(movie.Title))
		poster := strings.ToLower(strings.TrimSpace(strings.Split(movie.PosterURL, "?")[0]))
		titlePosterKey := ""
		if title != "" && poster != "" {
			titlePosterKey = title + "|" + poster
		}

		_, slugSeen := seenSlugs[slug]
		_, titlePosterSeen := seenTitlesAndPosters[titlePosterKey]
		if (slug != "" && slugSeen) || (titlePosterKey != "" && titlePosterSeen) {
			return
		}

		combined = append(combined, movie)
		if slug != "" {
			seenSlugs[slug] = struct{}{}
		}
		if titlePosterKey != "" {
			seenTitlesAndPosters[titlePosterKey] = struct{}{}
		}
	}

	// Prepend custom Admin Movies from PostgreSQL
	if db.DB != nil {
		var dbMovies []model.MovieDB
		db.DB.Order("created_at desc").Find(&dbMovies)
		for _, dbm := range dbMovies {
			var genresList []string
			if dbm.Genres != "" {
				for _, g := range strings.Split(dbm.Genres, ",") {
					genresList = append(genresList, strings.TrimSpace(g))
				}
			}
			addUnique(model.Movie{
				ID:         int(dbm.ID),
				Title:      dbm.Title,
				OriginName: dbm.OriginName,
				Slug:       dbm.Slug,
				Year:       dbm.Year,
				Rating:     dbm.Rating,
				Genres:     genresList,
				PosterURL:  dbm.PosterURL,
				ThumbURL:   dbm.ThumbURL,
			})
		}
	}

	for _, movie := range apiMovies {
		addUnique(movie)
	}
	c.JSON(http.StatusOK, combined)
}

func GetMovieByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	// Check DB Movies first
	if db.DB != nil {
		var dbm model.MovieDB
		if err := db.DB.First(&dbm, id).Error; err == nil {
			var genresList []string
			if dbm.Genres != "" {
				for _, g := range strings.Split(dbm.Genres, ",") {
					genresList = append(genresList, strings.TrimSpace(g))
				}
			}
			c.JSON(http.StatusOK, model.Movie{
				ID:         int(dbm.ID),
				Title:      dbm.Title,
				OriginName: dbm.OriginName,
				Slug:       dbm.Slug,
				Year:       dbm.Year,
				Rating:     dbm.Rating,
				Genres:     genresList,
				PosterURL:  dbm.PosterURL,
				ThumbURL:   dbm.ThumbURL,
			})
			return
		}
	}

	pageMovies, err := fetchMovies(1)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy dữ liệu phim", "detail": err.Error()})
		return
	}

	for _, m := range pageMovies {
		if m.ID == id {
			c.JSON(http.StatusOK, m)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phim"})
}

func CreateMovie(c *gin.Context) {
	var newMovie model.Movie

	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newMovie.ID = len(movies) + 1
	movies = append(movies, newMovie)

	c.JSON(http.StatusCreated, newMovie)
}
