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
            path: 'dynamic-table',
            component: () => import('@/views/table/dynamic-table/index'),
            name: 'DynamicTable',
            meta: { title: 'Dynamic Table' }
        },
    ]
}
export default weagentRouter
