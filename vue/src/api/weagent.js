import request from '@/utils/request'

var rootUrl = "http://127.0.0.1:3004"

export const GetoutStatusReview = 0  // 审核中
export const GetoutStatusRefused = 1 // 审核拒绝
export const GetoutStatusSuccess = 2 // 提现成功
export const GetoutStatusFailed = 3  // 提现失败

export function moneyGetoutRecord (data) {
    return request({
        url: rootUrl + '/money/getout/record',
        method: 'post',
        data
    })
}

export function moneyGetoutResult (data) {
    return request({
        url: rootUrl + '/money/getout/result',
        method: 'post',
        data
    })
}

export function moneyGetoutPlayerRecord (data) {
    return request({
        url: rootUrl + '/money/getout/playerrecord',
        method: 'post',
        data
    })
}

export function moneyGetoutCount () {
    return request({
        url: rootUrl + '/money/getout/count',
        method: 'get',
    })
}

export function login (data) {
    return request({
        url: rootUrl + '/auth/web/login',
        method: 'post',
        data
    })
}

export function getInfo () {
    return request({
        url: rootUrl + '/auth/web/getinfo',
        method: 'get'
    })
}

