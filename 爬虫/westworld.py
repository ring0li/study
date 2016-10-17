#!/usr/bin/env python
# -*- coding: utf-8 -*

import requests
from pyquery import PyQuery as pq
import time
import datetime
import json
import os
import sys


def get(i, e):
    e = pq(e)
    title = e.find('.name').text()
    pan = e.find('.links').find('a').eq(1).attr('href')
    print title
    if title.find('S01E03') != -1:
        send_mail_command = 'echo "' + pan + '" | mail -s "' + title + '" liuli@jindanlicai.com'
        os.system(send_mail_command.encode('utf-8'))
        print "%s %s %s" % (datetime.datetime.now(), title, pan)
        sys.exit(0)


url_login = 'http://www.zimuzu.tv/User/Login/ajaxLogin'
url_westworld = 'http://www.zimuzu.tv/gresource/33701'

s = requests.session()
data = {'account': 'bench1', 'password': 'benchi', 'remember': '1'}
resp = s.post(url_login, data=data)
ret = json.loads(resp.text)

print ret['info']
if ret['status'] == 1:
    while True:
        try:
            response = s.get(url_westworld)
            # Consider any status other than 2xx an error
            if not response.status_code // 100 == 2:
                print "Error: Unexpected response {}".format(response)
                continue
        except requests.exceptions.RequestException as e:
            # A serious problem happened, like an SSLError or InvalidURL
            print "Error: {}".format(e)
            continue

        d = pq(response.text)
        content = d(".view-res-list").find('li')
        content.each(get)

        time.sleep(60)
