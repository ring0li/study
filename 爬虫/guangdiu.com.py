# -*- coding: utf-8 -*-

import requests
from pyquery import PyQuery as pq
import os
import time
import datetime

domain = 'http://guangdiu.com'
keywords = (u'鸭脖', u'午餐肉')

with open('id.txt', 'r') as f:
    id = int(f.read())
print id

while True:
    r = requests.get(domain + '/detail.php', params={'id': id})
    r.encoding = 'utf-8'
    p = pq(r.text)
    url = r.url
    title = p(".dtitlelink").text().lower()

    if r.text.find(u'不存在') != -1:
        r = requests.get(domain)
        r.encoding = 'utf-8'
        p = pq(r.text)
        href = p(".gooditem").find('.goodname').attr('href').lower()
        new_id = int(href[14:])
        # print new_id, id
        if new_id <= id:
            time.sleep(5)
            continue

    for keyword in keywords:
        if title.find(keyword) != -1:
            send_mail_command = 'echo "' + url + '" | mail -s "' + title + '" liuli@jindanlicai.com'
            os.system(send_mail_command.encode('utf-8'))

    print "%s %d %s" % (datetime.datetime.now(), id, title)

    id += 1
    with open('id.txt', 'w') as f:
        f.write(str(id))
