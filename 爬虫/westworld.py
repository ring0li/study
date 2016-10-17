#!/usr/bin/env python
# -*- coding: utf-8 -*

import requests
from pyquery import PyQuery as pq
import time
import datetime
import json
import os

url_login = 'http://www.zimuzu.tv/User/Login/ajaxLogin'
url_westworld = 'http://www.zimuzu.tv/gresource/33701'

s = requests.session()

data = {'account': 'bench1', 'password': 'benchi', 'remember': '1'}
resp = s.post(url_login, data=data)
ret = json.loads(resp.text)


def get(i, e):
    e = pq(e)
    title = e.find('.name').text()
    pan = e.find('.links').find('a').eq(3).attr('href')
    if title.find('S01E03') != -1:
        send_mail_command = 'echo "' + pan + '" | mail -s "' + title + '" liuli@jindanlicai.com'
        os.system(send_mail_command.encode('utf-8'))
        print "%s %s %s" % (datetime.datetime.now(), title, pan)


print ret['info']
if ret['status'] == 1:
    while True:
        resp = s.get(url_westworld)

        d = pq(resp.text)
        content = d(".view-res-list").find('li')
        content.each(get)

        time.sleep(60)
