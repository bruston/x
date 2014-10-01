Simple shellshock vulnerability test.

<pre>
Usage of ./shocked (and its defaults):
  -stdin=false: If true, a list of URLs to check is read from Stdin.
  -timeout=10: How long to wait for a response before giving up.
  -url="": The URL to check.
  -verbose=false: If true, the response body will be printed.
  -workers=5: The number of concurrent requests to make.
</pre>

Example:
```bash
./shocked -url="http://dev.lan/cgi-bin/poc.cgi" -verbose
cat urls.txt | ./shocked -stdin -workers=3
```
