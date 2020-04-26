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

export function divinationRecordGet (data) {
    return request({
        url: rootUrl + '/divination/record/get',
        method: 'post',
        data
    })
}

export function divinationRecordAdd (data) {
    return request({
        url: rootUrl + '/divination/record/add',
        method: 'post',
        data
    })
}

export function divinationRecordDel (data) {
    return request({
        url: rootUrl + '/divination/record/del',
        method: 'post',
        data
    })
}

export function divinationRecordSetBest (data) {
    return request({
        url: rootUrl + '/divination/record/setbest',
        method: 'post',
        data
    })
}

export function divinationGetBest (data) {
    return request({
        url: rootUrl + '/divination/getbest',
        method: 'post',
        data
    })
}

export function divinationRecordCount (data) {
    return request({
        url: rootUrl + '/divination/record/count',
        method: 'post',
        data
    })
}

export function divinationConfigFirstGet () {
    return request({
        url: rootUrl + '/divination/config/first/get',
        method: 'get'
    })
}

export function divinationConfigFirstSet (data) {
    return request({
        url: rootUrl + '/divination/config/first/set',
        method: 'post',
        data
    })
}

export function ClientConfigGet () {
    return request({
        url: rootUrl + '/client/config/get',
        method: 'get'
    })
}

export function ClientConfigSet (data) {
    return request({
        url: rootUrl + '/client/config/set',
        method: 'post',
        data
    })
}