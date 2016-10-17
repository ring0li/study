#!/usr/bin/env python
# -*- coding: utf-8 -*

import requests
from pyquery import PyQuery as pq
import os
import time
import datetime

domain = 'http://guangdiu.com'
keywords = (u'午餐肉', u'元蹄', u'扣肉罐头', u'甘竹', u'hormel',  # u'罐头',
            u'生蚝', u'牛脆骨', u'冷冻生鲜',
            u'Mane‘n Tail',
            u'具良治',
            u'阿瑞娜',
            u'榴莲',
            u'露诗',
            u'nici',
            # u'内裤',
            u'钢琴',  # 钢琴
            u'竿', u'渔具', u'钓鱼', u'线组',
            )

with open('id.txt', 'r') as f:
    id = int(f.read())
print '继续执行id：', id

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
            body = u"匹配关键词：" + keyword + "\n" + url
            send_mail_command = 'echo "' + body + '" | mail -s "' + title + '" liuli@jindanlicai.com'
            os.system(send_mail_command.encode('utf-8'))
            print "%s %d %s %s" % (datetime.datetime.now(), id, title, url)

    id += 1
    with open('id.txt', 'w') as f:
        f.write(str(id))
