# zap-baseline.py の hook。spider は GET しか辿らないので POST /items を明示的に送る
import os
from urllib.parse import urlparse

def zap_started(zap, target):
    body = '{"name":"zap"}'
    res = zap.core.send_request(
        f"POST {target.rstrip('/')}/items HTTP/1.1\r\n"
        f"Host: {urlparse(target).netloc}\r\n"
        f"Authorization: {os.environ.get('ZAP_AUTH_HEADER_VALUE', '')}\r\n"
        "Content-Type: application/json\r\n"
        f"Content-Length: {len(body)}\r\n\r\n"
        + body
    )
    print("hook POST /items ->", res[0]["responseHeader"].split("\r\n")[0])
