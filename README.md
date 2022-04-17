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

