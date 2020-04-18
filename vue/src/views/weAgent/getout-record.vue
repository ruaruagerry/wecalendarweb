<template>
    <div class="app-container">
        <div class="filter-container">
            <el-input v-model="listQuery.playerid"
                      placeholder="PlayerID"
                      style="width: 200px;"
                      class="filter-item" />
            <el-button v-waves
                       class="filter-item"
                       type="primary"
                       icon="el-icon-search"
                       @click="handleSearchPlayerRecord">
                Search
            </el-button>
        </div>

        <el-table :key="tableKey"
                  v-loading="listLoading"
                  :data="list"
                  border
                  fit
                  highlight-current-row>
            <el-table-column label="rid"
                             width="80"
                             align="center">
                <template slot-scope="{row}">
                    <span>{{ row.rid }}</span>
                </template>
            </el-table-column>
            <el-table-column label="id"
                             prop="id"
                             align="center"
                             width="80">
                <template slot-scope="{row}">
                    <span>{{ row.id }}</span>
                </template>
            </el-table-column>
            <el-table-column label="name"
                             prop="name"
                             align="center"
                             width="180">
                <template slot-scope="{row}">
                    <span>{{ row.name }}</span>
                </template>
            </el-table-column>
            <el-table-column label="getoutmoney"
                             prop="getoutmoney"
                             align="center"
                             width="280">
                <template slot-scope="{row}">
                    <span>{{ row.getoutmoney }}</span>
                </template>
            </el-table-column>
            <el-table-column label="createtime"
                             width="300"
                             align="center">
                <template slot-scope="{row}">
                    <span>{{ row.createtime | parseTime('{y}-{m}-{d} {h}:{i}:{s}') }}</span>
                </template>
            </el-table-column>
            <el-table-column label="resulttime"
                             width="300"
                             align="center">
                <template slot-scope="{row}">
                    <span v-if='row.resulttime != "0001-01-01T00:00:00Z"'>{{ row.resulttime | parseTime('{y}-{m}-{d} {h}:{i}:{s}') }}</span>
                    <span v-else>{{ "" }}</span>
                </template>
            </el-table-column>
            <el-table-column label="status"
                             class-name="status-col"
                             width="118">
                <template slot-scope="{row}">
                    <el-tag :type="row.status | statusFilter">
                        {{ statusNameFilter(row.status) }}
                    </el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Actions"
                             align="center"
                             width="330"
                             class-name="small-padding fixed-width">
                <template slot-scope="{row}">
                    <el-button v-if="row.status == 0"
                               size="mini"
                               type="success"
                               @click="handleModifyStatus(row, 2)">
                        同意
                    </el-button>
                    <el-button v-if="row.status == 0"
                               size="mini"
                               type="danger"
                               @click="handleModifyStatus(row, 1)">
                        拒绝
                    </el-button>
                </template>
            </el-table-column>
        </el-table>

        <pagination v-show="total>0"
                    :total="total"
                    :page.sync="listQuery.page"
                    :limit.sync="listQuery.limit"
                    @pagination="getList" />
    </div>
</template>

<script>
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { moneyGetoutRecord, moneyGetoutResult, moneyGetoutCount, moneyGetoutPlayerRecord } from '@/api/weagent'

export default {
    name: 'GetoutRecord',
    components: { Pagination },
    directives: { waves },
    filters: {
        statusFilter (status) {
            const statusMap = {
                0: 'info',
                1: 'danger',
                2: 'success',
                3: 'danger'
            }
            return statusMap[status]
        }
    },
    data () {
        return {
            tableKey: 0,
            list: null,
            total: 0,
            listLoading: true,
            listQuery: {
                page: 1,
                limit: 20,
                playerid: '',
            },
        }
    },
    created () {
        this.getList()
    },
    methods: {
        statusNameFilter (status) {
            const statusMap = {
                0: '审核中',
                1: '拒绝',
                2: '提现成功',
                3: '提现失败'
            }

            return statusMap[status]
        },
        getList () {
            this.listLoading = true

            moneyGetoutCount().then(response => {
                this.total = response.count
            })

            var start = (this.listQuery.page - 1) * this.listQuery.limit
            var end = start + this.listQuery.limit - 1

            var params = {
                start: start,
                end: end
            }

            var data = JSON.stringify(params)
            moneyGetoutRecord(data).then(response => {
                this.list = response.getoutrecords

                // Just to simulate the time of the request
                setTimeout(() => {
                    this.listLoading = false
                }, 1.5 * 1000)
            })
        },
        handleModifyStatus (row, status) {
            var params = {
                rid: row.rid,
                status: status
            }
            var data = JSON.stringify(params)

            moneyGetoutResult(data).then(response => {
                this.$message({
                    message: '操作Success',
                    type: 'success'
                })
                row.status = status
                row.resulttime = response.resulttime
            })
        },
        handleSearchPlayerRecord () {
            this.listLoading = true

            var start = (this.listQuery.page - 1) * this.listQuery.limit
            var end = start + this.listQuery.limit - 1

            var params = {
                playerid: this.listQuery.playerid,
                start: start,
                end: end,
            }
            var data = JSON.stringify(params)

            moneyGetoutPlayerRecord(data).then(response => {
                this.list = response.getoutrecords

                // Just to simulate the time of the request
                setTimeout(() => {
                    this.listLoading = false
                }, 1.5 * 1000)
            })
        },
    }
}
</script>
