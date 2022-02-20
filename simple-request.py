import http.client
import sys

if len(sys.argv) != 3 and len(sys.argv) != 4:
	print('Usage: {} hostname resource [POST | DELETE]')
	print('')
	print('Send a HTTP request to the given hostname (e.g., localhost:8000)')
	print('at the given resource (e.g., /tracker/some-game/bool-value).')
	print('')
	print('If a method isn\'t specified, "POST" is assumed.')
	exit()

host = sys.argv[1]
res = sys.argv[2]
method = 'POST'
if len(sys.argv) == 4:
	method = sys.argv[3]

conn = http.client.HTTPConnection(host)
conn.request(method, res)
r = conn.getresponse()
print(r.status, r.reason)
