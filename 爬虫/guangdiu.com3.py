#!/usr/bin/env python3
# -*- coding: utf-8 -*

import requests
from pyquery import PyQuery as pq
import os
import time
from datetime import datetime
from datetime import timedelta
import configparser
import pymysql

from email.header import Header
from email.mime.text import MIMEText
from email.utils import parseaddr, formataddr

import smtplib


def get_keyword():
    conn = pymysql.connect(host='localhost',
                           user='root',
                           passwd='benchi',
                           db='liuli',
                           charset='utf8', )
    cursor = conn.cursor()
    cursor.execute("select * from guangdiu_keyword where is_del=0")
    data = cursor.fetchall()
    keywords = []
    for row in data:
        keywords.append(row[1])
    return keywords


def get_hotest(i, e):
    global hotest_body

    e = pq(e)
    no = e.find('.hourranknumtext').text()
    title = e.find('.hourrankimg').attr('alt')
    url = e.find('.hourrankimgdiv').attr('href')
    url = url[1:]

    hotest_body += "<br>" + "<a href='" + domain + url + "'>" + str(no) + title + "</a>"


def _format_addr(s):
    name, addr = parseaddr(s)
    return formataddr((Header(name, 'utf-8').encode(), addr))


def send_email(title, body):
    from_addr = 'alert@jindanlicai.com'
    password = 'e9E-Aek-V75-ewd'
    to_addr = 'liuli@jindanlicai.com'
    smtp_server = 'smtp.exmail.qq.com'

    msg = MIMEText(body, 'html', 'utf-8')
    msg['From'] = _format_addr('逛丢 <%s>' % from_addr)
    msg['To'] = _format_addr('主人 <%s>' % to_addr)
    msg['Subject'] = Header(title, 'utf-8').encode()

    server = smtplib.SMTP(smtp_server, 25)
    # server.set_debuglevel(1)
    server.login(from_addr, password)
    server.sendmail(from_addr, [to_addr], msg.as_string())
    server.quit()


domain = u'http://guangdiu.com'
url_hotest = u'http://guangdiu.com/rank.php'
keywords = []

config = configparser.ConfigParser()
config.read('guangdiu.ini')
next_query_id = int(config['DEFAULT']['next_query_id'])
next_hotest_YmdH = config['DEFAULT']['next_hotest_YmdH']

print('继续执行id：', next_query_id)
while True:
    if (datetime.now().second == 0 or len(keywords) == 0):
        keywords = get_keyword()
        print(keywords)

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
            time.sleep(1)
            continue

    title = p(".dtitlelink").text().lower()
    for keyword in keywords:
        if title.find(keyword) != -1:
            body = u"匹配关键词：" + keyword + "<br>" + "<a href='" + url + "'>" + title + "</a>"
            print("%s %s %s %s" % (datetime.now(), next_query_id, title, url))
            send_email(title, body)

    next_query_id = int(next_query_id) + 1
    config['DEFAULT']['next_query_id'] = str(next_query_id)
    with open('guangdiu.ini', 'w') as configfile:  # save
        config.write(configfile)

    current_Ymd = time.strftime("%Y%m%d")
    current_YmdH = time.strftime("%Y%m%d%H")

    next_hotest_Ymd = time.strftime('%Y%m%d', time.strptime(next_hotest_YmdH, "%Y%m%d%H"))
    next_hotest_H = time.strftime('%H', time.strptime(next_hotest_YmdH, "%Y%m%d%H"))

    if next_hotest_YmdH <= current_YmdH:
        data = {'date': next_hotest_Ymd, 'hour': next_hotest_H}
        r = requests.get(url_hotest, params=data)
        r.encoding = 'utf-8'
        p = pq(r.text)
        content = p('.hourrankitem')

        hotest_body = ''
        content.each(get_hotest)
        if hotest_body or current_YmdH >= next_hotest_YmdH:
            send_email(next_hotest_YmdH, hotest_body)
            datetime_next_hour = datetime.strptime(next_hotest_YmdH, "%Y%m%d%H") + timedelta(hours=1)
            next_hotest_YmdH = datetime_next_hour.strftime('%Y%m%d%H')
            config['DEFAULT']['next_hotest_YmdH'] = str(next_hotest_YmdH)
            with open('guangdiu.ini', 'w') as configfile:
                config.write(configfile)
