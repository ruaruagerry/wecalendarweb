import request from '@/utils/request'

export function logout () {
    return request({
        url: '/user/logout',
        method: 'post'
    })
}
