import axios from 'axios'
import { MessageBox, Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'

function base64_decode (input) { // 解码，配合decodeURIComponent使用
    var base64EncodeChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/='
    var output = ''
    var chr1, chr2, chr3
    var enc1, enc2, enc3, enc4
    var i = 0
    input = input.replace(/[^A-Za-z0-9\+\/\=]/g, '')
    while (i < input.length) {
        enc1 = base64EncodeChars.indexOf(input.charAt(i++))
        enc2 = base64EncodeChars.indexOf(input.charAt(i++))
        enc3 = base64EncodeChars.indexOf(input.charAt(i++))
        enc4 = base64EncodeChars.indexOf(input.charAt(i++))
        chr1 = (enc1 << 2) | (enc2 >> 4)
        chr2 = ((enc2 & 15) << 4) | (enc3 >> 2)
        chr3 = ((enc3 & 3) << 6) | enc4
        output = output + String.fromCharCode(chr1)
        if (enc3 != 64) {
            output = output + String.fromCharCode(chr2)
        }
        if (enc4 != 64) {
            output = output + String.fromCharCode(chr3)
        }
    }
    return utf8_decode(output)
}

function utf8_decode (utftext) { // utf-8解码
    var string = ''
    let i = 0
    let c = 0
    let c1 = 0
    let c2 = 0
    while (i < utftext.length) {
        c = utftext.charCodeAt(i)
        if (c < 128) {
            string += String.fromCharCode(c)
            i++
        } else if ((c > 191) && (c < 224)) {
            c1 = utftext.charCodeAt(i + 1)
            string += String.fromCharCode(((c & 31) << 6) | (c1 & 63))
            i += 2
        } else {
            c1 = utftext.charCodeAt(i + 1)
            c2 = utftext.charCodeAt(i + 2)
            string += String.fromCharCode(((c & 15) << 12) | ((c1 & 63) << 6) | (c2 & 63))
            i += 3
        }
    }
    return string
}

// create an axios instance
const service = axios.create({
    baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
    // withCredentials: true, // send cookies when cross-domain requests
    timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
    config => {
        // do something before request is sent

        if (store.getters.token) {
            // let each request carry token
            config.headers['Session'] = getToken()
        }
        return config
    },
    error => {
        // do something with request error
        console.log(error) // for debug
        return Promise.reject(error)
    }
)

// response interceptor
service.interceptors.response.use(
    /**
         * If you want to get http information such as headers or status
         * Please return  response => response
        */

    /**
         * Determine the request status by custom code
         * Here is just an example
         * You can also judge the status by HTTP Status Code
         */
    response => {
        const res = response.data
        const requrl = response.config.url

        if (res.result !== 0) {
            Message({
                message: res.msg || 'Error',
                type: 'error',
                duration: 5 * 1000
            })

            // 50008: Illegal token; 50012: Other clients logged in; 50014: Token expired;
            // if (res.result === 9 || res.result === 10 || res.result === 11) {
            //     // to re-login
            //     MessageBox.confirm('You have been logged out, you can cancel to stay on this page, or log in again', 'Confirm logout', {
            //         confirmButtonText: 'Re-Login',
            //         cancelButtonText: 'Cancel',
            //         type: 'warning'
            //     }).then(() => {
            //         store.dispatch('user/resetToken').then(() => {
            //             location.reload()
            //         })
            //     })
            // }

            return Promise.reject(new Error(res.msg || 'Error'))
        } else {
            // todo login
            if (requrl !== '/dev-api/user/logout') {
                var jsondata = new Object()
                if (res.data != undefined) {
                    res.data = base64_decode(res.data)
                    jsondata = JSON.parse(res.data)
                }

                return jsondata
            }

            return res
        }
    },
    error => {
        console.log('err' + error) // for debug
        Message({
            message: error.message,
            type: 'error',
            duration: 5 * 1000
        })
        return Promise.reject(error)
    }
)

export default service
