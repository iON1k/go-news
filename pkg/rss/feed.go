package rss

// Корневой объект RSS
type Feed struct {
	Chanel Channel `xml:"channel"`
}

// Канал RSS
type Channel struct {
	Items []Item `xml:"item"`
}

// Пост из RSS
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}
