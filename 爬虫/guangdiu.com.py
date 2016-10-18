#!/usr/bin/env python
# -*- coding: utf-8 -*

import requests
from pyquery import PyQuery as pq
import os
import time
import datetime
import configparser


def get_hotest(i, e):
    global hotest_body

    e = pq(e)
    no = e.find('.hourranknumtext').text()
    title = e.find('.hourrankimg').attr('alt')
    url = e.find('.hourrankimgdiv').attr('href')
    url = url[1:]

    hotest_body += "\n" + str(no) + "\n" + title + "\n" + domain + url


def send_email(title, body):
    print title
    print body
    print "\n\n"
    send_mail_command = 'echo "' + body + '" | mail -s "' + title + '" liuli@jindanlicai.com'
    os.system(send_mail_command.encode('utf-8'))


domain = u'http://guangdiu.com'
url_hotest = u'http://guangdiu.com/rank.php'

keywords = (u'午餐肉', u'元蹄', u'扣肉罐头', u'甘竹', u'hormel', u'罐头',
            u'生蚝', u'牛脆骨', u'冷冻生鲜',
            u'Mane‘n Tail',
            u'具良治',
            u'阿瑞娜',
            u'榴莲',
            u'露诗',
            u'nici',
            u'钢琴',
            u'竿', u'渔具', u'钓鱼', u'线组',
            )

config = configparser.ConfigParser()
config.read('guangdiu.ini')
next_hotest_hour = config['DEFAULT']['next_hotest_hour']
next_query_id = config['DEFAULT']['next_query_id']

print '继续执行id：', next_query_id
while True:
    r = requests.get(domain + '/detail.php', params={'id': next_query_id})
    r.encoding = 'utf-8'
    p = pq(r.text)
    url = r.url

    if r.text.find(u'不存在') != -1:
        # id可能不连续，需要判断是最新的
        r = requests.get(domain)
        r.encoding = 'utf-8'
        p = pq(r.text)
        href = p(".gooditem").find('.goodname').attr('href').lower()
        new_query_id = int(href[14:])
        # print new_id, id
        if new_query_id <= next_query_id:
            time.sleep(10)
            continue

    title = p(".dtitlelink").text().lower()
    for keyword in keywords:
        if title.find(keyword) != -1:
            body = u"匹配关键词：" + keyword + "\n" + url
            send_email(title, body)
            print "%s %s %s %s" % (datetime.datetime.now(), next_query_id, title, url)

    next_query_id = int(next_query_id) + 1
    config['DEFAULT']['next_query_id'] = str(next_query_id)
    with open('guangdiu.ini', 'w') as configfile:  # save
        config.write(configfile)



    current_Ymd = time.strftime("%Y%m%d")
    current_hour = time.strftime("%H")
    if not next_hotest_hour or int(next_hotest_hour) <= int(current_hour):
        r = requests.get(url_hotest, params={'date': current_Ymd, 'hour': next_hotest_hour})
        r.encoding = 'utf-8'
        p = pq(r.text)
        content = p('.hourrankitem')

        hotest_title = current_Ymd + str(next_hotest_hour)
        hotest_body = ''
        content.each(get_hotest)
        send_email(hotest_title, hotest_body)

        next_hotest_hour = int(next_hotest_hour) + 1
        config['DEFAULT']['next_hotest_hour'] = str(next_hotest_hour)
        with open('guangdiu.ini', 'w') as configfile:
            config.write(configfile)
