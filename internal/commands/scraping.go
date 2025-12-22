package cmnd

import(
	"fmt"
	"context"
	"time"
	"database/sql"
	"github.com/CoupDeGrace92/gator/internal/web"
	"github.com/CoupDeGrace92/gator/internal/config"
	"github.com/CoupDeGrace92/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func variableTimeParser(date string) (time.Time, error){
	layouts := []string{"Mon, 02 Jan 2006 15:04:05 -0700","Mon, 02 Jan 2006 15:04:05 MST","2006-01-02T15:04:05Z07:00","02 Jan 06 15:04 MST","02 Jan 06 15:04 -0700"}
	for _, layout := range layouts {
		t, err := time.Parse(layout, date)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("Could not parse %v with the available formats", date)
}

func ScrapeFeeds(s *config.State, ctx context.Context) error {
	url, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil{
		err = fmt.Errorf("Error getting next feed to fetch: %v\n", err)
		return err
	}

	RssFeed, err := web.FetchFeed(ctx, url)
	if err != nil{
		err = fmt.Errorf("Error getting next feed to fetch: %v\n", err)
		return err
	}

	err = s.Db.MarkFeedFetched(ctx, url)
	if err != nil {
		err = fmt.Errorf("Error marking feed (%v) as fetched: %v\n", url, err)
		return err
	}

	feedID, err := s.Db.GetFeedIds(ctx, url)
	if err != nil {
		err = fmt.Errorf("Error in getting feed id: %v\n", err)
		return err
	}

	for _, item := range RssFeed.Channel.Item{
		now := time.Now()
		var argParams database.CreatePostParams
		argParams.ID = uuid.New()
		argParams.CreatedAt = now
		argParams.UpdatedAt = now
		argParams.Title = item.Title
		argParams.Url = item.Link
		if item.Description == ""{
			argParams.Description = sql.NullString{String: "", Valid: false}
		} else {
			argParams.Description = sql.NullString{String: item.Description, Valid: true}
		}
		if item.PubDate == ""{
			argParams.PublishedAt = sql.NullTime{Valid: false}
		} else {
			t, err := variableTimeParser(item.PubDate)
			if err != nil {
				fmt.Printf("failed to parse %v from %v on article %v", item.PubDate, RssFeed.Channel.Title, item.Title)
				argParams.PublishedAt = sql.NullTime{Valid: false}
			} else {
				argParams.PublishedAt = sql.NullTime{Time: t, Valid: true}
			}
		}

		argParams.FeedID = feedID
		_, err := s.Db.CreatePost(ctx, argParams)
		if err != nil {
			if pqErr, ok :=err.(*pq.Error); ok && pqErr.Code == "23505"{
				fmt.Printf("Ignoring duplicate post for %v\n", item.Title)
				continue //23505 is the error code for a duplicate error
			}
			err = fmt.Errorf("Error creating post: %v\n", err)
			return err
		}
	}
	return nil
}