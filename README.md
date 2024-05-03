# GO-CONC

Musings of streaming json over chunked encoding

from server directory
```bash
go build && ./server
```

in another terminal:
```bash
curl --trace - "http://localhost:3000/foobars"
```

### Reference

https://instantdomains.com/engineering/streaming-json-jsons
https://andrewjesaitis.com/2016/08/25/streaming-comparison/
https://github.com/eBay/jsonpipe
https://gist.github.com/smoya/bfedc2dd38d52d9bcd384fcb12e77d30