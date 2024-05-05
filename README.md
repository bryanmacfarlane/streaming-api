# Steaming API

## Problem Statement

In a modern rich web application, a web client makes async requests to the server to dynamically render views.  It can either be done as an initial SSR page (SEO) followed by asynchronous calls to populate views, or it could be a completely client side rendered page.  They both need to solve similar problems.

The ideal case is your server calls are fast, lightweight and latency is low.  In that case, you can lean on AJAX browser requests to concurrently populate the page.  Google search type ahead is the perfect example where the server is fast enough to dynamically respond as you type.

However, what if that isn't the case?  What if latency isn't ideal for some regions?  What if the server has a non-trivial overhead per call?  (auth, per-request process forking, etc etc).  Now not only is your app not responsive, but you're incurring much more server load by making many requests.

The most common solution is some form of batching requests.  That's a valid solution.  However, if you batch too much data, now the client has to wait for all of the work to complete and all of the data to return.  Depending on how slow the server is, that can lead to a bad experience.  This also limits the responsiveness of the page to the slowest possible chunk of server work in that batch.

## One call, Consolidate Overhead, Concurrently Work, Stream Results

If the server can concurrently perform the fetches (with proper waits for state) and offer the results to the client as they are available, we consolidate the overhead and can still have a responsive client. 

Http chunked encoding is an http 1.1 feature.  Ebay leveraged and [even created a json pipe client](https://github.com/eBay/jsonpipe).  

We also need a concurrent server solution which writes the streaming chunks.

> Note: This is not a RESTful api.  This is a first party route per page-view

So for the "foobars" page, let's get all the data async.  Notice the overhead of the request and then the three pieces of data being made available to the client as the server has them.  Also notice that the third "chunk" of data relies on the first chunk (a foo has a baz id).

Foo and bar can also come out of order because they don't depend each other.  But baz will always come after Foo.  In this specific call, bar was offered to the client before foo.

```bash
curl --trace - "http://localhost:3000/foobars"
== Info:   Trying [::1]:3000...
== Info: Connected to localhost (::1) port 3000
=> Send header, 84 bytes (0x54)
0000: 47 45 54 20 2f 66 6f 6f 62 61 72 73 20 48 54 54 GET /foobars HTT
0010: 50 2f 31 2e 31 0d 0a 48 6f 73 74 3a 20 6c 6f 63 P/1.1..Host: loc
0020: 61 6c 68 6f 73 74 3a 33 30 30 30 0d 0a 55 73 65 alhost:3000..Use
0030: 72 2d 41 67 65 6e 74 3a 20 63 75 72 6c 2f 38 2e r-Agent: curl/8.
0040: 34 2e 30 0d 0a 41 63 63 65 70 74 3a 20 2a 2f 2a 4.0..Accept: */*
0050: 0d 0a 0d 0a                                     ....
<= Recv header, 17 bytes (0x11)
0000: 48 54 54 50 2f 31 2e 31 20 32 30 30 20 4f 4b 0d HTTP/1.1 200 OK.
0010: 0a                                              .
<= Recv header, 41 bytes (0x29)
0000: 43 6f 6e 74 65 6e 74 2d 54 79 70 65 3a 20 74 65 Content-Type: te
0010: 78 74 2f 70 6c 61 69 6e 3b 20 63 68 61 72 73 65 xt/plain; charse
0020: 74 3d 75 74 66 2d 38 0d 0a                      t=utf-8..
<= Recv header, 33 bytes (0x21)
0000: 58 2d 43 6f 6e 74 65 6e 74 2d 54 79 70 65 2d 4f X-Content-Type-O
0010: 70 74 69 6f 6e 73 3a 20 6e 6f 73 6e 69 66 66 0d ptions: nosniff.
0020: 0a                                              .
<= Recv header, 37 bytes (0x25)
0000: 44 61 74 65 3a 20 46 72 69 2c 20 30 33 20 4d 61 Date: Fri, 03 Ma
0010: 79 20 32 30 32 34 20 31 39 3a 31 37 3a 35 39 20 y 2024 19:17:59 
0020: 47 4d 54 0d 0a                                  GMT..
<= Recv header, 28 bytes (0x1c)
0000: 54 72 61 6e 73 66 65 72 2d 45 6e 63 6f 64 69 6e Transfer-Encodin
0010: 67 3a 20 63 68 75 6e 6b 65 64 0d 0a             g: chunked..
<= Recv header, 2 bytes (0x2)
0000: 0d 0a                                           ..
<= Recv data, 129 bytes (0x81)
0000: 37 62 0d 0a 7b 0a 20 20 22 69 64 65 6e 74 69 66 7b..{.  "identif
0010: 69 65 72 22 3a 20 22 62 61 72 22 2c 0a 20 20 22 ier": "bar",.  "
0020: 64 61 74 61 22 3a 20 7b 0a 20 20 20 20 22 69 64 data": {.    "id
0030: 22 3a 20 32 2c 0a 20 20 20 20 22 6d 65 73 73 61 ": 2,.    "messa
0040: 67 65 22 3a 20 22 62 61 72 20 6d 65 73 73 61 67 ge": "bar messag
0050: 65 20 66 6f 72 20 32 22 2c 0a 20 20 20 20 22 62 e for 2",.    "b
0060: 61 7a 69 64 22 3a 20 33 0a 20 20 7d 2c 0a 20 20 azid": 3.  },.  
0070: 22 65 72 72 6f 72 22 3a 20 22 22 0a 7d 0a 0a 0d "error": "".}...
0080: 0a                                              .
{
  "identifier": "bar",
  "data": {
    "id": 2,
    "message": "bar message for 2",
    "bazid": 3
  },
  "error": ""
}

<= Recv data, 109 bytes (0x6d)
0000: 36 37 0d 0a 7b 0a 20 20 22 69 64 65 6e 74 69 66 67..{.  "identif
0010: 69 65 72 22 3a 20 22 66 6f 6f 22 2c 0a 20 20 22 ier": "foo",.  "
0020: 64 61 74 61 22 3a 20 7b 0a 20 20 20 20 22 69 64 data": {.    "id
0030: 22 3a 20 31 2c 0a 20 20 20 20 22 74 69 74 6c 65 ": 1,.    "title
0040: 22 3a 20 22 66 6f 6f 20 74 69 74 6c 65 20 66 6f ": "foo title fo
0050: 72 20 31 22 0a 20 20 7d 2c 0a 20 20 22 65 72 72 r 1".  },.  "err
0060: 6f 72 22 3a 20 22 22 0a 7d 0a 0a 0d 0a          or": "".}....
{
  "identifier": "foo",
  "data": {
    "id": 1,
    "title": "foo title for 1"
  },
  "error": ""
}

<= Recv data, 118 bytes (0x76)
0000: 36 62 0d 0a 7b 0a 20 20 22 69 64 65 6e 74 69 66 6b..{.  "identif
0010: 69 65 72 22 3a 20 22 62 61 7a 22 2c 0a 20 20 22 ier": "baz",.  "
0020: 64 61 74 61 22 3a 20 7b 0a 20 20 20 20 22 69 64 data": {.    "id
0030: 22 3a 20 33 2c 0a 20 20 20 20 22 61 64 64 72 65 ": 3,.    "addre
0040: 73 73 22 3a 20 22 62 61 7a 20 61 64 64 72 65 73 ss": "baz addres
0050: 73 20 66 6f 72 20 33 22 0a 20 20 7d 2c 0a 20 20 s for 3".  },.  
0060: 22 65 72 72 6f 72 22 3a 20 22 22 0a 7d 0a 0a 0d "error": "".}...
0070: 0a 30 0d 0a 0d 0a                               .0....
{
  "identifier": "baz",
  "data": {
    "id": 3,
    "address": "baz address for 3"
  },
  "error": ""
}

== Info: Connection #0 to host localhost left intact
```

From the server, we want to make it productive to:

- Perform concurrent work safely
- Synchronize writing the chunks to the stream down to the client

The proof of concept code here offers a Chunked middleware and a concgroup go package to make concurrently doing work and synchronizing consistent.

In this sample code, we have a `/foobars` route for the `foobars` view.  It concurrently gets a foo, concurrently gets a bar, waits for both because we need state from a bar to get the final chunk, a baz.

```golang
	r := chi.NewRouter()
	r.Use(Chunked)

	r.Get("/foobars", func(w http.ResponseWriter, r *http.Request) {
		cw := NewChunkedWriter(w)

		cg, _ := concgroup.WithOptions(context.Background(), MAX_CONCURRENCY, func(identifier string, obj interface{}, err error) {
			cw.Send(NewChunk(identifier, obj, err))
		})

		cg.Go("foo", func() (interface{}, error) {
			return data.GetFoo(1)
		})

		bazId := 0
		cg.Go("bar", func() (interface{}, error) {
			bar, err := data.GetBar(2)
			if err != nil {
				return nil, err
			}
			bazId = bar.BazId
			return bar, err
		})

		cg.Wait()

		baz, err := data.GetBaz(bazId)
		cw.Send(NewChunk("baz", baz, err))
		// cw.Done()
	})
```

## Graceful Degradation

If the inital SSR "primary data" delivers the critical data for the view and we asynchronously get the secondary data (streaming and non blocking) then any of those chunks can fail (or we can circuit break) and the page can gracefully handle.

## Authentication and Authorization

The single request consolidates authentication.

For authorization, the initial SSR page (or request) can produce a signed JWT token with the resources and actors already authorized.  This should be an opaque and signed token to the client with a short (1 min?) TTL.

The subsequent AJAX request should supply the opaque token.  Ideally, a majority of the authorization checks should be short circuited without incurreing server work.  The token is signed, therefore trusted without calling the authz system. 

## Run

from server directory
```bash
go build && ./server
```

in another terminal:
```bash
curl --trace - "http://localhost:3000/foobars"
```

## Reference

https://instantdomains.com/engineering/streaming-json-jsons  
https://andrewjesaitis.com/2016/08/25/streaming-comparison/  
https://github.com/eBay/jsonpipe  
https://gist.github.com/smoya/bfedc2dd38d52d9bcd384fcb12e77d30  