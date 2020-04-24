<template>
    <div class="app-container">
        <el-date-picker v-model="time"
                        type="datetime"
                        format="yyyy-MM-dd"
                        placeholder="Release time" />
        <el-button type="primary"
                   @click="handleAddDivination">新建吐槽</el-button>
        <el-button type="primary"
                   @click="handleBestDivination">最佳吐槽</el-button>

        <el-dialog :visible.sync="dialogVisible"
                   :title="'New Divination'">
            <el-form :model="divination"
                     label-width="80px"
                     label-position="left">
                <el-form-item v-if="!confirmVisible"
                              label="用户ID">
                    <el-input v-model="divination.playerid"
                              v-bind:disabled="true" />
                </el-form-item>
                <el-form-item label="昵称">
                    <el-input v-model="divination.name"
                              v-bind:disabled="!confirmVisible" />
                </el-form-item>
                <el-form-item label="头像">
                    <el-input v-model="divination.portrait"
                              v-bind:disabled="!confirmVisible" />
                </el-form-item>
                <el-form-item label="文本">
                    <el-input v-model="divination.content"
                              :autosize="{ minRows: 2, maxRows: 4}"
                              v-bind:disabled="!confirmVisible"
                              type="textarea" />
                </el-form-item>
            </el-form>
            <div style="text-align:right;">
                <el-button type="danger"
                           @click="dialogVisible=false">Cancel</el-button>
                <el-button type="primary"
                           v-if="confirmVisible"
                           @click="confirmDivination">Confirm</el-button>
            </div>
        </el-dialog>

        <el-table :key="tableKey"
                  v-loading="listLoading"
                  :data="list"
                  border
                  fit
                  highlight-current-row>
            <el-table-column label="用户ID"
                             prop="id"
                             width="100"
                             align="center">
                <template slot-scope="{row}">
                    <span>{{ row.playerid }}</span>
                </template>
            </el-table-column>
            <el-table-column label="吐槽ID"
                             prop="divinationid"
                             align="center"
                             width="100">
                <template slot-scope="{row}">
                    <span>{{ row.divinationid }}</span>
                </template>
            </el-table-column>
            <el-table-column label="昵称"
                             prop="nickname"
                             align="center"
                             width="120">
                <template slot-scope="{row}">
                    <span>{{ row.nickname }}</span>
                </template>
            </el-table-column>
            <el-table-column label="文本"
                             prop="content"
                             align="center"
                             width="740">
                <template slot-scope="{row}">
                    <span>{{ row.content }}</span>
                </template>
            </el-table-column>
            <el-table-column label="生成时间"
                             width="278"
                             align="center">
                <template slot-scope="{row}">
                    <span>{{ row.time | parseTime('{y}-{m}-{d} {h}:{i}:{s}') }}</span>
                </template>
            </el-table-column>
            <el-table-column label="操作"
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
                               @click="handleSetBest(row)">
                        最佳
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
import { divinationRecordGet, divinationRecordDel, divinationRecordCount, divinationRecordSetBest, divinationRecordAdd, divinationGetBest } from '@/api/weagent'

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

            confirmVisible: true,
            dialogVisible: false,
            divination: {
                playerid: "",
                name: "",
                portrait: "",
                content: "",
            },

            tableKey: 0,
            list: null,
            total: 0,
            listLoading: true,
            lastpage: 0,
            listQuery: {
                page: 1,
                limit: 10,
                playerid: '',
            },
        }
    },
    created () {
        this.getList()
    },
    methods: {
        // 添加吐槽
        handleAddDivination () {
            this.divination = {
                name: "",
                portrait: "",
                content: "",
            }
            this.dialogVisible = true
            this.confirmVisible = true
        },
        // 提交新增吐槽
        confirmDivination () {
            var params = {
                nowdata: dateToString(this.time),
                content: this.divination.content,
                name: this.divination.name,
                portrait: this.divination.portrait,
            }
            var data = JSON.stringify(params)

            divinationRecordAdd(data).then(response => {
                this.dialogVisible = false

                this.$message({
                    message: '操作Success',
                    type: 'success'
                })

                this.total += 1
                this.refreshList()
            })
        },
        // 获取吐槽总数
        async getCount () {
            var params = {
                nowdata: dateToString(this.time),
            }
            var data = JSON.stringify(params)
            divinationRecordCount(data).then(response => {
                this.total = response.count
            })
        },
        // 刷新界面
        refreshList () {
            this.listQuery.page = Math.floor((this.total - 1) / this.listQuery.limit) + 1
            this.getList()
        },
        // 获取吐槽列表
        getList () {
            this.listLoading = true

            this.getCount()

            var start = Math.floor(this.listQuery.page - 1) * this.listQuery.limit
            var end = start + this.listQuery.limit - 1

            var params = {
                nowdata: dateToString(this.time),
                start: start,
                end: end
            }
            var data = JSON.stringify(params)
            divinationRecordGet(data).then(response => {
                this.list = response.records

                this.lastpage = 1
                if (this.list != null) {
                    this.lastpage = (this.list.length / this.listQuery.limit) + 1
                }

                // Just to simulate the time of the request
                setTimeout(() => {
                    this.listLoading = false
                }, 1.5 * 1000)
            })
        },
        // 删除吐槽
        handleDelDivination (row) {
            var params = {
                nowdata: dateToString(this.time),
                divinationid: row.divinationid,
            }
            var data = JSON.stringify(params)

            divinationRecordDel(data).then(response => {
                this.$message({
                    message: '操作Success',
                    type: 'success'
                })
                this.total -= 1
                this.refreshList()
            })
        },
        // 设为当日最佳吐槽
        handleSetBest (row) {
            var params = {
                nowdata: dateToString(this.time),
                divinationid: row.divinationid,
            }
            var data = JSON.stringify(params)

            divinationRecordSetBest(data).then(response => {
                this.$message({
                    message: '操作Success',
                    type: 'success'
                })
            })
        },
        // 获取当日最佳吐槽
        handleBestDivination () {
            var params = {
                nowdata: dateToString(this.time),
            }
            var data = JSON.stringify(params)

            divinationGetBest(data).then(response => {
                this.$message({
                    message: '操作Success',
                    type: 'success'
                })

                this.divination = {
                    name: response.nickname,
                    portrait: response.portrait,
                    content: response.content,
                    playerid: response.playerid,
                }

                this.dialogVisible = true
                this.confirmVisible = false
            })
        }
    }
}
</script>
