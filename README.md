# hackernews-api
An HTTP API micro-service that uses the YCombinator HackerNews API and upon request returns the top 50 of the latest 500 stories where the author has karma above 2413. The position is determined by the number of comments in relation to the top 50 stories. The story with the most comments should have position: 1, so on and so forth.

# Specification

## Request:

GET /stories

## Response:

```json
{
"stories": [
{
"author": "nick1",
"karma": 5341,
"comments": 192,
"title": "article title", "position": 1
},
{
"author": "nick2",
"karma": 7629,
"title": "article title", "comments": 12,
"position": 3
},
{
"author": "nick3",
"karma": 6293,
"title": "article title", "comments": 180,
"position": 2
}
]
}
```

# Implementation

As the ask is to get a curated list of the top 50 stories based on the authors karma, so we need to first set up different models for our requirements.
We will break down the application into 3 different parts

1. `Models` : We define the structs for the various items. The models which we have used are items(Stories, comments, jobs, Ask HNs and even polls are just items. They're identified by their ids, which are unique integers, and live under /v0/item/<id>), story(the response structure), topstories(the list of top stories) and user(which contains the varius fileds of an author)
2. `Handlers` : Where we mock and call the various services. This also enables us to extend our program and later ad additional functionalities instead of modifying the existing fuctions
3. `Hackernews services` : We define the various services as per our requirememts.
   1. GetItem returns item from hackernews API for given id
      ```
      GetItem(id int) (models.Item, error)
      ```
   2. GetUser returns user from hackernews API for given username
      ```
      GetUser(username string) (models.User, error)
      ```
   3. GetTopStories returns top 500 stories from hackernews API
      ```
      GetTopStories() ([]models.Item, error)
      ```
   4. Curated50 returns top 50 of the latest 500 stories where the author has karma above 2413 with most comments
      ```
      Curated50(minKarma int) ([]models.Story, error)
      ```

We first collect the curated50 stories based on the authors minimum karma i.e 2413 and then sort the list based on the comments. 

We use goroutines( 50 workers ) for parallel execution of our code as it takes around 3 mins to go through the entire 500 list if we are using a serialisation method.

We have also used a in memory caching mechanism that creates a cache with a default expiration time of 5 minutes, and which purges expired items every 10 minutes

```go
cache:      cache.New(5*time.Minute, 10*time.Minute),
```