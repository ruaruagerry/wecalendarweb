<template>
    <div class="app-container">
        <el-date-picker v-model="time"
                        type="datetime"
                        format="yyyy-MM-dd"
                        placeholder="Release time" />
        <el-button type="primary"
                   @click="handleAddDivination">New Divination</el-button>

        <el-dialog :visible.sync="dialogVisible"
                   :title="'New Divination'">
            <el-form :model="divination"
                     label-width="80px"
                     label-position="left">
                <el-form-item label="name">
                    <el-input v-model="divination.name"
                              placeholder="name" />
                </el-form-item>
                <el-form-item label="portrait">
                    <el-input v-model="divination.portrait"
                              placeholder="portrait" />
                </el-form-item>
                <el-form-item label="content">
                    <el-input v-model="divination.content"
                              :autosize="{ minRows: 2, maxRows: 4}"
                              type="textarea"
                              placeholder="content" />
                </el-form-item>
            </el-form>
            <div style="text-align:right;">
                <el-button type="danger"
                           @click="dialogVisible=false">Cancel</el-button>
                <el-button type="primary"
                           @click="confirmDivination">Confirm</el-button>
            </div>
        </el-dialog>

        <el-table :key="tableKey"
                  v-loading="listLoading"
                  :data="list"
                  border
                  fit
                  highlight-current-row>
            <el-table-column label="id"
                             prop="id"
                             width="80"
                             align="center">
                <template slot-scope="{row}">
                    <span>{{ row.playerid }}</span>
                </template>
            </el-table-column>
            <el-table-column label="divinationid"
                             prop="divinationid"
                             align="center"
                             width="80">
                <template slot-scope="{row}">
                    <span>{{ row.divinationid }}</span>
                </template>
            </el-table-column>
            <el-table-column label="nickname"
                             prop="nickname"
                             align="center"
                             width="180">
                <template slot-scope="{row}">
                    <span>{{ row.nickname }}</span>
                </template>
            </el-table-column>
            <el-table-column label="content"
                             prop="content"
                             align="center"
                             width="280">
                <template slot-scope="{row}">
                    <span>{{ row.content }}</span>
                </template>
            </el-table-column>
            <el-table-column label="time"
                             width="300"
                             align="center">
                <template slot-scope="{row}">
                    <span>{{ row.time | parseTime('{y}-{m}-{d} {h}:{i}:{s}') }}</span>
                </template>
            </el-table-column>
            <el-table-column label="Actions"
                             align="center"
                             width="330"
                             class-name="small-padding fixed-width">
                <template slot-scope="{row}">
                    <el-button size="mini"
                               type="success"
                               @click="handleDelDivination(row)">
                        删除
                    </el-button>
                    <el-button size="mini"
                               type="danger"
                               @click="handleSetBest(row, 1)">
                        设置最佳
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
import { parseTime, dateToString } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { divinationRecordGet, divinationRecordDel, divinationRecordCount, divinationRecordSetBest, divinationRecordAdd } from '@/api/weagent'

export default {
    name: 'DivinationRecord',
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
            time: new Date(),

            dialogVisible: false,
            divination: {
                name: "",
                portrait: "",
                content: "",
            },

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
        handleAddDivination () {
            this.dialogVisible = true
        },
        confirmDivination () {
            console.log("!!time:", dateToString(this.time))

            // var params = {
            //     nowdata: "",
            //     content: this.divination.content,
            //     name: this.divination.name,
            //     portrait: this.divination.portrait,
            // }
            // var data = JSON.stringify(params)

            // divinationRecordAdd(data).then(response => {
            //     this.dialogVisible = false

            //     this.$message({
            //         message: '操作Success',
            //         type: 'success'
            //     })
            // })
        },
        getList () {
            this.listLoading = false
            // this.listLoading = true

            // divinationRecordCount().then(response => {
            //     this.total = response.count
            // })

            // var start = (this.listQuery.page - 1) * this.listQuery.limit
            // var end = start + this.listQuery.limit - 1

            // var params = {
            //     start: start,
            //     end: end
            // }

            // var data = JSON.stringify(params)
            // divinationRecordGet(data).then(response => {
            //     this.list = response.records

            //     // Just to simulate the time of the request
            //     setTimeout(() => {
            //         this.listLoading = false
            //     }, 1.5 * 1000)
            // })
        },
        handleDelDivination (row) {
            var params = {
                nowdata: "",
                divinationid: row.divinationid,
            }
            var data = JSON.stringify(params)

            divinationRecordDel(data).then(response => {
                this.$message({
                    message: '操作Success',
                    type: 'success'
                })
                row.status = status
                row.resulttime = response.resulttime
            })
        },
    }
}
</script>
