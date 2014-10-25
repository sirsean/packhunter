package rank

import (
	"github.com/sirsean/friendly-ph/ph"
	"github.com/sirsean/friendly-ph/model"
	"log"
	"sort"
	"sync"
)

type Product struct {
	ph.Post
	TagVotes int `json:"tag_votes"`
}
type productSorter struct {
	products []Product
	by func(p1, p2 Product) bool
}
func (s productSorter) Len() int {
	return len(s.products)
}
func (s productSorter) Swap(i, j int) {
	s.products[i], s.products[j] = s.products[j], s.products[i]
}
func (s productSorter) Less(i, j int) bool {
	return s.by(s.products[i], s.products[j])
}

type syncVoteTracker struct {
	sync.RWMutex
	numVotes map[int]int
	products map[int]ph.Post
}

func newVoteTracker() syncVoteTracker {
	t := syncVoteTracker{}
	t.numVotes = make(map[int]int)
	t.products = make(map[int]ph.Post)
	return t
}

func (t syncVoteTracker) Add(v ph.Vote) {
	t.Lock()
	defer t.Unlock()
	t.numVotes[v.PostId] += 1
	t.products[v.PostId] = v.Post
}

func (t syncVoteTracker) Products() []Product {
	t.Lock()
	defer t.Unlock()

	products := make([]Product, 0)
	for k, v := range t.products {
		log.Printf("k %v", k)
		products = append(products, Product{
			TagVotes: t.numVotes[k],
			Post: v,
		})
	}
	sorter := productSorter{
		products: products,
		by: func(p1, p2 Product) bool {
			if p1.TagVotes == p2.TagVotes {
				return p1.CreatedAt.After(p2.CreatedAt)
			}
			return p1.TagVotes > p2.TagVotes
		},
	}
	sort.Sort(sorter)

	return products
}

func ForTag(accessToken string, tag model.Tag) []Product {
	voteTracker := newVoteTracker()
	var wg sync.WaitGroup
	wg.Add(len(tag.Users))
	for _, u := range tag.Users {
		go func(wg *sync.WaitGroup, u ph.User, accessToken string) {
			log.Printf("%v", u)
			user := ph.GetUserById(accessToken, u.Id)
			log.Printf("votes %v", user.VotesCount)
			for _, v := range user.Votes {
				voteTracker.Add(v)
				//productVotes[v.PostId] += 1
				//productMap[v.PostId] = v.Post
			}
			wg.Done()
		}(&wg, u, accessToken)
	}
	wg.Wait()

	products := voteTracker.Products()

	return products
}
