#!/usr/bin/env python
# -*- coding: utf-8 -*
'''
海淀驾校自动约车脚本
'''
import requests
from pyquery import PyQuery
from PIL import Image
import xmltodict
import time

username = '232103198611011733'  # 登陆用户名
password = 'benchi'  # 登录密码
yuyue_date = '20161116'  # 预约日期
yuyue_time = '2002'  # 2001上午，2002下午，2003晚上

session = requests.session()
# 登录页
response = session.get('http://haijia.bjxueche.net')
pq = PyQuery(response.text)

# 登陆接口必填项
__EVENTVALIDATION = pq("#__EVENTVALIDATION").val()
__VIEWSTATEGENERATOR = pq("#__VIEWSTATEGENERATOR").val()
__VIEWSTATE = pq("#__VIEWSTATE").val()

# 验证码
response = session.get(
    'http://haijia.bjxueche.net/tools/CreateCode.ashx?key=ImgCode&random=' + str(int(time.time() * 1000)))
with open('captcha.jpg', 'wb') as f:
    f.write(response.content)
    f.close()
im = Image.open('captcha.jpg')
im.show()
im.close()
captcha = raw_input("请输入验证码:")  # 自动识别准确率低

# 登陆
data = {
    'txtUserName': username,
    'txtPassword': password,
    'txtIMGCode': captcha,
    'BtnLogin': '登  录',
    '__EVENTVALIDATION': __EVENTVALIDATION,
    '__VIEWSTATEGENERATOR': __VIEWSTATEGENERATOR,
    '__VIEWSTATE': __VIEWSTATE,
}
print data
ret = session.post('http://haijia.bjxueche.net/Login.aspx?LoginOut=true', data=data)
# ret.encoding = 'utf-8'
if ret.content.find('alert') != -1:
    print '登陆失败'

print '登陆成功'
while True:
    # 可约车的列表
    data = {
        "yyrq": yuyue_date,
        "yysd": yuyue_time,
        "xllxID": "2",
        "pageSize": 35,
        "pageNum": 1,
    }
    response = session.post('http://haijia.bjxueche.net/Han/ServiceBooking.asmx/GetCars', data=data)
    response = xmltodict.parse(response.text)
    car_list = response['string']['#text']
    car_list = car_list[:-2]
    car_list = eval(car_list, {'false': False, 'true': True, 'null': None})

    # print response
    # print car_list
    for car in car_list:
        # 约车
        data = {
            "yyrq": yuyue_date,
            "xnsd": yuyue_time,
            "cnbh": car['CNBH'],
            "imgCode": "",
            "KMID": "2",
        }
        response = session.post('http://haijia.bjxueche.net/Han/ServiceBooking.asmx/BookingCar', data=data)
        response = xmltodict.parse(response.text)
        ret = response['string']['#text']
        ret = eval(ret, {'false': False, 'true': True, 'null': None})
        if ret[0]['Result'] == True:
            print '预约成功'
            exit()

    print '查询中'
    time.sleep(5)
