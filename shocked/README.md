Simple shellshock vulnerability test.

<pre>
Usage of ./shocked:
  -url="": The URL to check.
  -verbose=false: If true, the response body will be printed.
</pre>

Example:
```bash
./shocked -url="http://dev.lan/cgi-bin/poc.cgi" -verbose
```
