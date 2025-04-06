package main

type config struct {
	Urls       []string `json:"rss_urls"`
	SyncPeriod int      `json:"sync_period"`
}
