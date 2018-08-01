package models

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed/extensions"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID                bson.ObjectId   `bson:"_id"`
	Username          string          `bson:"username"`
	Lastname          string          `bson:"lastname"`
	Firstname         string          `bson:"firstname"`
	Photourl          string          `bson:"photourl"`
	Hash              string          `bson:"hash"`
	SubscribeFeedURLs []string        `bson:"subscribeFeedURLs"`
	UnReadItems       []bson.ObjectId `bson:"unReadItems"`
	StarItems         []bson.ObjectId `bson:"starItems"`
	ReadItLaterItems  []bson.ObjectId `bson:"readItLaterItems"`
}

// Unique field "FeedURL" in order to prevent feed polution
type Feed struct {
	ID              bson.ObjectId     `bson:"_id"`
	FeedURL         string            `bson:"feedURL"`
	Title           string            `bson:"title"`
	Description     string            `bson:"description"`
	Link            string            `bson:"link"`
	FeedLink        string            `bson:"feedLink"`
	Updated         string            `bson:"updated"`
	UpdatedParsed   *time.Time        `bson:"updatedParsed"`
	Published       string            `bson:"published"`
	PublishedParsed *time.Time        `bson:"publishedParsed"`
	Author          *Person           `bson:"author"`
	Language        string            `bson:"language"`
	Image           *Image            `bson:"image"`
	Copyright       string            `bson:"copyright"`
	Generator       string            `bson:"generator"`
	Categories      []string          `bson:"categories"`
	Extensions      ext.Extensions    `bson:"extensions"`
	Custom          map[string]string `bson:"custom"`
	FeedType        string            `bson:"feedType"`
	FeedVersion     string            `bson:"feedVersion"`
	SubscribeCount  int               `bson:"subscribeCount"`
}

// Item is the universal Item type that atom.Entry
// and rss.Item gets translated to.  It represents
// a single entry in a given feed.
type Item struct {
	FeedID          bson.ObjectId     `bson:"feedID"`
	Title           string            `bson:"title"`
	Description     string            `bson:"description"`
	Content         string            `bson:"content"`
	Link            string            `bson:"link"`
	Updated         string            `bson:"updated"`
	UpdatedParsed   *time.Time        `bson:"updatedParsed"`
	Published       string            `bson:"published"`
	PublishedParsed string            `bson:"publishedParsed"`
	Author          *Person           `bson:"author"`
	GUID            string            `bson:"guid"`
	Image           *Image            `bson:"image"`
	Categories      []string          `bson:"categories"`
	Enclosures      []*Enclosure      `bson:"enclosures"`
	Extensions      ext.Extensions    `bson:"extensions"`
	Custom          map[string]string `bson:"custom"`
	StarCount       int               `bson:"starCount"`
}

// Person is an individual specified in a feed
// (e.g. an author)
type Person struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

// Image is an image that is the artwork for a given
// feed or item.
type Image struct {
	URL   string `bson:"url"`
	Title string `bson:"title"`
}

// Enclosure is a file associated with a given Item.
type Enclosure struct {
	URL    string `bson:"url"`
	Length string `bson:"length"`
	Type   string `bson:"type"`
}

//
type Thing struct {
	ID         int
	Name       string
	Content    string
	CreateTime time.Time
}
type Note struct {
	ID         int
	Name       string
	Content    string
	CreateTime time.Time
}

type Session struct {
	ID        bson.ObjectId `bson:"_id"`
	SessionID string        `bson:"SessionId"`
}

var DBConfig = struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}{}

//= =!
var Users *mgo.Collection
var Feeds *mgo.Collection
var Items *mgo.Collection
var Sessions *mgo.Collection

func Init() {
	url := "mongodb://admin:feedall@127.0.0.1:27017/feedall"
	Session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	fmt.Println("Start dial mongodb!")

	Users = Session.DB("feedall").C("users")
	Feeds = Session.DB("feedall").C("feeds")
	Items = Session.DB("feedall").C("items")
	Sessions = Session.DB("feedall").C("sessions")
}

func Insert(collection *mgo.Collection, i interface{}) error {
	return collection.Insert(i)
}
func FindOne(collection *mgo.Collection, q interface{}, i interface{}) error {
	return collection.Find(q).One(i)
}
func FindAll(collection *mgo.Collection, q interface{}, i interface{}) error {
	return collection.Find(q).All(i)
}
func FindSortLimit(collection *mgo.Collection, q interface{}, s string, n int, i interface{}) error {
	return collection.Find(q).Sort(s).Limit(n).All(i)
}
func PipeAll(collection *mgo.Collection, q interface{}, i interface{}) error {
	return collection.Pipe(q).All(i)
}
func PipeOne(collection *mgo.Collection, q interface{}, i interface{}) error {
	return collection.Pipe(q).One(i)
}
func Update(collection *mgo.Collection, q interface{}, i interface{}) error {
	return collection.Update(q, i)
}
func Upsert(collection *mgo.Collection, q interface{}, i interface{}) (info *mgo.ChangeInfo, err error) {
	return collection.Upsert(q, i)
}
