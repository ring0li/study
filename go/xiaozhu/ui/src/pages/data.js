import axios from 'axios';
import qs from 'qs'

var apiUrl = "http://wx.iiwoo.com";

if(process.env.NODE_ENV=="development"){
    // apiUrl = "http://127.0.0.1";
}

axios.defaults.baseURL = apiUrl;
axios.defaults.withCredentials = true;
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';


export function fetch(url, params) { // 封装axios的post请求
    return new Promise((resolve, reject) => { // promise 用法,自行查阅
        axios.post(url, qs.stringify(params)).then(response => {
            resolve(response.data) // promise相关
        }).catch(error => {
            reject(error) // promise相关
        })
    })
}

export default {
    getTipData: (params) => {
        return fetch('/switchset/GetMsgReceiveSwitch', params)
    },
    setTipData: (params) => {
        return fetch('/switchset/SetMsgReceiveSwitch', params)
    },
    //
    getFriendTipData: (params) => {
        return fetch('/switchset/GetFriendList', params)
    },
    setFriendTipData: (params) => {
        return fetch('/switchset/SetFriendMsgReceiveSwitch', params)
    },
    delFriend: (params) => {
        return fetch('/switchset/DelFriend', params)
    },
    getWaterTipData: (params) => {
        return fetch('/switchset/GetWaterSwitch', params)
    },
    setWaterTipData: (params) => {
        return fetch('/switchset/SetWaterSwitch', params)
    },  
}


