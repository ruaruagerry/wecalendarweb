/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const weagentRouter = {
    path: '/weagent',
    component: Layout,
    redirect: '/weagent/getout-record',
    name: 'Weagent',
    meta: {
        title: 'Weagent',
        icon: 'table'
    },
    children: [
        {
            path: 'getout-record',
            component: () => import('@/views/weAgent/getout-record'),
            name: 'GetoutRecord',
            meta: { title: '收益审核' }
        },
        {
            path: 'divination-record',
            component: () => import('@/views/weAgent/divination-record'),
            name: 'DivinationRecord',
            meta: { title: '吐槽列表' }
        },
        {
            path: 'client-config',
            component: () => import('@/views/weAgent/client-config'),
            name: 'ClientConfig',
            meta: { title: '客户端配置' }
        },
    ]
}
export default weagentRouter
