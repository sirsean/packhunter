package ph

import (
	"fmt"
	"log"
	"github.com/sirsean/friendly-ph/config"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
)

/*
{
  "user": {
    "id": 1,
    "name": "Karl User the 7th",
    "headline": "Product Hunter",
    "created_at": "2014-10-25T08:07:57.425-07:00",
    "username": "producthunter7",
    "image_url": {
      "48px": "/assets/ph-logo.png",
      "73px": "/assets/ph-logo.png",
      "original": "/assets/ph-logo.png"
    },
    "profile_url": "http://www.producthunt.com/producthunter7",
    "votes_count": 0,
    "posts_count": 0,
    "maker_of_count": 0,
    "email": "email7@factory.com",
    "role": "user",
    "permissions": {
      "can_vote_posts": true,
      "can_comment": false,
      "can_post": false
    },
    "notifications": {
      "total": 0,
      "unseen": 0
    },
    "first_time_user": false,
	"votes": [
      {
        "id": 2,
        "created_at": "2014-10-25T08:08:15.102-07:00",
        "post_id": 1,
        "post": {
          "id": 1,
          "name": "Awesome Idea #44",
          "tagline": "Great new search engine",
          "created_at": "2014-10-25T08:08:15.038-07:00",
          "day": "2014-10-25",
          "comments_count": 0,
          "votes_count": 2,
          "discussion_url": "http://www.producthunt.com/posts/awesome-idea-44",
          "redirect_url": "http://www.producthunt.com/l/3c2b5bf46a/1",
          "screenshot_url": {
            "300px": "http://placehold.it/850x850.png",
            "850px": "http://placehold.it/850x850.png"
          },
          "maker_inside": false,
          "user": {
            "id": 2,
            "name": "Karl User the 244th",
            "headline": "Product Hunter",
            "created_at": "2014-10-25T08:08:15.024-07:00",
            "username": "producthunter243",
            "image_url": {
              "48px": "/assets/ph-logo.png",
              "73px": "/assets/ph-logo.png",
              "original": "/assets/ph-logo.png"
            },
            "profile_url": "http://www.producthunt.com/producthunter243"
          }
        }
      },
	],
    "posts": [],
    "maker_of": []
  }
}
*/

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Headline string `json:"headline"`
	CreatedAt time.Time `json:"created_at"`
	Username string `json:"username"`
	ImageUrl ImageUrl `json:"image_url"`
	ProfileUrl string `json:"profile_url"`
	VotesCount int `json:"votes_count"`
	PostsCount int `json:"posts_count"`
	MakerOfCount int `json:"maker_of_count"`
	Email string `json:"email"`
	Role string `json:"role"`
	Permissions Permissions `json:"permissions"`
	Votes []Vote `json:"votes"`
	//Posts []Post `json:"posts"`
	//MakerOf []MakerOf `json:"maker_of"`
}

type ImageUrl struct {
	Small string `json:"48px"`
	Medium string `json:"73px"`
	Original string `json:"original"`
}

type Permissions struct {
	CanVotePosts bool `json:"can_vote_posts"`
	CanComment bool `json:"can_comment"`
	CanPost bool `json:"can_post"`
}

type Vote struct {
	Id int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	PostId int `json:"post_id"`
	Post Post `json:"post"`
}

type Post struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Tagline string `json:"tagline"`
	CreatedAt time.Time `json:"created_at"`
	Day string `json:"day"`
	CommentsCount int `json:"comments_count"`
	VotesCount int `json:"votes_count"`
	DiscussionUrl string `json:"discussion_url"`
	RedirectUrl string `json:"redirect_url"`
	ScreenshotUrl ScreenshotUrl `json:"screenshot_url"`
	MakerInside bool `json:"maker_inside"`
}

type ScreenshotUrl struct {
	Small string `json:"300px"`
	Large string `json:"850px"`
}

type UserResponse struct {
	User User `json:"user"`
}

func Me(accessToken string) User {
	url := fmt.Sprintf("%v/me", config.Get().ProductHunt.Endpoint)
	resp, _ := get(accessToken, url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	u := UserResponse{}
	json.Unmarshal(body, &u)

	return u.User
}

func GetUserById(accessToken string, userId int) User {
	url := fmt.Sprintf("%v/users/%v", config.Get().ProductHunt.Endpoint, userId)
	resp, _ := get(accessToken, url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	r := UserResponse{}
	json.Unmarshal(body, &r)

	return r.User
}

type FollowingResponse struct {
	Following []FollowingUser `json:"following"`
}

func (r FollowingResponse) Users() []User {
	users := make([]User, len(r.Following))
	for i, u := range r.Following {
		users[i] = u.User
	}
	return users
}

type FollowingUser struct {
	Id int `json:"id"`
	User User `json:"user"`
}

// TODO: Handle pagination of many following
func Following(accessToken string, userId int) []User {
	url := fmt.Sprintf("%v/users/%v/following",
		config.Get().ProductHunt.Endpoint,
		userId)
	resp, _ := get(accessToken, url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//b := make(map[string]interface{})
	//json.Unmarshal(body, &b)
	//log.Printf("me: %v", b)

	r := FollowingResponse{}
	json.Unmarshal(body, &r)
	log.Printf("following: %v", r)

	return r.Users()
}

func get(accessToken string, url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed /following %v", err)
	}
	return resp, err
}
